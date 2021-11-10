// Package linkerd - for lifecycle management of Linkerd
package linkerd

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-adapter-library/status"
	"github.com/layer5io/meshery-linkerd/internal/config"
	"github.com/layer5io/meshery-linkerd/linkerd/cert"
	mesherykube "github.com/layer5io/meshkit/utils/kubernetes"
)

const (
	// LinkerdHelmStableRepo is the URL for linkerd stable helm repo
	LinkerdHelmStableRepo = "https://helm.linkerd.io/stable"
	// LinkerdHelmEdgeRepo is the URL for linkerd edge helm repo
	LinkerdHelmEdgeRepo = "https://helm.linkerd.io/edge"
)

var (
	//Namespace in which Linkerd is installed, (addons need to know this)
	linkerdNamespace = "linkerd"
)

func (linkerd *Linkerd) installLinkerd(del bool, version, namespace string) (string, error) {
	linkerdNamespace = namespace
	linkerd.Log.Info(fmt.Sprintf("Requested install of version: %s", version))
	linkerd.Log.Info(fmt.Sprintf("Requested action is delete: %v", del))
	linkerd.Log.Info(fmt.Sprintf("Requested action is in namespace: %s", namespace))

	// Overiding the namespace to be empty
	// This is intentional as deploying linkerd on custom namespace
	// is a bit tricky
	st := status.Installing

	if del {
		st = status.Removing
	}

	err := linkerd.Config.GetObject(adapter.MeshSpecKey, linkerd)
	if err != nil {
		return st, ErrMeshConfig(err)
	}

	if err := linkerd.applyHelmChart(version, namespace, del); err != nil {
		linkerd.Log.Error(ErrInstallLinkerd(err))

		linkerd.Log.Info("Attempting manifest installation...")

		// Attempt manifest installation
		manifest, err := linkerd.fetchManifest(version, namespace, del)
		if err != nil {
			linkerd.Log.Error(ErrInstallLinkerd(err))
			return st, ErrInstallLinkerd(err)
		}

		err = linkerd.applyManifest([]byte(manifest), del, namespace)
		if err != nil {
			linkerd.Log.Error(ErrInstallLinkerd(err))
			return st, ErrInstallLinkerd(err)
		}

		return st, nil
	}

	if del {
		return status.Removed, nil
	}
	return status.Installed, nil
}

func (linkerd *Linkerd) applyHelmChart(version string, namespace string, isDel bool) error {
	loc, cver := getChartLocationAndVersion(version)
	if loc == "" || cver == "" {
		return ErrInvalidVersionForMeshInstallation
	}

	// Generate certificates for linkerd
	c, pk, err := cert.GenerateRootCAWithDefaults("cluster.local")
	if err != nil {
		return ErrApplyHelmChart(err)
	}

	// Encode private key
	keyPEM, err := cert.EncodePrivateKeyPEM(pk)
	if err != nil {
		return ErrApplyHelmChart(err)
	}

	// Encode certificate
	certPEM, err := cert.EncodeCertificatesPEM(c)
	if err != nil {
		return ErrApplyHelmChart(err)
	}

	// Get expiry
	exp := c.NotAfter.Format(time.RFC3339)

	// Create namespace in which the installation was requested - Both
	// Helm and Linkerd to are too picky about this
	createHelmNS(linkerd.MesheryKubeclient, namespace, "linkerd2")

	err = linkerd.AnnotateNamespace(namespace, isDel, map[string]string{
		"app.kubernetes.io/managed-by":   "helm",
		"meta.helm.sh/release-name":      "linkerd2",
		"meta.helm.sh/release-namespace": namespace,
	})
	if err != nil {
		return ErrAnnotatingNamespace(err)
	}

	var act mesherykube.HelmChartAction
	if isDel {
		act = mesherykube.UNINSTALL
	} else {
		act = mesherykube.INSTALL
	}
	err = linkerd.MesheryKubeclient.ApplyHelmChart(mesherykube.ApplyHelmChartConfig{
		ChartLocation: mesherykube.HelmChartLocation{
			Repository: loc,
			Chart:      "linkerd2",
			Version:    cver,
		},
		Namespace: namespace,
		// CreateNamespace: true, // Don't use this => Linkerd NS has "special" requirements
		Action: act,
		OverrideValues: map[string]interface{}{
			"namespace":        namespace,
			"installNamespace": false,
			"global": map[string]interface{}{
				"identityTrustAnchorsPEM": string(certPEM),
			},
			"identityTrustAnchorsPEM": string(certPEM),
			"identity": map[string]interface{}{
				"issuer": map[string]interface{}{
					"crtExpiry": exp,
					"tls": map[string]interface{}{
						"keyPEM": string(keyPEM),
						"crtPEM": string(certPEM),
					},
				},
			},
		},
	})
	return err
}

