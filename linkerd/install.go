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
	"io/ioutil"
	"os"
	"path"
	"time"

	getter "github.com/hashicorp/go-getter"
	"github.com/layer5io/meshery-linkerd/pkg/util"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	emojivotoInstallFile = "https://run.linkerd.io/emojivoto.yml"
	booksAppInstallFile  = "https://run.linkerd.io/booksapp.yml"
	cachePeriod          = 1 * time.Hour
)

var (
	emojivotoLocalFile = path.Join(os.TempDir(), "emojivoto.yml")
	booksAppLocalFile  = path.Join(os.TempDir(), "booksapp.yml")
)

func (iClient *Client) downloadFile(urlToDownload, localFile string) error {
	dFile, err := os.Create(localFile)
	if err != nil {
		err = errors.Wrapf(err, "unable to create a file on the filesystem at %s", localFile)
		logrus.Error(err)
		return err
	}

	defer util.SafeClose(dFile, &err)
	err = getter.GetFile(localFile, urlToDownload)
	if err != nil {
		err = errors.Wrapf(err, "Download the file failed %s", localFile)
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
		// TODO Change to the HashiCorp tool which uses in the shipyard-run repo
		if err = iClient.downloadFile(remoteURL, localFile); err != nil {
			return "", err
		}
		logrus.Debug("file successfully downloaded . . .")
	}
	/* #nosec */
	b, err := ioutil.ReadFile(localFile)
	return string(b), err
}
