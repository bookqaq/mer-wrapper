package analysis

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"reflect"
	"testing"
)

func Test_buildECPrivKeys(t *testing.T) {
	type args struct {
		DBase64 string
		XBase64 string
		YBase64 string
	}
	tests := []struct {
		name    string
		args    args
		want    *ecdsa.PrivateKey
		wantErr bool
	}{
		{"sample private key",
			args{
				DBase64: "kRAB_FrWLQdSWnCiBFnjrhG6Jbqg6Fwtm7QrDYPC0mg",
				XBase64: "e5zwA7scbeI2653Z6hKV-ktJ9fDXAIce8GLDWcCl-Z0",
				YBase64: "9QXtjg2OCbLmesYFDl50u7dI690wMINClVz2HgZQNng",
			}, nil, false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := buildECPrivKeys(tt.args.DBase64, tt.args.XBase64, tt.args.YBase64)
			if (err != nil) != tt.wantErr {
				t.Errorf("buildECPrivKeys() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				priv_der, err := x509.MarshalPKCS8PrivateKey(got)
				if err != nil {
					t.Error(err)
				}

				pub_der, err := x509.MarshalPKIXPublicKey(&got.PublicKey)
				if err != nil {
					t.Error(err)
				}

				priv_pem := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: priv_der})
				pub_pem := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pub_der})

				fmt.Printf("%s\n", priv_pem)
				fmt.Printf("%s\n", pub_pem)

				t.Errorf("buildECPrivKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_tryVerifySignature(t *testing.T) {
	priv, err := buildECPrivKeys("kRAB_FrWLQdSWnCiBFnjrhG6Jbqg6Fwtm7QrDYPC0mg", "e5zwA7scbeI2653Z6hKV-ktJ9fDXAIce8GLDWcCl-Z0", "9QXtjg2OCbLmesYFDl50u7dI690wMINClVz2HgZQNng")
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		json                  string
		targetSignatureBase64 string
		key                   *ecdsa.PrivateKey
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"sample token", args{
			json:                  "eyJ0eXAiOiJkcG9wK2p3dCIsImFsZyI6IkVTMjU2IiwiandrIjp7ImNydiI6IlAtMjU2Iiwia3R5IjoiRUMiLCJ4IjoiZTV6d0E3c2NiZUkyNjUzWjZoS1Yta3RKOWZEWEFJY2U4R0xEV2NDbC1aMCIsInkiOiI5UVh0amcyT0NiTG1lc1lGRGw1MHU3ZEk2OTB3TUlOQ2xWejJIZ1pRTm5nIn19.eyJpYXQiOjE2ODk1NTg4OTIsImp0aSI6Ijg3MGNjZmNkLWFhZDAtNDY5ZC1hM2FjLWE1YTg4YWRmZDlhZiIsImh0dSI6Imh0dHBzOi8vYXBpLm1lcmNhcmkuanAvdjIvZW50aXRpZXM6c2VhcmNoIiwiaHRtIjoiUE9TVCIsInV1aWQiOiJmMjVhZWI0My01MTNkLTQwNDUtOTRhOC0yM2ZhNDYxOGIyNjUifQ",
			targetSignatureBase64: "J4EPhmNia_4AQfKenUi8xtSV94ru9DpXesx-1F-mh5-q1zYNpSOYvR7d7ERl9OcZGFDj9PYu51UzBSmVVhhPgA",
			key:                   priv,
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tryVerifySignature(tt.args.json, tt.args.targetSignatureBase64, tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("tryVerifySignature() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
