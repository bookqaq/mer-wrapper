package v2

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/bookqaq/mer-wrapper/common"
	"github.com/google/uuid"
)

func searchParse(p SearchData) ([]byte, error) {
	sp := MercariV2Search{}
	sp.DefaultDatabases = []string{"DATASET_TYPE_MERCARI", "DATASET_TYPE_BEYOND"}
	sp.IndexRouting = "INDEX_ROUTING_UNSPECIFIED"
	sp.PageSize = p.Limit
	// *sp.PageToken = "v1:0"	// only search one page for all tasks

	// searchCondition
	sp.SearchCondition.HasCoupon = false
	sp.SearchCondition.Status = []string{"STATUS_ON_SALE"}
	// sp.SearchCondition.Sort = p.Sort
	sp.SearchCondition.Sort = SearchOptionSortCreatedTime
	sp.SearchCondition.Order = SearchOptionOrderDESC
	sp.SearchCondition.Keyword = p.Keyword
	if len(p.TargetPrice) == 2 && p.TargetPrice[0] >= 0 && p.TargetPrice[0] <= p.TargetPrice[1] {
		sp.SearchCondition.PriceMin = p.TargetPrice[0]
		sp.SearchCondition.PriceMax = p.TargetPrice[1]
	}

	sp.SearchSessionId = generateSearchSessionId(DefaultLengthSearchSessionId)
	sp.ServiceFrom = "suruga"
	sp.WithItemBrand = true
	sp.WithItemPromotions = true
	sp.WithItemSize = false
	sp.WithItemSizes = true
	sp.WithShopName = false

	res, err := json.Marshal(sp)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func fetch(req *http.Request) ([]byte, error) {
	resp, err := common.Client.Content.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	//fmt.Println(resp.Status)

	gzReader, err := gzip.NewReader(resp.Body)
	if err != nil {
		return nil, err
	}

	result, err := io.ReadAll(gzReader)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func Search(data SearchData) ([]MercariV2Item, error) {
	sdata, err := searchParse(data)
	if err != nil {
		return nil, err
	}
	dPoP := dPoPGenerator(uuid.NewString(), searchParams.searchMethod, searchParams.searchURLv2)
	req, err := http.NewRequest(searchParams.searchMethod, searchParams.searchURLv2, bytes.NewReader(sdata))
	if err != nil {
		return nil, err
	}
	req.Header.Add("dpop", dPoP)
	req.Header.Add("accept-encoding", "gzip, deflate, br")
	req.Header.Add("accept", "application/json, text/plain, */*")
	req.Header.Add("x-platform", "web")

	res, err := fetch(req)
	if err != nil {
		return nil, err
	}

	var result MercariV2SearchResponse
	err = json.Unmarshal(res, &result)
	if err != nil {
		return nil, err
	}
	return result.Items, nil
}

func Item(item string) (MercariDetail, error) {
	reqVal := url.Values{}
	reqVal.Add("id", item)
	url := fmt.Sprintf("%s?%s", itemParams.itemURL, reqVal.Encode())
	dPoP := dPoPGenerator(uuid.NewString(), itemParams.itemMethod, itemParams.itemURL)
	req, err := http.NewRequest(itemParams.itemMethod, url, nil)
	if err != nil {
		return MercariDetail{}, err
	}
	req.Header.Add("dpop", dPoP)
	req.Header.Add("accept-encoding", "gzip, deflate, br")
	req.Header.Add("accept", "application/json, text/plain, */*")
	req.Header.Add("x-platform", "web")
	res, err := fetch(req)
	if err != nil {
		return MercariDetail{}, err
	}

	var result ItemResultData
	err = json.Unmarshal(res, &result)
	if err != nil {
		return MercariDetail{}, err
	}
	return result.Data, nil
}
