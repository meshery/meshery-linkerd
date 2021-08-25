package cert

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"time"
)

const (
	// DefaultLifetime configures certificate validity.
	DefaultLifetime = (24 * 365) * time.Hour

	// DefaultClockSkewAllowance indicates the maximum allowed difference in clocks
	// in the network.
	DefaultClockSkewAllowance = 10 * time.Second
)

// CreateRootCA generates root CA
func CreateRootCA(name string, key *ecdsa.PrivateKey, validFrom *time.Time) (*x509.Certificate, error) {
	dc := GetDefaultX509Cert(1, &key.PublicKey, validFrom)
	dc.Subject = pkix.Name{CommonName: name}
	dc.IsCA = true
	dc.MaxPathLen = -1
	dc.BasicConstraintsValid = true
	dc.KeyUsage = x509.KeyUsageCertSign | x509.KeyUsageCRLSign

	crtb, err := x509.CreateCertificate(rand.Reader, dc, dc, key.Public(), key)
	if err != nil {
		return nil, ErrCreateRootCA(err)
	}

	pc, err := x509.ParseCertificate(crtb)
	if err != nil {
		return nil, ErrCreateRootCA(err)
	}

	return pc, nil
}

// GenerateKey creates a new P-256 ECDSA private key from the default random source.
func GenerateKey() (*ecdsa.PrivateKey, error) {
	pk, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, ErrGeneratePK(err)
	}

	return pk, nil
}

// GenerateRootCAWithDefaults generates a new root CA with default settings.
func GenerateRootCAWithDefaults(name string) (*x509.Certificate, *ecdsa.PrivateKey, error) {
	key, err := GenerateKey()
	if err != nil {
		return nil, nil, ErrGenerateDefaultRootCA(err)
	}

	now := time.Now()

	res, err := CreateRootCA(name, key, &now)
	if err != nil {
		return nil, nil, ErrGenerateDefaultRootCA(err)
	}

	return res, key, nil
}

// GetDefaultX509Cert returns x509 cert with some defaults
func GetDefaultX509Cert(serialNumber uint64, k *ecdsa.PublicKey, validFrom *time.Time) *x509.Certificate {
	const SignatureAlgorithm = x509.ECDSAWithSHA256

	if validFrom == nil {
		now := time.Now()
		validFrom = &now
	}
	notBefore, notAfter := GetWindow(*validFrom, DefaultLifetime, DefaultClockSkewAllowance)

	return &x509.Certificate{
		SerialNumber:       big.NewInt(int64(serialNumber)),
		SignatureAlgorithm: SignatureAlgorithm,
		NotBefore:          notBefore,
		NotAfter:           notAfter,
		PublicKey:          k,
		KeyUsage:           x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
			x509.ExtKeyUsageClientAuth,
		},
	}
}

// GetWindow returns cert validity window based on the arguments given
func GetWindow(t time.Time, lifetime, clockSkewAllowance time.Duration) (time.Time, time.Time) {
	return t.Add(-clockSkewAllowance), t.Add(lifetime).Add(clockSkewAllowance)
}
