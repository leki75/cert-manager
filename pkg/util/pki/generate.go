/*
Copyright 2019 The Jetstack cert-manager contributors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package pki

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/leki75/cert-manager/pkg/apis/certmanager/v1alpha1"
)

const (
	// MinRSAKeySize is the minimum RSA keysize allowed to be generated by the
	// generator functions in this package.
	MinRSAKeySize = 2048

	// MaxRSAKeySize is the maximum RSA keysize allowed to be generated by the
	// generator functions in this package.
	MaxRSAKeySize = 8192

	// ECCurve256 represents a 256bit ECDSA key.
	ECCurve256 = 256
	// ECCurve384 represents a 384bit ECDSA key.
	ECCurve384 = 384
	// ECCurve521 represents a 521bit ECDSA key.
	ECCurve521 = 521
)

// GeneratePrivateKeyForCertificate will generate a private key suitable for
// the provided cert-manager Certificate resource, taking into account the
// parameters on the provided resource.
// The returned key will either be RSA or ECDSA.
func GeneratePrivateKeyForCertificate(crt *v1alpha1.Certificate) (crypto.Signer, error) {
	switch crt.Spec.KeyAlgorithm {
	case v1alpha1.KeyAlgorithm(""), v1alpha1.RSAKeyAlgorithm:
		keySize := MinRSAKeySize

		if crt.Spec.KeySize > 0 {
			keySize = crt.Spec.KeySize
		}

		return GenerateRSAPrivateKey(keySize)
	case v1alpha1.ECDSAKeyAlgorithm:
		keySize := ECCurve256

		if crt.Spec.KeySize > 0 {
			keySize = crt.Spec.KeySize
		}

		return GenerateECPrivateKey(keySize)
	default:
		return nil, fmt.Errorf("unsupported private key algorithm specified: %s", crt.Spec.KeyAlgorithm)
	}
}

// GenerateRSAPrivateKey will generate a RSA private key of the given size.
// It places restrictions on the minimum and maximum RSA keysize.
func GenerateRSAPrivateKey(keySize int) (*rsa.PrivateKey, error) {
	// Do not allow keySize < 2048
	// https://en.wikipedia.org/wiki/Key_size#cite_note-twirl-14
	if keySize < MinRSAKeySize {
		return nil, fmt.Errorf("weak rsa key size specified: %d. minimum key size: %d", keySize, MinRSAKeySize)
	}
	if keySize > MaxRSAKeySize {
		return nil, fmt.Errorf("rsa key size specified too big: %d. maximum key size: %d", keySize, MaxRSAKeySize)
	}

	return rsa.GenerateKey(rand.Reader, keySize)
}

// GenerateECPrivateKey will generate an ECDSA private key of the given size.
// It can be used to generate 256, 384 and 521 sized keys.
func GenerateECPrivateKey(keySize int) (*ecdsa.PrivateKey, error) {
	var ecCurve elliptic.Curve

	switch keySize {
	case ECCurve256:
		ecCurve = elliptic.P256()
	case ECCurve384:
		ecCurve = elliptic.P384()
	case ECCurve521:
		ecCurve = elliptic.P521()
	default:
		return nil, fmt.Errorf("unsupported ecdsa key size specified: %d", keySize)
	}

	return ecdsa.GenerateKey(ecCurve, rand.Reader)
}

// EncodePrivateKey will encode a given crypto.PrivateKey by first inspecting
// the type of key encoding and then inspecting the type of key provided.
// It only supports encoding RSA or ECDSA keys.
func EncodePrivateKey(pk crypto.PrivateKey, keyEncoding v1alpha1.KeyEncoding) ([]byte, error) {
	switch keyEncoding {
	case v1alpha1.KeyEncoding(""), v1alpha1.PKCS1:
		switch k := pk.(type) {
		case *rsa.PrivateKey:
			return EncodePKCS1PrivateKey(k), nil
		case *ecdsa.PrivateKey:
			return EncodeECPrivateKey(k)
		default:
			return nil, fmt.Errorf("error encoding private key: unknown key type: %T", pk)
		}
	case v1alpha1.PKCS8:
		return EncodePKCS8PrivateKey(pk)
	default:
		return nil, fmt.Errorf("error encoding private key: unknown key encoding: %s", keyEncoding)
	}
}

// EncodePKCS1PrivateKey will marshal a RSA private key into x509 PEM format.
func EncodePKCS1PrivateKey(pk *rsa.PrivateKey) []byte {
	block := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(pk)}

	return pem.EncodeToMemory(block)
}

// EncodePKCS8PrivateKey will marshal a private key into x509 PEM format.
func EncodePKCS8PrivateKey(pk interface{}) ([]byte, error) {
	keyBytes, err := x509.MarshalPKCS8PrivateKey(pk)
	if err != nil {
		return nil, err
	}
	block := &pem.Block{Type: "PRIVATE KEY", Bytes: keyBytes}

	return pem.EncodeToMemory(block), nil
}

// EncodeECPrivateKey will marshal an ECDSA private key into x509 PEM format.
func EncodeECPrivateKey(pk *ecdsa.PrivateKey) ([]byte, error) {
	asnBytes, err := x509.MarshalECPrivateKey(pk)
	if err != nil {
		return nil, fmt.Errorf("error encoding private key: %s", err.Error())
	}

	block := &pem.Block{Type: "EC PRIVATE KEY", Bytes: asnBytes}
	return pem.EncodeToMemory(block), nil
}

// PublicKeyForPrivateKey will return the crypto.PublicKey for the given
// crypto.PrivateKey. It only supports RSA and ECDSA keys.
func PublicKeyForPrivateKey(pk crypto.PrivateKey) (crypto.PublicKey, error) {
	switch k := pk.(type) {
	case *rsa.PrivateKey:
		return k.Public(), nil
	case *ecdsa.PrivateKey:
		return k.Public(), nil
	default:
		return nil, fmt.Errorf("unknown private key type: %T", pk)
	}
}

// PublicKeyMatchesCertificate can be used to verify the given public key
// is the correct counter-part to the given x509 Certificate.
// It will return false and no error if the public key is *not* valid for the
// given Certificate.
// It will return true if the public key *is* valid for the given Certificate.
// It will return an error if either of the passed parameters are of an
// unrecognised type (i.e. non RSA/ECDSA)
func PublicKeyMatchesCertificate(check crypto.PublicKey, crt *x509.Certificate) (bool, error) {
	switch pub := crt.PublicKey.(type) {
	case *rsa.PublicKey:
		rsaCheck, ok := check.(*rsa.PublicKey)
		if !ok {
			return false, nil
		}
		if pub.N.Cmp(rsaCheck.N) != 0 {
			return false, nil
		}
		return true, nil
	case *ecdsa.PublicKey:
		ecdsaCheck, ok := check.(*ecdsa.PublicKey)
		if !ok {
			return false, nil
		}
		if pub.X.Cmp(ecdsaCheck.X) != 0 || pub.Y.Cmp(ecdsaCheck.Y) != 0 {
			return false, nil
		}
		return true, nil
	default:
		return false, fmt.Errorf("unrecognised Certificate public key type")
	}
}

// PublicKeyMatchesCSR can be used to verify the given public key is the correct
// counter-part to the given x509 CertificateRequest.
// It will return false and no error if the public key is *not* valid for the
// given CertificateRequest.
// It will return true if the public key *is* valid for the given CertificateRequest.
// It will return an error if either of the passed parameters are of an
// unrecognised type (i.e. non RSA/ECDSA)
func PublicKeyMatchesCSR(check crypto.PublicKey, csr *x509.CertificateRequest) (bool, error) {
	switch pub := csr.PublicKey.(type) {
	case *rsa.PublicKey:
		rsaCheck, ok := check.(*rsa.PublicKey)
		if !ok {
			return false, nil
		}
		if pub.N.Cmp(rsaCheck.N) != 0 {
			return false, nil
		}
		return true, nil
	case *ecdsa.PublicKey:
		ecdsaCheck, ok := check.(*ecdsa.PublicKey)
		if !ok {
			return false, nil
		}
		if pub.X.Cmp(ecdsaCheck.X) != 0 || pub.Y.Cmp(ecdsaCheck.Y) != 0 {
			return false, nil
		}
		return true, nil
	default:
		return false, fmt.Errorf("unrecognised Certificate public key type")
	}
}
