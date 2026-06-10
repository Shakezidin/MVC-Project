package tools

import (
	"context"
	"encoding/json"

	"github.com/banking/bank-server/mcp/client"
	"github.com/banking/bank-server/mcp/types"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type GetAccountsTool struct {
	BankClient *client.BankClient
}

type GetAccountsInput struct{}

type GetAccountsOutput struct {
	Content string `json:"content"`
}

func NewGetAccountsTool(bankClient *client.BankClient) *GetAccountsTool {
	return &GetAccountsTool{
		BankClient: bankClient,
	}

}

func (t *GetAccountsTool) Definition() *mcp.Tool {
	return &mcp.Tool{
		Name:        "get_accounts",
		Description: "Get all bank accounts of logged in user",
	}
}

func (t *GetAccountsTool) Handler(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input GetAccountsInput,
) (*mcp.CallToolResult, GetAccountsOutput, error) {

	var response types.MCPResponse

	err := t.BankClient.Get(
		"/api/v1/accounts",
		&response,
	)

	if err != nil {
		return nil, GetAccountsOutput{}, err
	}

	bytes, err := json.MarshalIndent(
		response,
		"",
		"  ",
	)

	if err != nil {
		return nil, GetAccountsOutput{}, err
	}

	return nil,
		GetAccountsOutput{
			Content: string(bytes),
		},
		nil

}
