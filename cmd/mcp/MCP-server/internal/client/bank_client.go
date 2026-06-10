package client

import (
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

type BankClient struct {
	client  *resty.Client
	baseURL string
}

func NewBankClient(
	baseURL string,
	token string,
) *BankClient {

	client := resty.New()

	client.SetTimeout(10 * time.Second)

	client.SetRetryCount(3)

	client.SetRetryWaitTime(1 * time.Second)

	client.SetHeader(
		"Content-Type",
		"application/json",
	)

	if token != "" {
		client.SetAuthToken(token)
	}

	return &BankClient{
		client:  client,
		baseURL: baseURL,
	}
}

func (b *BankClient) Get(
	path string,
	result interface{},
) error {

	resp, err := b.client.R().
		SetResult(result).
		Get(fmt.Sprintf("%s%s", b.baseURL, path))

	if err != nil {
		return err
	}

	if resp.StatusCode() >= 400 {
		return fmt.Errorf(
			"bank API returned status code %d",
			resp.StatusCode(),
		)
	}

	return nil
}
