package v2

import "net/http"

type SearchData struct {
	Keyword     string
	Limit       int
	Sort        string
	TargetPrice *[2]int
}

// const data
const (
	rootURL = "https://api.mercari.jp/"
)

const DefaultLengthSearchSessionId = 32

const (
	SearchOptionOrderDESC = "ORDER_DESC"
	SearchOptionOrderASC  = "ORDER_ASC"
)

const (
	SearchOptionSortScore       = "SORT_SCORE"
	SearchOptionSortCreatedTime = "SORT_CREATED_TIME"
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

type MerProduct interface {
	GetProductURL() string
}

// entities:search request body
type MercariV2Search struct {
	DefaultDatabases   []string                     `json:"defaultDatabases"` // have default value
	IndexRouting       string                       `json:"indexRouting"`     // have default value
	PageSize           int                          `json:"pageSize"`
	PageToken          *string                      `json:"pageToken,omitempty"` // pagination, stucture like "v1:1" "v1:2", not appear in first page
	SearchCondition    MercariV2SearchRequestDetail `json:"searchCondition"`
	SearchSessionId    string                       `json:"searchSessionId"`
	ServiceFrom        string                       `json:"serviceFrom"`    // have default value
	ThumbnailTypes     []any                        `json:"thumbnailTypes"` // default empty
	UserId             string                       `json:"userId"`
	WithItemBrand      bool                         `json:"withItemBrand"`
	WithItemPromotions bool                         `json:"withItemPromotions"`
	WithItemSize       bool                         `json:"withItemSize"`
	WithItemSizes      bool                         `json:"withItemSizes"`
	WithShopName       bool                         `json:"withShopName"`
}

// searchCondition part of request body
type MercariV2SearchRequestDetail struct {
	Attributes       []any    `json:"attributes"`      // default empty
	BrandId          []any    `json:"brandId"`         // default empty
	CategoryId       []any    `json:"categoryId"`      // default empty
	ColorId          []any    `json:"colorId"`         // default empty
	ExcludeKeyword   string   `json:"excludeKeyword"`  // TODO: check if this can achieve what it means
	HasCoupon        bool     `json:"hasCoupon"`       // default false
	ItemConditionId  []int    `json:"itemConditionId"` // default empty
	ItemTypes        []any    `json:"itemTypes"`       // default empty
	Keyword          string   `json:"keyword"`
	Order            string   `json:"order"`
	PriceMax         int      `json:"priceMax"`         // default empty
	PriceMin         int      `json:"priceMin"`         // default empty
	SellerId         []string `json:"sellerId"`         // default empty
	ShippingFromArea []any    `json:"shippingFromArea"` // default empty
	ShippingMethod   []any    `json:"shippingMethod"`   // default empty
	ShippingPayerId  []any    `json:"shippingPayerId"`  // default empty
	SizeId           []any    `json:"sizeId"`           // default empty
	SKUIds           []any    `json:"skuIds"`           // default empty
	Sort             string   `json:"sort"`
	Status           []string `json:"status"` // default empty
}

// entities:search response body
type MercariV2SearchResponse struct {
	Meta       any             `json:"meta"`
	Components any             `json:"components"`
	Items      []MercariV2Item `json:"items"`
}

// Search() result item that function return
type MercariV2Item struct {
	ProductId     string   `json:"id"`
	ProductName   string   `json:"name"`
	Price         string   `json:"price"`
	Created       string   `json:"created"`
	Updated       string   `json:"updated"`
	ImageURL      []string `json:"thumbnails"`
	ItemType      string   `json:"itemType"` // "ITEM_TYPE_MERCARI"
	Condition     string   `json:"itemConditionId"`
	ShippingPayer string   `json:"shippingPayerId"` // 0(or 2): by seller
	Status        string   `json:"status"`
	Seller        string   `json:"sellerId"`
	Buyer         string   `json:"buyerId"`
}

type Name_Id_Unit struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

// Item() response body
type ItemResultData struct {
	Result string                 `json:"result"`
	Meta   map[string]interface{} `json:"meta"`
	Data   MercariDetail          `json:"items"`
}

// Item() response body item part
type MercariDetail struct {
	ProductId         string         `json:"id"`
	ProductName       string         `json:"name"`
	Price             int            `json:"price"`
	Seller            ItemSellerInfo `json:"seller"`
	Status            string         `json:"status"`
	Description       string         `json:"description"`
	Condition         Name_Id_Unit   `json:"condition"`
	Like              int            `json:"num_likes"`
	Comment           int            `json:"num_comments"`
	ImageURL          []string       `json:"photos"`
	Created           int64          `json:"created"`
	Updated           int64          `json:"updated"`
	AnonymousShipping bool           `json:"is_anonymous_shipping"`
	ShippingDuration  Name_Id_Unit   `json:"shipping_duration"`
	ShippingFrom      Name_Id_Unit   `json:"shipping_from_area"`
	ShippingMethod    Name_Id_Unit   `json:"shipping_method"`
	ShippingPayer     Name_Id_Unit   `json:"shipping_payer"`
}

type ItemSellerInfo struct {
	ID           int64   `json:"id"`
	Name         string  `json:"name"`
	QuickShipper bool    `json:"quick_shipper"`
	NumSell      int32   `json:"num_sell_items"`
	Avatar       string  `json:"photo_thumbnail_url"`
	Created      int64   `json:"created"`
	SmsAuth      string  `json:"register_sms_confirmation"`
	SmsAuthAt    string  `json:"register_sms_confirmation_at"`
	Score        int     `json:"score"`
	Rating       float32 `json:"star_rating_score"`
}

type ItemSellerRating struct {
	Good   int32
	Bad    int32
	Normal int32
}
