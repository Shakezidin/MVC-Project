package tools

import (
	"context"

	"bank-mcp-server/internal/client"
	"bank-mcp-server/internal/types"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type GetAccountsTool struct {
	BankClient *client.BankClient
}

func NewGetAccountsTool(
	bankClient *client.BankClient,
) *GetAccountsTool {

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
	request *mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {

	var response types.MCPResponse

	err := t.BankClient.Get(
		"/api/v1/accounts",
		&response,
	)

	if err != nil {
		return nil, err
	}

	return JSONResponse(response)
}
