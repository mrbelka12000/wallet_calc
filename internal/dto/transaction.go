package dto

type BaseTransaction struct {
	Date        string  `json:"date,omitempty"`
	Description string  `json:"description,omitempty"`
	Amount      float64 `json:"amount,omitempty"`
	Transaction string  `json:"transaction,omitempty"`
	Category    string  `json:"category,omitempty"`
	Confidence  float64 `json:"confidence,omitempty"`
}
