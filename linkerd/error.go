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
	// ErrSampleAppCode is the error code for ErrSampleApp
	ErrSampleAppCode = "1011"
	// ErrCustomOperationCode is the error code for ErrCustomOperation
	ErrCustomOperationCode = "1012"
	// ErrOpInvalidCode is the error code for ErrOpInvalid
	ErrOpInvalidCode = "1013"
	// ErrInvalidOAMComponentTypeCode represents the error code which is
	// generated when an invalid oam component is requested
	ErrInvalidOAMComponentTypeCode = "1014"
	// ErrLinkerdCoreComponentFailCode represents the error code which is
	// generated when an linkerd core operations fails
	ErrLinkerdCoreComponentFailCode = "1015"
	// ErrProcessOAMCode represents the error code which is
	// generated when an OAM operations fails
	ErrProcessOAMCode = "1016"
	// ErrParseLinkerdCoreComponentCode represents the error code which is
	// generated when linkerd core component manifest parsing fails
	ErrParseLinkerdCoreComponentCode = "1017"
	// ErrParseOAMComponentCode represents the error code which is
	// generated during the OAM component parsing
	ErrParseOAMComponentCode = "1018"
	// ErrParseOAMConfigCode represents the error code which is
	// generated during the OAM configuration parsing
	ErrParseOAMConfigCode = "1019"

	// ErrApplyHelmChartCode represents the error code which is
	// generated during the Helm Chart installation
	ErrApplyHelmChartCode = "1020"

	// ErrNilClientCode represents the error code which is
	// generated when Kubernetes client is nil
	ErrNilClientCode = "1019"

	//ErrAddonFromHelmCode represents the error while installing addons through helm charts
	ErrAddonFromHelmCode = "1014"

	//ErrInvalidVersionForMeshInstallationCode represents the error while installing mesh through helm charts with invalid version
	ErrInvalidVersionForMeshInstallationCode = "1015"

	//ErrAnnotatingNamespaceCode represents the error while annotating namespace
	ErrAnnotatingNamespaceCode = "1016"
	//ErrInvalidVersionForMeshInstallation represents the error while installing mesh through helm charts with invalid version
	ErrInvalidVersionForMeshInstallation = errors.New(ErrInvalidVersionForMeshInstallationCode, errors.Alert, []string{"Invalid version passed for helm based installation"}, []string{"Version passed is invalid"}, []string{"Version might not be prefixed with \"stable-\" or \"edge-\""}, []string{"Version should be prefixed with \"stable-\" or \"edge-\"", "Version might be empty"})
	// ErrOpInvalid is the error for invalid operation
	ErrOpInvalid = errors.New(ErrOpInvalidCode, errors.Alert, []string{"Invalid operation"}, []string{}, []string{}, []string{})

	// ErrParseOAMComponent represents the error which is
	// generated during the OAM component parsing
	ErrParseOAMComponent = errors.New(ErrParseOAMComponentCode, errors.Alert, []string{"error parsing the component"}, []string{"Error occurred while parsing application component in the OAM request made by Meshery server"}, []string{"Could not unmarshall configuration component received via ProcessOAM gRPC call into a valid Component struct"}, []string{"Check if Meshery Server is creating valid component for ProcessOAM gRPC call. This error should never happen and can be reported as a bug in Meshery Server. Also check if Meshery Server and adapters are referring to same component struct provided in MeshKit."})

	// ErrParseOAMConfig represents the error which is
	// generated during the OAM configuration parsing
	ErrParseOAMConfig = errors.New(ErrParseOAMConfigCode, errors.Alert, []string{"error parsing the configuration"}, []string{"Error occured while parsing configuration in the request made by Meshery Server"}, []string{"Could not unmarshall OAM config recieved via ProcessOAM gRPC call into a valid Config struct"}, []string{"Check if Meshery Server is creating valid config for ProcessOAM gRPC call. This error should never happen and can be reported as a bug in Meshery Server. Also, confirm that Meshery Server and Adapters are referring to same config struct provided in MeshKit"})

	// ErrNilClient represents the error which is
	// generated when Kubernetes client is nil
	ErrNilClient = errors.New(ErrNilClientCode, errors.Alert, []string{"Kubernetes client not initialized"}, []string{"Kubernetes client is nil"}, []string{"Kubernetes client not initialized"}, []string{"Reconnect the Meshery Adapter to Meshery Server"})
)

// ErrInstallLinkerd is the error for install mesh
func ErrInstallLinkerd(err error) error {
	return errors.New(ErrInstallLinkerdCode, errors.Alert, []string{"Error with Linkerd operation: "}, []string{err.Error()}, []string{}, []string{})
}

