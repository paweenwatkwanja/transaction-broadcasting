package models

type BroadcastRequest struct {
	Symbol    string `json:"symbol"`
	Price     uint64 `json:"price"`
	Timestamp uint64 `json:"timestamp"`
}

type BroadcastResponse struct {
	TxHash   string `json:"tx_hash"`
	TxStatus string `json:"tx_status"`
}
