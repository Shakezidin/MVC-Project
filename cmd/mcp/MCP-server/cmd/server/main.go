package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type BankClient struct {
	client  *resty.Client
	baseURL string
	token   string
}

func NewBankClient() *BankClient {
	baseURL := os.Getenv("BANK_API_BASE_URL")
	token := os.Getenv("JWT_TOKEN")

	client := resty.New()

	client.SetTimeout(10 * time.Second)
	client.SetRetryCount(3)
	client.SetRetryWaitTime(1 * time.Second)

	client.SetHeader("Content-Type", "application/json")

	if token != "" {
		client.SetAuthToken(token)
	}

	return &BankClient{
		client:  client,
		baseURL: baseURL,
		token:   token,
	}
}

func (b *BankClient) Get(path string, result interface{}) error {
	resp, err := b.client.R().
		SetResult(result).
		Get(fmt.Sprintf("%s%s", b.baseURL, path))

	if err != nil {
		return err
	}

	if resp.StatusCode() >= 400 {
		return fmt.Errorf("bank API returned status: %d", resp.StatusCode())
	}

	return nil
}

type MCPResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func jsonResponse(data interface{}) (*mcp.CallToolResult, error) {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return nil, err
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: string(bytes),
			},
		},
	}, nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	bankClient := NewBankClient()

	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    "bank-mcp-server",
			Version: "1.0.0",
		},
		nil,
	)

	/*
		==================================================
		GET ACCOUNTS TOOL
		==================================================
	*/
	getAccountsTool := &mcp.Tool{
		Name:        "get_accounts",
		Description: "Get all bank accounts of the logged in user",
	}

	mcp.AddTool(
		server,
		getAccountsTool,
		func(
			ctx context.Context,
			request *mcp.CallToolRequest,
		) (*mcp.CallToolResult, error) {

			var response MCPResponse

			err := bankClient.Get("/api/v1/accounts", &response)
			if err != nil {
				return nil, err
			}

			return jsonResponse(response)
		},
	)

	/*
		==================================================
		GET ALL BALANCES TOOL
		==================================================
	*/
	getAllBalancesTool := &mcp.Tool{
		Name:        "get_all_balances",
		Description: "Get balances of all user bank accounts",
	}

	mcp.AddTool(
		server,
		getAllBalancesTool,
		func(
			ctx context.Context,
			request *mcp.CallToolRequest,
		) (*mcp.CallToolResult, error) {

			var response MCPResponse

			err := bankClient.Get("/api/v1/accounts/balances", &response)
			if err != nil {
				return nil, err
			}

			return jsonResponse(response)
		},
	)

	/*
		==================================================
		GET ACCOUNT BALANCE TOOL
		==================================================
	*/
	getAccountBalanceTool := &mcp.Tool{
		Name:        "get_account_balance",
		Description: "Get balance of a specific account",

		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"account_id": map[string]interface{}{
					"type":        "string",
					"description": "Bank account UUID",
				},
			},
			Required: []string{"account_id"},
		},
	}

	mcp.AddTool(
		server,
		getAccountBalanceTool,
		func(
			ctx context.Context,
			request *mcp.CallToolRequest,
		) (*mcp.CallToolResult, error) {

			accountID, ok := request.Params.Arguments["account_id"].(string)
			if !ok || accountID == "" {
				return nil, fmt.Errorf("account_id is required")
			}

			var response MCPResponse

			path := fmt.Sprintf(
				"/api/v1/accounts/%s/balance",
				accountID,
			)

			err := bankClient.Get(path, &response)
			if err != nil {
				return nil, err
			}

			return jsonResponse(response)
		},
	)

	/*
		==================================================
		GET BENEFICIARIES TOOL
		==================================================
	*/
	getBeneficiariesTool := &mcp.Tool{
		Name:        "get_beneficiaries",
		Description: "Get all beneficiaries of logged in user",
	}

	mcp.AddTool(
		server,
		getBeneficiariesTool,
		func(
			ctx context.Context,
			request *mcp.CallToolRequest,
		) (*mcp.CallToolResult, error) {

			var response MCPResponse

			err := bankClient.Get("/api/v1/beneficiaries", &response)
			if err != nil {
				return nil, err
			}

			return jsonResponse(response)
		},
	)

	/*
		==================================================
		GET TRANSFER MODES TOOL
		==================================================
	*/
	getTransferModesTool := &mcp.Tool{
		Name:        "get_transfer_modes",
		Description: "Get available bank transfer modes like UPI, NEFT, RTGS, IMPS",
	}

	mcp.AddTool(
		server,
		getTransferModesTool,
		func(
			ctx context.Context,
			request *mcp.CallToolRequest,
		) (*mcp.CallToolResult, error) {

			var response MCPResponse

			err := bankClient.Get("/api/v1/transfer-modes", &response)
			if err != nil {
				return nil, err
			}

			return jsonResponse(response)
		},
	)

	log.Println("Starting Bank MCP Server...")

	err = server.Run(
		context.Background(),
		&mcp.StdioTransport{},
	)

	if err != nil {
		log.Fatalf("failed to start MCP server: %v", err)
	}
}
