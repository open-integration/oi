package service

import (
	"context"
	"io"

	v1 "github.com/open-integration/core/pkg/api/v1"
	"github.com/open-integration/core/pkg/logger"
	"google.golang.org/grpc"
)

const (
	// Local type using gGRPC to connect to
	Local = "local"

	// Kubernetes type
	Kubernetes = "kubernetes"
)

type (
	// Service expose an interface to manage services
	Service interface {
		Caller
		RunnerKiller
		API
	}

	// Caller to call the service
	Caller interface {
		Call(context context.Context, req *v1.CallRequest) (*v1.CallResponse, error)
	}
	// RunnerKiller to start and stop the the service
	RunnerKiller interface {
		Run() error
		Kill() error
	}

	// API expose general methods
	API interface {
		Schemas() map[string]string // TODO: cache!
	}

	// Options shows all the available options to build runner
	Options struct {
		Type                 string
		Name                 string
		ID                   string
		Version              string
		Logger               logger.Logger
		Dailer               dialer
		PortGenerator        portGenerator
		ServiceClientCreator serviceClientCreator
		LogsDirectory        string

		// Local runner options
		LocalLogFileCreator logFileCreator
		LocalCommandCreator cmdCreator
		LocalPathToBinary   string

		// Kubernetes runner options
		KubernetesKubeConfigPath   string
		KubernetesContext          string
		KubernetesNamespace        string
		KubeconfigHost             string
		KubeconfigB64Crt           string
		KubeconfigToken            string
		Kube                       kube
		KubernetesGrpcDialViaPodIP bool
		KubernetesVolumeClaimName  string
		KubernetesVolumeName       string
	}

	dialer interface {
		Dial(target string, opts ...grpc.DialOption) (*grpc.ClientConn, error)
	}

	logFileCreator interface {
		Create(dir string, name string) (io.Writer, error)
	}

	serviceClientCreator interface {
		New(cc *grpc.ClientConn) v1.ServiceClient
	}

	portGenerator interface {
		GetAvailable() (string, error)
	}
)

// New builds new runner based on Options.Type
func New(opt *Options) Service {
	if opt.Type == Local {
		return &localService{
			Logger:               opt.Logger,
			name:                 opt.Name,
			id:                   opt.ID,
			logFileCreator:       opt.LocalLogFileCreator,
			logsFileDirectory:    opt.LogsDirectory,
			serviceClientCreator: opt.ServiceClientCreator,
			portGenerator:        opt.PortGenerator,
			dialer:               opt.Dailer,
			cmdCreator:           opt.LocalCommandCreator,
			path:                 opt.LocalPathToBinary,
		}
	}

	if opt.Type == Kubernetes {
		runner := &kubernetesService{
			Logger:               opt.Logger,
			name:                 opt.Name,
			version:              opt.Version,
			id:                   opt.ID,
			kubeconfigPath:       opt.KubernetesKubeConfigPath,
			kubeconfigNamespace:  opt.KubernetesNamespace,
			kubeconfigContext:    opt.KubernetesContext,
			kubeconfigHost:       opt.KubeconfigHost,
			kubeconfigToken:      opt.KubeconfigToken,
			kubeconfigB64Crt:     opt.KubeconfigB64Crt,
			kube:                 opt.Kube,
			dialer:               opt.Dailer,
			portGenerator:        opt.PortGenerator,
			serviceClientCreator: opt.ServiceClientCreator,
			hostname:             "localhost",
			grpcDialViaPodIP:     opt.KubernetesGrpcDialViaPodIP,
			logsDirectory:        opt.LogsDirectory,
			kubeVolumeClaimName:  opt.KubernetesVolumeClaimName,
			kubeVolumeName:       opt.KubernetesVolumeName,
		}

		return runner
	}
	return nil
}
