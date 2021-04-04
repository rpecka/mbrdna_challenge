package requests

type GetLocationResponse struct {
	Longitude LocValue `json:"longitude"`
	Latitude LocValue `json:"latitude"`
	Heading LocValue `json:"heading"`
}

type LocValue struct {
	Value float64 `json:"value"`
	RetrievalStatus string `json:"retrievalstatus"`
	Timestamp int `json:"timestamp"`
}
