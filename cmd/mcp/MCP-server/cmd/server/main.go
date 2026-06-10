package main

import (
	"context"
	"log"

	"bank-mcp-server/internal/client"
	"bank-mcp-server/internal/config"
	"bank-mcp-server/internal/tools"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	/*
		==================================================
		LOAD CONFIG
		==================================================
	*/
	cfg := config.Load()

	/*
		==================================================
		INITIALIZE BANK CLIENT
		==================================================
	*/
	bankClient := client.NewBankClient(
		cfg.BankAPIBaseURL,
		cfg.JWTToken,
	)

	/*
		==================================================
		CREATE MCP SERVER
		==================================================
	*/
	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    "bank-mcp-server",
			Version: "1.0.0",
		},
		nil,
	)

	/*
		==================================================
		INITIALIZE TOOLS
		==================================================
	*/
	getAccountsTool := tools.NewGetAccountsTool(
		bankClient,
	)

	getAllBalancesTool := tools.NewGetAllBalancesTool(
		bankClient,
	)

	getAccountBalanceTool := tools.NewGetAccountBalanceTool(
		bankClient,
	)

	getBeneficiariesTool := tools.NewGetBeneficiariesTool(
		bankClient,
	)

	getTransferModesTool := tools.NewGetTransferModesTool(
		bankClient,
	)

	/*
		==================================================
		REGISTER TOOLS
		==================================================
	*/
	mcp.AddTool(
		server,
		getAccountsTool.Definition(),
		getAccountsTool.Handler,
	)

	mcp.AddTool(
		server,
		getAllBalancesTool.Definition(),
		getAllBalancesTool.Handler,
	)

	mcp.AddTool(
		server,
		getAccountBalanceTool.Definition(),
		getAccountBalanceTool.Handler,
	)

	mcp.AddTool(
		server,
		getBeneficiariesTool.Definition(),
		getBeneficiariesTool.Handler,
	)

	mcp.AddTool(
		server,
		getTransferModesTool.Definition(),
		getTransferModesTool.Handler,
	)

	/*
		==================================================
		START SERVER
		==================================================
	*/
	log.Println("Starting Bank MCP Server...")

	err := server.Run(
		context.Background(),
		&mcp.StdioTransport{},
	)

	if err != nil {
		log.Fatalf(
			"failed to start MCP server: %v",
			err,
		)
	}
}
