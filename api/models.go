package api

type AmountOpRequestBody struct {
	UserID string  `json:"user_id"`
	Amount float64 `json:"amount"`
}

type GetBalanceResponse struct {
	UserID  string  `json:"user_id"`
	Balance float64 `json:"balance"`
}
