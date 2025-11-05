package models

type BroadcastRequest struct {
	Symbol    string `json:"symbol" validate:"required"`
	Price     uint64 `json:"price" validate:"required"`
	Timestamp uint64 `json:"timestamp" validate:"required"`
}

type BroadcastResponse struct {
	TxHash   string `json:"tx_hash"`
	TxStatus string `json:"tx_status"`
}
