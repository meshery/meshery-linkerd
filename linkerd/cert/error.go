package cert

import (
	"github.com/layer5io/meshkit/errors"
)

var (
	// ErrCertEncodeCode represents the error code which is
	// generated when an encode operation fails
	ErrCertEncodeCode = "1101"
	// ErrEncodeCertificatesPEMCode represents the error code which is
	// generated when an certificate PEM encode operations fails
	ErrEncodeCertificatesPEMCode = "1102"
	// ErrEncodePrivateKeyPEMCode represents the error code which is
	// generated when an private key PEM encode operations fails
	ErrEncodePrivateKeyPEMCode = "1103"
	// ErrCreateRootCACode represents the error code which is
	// generated when root CA generation fails
	ErrCreateRootCACode = "1104"
	// ErrGeneratePKCode represents the error code which is
	// generated when private key generation fails
	ErrGeneratePKCode = "1105"
	// ErrGenerateDefaultRootCACode represents the error code which is
	// generated when defaut root CA generation fails
	ErrGenerateDefaultRootCACode = "1106"
)

// ErrCertEncode is the error for encode failure
func ErrCertEncode(err error) error {
	return errors.New(ErrCertEncodeCode, errors.Fatal, []string{"Failed to encode: "}, []string{err.Error()}, []string{}, []string{})
}

// ErrEncodeCertificatesPEM is the error for certificate encode failure
func ErrEncodeCertificatesPEM(err error) error {
	return errors.New(ErrEncodeCertificatesPEMCode, errors.Fatal, []string{"Failed to encode certificate PEM: "}, []string{err.Error()}, []string{}, []string{})
}

// ErrEncodePrivateKeyPEM is the error for private key PEM encode failure
func ErrEncodePrivateKeyPEM(err error) error {
	return errors.New(ErrEncodePrivateKeyPEMCode, errors.Fatal, []string{"Failed to encode private key PEM: "}, []string{err.Error()}, []string{}, []string{})
}

// ErrCreateRootCA is the error for root ca generation failure
func ErrCreateRootCA(err error) error {
	return errors.New(ErrCreateRootCACode, errors.Alert, []string{"Failed to create Root CA: "}, []string{err.Error()}, []string{}, []string{})
}

// ErrGeneratePK is the error for private key generation failure
func ErrGeneratePK(err error) error {
	return errors.New(ErrGeneratePKCode, errors.Alert, []string{"Failed to create Private Key: "}, []string{err.Error()}, []string{}, []string{})
}

// ErrGenerateDefaultRootCA is the error for default root ca generation failure
func ErrGenerateDefaultRootCA(err error) error {
	return errors.New(ErrGenerateDefaultRootCACode, errors.Alert, []string{"Failed to create default Root CA: "}, []string{err.Error()}, []string{}, []string{})
}
