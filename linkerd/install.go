// Copyright 2019 Layer5.io
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package linkerd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/layer5io/meshery-linkerd/pkg/util"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	repoURL              = "https://api.github.com/repos/linkerd/linkerd2/releases"
	emojivotoInstallFile = "https://run.linkerd.io/emojivoto.yml"
	booksAppInstallFile  = "https://run.linkerd.io/booksapp.yml"

	cachePeriod = 1 * time.Hour
)

var (
	urlsuffix          = "-" + runtime.GOOS //defining
	localFile          = path.Join(os.TempDir(), "linkerd-cli")
	emojivotoLocalFile = path.Join(os.TempDir(), "emojivoto.yml")
	booksAppLocalFile  = path.Join(os.TempDir(), "booksapp.yml")
)

// Asset is used to store the individual asset data as part of a release
type Asset struct {
	Name        string `json:"name,omitempty"`
	State       string `json:"state,omitempty"`
	DownloadURL string `json:"browser_download_url,omitempty"`
}

// Release is used to save the release informations
type Release struct {
	ID      int      `json:"id,omitempty"`
	TagName string   `json:"tag_name,omitempty"`
	Name    string   `json:"name,omitempty"`
	Draft   bool     `json:"draft,omitempty"`
	Assets  []*Asset `json:"assets,omitempty"`
}

func (iClient *Client) getLatestReleaseURL() error {

	if iClient.linkerdReleaseDownloadURL == "" || time.Since(iClient.linkerdReleaseUpdatedAt) > cachePeriod {
		logrus.Debugf("API info url: %s", repoURL)
		resp, err := http.Get(repoURL)
		if err != nil {
			err = errors.Wrapf(err, "error getting latest version info")
			logrus.Error(err)
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			err = fmt.Errorf("unable to fetch release info due to an unexpected status code: %d", resp.StatusCode)
			logrus.Error(err)
			return err
		}

		// TODO Need to confirm that the github APIv3 limit the number of request
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			err = errors.Wrapf(err, "error parsing response body")
			logrus.Error(err)
			return err
		}
		// TODO There may have a consider if the top 10 release did not includes the stable version
		releaseList := make([]*Release, 10)

		err = json.Unmarshal(body, &releaseList)
		if err != nil {
			err = errors.Wrapf(err, "error unmarshalling response body")
			logrus.Error(err)
			return err
		}

		for _, v := range releaseList {
			if strings.HasPrefix(v.TagName, "stable") && !v.Draft {
				for _, asset := range v.Assets {
					if strings.HasSuffix(asset.Name, urlsuffix) {
						iClient.linkerdReleaseVersion = strings.Replace(asset.Name, urlsuffix, "", -1)
						iClient.linkerdReleaseDownloadURL = asset.DownloadURL
						iClient.linkerdReleaseUpdatedAt = time.Now()
						return nil
					}
				}
			}
		}
		err = errors.New("unable to extract the download URL")
		logrus.Error(err)
		return err
	}
	return nil
}

func (iClient *Client) downloadFile(urlToDownload, localFile string) error {
	dFile, err := os.Create(localFile)
	if err != nil {
		err = errors.Wrapf(err, "unable to create a file on the filesystem at %s", localFile)
		logrus.Error(err)
		return err
	}

	defer util.SafeClose(dFile, &err)

	/* #nosec */
	resp, err := http.Get(urlToDownload)
	if err != nil {
		err = errors.Wrapf(err, "unable to download the file from URL: %s", iClient.linkerdReleaseDownloadURL)
		logrus.Error(err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("unable to download the file from URL: %s, status: %s", iClient.linkerdReleaseDownloadURL, resp.Status)
		logrus.Error(err)
		return err
	}

	_, err = io.Copy(dFile, resp.Body)
	if err != nil {
		err = errors.Wrapf(err, "unable to write the downloaded file to the file system at %s", localFile)
		logrus.Error(err)
		return err
	}
	/* #nosec */
	err = os.Chmod(localFile, 0755)
	if err != nil {
		err = errors.Wrapf(err, "unable to change permission on %s", localFile)
		logrus.Error(err)
		return err
	}
	return nil
}

func (iClient *Client) downloadLinkerd() error {
	logrus.Debug("preparing to download the latest linkerd release")
	err := iClient.getLatestReleaseURL()
	if err != nil {
		return err
	}
	fileName := iClient.linkerdReleaseVersion
	downloadURL := iClient.linkerdReleaseDownloadURL
	logrus.Debugf("retrieved latest file name: %s and download url: %s", fileName, downloadURL)

	proceedWithDownload := true

	lFileStat, err := os.Stat(localFile)
	if err == nil {
		if time.Since(lFileStat.ModTime()) > cachePeriod {
			proceedWithDownload = true
		} else {
			proceedWithDownload = false
		}
	}

	if proceedWithDownload {
		if err = iClient.downloadFile(downloadURL, localFile); err != nil {
			return err
		}
		logrus.Debug("package successfully downloaded . . .")
	}
	return nil
}

func (iClient *Client) execute(command ...string) (string, string, error) {
	err := iClient.downloadLinkerd()
	if err != nil {
		return "", "", err
	}
	logrus.Debugf("checking if install file exists at path: %s", localFile)
	_, err = os.Stat(localFile)
	if err != nil {
		err = errors.Wrap(err, "path not found")
		logrus.Error(err)
	}

	// TODO: execute
	logrus.Debugf("command to be executed: %s %v", localFile, command)
	/* #nosec */
	cmd := exec.Command(localFile, command...)
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb

	err = cmd.Run()
	if err != nil {
		err = errors.Wrapf(err, "error while executing requested command")
		logrus.Error(err)
	}
	logrus.Debugf("Received output: %s", outb.String())
	logrus.Debugf("Received error: %s", errb.String())
	return outb.String(), errb.String(), nil
}

func (iClient *Client) getYAML(remoteURL, localFile string) (string, error) {

	proceedWithDownload := true

	lFileStat, err := os.Stat(localFile)
	if err == nil {
		if time.Since(lFileStat.ModTime()) > cachePeriod {
			proceedWithDownload = true
		} else {
			proceedWithDownload = false
		}
	}

	if proceedWithDownload {
		if err = iClient.downloadFile(remoteURL, localFile); err != nil {
			return "", err
		}
		logrus.Debug("file successfully downloaded . . .")
	}
	/* #nosec */
	b, err := ioutil.ReadFile(localFile)
	return string(b), err
}
