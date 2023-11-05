package model

type WebResponse struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
	Data   string `json:"data"`
}
