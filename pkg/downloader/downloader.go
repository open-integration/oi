package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"runtime"

	"github.com/open-integration/core/pkg/logger"
)

type (
	// Downloader to use to fetch services from catalog
	Downloader interface {
		Download(name string, version string) error
		Store() string
	}

	// Options to create new Downloader
	Options struct {
		Store  string
		Logger logger.Logger
	}

	downloader struct {
		store  string
		logger logger.Logger
	}
)

// New creates new Downloader
func New(opt Options) Downloader {
	return &downloader{
		store:  opt.Store,
		logger: opt.Logger,
	}
}

func (d *downloader) Download(name string, version string) error {
	url := fmt.Sprintf("https://github.com/open-integration/core-services/releases/download/%s/%s-%s-%s-%s", version, name, version, runtime.GOOS, runtime.GOARCH)
	d.logger.Debug("Downloading", "url", url)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("%s:%s not found", name, version)
	}

	fileName := path.Join(d.Store(), name)

	// Create the file
	out, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	// Add prmission to exec file
	err = os.Chmod(fileName, os.ModePerm)
	return err
}

func (d *downloader) Store() string {
	return d.store
}