func getChartLocationAndVersion(version string) (string, string) {
	if strings.HasPrefix(version, "edge-") {
		return LinkerdHelmEdgeRepo, strings.TrimPrefix(version, "edge-")
	}

	if strings.HasPrefix(version, "stable-") {
		return LinkerdHelmStableRepo, strings.TrimPrefix(version, "stable-")
	}

	return "", ""
}

func (linkerd *Linkerd) fetchManifest(version string, namespace string, isDel bool) (string, error) {
	var (
		out bytes.Buffer
		er  bytes.Buffer
	)

	Executable, err := linkerd.getExecutable(version)
	if err != nil {
		return "", ErrFetchManifest(err, err.Error())
	}
	execCmd := []string{"install", "--ignore-cluster", "--linkerd-namespace", namespace}
	if isDel {
		execCmd = []string{"uninstall", "--linkerd-namespace", namespace}
	}

	// We need a variable executable here hence using nosec
	// #nosec
	command := exec.Command(Executable, execCmd...)
	command.Stdout = &out
	command.Stderr = &er
	err = command.Run()
	if err != nil {
		return "", ErrFetchManifest(err, er.String())
	}

	return out.String(), nil
}

func (linkerd *Linkerd) applyManifest(contents []byte, isDel bool, namespace string) error {
	err := linkerd.MesheryKubeclient.ApplyManifest(contents, mesherykube.ApplyOptions{
		Namespace:    namespace,
		Update:       true,
		Delete:       isDel,
		IgnoreErrors: true,
	})
	if err != nil {
		return err
	}

	return nil
}

// getExecutable looks for the executable in
// 1. $PATH
// 2. Root config path
//
// If it doesn't find the executable in the path then it proceeds
// to download the binary from github releases and installs it
// in the root config path
func (linkerd *Linkerd) getExecutable(release string) (string, error) {
	const binaryName = "linkerd"
	alternateBinaryName := "linkerd-" + release

	// Look for the executable in the path
	linkerd.Log.Info("Looking for linkerd in the path...")
	executable, err := exec.LookPath(binaryName)
	if err == nil {
		return executable, nil
	}
	executable, err = exec.LookPath(alternateBinaryName)
	if err == nil {
		return executable, nil
	}

	// Look for config in the root path
	binPath := path.Join(config.RootPath(), "bin")
	linkerd.Log.Info("Looking for linkerd in", binPath, "...")
	executable = path.Join(binPath, alternateBinaryName)
	if _, err := os.Stat(executable); err == nil {
		return executable, nil
	}

	// Proceed to download the binary in the config root path
	linkerd.Log.Info("linkerd not found in the path, downloading...")
	res, err := downloadBinary(runtime.GOOS, runtime.GOARCH, release)
	if err != nil {
		return "", err
	}
	// Install the binary
	linkerd.Log.Info("Installing...")
	if err = installBinary(path.Join(binPath, alternateBinaryName), runtime.GOOS, res); err != nil {
		return "", err
	}

	linkerd.Log.Info("Done")
	return path.Join(binPath, alternateBinaryName), nil
}

func downloadBinary(platform, arch, release string) (*http.Response, error) {
	var url = "https://github.com/linkerd/linkerd2/releases/download"
	switch platform {
	case "darwin":
		fallthrough
	case "windows":
		url = fmt.Sprintf("%s/%s/linkerd2-cli-%s-%s", url, release, release, platform)
	case "linux":
		url = fmt.Sprintf("%s/%s/linkerd2-cli-%s-%s-%s", url, release, release, platform, arch)
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, ErrDownloadBinary(err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, ErrDownloadBinary(fmt.Errorf("bad status: %s", resp.Status))
	}

	return resp, nil
}

func installBinary(location, platform string, res *http.Response) error {
	// Close the response body
	defer func() {
		if err := res.Body.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	out, err := os.Create(location)
	if err != nil {
		return err
	}
	defer func() {
		if err := out.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	switch platform {
	case "darwin":
		fallthrough
	case "linux":
		_, err = io.Copy(out, res.Body)
		if err != nil {
			return ErrInstallBinary(err)
		}

		if err = out.Chmod(0750); err != nil {
			return ErrInstallBinary(err)
		}
	case "windows":
	}
	return nil
}

func createHelmNS(c *mesherykube.Client, ns, relName string) {
	const linkerdNS = `apiVersion: v1
kind: Namespace
metadata:
  name: %s
  labels:
    app.kubernetes.io/managed-by: Helm
  annotations:
    config.linkerd.io/admission-webhooks: disabled
    meta.helm.sh/release-name: %s
    meta.helm.sh/release-namespace: %s`

	_ = c.ApplyManifest([]byte(fmt.Sprintf(linkerdNS, ns, relName, ns)), mesherykube.ApplyOptions{
		Update:       true,
		IgnoreErrors: true,
	})
}
