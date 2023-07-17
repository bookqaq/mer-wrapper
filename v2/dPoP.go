package v2

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/bookqaq/mer-wrapper/common"
)

type ECDSASignature struct {
	R, S *big.Int
}

type payload struct {
	Iat int64  `json:"iat"`
	Jti string `json:"jti"`
	Htu string `json:"htu"`
	Htm string `json:"htm"`
	Uid string `json:"uuid"`
}

type pkey_jwk struct {
	Crv string `json:"crv"`
	Kty string `json:"kty"`
	X   string `json:"x"`
	Y   string `json:"y"`
}

type pkey_header struct {
	Typ string   `json:"typ"`
	Alg string   `json:"alg"`
	Jwk pkey_jwk `json:"jwk"`
}

func byteToBase64URL(target []byte) string {
	return base64.RawURLEncoding.EncodeToString(target)
}

// a function that shouldn't fail
func dPoPGenerator(uuid_ string, method string, url_ string) string { //因为有 url和uuid 包了
	private_key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		fmt.Println("error at dPoPGenerator/ecdsa.GenerateKey():\n", err)
		os.Exit(60)
	}

	pl := payload{time.Now().Unix(), uuid_, url_, strings.ToUpper(method), common.Client.ClientID}
	pkjwk := pkey_jwk{"P-256", "EC", byteToBase64URL(private_key.PublicKey.X.Bytes()), byteToBase64URL(private_key.PublicKey.Y.Bytes())}
	pkh := pkey_header{"dpop+jwt", "ES256", pkjwk}

	headerString, err := json.Marshal(pkh)
	if err != nil {
		fmt.Println("error at dPoPGenerator/json.Marshal(pkh):\n", err)
		os.Exit(61)
	}
	payloadString, err := json.Marshal(pl)
	if err != nil {
		fmt.Println("error at dPoPGenerator/json.Marshal(pl):\n", err)
		os.Exit(62)
	}

	data_unsigned := fmt.Sprintf("%s.%s", byteToBase64URL(headerString), byteToBase64URL(payloadString))

	hval := sha256.Sum256([]byte(data_unsigned))

	r, s, err := ecdsa.Sign(rand.Reader, private_key, hval[:])
	if err != nil {
		fmt.Println(err)
		os.Exit(63)
	}
	// sig := &ECDSASignature{}
	// if _, err := asn1.Unmarshal(signature, sig); err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(64)
	// }

	signatured := append(r.Bytes(), s.Bytes()...)

	signaturedString := byteToBase64URL(signatured)

	result := fmt.Sprintf("%s.%s", data_unsigned, signaturedString)
	return result
}

func generateSearchSessionId(length int) string {
	buflen := length
	if buflen%2 != 0 {
		buflen += 1
	}
	buf := make([]byte, buflen/2)
	rand.Read(buf)
	return hex.EncodeToString(buf)[:length]
}
