package aws

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type Policy struct {
	Name string
}

func InitTools() server.ServerTool {
	iam_list_policies := mcp.NewTool("list", mcp.WithDescription("Lists all IAM policies"))

	return server.ServerTool{
		Tool:    iam_list_policies,
		Handler: handleListTool,
	}
}

func handleListTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to load aws config: %v", err)), err
	}

	client := iam.NewFromConfig(cfg)

	var policies []Policy

	var marker *string
	for {
		output, err := client.ListPolicies(ctx, &iam.ListPoliciesInput{
			Scope:  types.PolicyScopeTypeAll,
			Marker: marker,
		})

		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("got an error: %v", err)), err
		}

		for _, p := range output.Policies {
			policies = append(policies, Policy{Name: *p.PolicyName})
		}

		if output.IsTruncated {
			marker = output.Marker
		} else {
			break
		}

	}

	r, err := json.Marshal(policies)
	if err != nil {
		return mcp.NewToolResultError("failed to marshal response"), err
	}

	return mcp.NewToolResultText(string(r)), nil

}
