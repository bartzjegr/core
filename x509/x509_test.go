package x509

import (
	"crypto/x509/pkix"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestX509SignCSR(t *testing.T) {
	rootCA, _ := NewCA(nil)
	rootCA.Data.Body.Name = "RootCA"
	rootCA.Data.Body.DNScope.Country = "UK"
	rootCA.Data.Body.DNScope.Organization = "pki.io"
	rootCA.GenerateRoot()

	subCA, _ := NewCA(nil)
	subCA.Data.Body.Name = "DevCA"
	subCA.Data.Body.DNScope.OrganizationalUnit = "Development"
	subCA.GenerateSub(rootCA)

	csr, _ := NewCSR(nil)
	csr.Data.Body.Name = "Server1"
	subject := pkix.Name{CommonName: csr.Data.Body.Name}
	csr.Generate(&subject)

	csrPublic, _ := csr.Public()

	cert, err := subCA.Sign(csrPublic, false)
	assert.Nil(t, err)
	assert.NotNil(t, cert)
	assert.NotEqual(t, cert.Data.Body.Certificate, "")

	certificate, err := PemDecodeX509Certificate([]byte(cert.Data.Body.Certificate))
	assert.Nil(t, err)
	assert.True(t, certificate.NotBefore.After(time.Now().AddDate(0, 0, -1)))
	assert.True(t, certificate.NotAfter.Before(time.Now().AddDate(0, 0, 1)))
}

func TestX509SignCSRKeepSubject(t *testing.T) {
	rootCA, _ := NewCA(nil)
	rootCA.Data.Body.Name = "RootCA"
	rootCA.GenerateRoot()

	csr, _ := NewCSR(nil)
	csr.Data.Body.Name = "Server1"
	subject := pkix.Name{CommonName: csr.Data.Body.Name}
	csr.Generate(&subject)

	csrPublic, _ := csr.Public()

	cert, err := rootCA.Sign(csrPublic, true)
	assert.Nil(t, err)
	assert.NotNil(t, cert)
	assert.NotEqual(t, cert.Data.Body.Certificate, "")

	certificate, err := PemDecodeX509Certificate([]byte(cert.Data.Body.Certificate))
	assert.Nil(t, err)
	assert.Equal(t, certificate.Subject.CommonName, subject.CommonName)
	assert.True(t, certificate.NotBefore.After(time.Now().AddDate(0, 0, -1)))
	assert.True(t, certificate.NotAfter.Before(time.Now().AddDate(0, 0, 1)))
}
