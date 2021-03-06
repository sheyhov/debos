package actions

import (
	"fmt"
	"github.com/go-debos/debos"
	"net/url"
	"path"
)

type DownloadAction struct {
	debos.BaseAction `yaml:",inline"`
	Url              string // URL for downloading
	Filename         string // File name, overrides the name from URL.
	Unpack           bool   // Unpack downloaded file to directory dedicated for download
	Compression      string // compression type
	Name             string // exporting path to file or directory(in case of unpack)
}

// validateUrl checks if supported URL is passed from recipe
// Return:
// - parsed URL
// - nil in case of success
func (d *DownloadAction) validateUrl() (*url.URL, error) {

	url, err := url.Parse(d.Url)
	if err != nil {
		return url, err
	}

	switch url.Scheme {
	case "http", "https":
		// Supported scheme
	default:
		return url, fmt.Errorf("Unsupported URL is provided: '%s'", url.String())
	}

	return url, nil
}

func (d *DownloadAction) Verify(context *debos.DebosContext) error {

	if len(d.Name) == 0 {
		return fmt.Errorf("Property 'name' is mandatory for download action\n")
	}
	_, err := d.validateUrl()
	return err
}

func (d *DownloadAction) Run(context *debos.DebosContext) error {
	var filename string
	d.LogStart()

	url, err := d.validateUrl()
	if err != nil {
		return err
	}

	if len(d.Filename) == 0 {
		// Trying to guess the name from URL Path
		filename = path.Base(url.Path)
	} else {
		filename = path.Base(d.Filename)
	}
	if len(filename) == 0 {
		return fmt.Errorf("Incorrect filename is provided for '%s'", d.Url)
	}
	filename = path.Join(context.Scratchdir, filename)
	originPath := filename

	switch url.Scheme {
	case "http", "https":
		err := debos.DownloadHttpUrl(url.String(), filename)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("Unsupported URL is provided: '%s'", url.String())
	}

	if d.Unpack == true {
		targetdir := filename + ".d"
		err := debos.UnpackTarArchive(filename, targetdir, d.Compression, "--no-same-owner", "--no-same-permissions")
		if err != nil {
			return err
		}
		originPath = targetdir
	}

	context.Origins[d.Name] = originPath

	return nil
}
