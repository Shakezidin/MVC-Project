package tools

import (
	"context"
	"encoding/json"

	"github.com/banking/bank-server/mcp/client"
	"github.com/banking/bank-server/mcp/types"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type GetBeneficiariesTool struct {
	BankClient *client.BankClient
}

type GetBeneficiariesInput struct{}

type GetBeneficiariesOutput struct {
	Content string `json:"content"`
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
	req *mcp.CallToolRequest,
	input GetBeneficiariesInput,
) (*mcp.CallToolResult, GetBeneficiariesOutput, error) {

	var response types.MCPResponse

	err := t.BankClient.Get(
		"/api/v1/beneficiaries",
		&response,
	)

	if err != nil {
		return nil, GetBeneficiariesOutput{}, err
	}

	bytes, err := json.MarshalIndent(
		response,
		"",
		"  ",
	)

	if err != nil {
		return nil, GetBeneficiariesOutput{}, err
	}

	return nil,
		GetBeneficiariesOutput{
			Content: string(bytes),
		},
		nil
}
