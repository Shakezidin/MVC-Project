package tools

import (
	"context"

	"bank-mcp-server/internal/client"
	"bank-mcp-server/internal/types"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type GetTransferModesTool struct {
	BankClient *client.BankClient
}

func NewGetTransferModesTool(
	bankClient *client.BankClient,
) *GetTransferModesTool {

	return &GetTransferModesTool{
		BankClient: bankClient,
	}
}

func (t *GetTransferModesTool) Definition() *mcp.Tool {
	return &mcp.Tool{
		Name:        "get_transfer_modes",
		Description: "Get available transfer modes",
	}
}

func (t *GetTransferModesTool) Handler(
	ctx context.Context,
	request *mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {

	var response types.MCPResponse

	err := t.BankClient.Get(
		"/api/v1/transfer-modes",
		&response,
	)

	if err != nil {
		return nil, err
	}

	return JSONResponse(response)
}
