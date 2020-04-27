package linkerd

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestClient(t *testing.T) {

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

	// TODO Could out more information about the client if we need
	if client == nil {
		t.Fatal("Client is nil")
	} else {
		t.Skip()
	}
}
