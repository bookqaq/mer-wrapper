package analysis

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"
)

//go:nolint
func buildECPrivKeys(DBase64 string, XBase64 string, YBase64 string) (*ecdsa.PrivateKey, error) {
	d, err := base64.RawURLEncoding.DecodeString(DBase64)
	if err != nil {
		return nil, err
	}
	x, err := base64.RawURLEncoding.DecodeString(XBase64)
	if err != nil {
		return nil, err
	}
	y, err := base64.RawURLEncoding.DecodeString(YBase64)
	if err != nil {
		return nil, err
	}

	return &ecdsa.PrivateKey{
		D: new(big.Int).SetBytes(d),
		PublicKey: ecdsa.PublicKey{
			Curve: elliptic.P256(),
			X:     new(big.Int).SetBytes(x),
			Y:     new(big.Int).SetBytes(y),
		},
	}, nil
}

func tryVerifySignature(json string, targetSignatureBase64 string, key *ecdsa.PrivateKey) error {
	targetSignature, err := base64.RawURLEncoding.DecodeString(targetSignatureBase64)
	if err != nil {
		return err
	}

	hashed := sha256.Sum256([]byte(json))

	// Try signature with ASN.1
	signatureWithAsn1, err := key.Sign(rand.Reader, hashed[:], nil)
	if err != nil {
		return err
	}

	fmt.Println(len(targetSignature), len(signatureWithAsn1))

	fmt.Printf("%x\n", signatureWithAsn1)
	if len(targetSignature) != len(signatureWithAsn1) {
		fmt.Println("not sha256-p256-asn1 signature, not the same length")
	}

	// Try signature without ASN.1
	r, s, err := ecdsa.Sign(rand.Reader, key, hashed[:])
	if err != nil {
		return err
	}

	signature := append(r.Bytes(), s.Bytes()...)

	fmt.Printf("%d %x\n", len(signature), signature)
	if len(targetSignature) != len(signature) {
		fmt.Println("not sha256-p256 signature, not the same length")
	}

	targetR, targetS := new(big.Int).SetBytes(targetSignature[:32]), new(big.Int).SetBytes(targetSignature[32:])

	// Try verify signature to determine what part of jwt is signed
	if ecdsa.Verify(&key.PublicKey, hashed[:], targetR, targetS) {
		fmt.Println("target signature verify success, can confirm it's a standard jwt (hash first two part base64 and then sign)")
	}

	return errors.New("an error that not letting go test cache result")
}
