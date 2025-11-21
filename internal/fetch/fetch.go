package fetch

import "example.com/readxapidb/internal/args"

func DB(a args.Args) ([]byte, error) {
	if a.Hostname == "" {
		// Local database is used
		return Local(a.FileName)
	}

	return FileSFTP(a.Username, a.Password, a.Hostname, a.FileName)

}
