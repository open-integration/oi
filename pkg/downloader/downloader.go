package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"runtime"

	"github.com/open-integration/oi/pkg/logger"
)

type (
	// Downloader to use to fetch services from catalog
	Downloader interface {
		Download(name string, version string) (string, error)
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

const base = "https://storage.googleapis.com/open-integration-service-catalog"

// New creates new Downloader
func New(opt Options) Downloader {
	return &downloader{
		store:  opt.Store,
		logger: opt.Logger,
	}
}

func (d *downloader) Download(name string, version string) (string, error) {
	var err error
	candidateFileName := fmt.Sprintf("%s-%s-%s-%s", name, version, runtime.GOOS, runtime.GOARCH)
	fullPath := path.Join(d.store, candidateFileName)
	if fileExists(fullPath) {
		d.logger.Info("Skipping download, service exist", "path", fullPath)
		return fullPath, nil
	}

	d.logger.Info("Downloading", "name", name, "url", fmt.Sprintf("%s/%s", base, candidateFileName))
	resp, err := http.Get(fmt.Sprintf("%s/%s", base, candidateFileName))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return "", fmt.Errorf("%s:%s not found", name, version)
	}

	// Create the file
	out, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}
	if out != nil {
		defer func() {
			if e := out.Close(); e != nil {
				err = e
			}
		}()
	}

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}

	// Add prmission to exec file
	err = os.Chmod(fullPath, os.ModePerm)
	if err != nil {
		return "", err
	}
	d.logger.Info("Downloaded", "name", name, "code", resp.StatusCode)
	return fullPath, err
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
