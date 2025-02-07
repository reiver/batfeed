package cfg

import (
	"fmt"

	"github.com/reiver/batfeed/env"
)

func WebServerTCPAddress() string {
	return fmt.Sprintf(":%s", env.TcpPort)
}
