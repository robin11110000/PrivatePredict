package aleo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type BetRequest struct {
	MarketID uint   `json:"market_id"`
	Amount   int64  `json:"amount"`
	Outcome  string `json:"outcome"`
}

type BetResponse struct {
	TxID string `json:"tx_id"`
}

func PlacePrivateBet(marketID uint, amount int64, outcome string) (string, error) {
	payload := BetRequest{
		MarketID: marketID,
		Amount:   amount,
		Outcome:  outcome,
	}

	body, _ := json.Marshal(payload)

	resp, err := http.Post(
		"http://localhost:3001/prove",
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return "", fmt.Errorf("aleo prover unreachable: %w", err)
	}
	defer resp.Body.Close()

	var result BetResponse
	json.NewDecoder(resp.Body).Decode(&result)

	return result.TxID, nil
}
