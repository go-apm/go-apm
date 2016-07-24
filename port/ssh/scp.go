package ssh

import (
	"fmt"
	"github.com/go-apm/go-apm/util/xhttp"
	"github.com/uber-go/zap"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"golang.org/x/net/context"
	"net"
	"os"
)

func FetchBinary(c context.Context, host string, src string, dst string) error {
	agent, err := getAgent()
	if err != nil {
		xhttp.CurrentLogger(c).Fatal("Failed to connect to SSH_AUTH_SOCK", zap.Error(err))
	}
	sshClient, err := ssh.Dial("tcp", host+":22", &ssh.ClientConfig{
		User: os.Getenv("USER"),
		Auth: []ssh.AuthMethod{
			ssh.PublicKeysCallback(agent.Signers),
		},
	})
	if err != nil {
		xhttp.CurrentLogger(c).Fatal("Failed to dial ssh", zap.Error(err))
	}
	session, err := sshClient.NewSession()
	if err != nil {
		xhttp.CurrentLogger(c).Fatal("Failed to create session", zap.Error(err))
	}
	xhttp.CurrentLogger(c).Info("Start scp -r", zap.String("from", src), zap.String("to", dst))
	cmd := fmt.Sprintf("scp -r %s %s", src, dst)
	if err := session.Run(cmd); err != nil {
		xhttp.CurrentLogger(c).Error("Execute scp error", zap.Error(err))
		return err
	}
	return nil
}

func getAgent() (agent.Agent, error) {
	agentConn, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK"))
	return agent.NewClient(agentConn), err
}
