package cert

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
)

// EncodeCertificatesPEM encodes the collection of provided certificates as
// a text blob of PEM-encoded certificates.
func EncodeCertificatesPEM(crts ...*x509.Certificate) ([]byte, error) {
	buf := bytes.Buffer{}
	for _, c := range crts {
		if err := encode(&buf, &pem.Block{Type: "CERTIFICATE", Bytes: c.Raw}); err != nil {
			return nil, ErrEncodeCertificatesPEM(err)
		}
	}

	return buf.Bytes(), nil
}

// EncodePrivateKeyPEM encodes the provided key as PEM-encoded text
func EncodePrivateKeyPEM(k *ecdsa.PrivateKey) ([]byte, error) {
	der, err := x509.MarshalECPrivateKey(k)
	if err != nil {
		return nil, ErrEncodePrivateKeyPEM(err)
	}

	return pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: der}), nil
}

func encode(buf *bytes.Buffer, blk *pem.Block) error {
	if err := pem.Encode(buf, blk); err != nil {
		return ErrCertEncode(err)
	}

	return nil
}
