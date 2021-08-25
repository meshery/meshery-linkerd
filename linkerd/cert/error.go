package cert

import (
	"github.com/layer5io/meshkit/errors"
)

var (
	ErrCertEncodeCode            = "1101"
	ErrEncodeCertificatesPEMCode = "1102"
	ErrEncodePrivateKeyPEMCode   = "1103"
	ErrCreateRootCACode          = "1104"
	ErrGeneratePKCode            = "1105"
	ErrGenerateDefaultRootCACode = "1106"
)

func ErrCertEncode(err error) error {
	return errors.New(ErrCertEncodeCode, errors.Fatal, []string{"Failed to encode: "}, []string{err.Error()}, []string{}, []string{})
}

func ErrEncodeCertificatesPEM(err error) error {
	return errors.New(ErrEncodeCertificatesPEMCode, errors.Fatal, []string{"Failed to encode certificate PEM: "}, []string{err.Error()}, []string{}, []string{})
}

func ErrEncodePrivateKeyPEM(err error) error {
	return errors.New(ErrEncodePrivateKeyPEMCode, errors.Fatal, []string{"Failed to encode private key PEM: "}, []string{err.Error()}, []string{}, []string{})
}

func ErrCreateRootCA(err error) error {
	return errors.New(ErrCreateRootCACode, errors.Alert, []string{"Failed to create Root CA: "}, []string{err.Error()}, []string{}, []string{})
}

func ErrGeneratePK(err error) error {
	return errors.New(ErrGeneratePKCode, errors.Alert, []string{"Failed to create Private Key: "}, []string{err.Error()}, []string{}, []string{})
}

func ErrGenerateDefaultRootCA(err error) error {
	return errors.New(ErrGenerateDefaultRootCACode, errors.Alert, []string{"Failed to create default Root CA: "}, []string{err.Error()}, []string{}, []string{})
}
