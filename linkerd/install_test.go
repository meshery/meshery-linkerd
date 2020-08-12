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
	"context"
	"io/ioutil"
	"os"
	"testing"
)

func TestPreck(t *testing.T) {

	//os.Setenv("KUBECONFIG", "/Users/aisuko/Documents/rke/kube_config_cluster.yml")
	//os.Setenv("CURRENTCONTEXT", "local")
	kubectlConfig := os.Getenv("KUBECONFIG")
	contextName := os.Getenv("CURRENTCONTEXT")
	byteKubeconfig, err := ioutil.ReadFile(kubectlConfig)
	if err != nil {
		t.Fatal(err)
	}
	client, err := newClient(byteKubeconfig, contextName)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	err = client.preCheck(ctx, "linkerd")
	if err != nil {
		t.Fatal(err)
	}

	//bol := assert.NotEmpty(t, deployMainyaml)
	//if !bol && err != nil {
	//	t.Fatal("deployMainyaml is blank")
	//}
	//
	//err = client.deployment(deployMainyaml)
	//
	//assert.NoError(t, err)
}
