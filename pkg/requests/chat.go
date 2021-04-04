package requests

type ChatRequest struct {
	Text string `json:"text"`
}

type ChatRespnse struct {
	Text string `json:"text"`
}
