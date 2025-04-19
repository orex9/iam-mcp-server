package aws

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/mark3labs/mcp-go/mcp"
)

type Policy struct {
	Name string
	Type string
	Arn  string
}

type Role struct {
	Name string
}

func handlePolicyListTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

func handleRolesListTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to load aws config: %v", err)), err
	}

	client := iam.NewFromConfig(cfg)

	var roles []Role

	var marker *string

	for {
		output, err := client.ListRoles(ctx, &iam.ListRolesInput{
			PathPrefix: aws.String("/"),
			Marker:     marker,
		})
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to list iam roles: %v", err)), err
		}
		for _, r := range output.Roles {
			roles = append(roles, Role{Name: *r.RoleName})
		}

		if output.IsTruncated {
			marker = output.Marker
		} else {
			break
		}

	}

	res, err := json.Marshal(roles)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to marshal roles struct")), err
	}

	return mcp.NewToolResultText(string(res)), nil
}

func handleGetRolePoliciesTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	arguments := request.Params.Arguments
	role, ok := arguments["role"].(string)
	if !ok {
		return mcp.NewToolResultError("failed to get argument: role"), fmt.Errorf("failed to get argument: role")
	}
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to load aws config: %v", err)), err
	}

	client := iam.NewFromConfig(cfg)

	var policies []Policy

	var inline_marker *string
	for {
		output, err := client.ListRolePolicies(ctx, &iam.ListRolePoliciesInput{
			RoleName: aws.String(role),
			Marker:   inline_marker,
		})
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to get role inline policies: %v", err)), err
		}
		for _, p := range output.PolicyNames {
			policies = append(policies, Policy{Name: p, Type: "inline"})
		}
		if output.IsTruncated {
			inline_marker = output.Marker
		} else {
			break
		}

	}

	var attached_marker *string
	for {
		output, err := client.ListAttachedRolePolicies(ctx, &iam.ListAttachedRolePoliciesInput{
			RoleName: aws.String(role),
			Marker:   attached_marker,
		})

		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to get role attached policies: %v", err)), err
		}

		for _, p := range output.AttachedPolicies {
			policies = append(policies, Policy{Name: *p.PolicyName, Type: "attached", Arn: *p.PolicyArn})
		}
		if output.IsTruncated {
			attached_marker = output.Marker
		} else {
			break
		}
	}

	res, err := json.Marshal(policies)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to marshal policies struct")), err
	}

	return mcp.NewToolResultText(string(res)), nil

}

func handleGetInlinePolicy(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	arguments := request.Params.Arguments
	policy, ok := arguments["policy"].(string)
	if !ok {
		return mcp.NewToolResultError("failed to get argument: policy"), fmt.Errorf("failed to get argument: policy")
	}
	role, ok := arguments["role"].(string)
	if !ok {
		return mcp.NewToolResultError("failed to get argument: role"), fmt.Errorf("failed to get argument: role")
	}
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to load aws config: %v", err)), err
	}

	client := iam.NewFromConfig(cfg)

	output, err := client.GetRolePolicy(ctx, &iam.GetRolePolicyInput{
		RoleName:   aws.String(role),
		PolicyName: aws.String(policy),
	})
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to get role inline policy: %v", err)), err
	}

	return mcp.NewToolResultText(string(*output.PolicyDocument)), nil

}

func handleGetAttachedPolicy(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	arguments := request.Params.Arguments
	policyArn, ok := arguments["policy_arn"].(string)
	if !ok {
		return mcp.NewToolResultError("failed to get argument: policy"), fmt.Errorf("failed to get argument: policy")
	}
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to load aws config: %v", err)), err
	}

	client := iam.NewFromConfig(cfg)

	output, err := client.GetPolicyVersion(ctx, &iam.GetPolicyVersionInput{
		PolicyArn: aws.String(policyArn),
		VersionId: aws.String("v1"),
	})
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to get role inline policy: %v", err)), err
	}

	return mcp.NewToolResultText(string(*output.PolicyVersion.Document)), nil

}
