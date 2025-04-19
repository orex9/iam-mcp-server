package stdio

import (
	"github.com/mark3labs/mcp-go/server"
	awss "github.com/orex9/iam-mcp-server/pkg/aws"
)

func RunStdioServer() error {
	s := server.NewMCPServer("AWS IAM list", "1.0.0", server.WithLogging(), server.WithRecovery())

	tools := awss.InitTools()
	s.AddTools(tools...)

	if err := server.ServeStdio(s); err != nil {
		return err
	}

	return nil
}
