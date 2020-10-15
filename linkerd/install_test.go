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
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestDownloadLinkerd(t *testing.T) {
	pathOfKubeconfig := os.Getenv("KUBECONFIG")
	contextName := os.Getenv("CURRENTCONTEXT")

	byteKubeconfig, err := ioutil.ReadFile(pathOfKubeconfig)

	if err != nil {
		t.Fatalf("Load kubeconfig err %s", err)
	}

	client, err := newClient(byteKubeconfig, contextName)

	if err != nil {
		t.Fatalf("NewClient function was failed %s", err)
	}

	err = client.downloadLinkerd()

	if err != nil {
		t.Fatalf("DownloadLinkerd function execution failed %s", err)
	}
}

func TestExecute(t *testing.T) {
	pathOfKubeconfig := os.Getenv("KUBECONFIG")
	contextName := os.Getenv("CURRENTCONTEXT")
	byteKubeconfig, err := ioutil.ReadFile(pathOfKubeconfig)

	if err != nil {
		t.Fatalf("Load kubeconfig err %s", err)
	}

	args := []string{
		"--context", contextName,
		"--kubeconfig", pathOfKubeconfig, "check", "--pre"}

	client, err := newClient(byteKubeconfig, contextName)

	if err != nil {
		t.Fatal(err)
	}

	outs, errs, err := client.execute(args...)

	if err != nil {
		fmt.Println(errs)
		t.Fatal(err)
	}
	t.Log(outs)
}
