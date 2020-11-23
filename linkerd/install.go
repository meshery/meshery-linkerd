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

	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-adapter-library/status"
	"github.com/layer5io/meshery-linkerd/internal/config"
	mesherykube "github.com/layer5io/meshkit/utils/kubernetes"
)

func (linkerd *Linkerd) installLinkerd(del bool, version string) (string, error) {
	linkerd.Log.Info(fmt.Sprintf("Requested install of version: %s", version))
	linkerd.Log.Info(fmt.Sprintf("Requested action is delete: %v", del))

	st := status.Installing

	if del {
		st = status.Removing
	}

	err := linkerd.Config.GetObject(adapter.MeshSpecKey, linkerd)
	if err != nil {
		return st, ErrMeshConfig(err)
	}

	manifest, err := linkerd.fetchManifest(version)
	if err != nil {
		linkerd.Log.Error(ErrInstallLinkerd(err))
		return st, ErrInstallLinkerd(err)
	}

	err = linkerd.applyManifest([]byte(manifest), del)
	if err != nil {
		linkerd.Log.Error(ErrInstallLinkerd(err))
		return st, ErrInstallLinkerd(err)
	}

	if del {
		return status.Removed, nil
	}
	return status.Installed, nil
}

func (linkerd *Linkerd) fetchManifest(version string) (string, error) {
	var (
		out bytes.Buffer
		er  bytes.Buffer
	)

	Executable, err := linkerd.getExecutable(version)
	if err != nil {
		return "", ErrFetchManifest(err, err.Error())
	}

	// We need a variable executable here hence using nosec
	// #nosec
	command := exec.Command(Executable, "install")
	command.Stdout = &out
	command.Stderr = &er
	err = command.Run()
	if err != nil {
		return "", ErrFetchManifest(err, er.String())
	}

	return out.String(), nil
}

func (linkerd *Linkerd) applyManifest(contents []byte, isDel bool) error {
	kclient, err := mesherykube.New(linkerd.KubeClient, linkerd.RestConfig)
	if err != nil {
		return err
	}

	err = kclient.ApplyManifest(contents, mesherykube.ApplyOptions{Delete: isDel})
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
	configRootPath := config.RootPath()
	linkerd.Log.Info("Looking for linkerd in", configRootPath, "...")
	executable = path.Join(configRootPath, alternateBinaryName)
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
	if err = installBinary(path.Join(configRootPath, alternateBinaryName), runtime.GOOS, res); err != nil {
		return "", err
	}

	linkerd.Log.Info("Done")
	return path.Join(configRootPath, alternateBinaryName), nil
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
	defer res.Body.Close()

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

		if err = out.Chmod(0755); err != nil {
			return ErrInstallBinary(err)
		}
	case "windows":
	}
	return nil
}
