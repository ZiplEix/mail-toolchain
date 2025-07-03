package logger

import (
	"fmt"
	"net"
	"time"
)

func Event(conn net.Conn, event string) {
	remote := conn.RemoteAddr().String()
	t := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("[%s] [%s] %s\n", t, remote, event)
}

func Info(msg string) {
	t := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("[%s] [INFO] %s\n", t, msg)
}

func Error(msg string) {
	t := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("[%s] [ERROR] %s\n", t, msg)
}