// ErrMeshConfig is the error for mesh config
func ErrMeshConfig(err error) error {
	return errors.New(ErrMeshConfigCode, errors.Alert, []string{"Error configuration mesh: "}, []string{err.Error()}, []string{}, []string{})
}

// ErrFetchManifest is the error for mesh port forward
func ErrFetchManifest(err error, des string) error {
	return errors.New(ErrFetchManifestCode, errors.Alert, []string{"Error fetching mesh manifest: %s", des}, []string{err.Error()}, []string{}, []string{})
}

// ErrDownloadBinary is the error while downloading linkerd binary
func ErrDownloadBinary(err error) error {
	return errors.New(ErrDownloadBinaryCode, errors.Alert, []string{"Error downloading Linkerd binary: "}, []string{err.Error()}, []string{"Checkout https://docs.github.com/en/rest/reference/repos#releases for more details"}, []string{})
}

// ErrInstallBinary is the error while downloading linkerd binary
func ErrInstallBinary(err error) error {
	return errors.New(ErrInstallBinaryCode, errors.Alert, []string{"Error installing Linkerd binary: "}, []string{err.Error()}, []string{}, []string{})
}

// ErrClientConfig is the error for setting client config
func ErrClientConfig(err error) error {
	return errors.New(ErrClientConfigCode, errors.Alert, []string{"Error setting client config: "}, []string{err.Error()}, []string{}, []string{})
}

// ErrClientSet is the error for setting clientset
func ErrClientSet(err error) error {
	return errors.New(ErrClientSetCode, errors.Alert, []string{"Error setting clientset: "}, []string{err.Error()}, []string{}, []string{})
}

// ErrSampleApp is the error for streaming event
func ErrSampleApp(err error) error {
	return errors.New(ErrSampleAppCode, errors.Alert, []string{"Error with sample app operation"}, []string{err.Error(), "Error occurred while trying to install a sample application using manifests"}, []string{"Invalid kubeclient config", "Invalid manifest"}, []string{"Reconnect your adapter to meshery server to refresh the kubeclient"})
}

// ErrCustomOperation is the error for streaming event
func ErrCustomOperation(err error) error {
	return errors.New(ErrCustomOperationCode, errors.Alert, []string{"Error with custom operation"}, []string{"Error occurred while applying custom manifest to the cluster", err.Error()}, []string{"Invalid kubeclient config", "Invalid manifest"}, []string{"Make sure to apply a valid Kubernetes manifest"})
}

// ErrParseLinkerdCoreComponent is the error when linkerd core component manifest parsing fails
func ErrParseLinkerdCoreComponent(err error) error {
	return errors.New(ErrParseLinkerdCoreComponentCode, errors.Alert, []string{"linkerd core component manifest parsing failing"}, []string{err.Error()}, []string{}, []string{})
}

// ErrInvalidOAMComponentType is the error when the OAM component name is not valid
func ErrInvalidOAMComponentType(compName string) error {
	return errors.New(ErrInvalidOAMComponentTypeCode, errors.Alert, []string{"invalid OAM component name: ", compName}, []string{}, []string{}, []string{})
}

// ErrLinkerdCoreComponentFail is the error when core linkerd component processing fails
func ErrLinkerdCoreComponentFail(err error) error {
	return errors.New(ErrLinkerdCoreComponentFailCode, errors.Alert, []string{"error in linkerd core component"}, []string{err.Error()}, []string{}, []string{})
}

// ErrProcessOAM is a generic error which is thrown when an OAM operations fails
func ErrProcessOAM(err error) error {
	return errors.New(ErrProcessOAMCode, errors.Alert, []string{"error performing OAM operations"}, []string{err.Error()}, []string{}, []string{})
}

// ErrApplyHelmChart is an error which is thrown when apply helm chart fails
func ErrApplyHelmChart(err error) error {
	return errors.New(ErrApplyHelmChartCode, errors.Alert, []string{"error applying helm chart"}, []string{err.Error()}, []string{}, []string{})
}

// ErrAddonFromHelm is the error for installing addons through helm chart
func ErrAddonFromHelm(err error) error {
	return errors.New(ErrAddonFromHelmCode, errors.Alert, []string{"Error with addon install operation by helm chart"}, []string{err.Error()}, []string{"The helm chart URL in additional properties of addon operation might be incorrect", "Could not apply service patch file for the given addon"}, []string{})
}

//ErrAnnotatingNamespace is the error while annotating the namespace
func ErrAnnotatingNamespace(err error) error {
	return errors.New(ErrAddonFromHelmCode, errors.Alert, []string{"Error with annotating namespace"}, []string{err.Error()}, []string{"Could not get the namespace in cluster", "Could not update namespace in cluster"}, []string{"Make sure the cluster is reachable"})
}
