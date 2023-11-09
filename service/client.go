package service

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func (service *serviceImpl) JSONRequest(reqModel, resModel interface{}, url string) (interface{}, error) {
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
