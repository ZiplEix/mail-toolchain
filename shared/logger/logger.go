package logger

import (
	"fmt"
	"net"
	"time"
)

func Event(conn net.Conn, message string) {
	ts := time.Now().Format("2006-01-02 15:04:05")
	addr := conn.RemoteAddr().String()
	fmt.Printf("[%s] [%s] %s\n", ts, addr, message)
}

func Info(msg string) {
	fmt.Printf("[INFO] %s\n", msg)
}

func Error(msg string) {
	fmt.Printf("[ERROR] %s\n", msg)
}
