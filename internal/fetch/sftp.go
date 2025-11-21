package fetch

import (
	"io"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func FileSFTP(username, password, host, filePath string) ([]byte, error) {
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}

	conn, err := ssh.Dial("tcp", host+":22", config)
	if err != nil {
		return nil, err
	}

	sftpClient, err := sftp.NewClient(conn)
	if err != nil {
		return nil, err
	}
	defer sftpClient.Close()

	f, err := sftpClient.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return io.ReadAll(f)
}
