package merwrapper

import "net/http"

type SearchData struct {
	Keyword     string
	Limit       int
	Sort        string
	TargetPrice [2]int
}

type MerClient struct {
	ClientID string
	Content  *http.Client
}

const (
	rootURL = "https://api.mercari.jp/"
)

var searchParams = struct {
	searchURLv2  string
	searchMethod string
}{
	searchURLv2:  rootURL + "v2/entities:search",
	searchMethod: http.MethodPost,
}

var itemParams = struct {
	itemURL    string
	itemMethod string
}{
	itemURL:    rootURL + "items/get",
	itemMethod: http.MethodGet,
}

type MercariV2Search struct {
	IndexRouting    string                `json:"indexRouting"`
	SearchSessionId string                `json:"searchSessionId"`
	Page            int                   `json:"pageSize"`
	SearchCondition MercariV2SearchDetail `json:"searchCondition"`
}

type MercariV2SearchDetail struct {
	Keyword   string   `json:"keyword"`
	Sort      string   `json:"sort"`
	PriceMax  int      `json:"priceMax"`
	PriceMin  int      `json:"priceMin"`
	HasCoupon bool     `json:"hasCoupon"`
	Status    []string `json:"status"`
	Order     string   `json:"order"`
}

type MercariV2Item struct {
	ProductId     string   `json:"id"`
	ProductName   string   `json:"name"`
	Price         int      `json:"price"`
	Created       string   `json:"created"`
	Updated       string   `json:"updated"`
	Condition     string   `json:"itemConditionId"`
	ImageURL      []string `json:"thumbnails"`
	Status        string   `json:"status"`
	Seller        string   `json:"sellerId"`
	Buyer         string   `json:"buyerId"`
	ShippingPayer string   `json:"shippingPayerId"` // 0: by seller
}

type Name_Id_Unit struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type ItemResultData struct {
	Result string                 `json:"result"`
	Meta   map[string]interface{} `json:"meta"`
	Data   MercariDetail          `json:"data"`
}

type MercariDetail struct {
	ProductId   string         `json:"id"`
	ProductName string         `json:"name"`
	Price       int            `json:"price"`
	Seller      ItemSellerInfo `json:"seller"`
	Status      string         `json:"status"`
}

type ItemSellerInfo struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
