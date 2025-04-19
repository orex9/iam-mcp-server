package aws

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func InitTools() []server.ServerTool {
	iam_list_policies := mcp.NewTool("list_policies", mcp.WithDescription("Lists all IAM policies"))
	iam_list_roles := mcp.NewTool("list_roles", mcp.WithDescription("Lists all IAM roles"))
	iam_get_role_policies := mcp.NewTool("get_role_policies", mcp.WithDescription("Gets policies attached to role"), mcp.WithString("role", mcp.Required()))
	iam_get_inline_policy := mcp.NewTool("get_inline_policy", mcp.WithDescription("Get role inline policy"), mcp.WithString("role", mcp.Required()), mcp.WithString("policy", mcp.Required()))
	iam_get_attached_policy := mcp.NewTool("get_attached_policy", mcp.WithDescription("Get role attached policy"), mcp.WithString("policy_arn", mcp.Required()))

	return []server.ServerTool{
		{
			Tool:    iam_list_policies,
			Handler: handlePolicyListTool,
		},
		{
			Tool:    iam_list_roles,
			Handler: handleRolesListTool,
		},
		{
			Tool:    iam_get_role_policies,
			Handler: handleGetRolePoliciesTool,
		},
		{
			Tool:    iam_get_inline_policy,
			Handler: handleGetInlinePolicy,
		},
		{
			Tool:    iam_get_attached_policy,
			Handler: handleGetAttachedPolicy,
		},
	}

}
