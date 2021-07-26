// Package linkerd - Error codes for the adapter
package linkerd

import (
	"github.com/layer5io/meshkit/errors"
)

var (
	// ErrInstallLinkerdCode is the error code for ErrInstallLinkerd
	ErrInstallLinkerdCode = "1003"
	// ErrMeshConfigCode is the error code for ErrMeshConfig
	ErrMeshConfigCode = "1004"
	// ErrFetchManifestCode is the error code for ErrFetchManifest
	ErrFetchManifestCode = "1005"
	// ErrDownloadBinaryCode is the error code for ErrDownloadBinary
	ErrDownloadBinaryCode = "1006"
	// ErrInstallBinaryCode is the error code for ErrInstallBinary
	ErrInstallBinaryCode = "1007"
	// ErrClientConfigCode is the error code for ErrClientConfig
	ErrClientConfigCode = "1008"
	// ErrClientSetCode is the error code for ErrClientSet
	ErrClientSetCode = "1009"
	// ErrStreamEventCode is the error code for ErrStreamEvent
	ErrStreamEventCode = "1010"
	// ErrSampleAppCode is the error code for ErrSampleApp
	ErrSampleAppCode = "1011"
	// ErrCustomOperationCode is the error code for ErrCustomOperation
	ErrCustomOperationCode = "1012"
	// ErrOpInvalidCode is the error code for ErrOpInvalid
	ErrOpInvalidCode = "1013"

	// ErrOpInvalid is the error for invalid operation
	ErrOpInvalid = errors.New(ErrOpInvalidCode, errors.Alert, []string{"Invalid operation"}, []string{}, []string{}, []string{})
)

// ErrInstallLinkerd is the error for install mesh
func ErrInstallLinkerd(err error) error {
	return errors.New(ErrInstallLinkerdCode, errors.Alert, []string{"Error with Linkerd operation: ", err.Error()}, []string{}, []string{}, []string{})
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
	return errors.New(ErrDownloadBinaryCode, errors.Alert, []string{"Error downloading Linkerd binary: ", err.Error()}, []string{}, []string{}, []string{})
}

// ErrInstallBinary is the error while downloading linkerd binary
func ErrInstallBinary(err error) error {
	return errors.New(ErrInstallBinaryCode, errors.Alert, []string{"Error installing Linkerd binary: ", err.Error()}, []string{}, []string{}, []string{})
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
