package v1

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/bookqaq/mer-wrapper/common"
)

func fetch(baseURL string, data searchData) (ResultData, error) {
	url_ := fmt.Sprintf("%s?%s", baseURL, data.paramGet().Encode())

	DPOP := dPoPGenerator(common.Client.ClientID, "get", searchURL)

	req, err := http.NewRequest("GET", url_, nil)
	if err != nil {
		err = fmt.Errorf("error creating Request at fetch:\n%s", err)
		return ResultData{}, err
	}
	req.Header.Add("DPOP", DPOP)
	req.Header.Add("X-Platform", "web")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Accept-Encoding", "deflate, gzip")

	resp, err := common.Client.Content.Do(req)
	if err != nil {
		err = fmt.Errorf("error accessing page at fetch:\n%s", err)
		return ResultData{}, err
	}
	defer resp.Body.Close()

	gzReader, err := gzip.NewReader(resp.Body)
	if err != nil {
		err = fmt.Errorf("creating gzip reader fail at fetch:\n%s", err)
		return ResultData{}, err
	}

	result, err := io.ReadAll(gzReader)
	if err != nil {
		err = fmt.Errorf("decode fail at fetch:\n%s", err)
		return ResultData{}, err
	}

	var content ResultData
	err = json.Unmarshal(result, &content)
	if err != nil {
		err = fmt.Errorf("json parse fail at fetch:\n%s", err)
		return ResultData{}, err
	}

	return content, nil
}

func Mercari_search(name string, sort string, order string, status string, limit int, times int) ([]MercariItem, error) {
	search := searchData{
		Keyword: name,
		Limit:   limit,
		Sort:    sort,
		Page:    0,
		Order:   order,
		Status:  status,
	}

	results := make([]MercariItem, 0)

	for search.Page < times {
		items, err := fetch(searchURL, search)
		if err != nil {
			return nil, err
		}
		if len(items.Data) <= 0 {
			break
		}
		results = append(results, items.Data...)
		if !items.Meta.HasNext {
			break
		}
		search.Page += 1
	}

	return results, nil
}
