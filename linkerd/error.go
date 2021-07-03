// Package linkerd - Error codes for the adapter
package linkerd

import (
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
	return errors.New(ErrInstallLinkerdCode, errors.Alert, []string{"Error with linkerd operation: ", err.Error()}, []string{}, []string{}, []string{})
}

// ErrMeshConfig is the error for mesh config
func ErrMeshConfig(err error) error {
	return errors.New(ErrMeshConfigCode, errors.Alert, []string{"Error configuration mesh: ", err.Error()}, []string{}, []string{}, []string{})
}

// ErrFetchManifest is the error for mesh port forward
func ErrFetchManifest(err error, des string) error {
	return errors.New(ErrFetchManifestCode, errors.Alert, []string{"Error fetching mesh manifest: %s", des}, []string{err.Error()}, []string{}, []string{})
}

// ErrDownloadBinary is the error while downloading linkerd binary
func ErrDownloadBinary(err error) error {
	return errors.New(ErrDownloadBinaryCode, errors.Alert, []string{"Error downloading linkerd binary: ", err.Error()}, []string{}, []string{}, []string{})
}

// ErrInstallBinary is the error while downloading linkerd binary
func ErrInstallBinary(err error) error {
	return errors.New(ErrInstallBinaryCode, errors.Alert, []string{"Error installing linkerd binary: ", err.Error()}, []string{}, []string{}, []string{})
}

// ErrClientConfig is the error for setting client config
func ErrClientConfig(err error) error {
	return errors.New(ErrClientConfigCode, errors.Alert, []string{"Error setting client config: ", err.Error()}, []string{}, []string{}, []string{})
}

// ErrClientSet is the error for setting clientset
func ErrClientSet(err error) error {
	return errors.New(ErrClientSetCode, errors.Alert, []string{"Error setting clientset: ", err.Error()}, []string{}, []string{}, []string{})
}

// ErrStreamEvent is the error for streaming event
func ErrStreamEvent(err error) error {
	return errors.New(ErrStreamEventCode, errors.Alert, []string{"Error streaming event: ", err.Error()}, []string{}, []string{}, []string{})
}

// ErrSampleApp is the error for streaming event
func ErrSampleApp(err error) error {
	return errors.New(ErrSampleAppCode, errors.Alert, []string{"Error with sample app operation: ", err.Error()}, []string{}, []string{}, []string{})
}

// ErrCustomOperation is the error for streaming event
func ErrCustomOperation(err error) error {
	return errors.New(ErrCustomOperationCode, errors.Alert, []string{"Error with custom operation: ", err.Error()}, []string{}, []string{}, []string{})
}
