package tools

import (
	"context"
	"encoding/json"

	"bank-mcp-server/internal/client"
	"bank-mcp-server/internal/types"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type GetAllBalancesTool struct {
	BankClient *client.BankClient
}

type GetAllBalancesInput struct{}

type GetAllBalancesOutput struct {
	Content string `json:"content"`
}

func NewGetAllBalancesTool(
	bankClient *client.BankClient,
) *GetAllBalancesTool {

	return &GetAllBalancesTool{
		BankClient: bankClient,
	}
}

func (t *GetAllBalancesTool) Definition() *mcp.Tool {
	return &mcp.Tool{
		Name:        "get_all_balances",
		Description: "Get balances of all user accounts",
	}
}

func (t *GetAllBalancesTool) Handler(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input GetAllBalancesInput,
) (*mcp.CallToolResult, GetAllBalancesOutput, error) {

	var response types.MCPResponse

	err := t.BankClient.Get(
		"/api/v1/accounts/balances",
		&response,
	)

	if err != nil {
		return nil, GetAllBalancesOutput{}, err
	}

	bytes, err := json.MarshalIndent(
		response,
		"",
		"  ",
	)

	if err != nil {
		return nil, GetAllBalancesOutput{}, err
	}

	return nil, GetAllBalancesOutput{
		Content: string(bytes),
	}, nil
}
