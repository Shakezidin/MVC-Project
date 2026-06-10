package tools

import (
	"context"
	"encoding/json"

	"github.com/banking/bank-server/mcp/client"
	"github.com/banking/bank-server/mcp/types"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type GetTransferModesTool struct {
	BankClient *client.BankClient
}

type GetTransferModesInput struct{}

type GetTransferModesOutput struct {
	Content string `json:"content"`
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
	req *mcp.CallToolRequest,
	input GetTransferModesInput,
) (*mcp.CallToolResult, GetTransferModesOutput, error) {
	var response types.MCPResponse

	err := t.BankClient.Get(
		"/api/v1/transfer-modes",
		&response,
	)

	if err != nil {
		return nil, GetTransferModesOutput{}, err
	}

	bytes, err := json.MarshalIndent(
		response,
		"",
		"  ",
	)

	if err != nil {
		return nil, GetTransferModesOutput{}, err
	}
	return nil,
		GetTransferModesOutput{
			Content: string(bytes),
		},
		nil
}
