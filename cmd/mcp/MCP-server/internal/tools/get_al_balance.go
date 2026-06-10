package tools

import (
	"context"

	"bank-mcp-server/internal/client"
	"bank-mcp-server/internal/types"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type GetAllBalancesTool struct {
	BankClient *client.BankClient
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
	request *mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {

	var response types.MCPResponse

	err := t.BankClient.Get(
		"/api/v1/accounts/balances",
		&response,
	)

	if err != nil {
		return nil, err
	}

	return JSONResponse(response)
}
