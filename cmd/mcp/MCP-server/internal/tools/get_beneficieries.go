package tools

import (
	"context"

	"bank-mcp-server/internal/client"
	"bank-mcp-server/internal/types"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type GetBeneficiariesTool struct {
	BankClient *client.BankClient
}

func NewGetBeneficiariesTool(
	bankClient *client.BankClient,
) *GetBeneficiariesTool {

	return &GetBeneficiariesTool{
		BankClient: bankClient,
	}
}

func (t *GetBeneficiariesTool) Definition() *mcp.Tool {
	return &mcp.Tool{
		Name:        "get_beneficiaries",
		Description: "Get all beneficiaries",
	}
}

func (t *GetBeneficiariesTool) Handler(
	ctx context.Context,
	request *mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {

	var response types.MCPResponse

	err := t.BankClient.Get(
		"/api/v1/beneficiaries",
		&response,
	)

	if err != nil {
		return nil, err
	}

	return JSONResponse(response)
}
