package tools

import (
	"context"
	"fmt"

	"bank-mcp-server/internal/client"
	"bank-mcp-server/internal/types"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type GetAccountBalanceTool struct {
	BankClient *client.BankClient
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

		InputSchema: mcp.ToolInputSchema{
			Type: "object",

			Properties: map[string]interface{}{
				"account_id": map[string]interface{}{
					"type": "string",
				},
			},

			Required: []string{"account_id"},
		},
	}
}

func (t *GetAccountBalanceTool) Handler(
	ctx context.Context,
	request *mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {

	accountID, ok := request.Params.Arguments["account_id"].(string)

	if !ok || accountID == "" {
		return nil, fmt.Errorf("account_id is required")
	}

	var response types.MCPResponse

	path := fmt.Sprintf(
		"/api/v1/accounts/%s/balance",
		accountID,
	)

	err := t.BankClient.Get(
		path,
		&response,
	)

	if err != nil {
		return nil, err
	}

	return JSONResponse(response)
}
