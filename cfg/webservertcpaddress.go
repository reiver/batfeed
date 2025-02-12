package cfg

import (
	"fmt"

	"github.com/reiver/socialfed-api/env"
)

func WebServerTCPAddress() string {
	return fmt.Sprintf(":%s", env.TcpPort)
}
