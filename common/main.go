package common

import "net/http"

type MerClient struct {
	ClientID string
	Content  *http.Client
}

var Client = MerClient{}
