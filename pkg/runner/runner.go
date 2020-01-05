package runner

import (
	"context"
	"io"

	v1 "github.com/open-integration/core/pkg/api/v1"
	"github.com/open-integration/core/pkg/logger"
	"google.golang.org/grpc"
)

const (
	// LocalRunner type
	LocalRunner = "local"

	// KubernetesRunner type
	KubernetesRunner = "kubernetes"
)

type (
	// Runner expose an interface to run services
	Runner interface {
		Run() error
		Kill() error
		Call(context context.Context, req *v1.CallRequest) (*v1.CallResponse, error)
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

		// Local runner options
		LocalLogFileCreator logFileCreator
		LocalLogsDirectory  string
		LocalCommandCreator cmdCreator
		LocalPathToBinary   string

		// Kubernetes runner options
		KubernetesKubeConfigPath string
		KubernetesContext        string
		KubernetesNamespace      string
		Kube                     kube
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
func New(opt *Options) Runner {
	if opt.Type == LocalRunner {
		return &localRunner{
			Logger:               opt.Logger,
			name:                 opt.Name,
			id:                   opt.ID,
			logFileCreator:       opt.LocalLogFileCreator,
			logsFileDirectory:    opt.LocalLogsDirectory,
			serviceClientCreator: opt.ServiceClientCreator,
			portGenerator:        opt.PortGenerator,
			dialer:               opt.Dailer,
			cmdCreator:           opt.LocalCommandCreator,
			path:                 opt.LocalPathToBinary,
		}
	}

	if opt.Type == KubernetesRunner {
		return &kubernetesRunner{
			Logger:               opt.Logger,
			name:                 opt.Name,
			version:              opt.Version,
			id:                   opt.ID,
			kubeconfigPath:       opt.KubernetesKubeConfigPath,
			kubeconfigNamespace:  opt.KubernetesNamespace,
			kubeconfigContext:    opt.KubernetesContext,
			kube:                 opt.Kube,
			dialer:               opt.Dailer,
			portGenerator:        opt.PortGenerator,
			serviceClientCreator: opt.ServiceClientCreator,
		}
	}
	return nil
}
