package config

import "fmt"

func ListenAddr() string {
	return fmt.Sprintf("%s:%d", HOST, PORT)
}
