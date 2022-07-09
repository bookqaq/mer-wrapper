package v1

import (
	"fmt"
	"net/url"
)

const (
	rootURL        = "https://api.mercari.jp/"
	rootProductURL = "https://jp.mercari.com/item/"
	searchURL      = rootURL + "search_index/search"
)

type ResultData struct {
	Meta ResultMetaData `json:"meta"`
	Data []MercariItem  `json:"data"`
}

type ResultMetaData struct {
	HasNext  bool   `json:"has_next"`
	NextPage int    `json:"next_page"`
	Sort     string `json:"sort"`
	Order    string `json:"order"`
}

type MercariItem struct {
	ProductId   string       `json:"id"`
	ProductName string       `json:"name"`
	Price       int          `json:"price"`
	Created     int64        `json:"created"`
	Updated     int64        `json:"updated"`
	Condition   Name_Id_Unit `json:"item_condition"`
	ImageURL    []string     `json:"thumbnails"`
	Status      string       `json:"status"` // on_sale / trading / sold_out
	Seller      Name_Id_Unit `json:"seller"`
	Buyer       Name_Id_Unit `json:"buyer"`
	Shipping    Name_Id_Unit `json:"shipping_from_area"`
}

type Name_Id_Unit struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type searchData struct {
	Keyword string
	Limit   int
	Page    int
	Sort    string
	Order   string
	Status  string
}

func (item *MercariItem) GetProductURL() string {
	return rootProductURL + item.ProductId
}

func (data *searchData) paramGet() url.Values {
	params := url.Values{}
	params.Add("keyword", data.Keyword)
	params.Add("limit", fmt.Sprint(data.Limit))
	params.Add("page", fmt.Sprint(data.Page))
	params.Add("sort", data.Sort)
	params.Add("order", data.Order)
	params.Add("status", data.Status)
	return params
}
