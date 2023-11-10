package service

import (
	"bytes"
	"encoding/json"
	"net/http"
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name=Client
type Client interface {
	JSONRequest(reqModel, resModel interface{}, url string) (interface{}, error)
}

type ClientJSON struct {
}

func (c ClientJSON) JSONRequest(reqModel, resModel interface{}, url string) (interface{}, error) {
	body, err := json.Marshal(&reqModel)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&resModel)
	if err != nil {
		return nil, err
	}

	return resModel, nil
}
