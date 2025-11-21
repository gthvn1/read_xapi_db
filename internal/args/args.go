package args

import (
	"flag"
	"fmt"
	"os"
)

type Args struct {
	Username string
	Password string
	Hostname string
	FileName string
}

func GetArgs() Args {
	fileName := flag.String("file", "", "Local database file path (used if -hostname is not provided)")
	username := flag.String("username", "", "SSH username (for remote fetch)")
	password := flag.String("password", "", "SSH password (for remote fetch)")
	hostname := flag.String("hostname", "", "Remote host (leave empty for local file)")

	flag.Parse()

	if *fileName == "" {
		fmt.Println("Error: -file is required")
		flag.Usage()
		os.Exit(1)
	}
	return Args{
		FileName: *fileName,
		Username: *username,
		Password: *password,
		Hostname: *hostname,
	}
}
