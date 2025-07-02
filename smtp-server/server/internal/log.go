package internal

import (
	"fmt"
	"net"
	"time"
)

func LogEvent(conn net.Conn, event string) {
	remoteAddr := conn.RemoteAddr().String()
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("[%s] [%s] %s\n", timestamp, remoteAddr, event)
}
