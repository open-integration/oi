package downloader

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"runtime"

	"github.com/open-integration/core/pkg/logger"
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

// New creates new Downloader
func New(opt Options) Downloader {
	return &downloader{
		store:  opt.Store,
		logger: opt.Logger,
	}
}

func (d *downloader) Download(name string, version string) (string, error) {
	candidateFileName := fmt.Sprintf("%s-%s-%s-%s", name, version, runtime.GOOS, runtime.GOARCH)
	fullPath := path.Join(d.store, candidateFileName)
	_, err := ioutil.ReadFile(fullPath)
	if os.IsExist(err) {
		d.logger.Debug("Skipping download, service exist", "path", fullPath)
		return fullPath, nil
	}

	url := fmt.Sprintf("https://storage.googleapis.com/open-integration-service-catalog/%s", candidateFileName)
	d.logger.Debug("Downloading", "name", name, "url", url)
	resp, err := http.Get(url)
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
	defer out.Close()

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
	d.logger.Debug("Downloaded", "name", name, "code", resp.StatusCode)
	return fullPath, nil
}
