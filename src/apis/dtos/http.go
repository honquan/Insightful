package dtos

type HttpResponse struct {
	Meta *MetaResp   `json:"meta"`
	Data interface{} `json:"data,omitempty"`
}
type MetaResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
