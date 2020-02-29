package utils

import (
	b64 "encoding/base64"
	"fmt"
	"strconv"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type (
	// Kubernetes expose abilities on run on kube cluster
	Kubernetes struct{}

	BuildKubeClientOptions struct {
		KubeconfigPath string
		Host           string
		B64CRT         string
		Token          string
	}
)

// BuildClient returns a kubernetes client based on path to kubeconfig
func (_k *Kubernetes) BuildClient(opt BuildKubeClientOptions) (*kubernetes.Clientset, error) {
	if opt.Host != "" && opt.B64CRT != "" && opt.Token != "" {
		ca, err := b64.StdEncoding.DecodeString(opt.B64CRT)
		if err != nil {
			return nil, err
		}
		return kubernetes.NewForConfig(&rest.Config{
			Host:        opt.Host,
			BearerToken: opt.Token,
			TLSClientConfig: rest.TLSClientConfig{
				CAData: ca,
			},
		})
	}
	var config *rest.Config
	if opt.KubeconfigPath == "" {
		var err error
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	} else {
		var err error
		config, err = clientcmd.BuildConfigFromFlags("", opt.KubeconfigPath)
		if err != nil {
			return nil, err
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}

// BuildClient returns a kubernetes client based on path to kubeconfig
func (_k *Kubernetes) BuildClientWithToken(host string, ca string, token string) (*kubernetes.Clientset, error) {
	config := &rest.Config{
		Host: host,
		TLSClientConfig: rest.TLSClientConfig{
			CAData: []byte(ca),
		},
		BearerToken: token,
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}

func (_k *Kubernetes) BuildPodDefinition(namespace string, name string, version string, id string, port string, volume string, volumeClaimName string, volumeMoutnPath string) (*v1.Pod, error) {
	if version == "" {
		version = "latest"
	}
	portInt, err := strconv.Atoi(port)
	if err != nil {
		return nil, err
	}
	p := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-%s", name, id),
			Namespace: namespace,
			Labels: map[string]string{
				"id": id,
			},
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				v1.Container{
					Name:            name,
					Image:           fmt.Sprintf("openintegration/service_catalog-%s:%s", name, version),
					ImagePullPolicy: v1.PullAlways,
					Ports: []v1.ContainerPort{
						v1.ContainerPort{
							Name:          "http",
							ContainerPort: int32(portInt),
							Protocol:      v1.ProtocolTCP,
						},
					},
					Env: []v1.EnvVar{
						v1.EnvVar{
							Name:  "PORT",
							Value: port,
						},
					},
				},
			},
		},
	}
	if volume != "" {
		p.Spec.Volumes = []v1.Volume{
			v1.Volume{
				Name: volume,
				VolumeSource: v1.VolumeSource{
					PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
						ClaimName: volumeClaimName,
					},
				},
			},
		}
		p.Spec.Containers[0].VolumeMounts = []v1.VolumeMount{
			v1.VolumeMount{
				Name:      volume,
				MountPath: volumeMoutnPath,
			},
		}
	}
	return p, nil
}

func (_k Kubernetes) BuildServiceDefinition(namespace string, name string, id string, port string, serviceType string) (*v1.Service, error) {
	t := v1.ServiceTypeClusterIP

	if serviceType == "LoadBalancer" {
		t = v1.ServiceTypeLoadBalancer
	}

	portInt, err := strconv.Atoi(port)
	if err != nil {
		return nil, err
	}
	return &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-%s", name, id),
			Namespace: namespace,
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{
				v1.ServicePort{
					Name:     "http",
					Protocol: v1.ProtocolTCP,
					Port:     int32(portInt),
					TargetPort: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: 8080,
					},
				},
			},
			Type: t,
			Selector: map[string]string{
				"id": id,
			},
		},
	}, nil
}

// CreatePod applies a pod definitions on given client
func (_k Kubernetes) CreatePod(client *kubernetes.Clientset, def *v1.Pod) (*v1.Pod, error) {
	ns := "default"
	if def.ObjectMeta.Namespace != "" {
		ns = def.ObjectMeta.Namespace
	}
	return client.CoreV1().Pods(ns).Create(def)
}

// WaitForPod waits til pod reaches given phase
func (_k Kubernetes) WaitForPod(client *kubernetes.Clientset, pod *v1.Pod, phase string) error {
	ns := "default"
	if pod.ObjectMeta.Namespace != "" {
		ns = pod.ObjectMeta.Namespace
	}
	w, err := client.CoreV1().Pods(ns).Watch(metav1.ListOptions{
		TypeMeta: metav1.TypeMeta{
			Kind: "Pod",
		},
	})
	if err != nil {
		return err
	}
	defer w.Stop()
	stopChan := make(chan bool)
	go func() {
		time.Sleep(30 * time.Second)
		stopChan <- true
	}()
	for {
		select {
		case ev := <-w.ResultChan():
			if ev.Object == nil {
				continue
			}

			p, ok := ev.Object.(*v1.Pod)
			if !ok {
				continue
			}
			if p.ObjectMeta.GetUID() == pod.ObjectMeta.GetUID() {
				if pod.Status.Phase == v1.PodRunning {
					return nil
				}
			}
		case <-stopChan:
			return nil
		}
	}
}

// CreateService applies a pod definitions on given client
func (_k Kubernetes) CreateService(client *kubernetes.Clientset, def *v1.Service) (*v1.Service, error) {
	ns := "default"
	if def.ObjectMeta.Namespace != "" {
		ns = def.ObjectMeta.Namespace
	}
	return client.CoreV1().Services(ns).Create(def)
}

// KillService deletes kubernetes service
func (_k Kubernetes) KillService(client *kubernetes.Clientset, namespace string, name string) error {
	return client.CoreV1().Services(namespace).Delete(name, nil)
}

// KillPod deletes kubernetes service
func (_k Kubernetes) KillPod(client *kubernetes.Clientset, namespace string, name string) error {
	return client.CoreV1().Pods(namespace).Delete(name, nil)
}
