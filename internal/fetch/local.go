package fetch

import "os"

func Local(path string) ([]byte, error) {
	return os.ReadFile(path)
}
