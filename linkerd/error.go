// Package linkerd - Error codes for the adapter
package linkerd

import (
	"fmt"

	"github.com/layer5io/meshkit/errors"
)

var (
	// ErrInstallLinkerdCode is the error code for ErrInstallLinkerd
	ErrInstallLinkerdCode = "linkerd_test_code"
	// ErrMeshConfigCode is the error code for ErrMeshConfig
	ErrMeshConfigCode = "linkerd_test_code"
	// ErrFetchManifestCode is the error code for ErrFetchManifest
	ErrFetchManifestCode = "linkerd_test_code"
	// ErrDownloadBinaryCode is the error code for ErrDownloadBinary
	ErrDownloadBinaryCode = "linkerd_test_code"
	// ErrInstallBinaryCode is the error code for ErrInstallBinary
	ErrInstallBinaryCode = "linkerd_test_code"
	// ErrClientConfigCode is the error code for ErrClientConfig
	ErrClientConfigCode = "linkerd_test_code"
	// ErrClientSetCode is the error code for ErrClientSet
	ErrClientSetCode = "linkerd_test_code"
	// ErrStreamEventCode is the error code for ErrStreamEvent
	ErrStreamEventCode = "linkerd_test_code"
	// ErrSampleAppCode is the error code for ErrSampleApp
	ErrSampleAppCode = "linkerd_test_code"
	// ErrCustomOperationCode is the error code for ErrCustomOperation
	ErrCustomOperationCode = "linkerd_test_code"

	// ErrOpInvalid is the error for invalid operation
	ErrOpInvalid = errors.NewDefault(errors.ErrOpInvalid, "Invalid operation")
)

// ErrInstallLinkerd is the error for install mesh
func ErrInstallLinkerd(err error) error {
	return errors.NewDefault(ErrInstallLinkerdCode, fmt.Sprintf("Error with linkerd operation: %s", err.Error()))
}

// ErrMeshConfig is the error for mesh config
func ErrMeshConfig(err error) error {
	return errors.NewDefault(ErrMeshConfigCode, fmt.Sprintf("Error configuration mesh: %s", err.Error()))
}

// ErrFetchManifest is the error for mesh port forward
func ErrFetchManifest(err error, des string) error {
	return errors.NewDefault(ErrFetchManifestCode, fmt.Sprintf("Error fetching mesh manifest: %s", des))
}

// ErrDownloadBinary is the error while downloading linkerd binary
func ErrDownloadBinary(err error) error {
	return errors.NewDefault(ErrDownloadBinaryCode, fmt.Sprintf("Error downloading linkerd binary: %s", err.Error()))
}

// ErrInstallBinary is the error while downloading linkerd binary
func ErrInstallBinary(err error) error {
	return errors.NewDefault(ErrInstallBinaryCode, fmt.Sprintf("Error installing linkerd binary: %s", err.Error()))
}

// ErrClientConfig is the error for setting client config
func ErrClientConfig(err error) error {
	return errors.NewDefault(ErrClientConfigCode, fmt.Sprintf("Error setting client config: %s", err.Error()))
}

// ErrClientSet is the error for setting clientset
func ErrClientSet(err error) error {
	return errors.NewDefault(ErrClientSetCode, fmt.Sprintf("Error setting clientset: %s", err.Error()))
}

// ErrStreamEvent is the error for streaming event
func ErrStreamEvent(err error) error {
	return errors.NewDefault(ErrStreamEventCode, fmt.Sprintf("Error streaming event: %s", err.Error()))
}

// ErrSampleApp is the error for streaming event
func ErrSampleApp(err error) error {
	return errors.NewDefault(ErrSampleAppCode, fmt.Sprintf("Error with sample app operation: %s", err.Error()))
}

// ErrCustomOperation is the error for streaming event
func ErrCustomOperation(err error) error {
	return errors.NewDefault(ErrCustomOperationCode, fmt.Sprintf("Error with custom operation: %s", err.Error()))
}
