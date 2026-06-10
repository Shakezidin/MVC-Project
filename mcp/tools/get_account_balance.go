package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/banking/bank-server/mcp/client"
	"github.com/banking/bank-server/mcp/types"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type GetAccountBalanceTool struct {
	BankClient *client.BankClient
}

type GetAccountBalanceInput struct {
	AccountID string `json:"account_id"`
}

type GetAccountBalanceOutput struct {
	Content string `json:"content"`
}

func NewGetAccountBalanceTool(
	bankClient *client.BankClient,
) *GetAccountBalanceTool {

	return &GetAccountBalanceTool{
		BankClient: bankClient,
	}

}

func (t *GetAccountBalanceTool) Definition() *mcp.Tool {
	return &mcp.Tool{
		Name:        "get_account_balance",
		Description: "Get balance of specific account",
	}
}

func (t *GetAccountBalanceTool) Handler(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input GetAccountBalanceInput,
) (*mcp.CallToolResult, GetAccountBalanceOutput, error) {

	if input.AccountID == "" {
		return nil, GetAccountBalanceOutput{}, fmt.Errorf(
			"account_id is required",
		)
	}

	var response types.MCPResponse

	path := fmt.Sprintf(
		"/api/v1/accounts/%s/balance",
		input.AccountID,
	)

	authTocken := req.Extra.Header.Get("Authorization")
	err := t.BankClient.Get(
		path,
		authTocken,
		&response,
	)

	if err != nil {
		return nil, GetAccountBalanceOutput{}, err
	}

	bytes, err := json.MarshalIndent(
		response,
		"",
		"  ",
	)

	if err != nil {
		return nil, GetAccountBalanceOutput{}, err
	}

	return nil,
		GetAccountBalanceOutput{
			Content: string(bytes),
		},
		nil
}
