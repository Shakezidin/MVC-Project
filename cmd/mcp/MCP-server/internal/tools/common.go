package tools

import (
	"encoding/json"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func JSONResponse(
	data interface{},
) (*mcp.CallToolResult, error) {

	bytes, err := json.MarshalIndent(
		data,
		"",
		"  ",
	)

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
