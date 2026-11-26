package paystack

/*
	Putting this in its own package (like internal/paystack) keeps external API
	calls completely decoupled from our HTTP handlers and database repositories.
*/

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)
const baseURL = "https://api.paystack.co"
// Client is defined right here in package paystack
type Client struct {
	secretKey  string
	httpClient *http.Client
}
type InitializeRequest struct {
	Email     string `json:"email"`
	Amount    int    `json:"amount"` // in kobo (naira * 100)
	Reference string `json:"reference"`
}
type InitializeResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		AuthorizationURL string `json:"authorization_url"`
		AccessCode       string `json:"access_code"`
		Reference        string `json:"reference"`
	} `json:"data"`
}

func NewClient(secretKey string) *Client {
	return &Client{
		secretKey: secretKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) InitializeTransaction(ctx context.Context, email string, amountInNaira int, reference string) (*InitializeResponse, error) {
	if email == "" {
		return nil, errors.New("email is required for Paystack payments")
	}

	payload := InitializeRequest{
		Email: email,
		Amount: amountInNaira * 100, // naira -> kobo for paystack
		Reference: reference,
	}

	// Converts the payload to json
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal paystack payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+"/transaction/initialize", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create paystack requests: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.secretKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call paystack api: %w", err)
	}
	defer resp.Body.Close()

	var initResp InitializeResponse
	if err := json.NewDecoder(resp.Body).Decode(&initResp);
	err != nil {
		return nil, fmt.Errorf("failed to decode paystack response: %w", err)
	}

	if resp.StatusCode != http.StatusOK || !initResp.Status {
		return nil, fmt.Errorf("paystack rejected transaction: %s", initResp.Message)
	}

	return &initResp, nil
}
