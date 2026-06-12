package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/banking/bank-server/internal/observability"
	"github.com/banking/bank-server/mcp/client"
	"github.com/banking/bank-server/mcp/config"
	"github.com/banking/bank-server/mcp/tools"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"go.uber.org/zap"
)

func main() {
	/*
		==================================================
		LOAD CONFIG
		==================================================
	*/
	cfg := config.Load()

	log, err := observability.NewLogger(
		"MCP-server",
		"tempBankLogs",
		"MCPServerLogs",
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fatal error: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()

	log.Sugar().Info("starting bank-server", zap.String("env", cfg.Log.Level))

	/*
		==================================================
		INITIALIZE BANK CLIENT
		==================================================
	*/
	bankClient := client.NewBankClient(
		cfg.BankAPIBaseURL,
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
	log.Sugar().Info("Starting Bank MCP Server...")

	handler := http.Handler(mcp.NewStreamableHTTPHandler(func(req *http.Request) *mcp.Server {
		return server
	}, nil))

	// Run the server over stdin/stdout, until the client disconnects.
	if err := http.ListenAndServe(":8082", handler); err != nil {
		log.Info("Server failed: %v")
	}
}
