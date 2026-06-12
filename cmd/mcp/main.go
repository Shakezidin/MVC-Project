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
	cfg := config.Load()

	log, err := observability.NewLogger(
		"MCP-server",
		cfg.PubSub.ProjectID,
		cfg.PubSub.TopicID,
		cfg.Log.Level,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fatal error creating logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()

	log.Info("starting MCP server", zap.String("config", cfg.String()))

	bankClient := client.NewBankClient(
		cfg.BankAPIBaseURL,
	)

	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    "bank-mcp-server",
			Version: "1.0.0",
		},
		nil,
	)

	getAccountsTool := tools.NewGetAccountsTool(bankClient)
	getAllBalancesTool := tools.NewGetAllBalancesTool(bankClient)
	getAccountBalanceTool := tools.NewGetAccountBalanceTool(bankClient)
	getBeneficiariesTool := tools.NewGetBeneficiariesTool(bankClient)
	getTransferModesTool := tools.NewGetTransferModesTool(bankClient)

	mcp.AddTool(server, getAccountsTool.Definition(), getAccountsTool.Handler)
	mcp.AddTool(server, getAllBalancesTool.Definition(), getAllBalancesTool.Handler)
	mcp.AddTool(server, getAccountBalanceTool.Definition(), getAccountBalanceTool.Handler)
	mcp.AddTool(server, getBeneficiariesTool.Definition(), getBeneficiariesTool.Handler)
	mcp.AddTool(server, getTransferModesTool.Definition(), getTransferModesTool.Handler)

	log.Info("Starting Bank MCP Server...")

	handler := http.Handler(mcp.NewStreamableHTTPHandler(func(req *http.Request) *mcp.Server {
		return server
	}, nil))

	if err := http.ListenAndServe(":8082", handler); err != nil {
		log.Fatal("Server failed", zap.Error(err))
	}
}
