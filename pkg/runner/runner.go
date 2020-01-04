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
		Type                      string
		Name                      string
		ID                        string
		Logger                    logger.Logger
		LocalDailer               dialer
		LocalLogFileCreator       logFileCreator
		LocalServiceClientCreator serviceClientCreator
		LocalPortGenerator        portGenerator
		LocalLogsDirectory        string
		LocalCommandCreator       cmdCreator
		LocalPathToBinary         string
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
			serviceClientCreator: opt.LocalServiceClientCreator,
			portGenerator:        opt.LocalPortGenerator,
			dialer:               opt.LocalDailer,
			cmdCreator:           opt.LocalCommandCreator,
			path:                 opt.LocalPathToBinary,
		}
	}

	if opt.Type == KubernetesRunner {
		return &kubernetesRunner{}
	}
	return nil
}
