package logger

import (
	"fmt"
	"net"
	"time"
)

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorGray   = "\033[90m"
	colorBlue   = "\033[94m"
	colorGreen  = "\033[92m"
	colorYellow = "\033[93m"
	colorRed    = "\033[91m"
)

// Global flags
var EnableDebug = true

func timestamp() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func Event(conn net.Conn, event string) {
	remote := conn.RemoteAddr().String()
	fmt.Printf("%s[%s] [%s] %s%s\n", colorGray, timestamp(), remote, event, colorReset)
}

func Eventf(conn net.Conn, format string, args ...interface{}) {
	remote := conn.RemoteAddr().String()
	fmt.Printf("%s[%s] [%s] %s%s\n", colorGray, timestamp(), remote, fmt.Sprintf(format, args...), colorReset)
}

func Info(msg string) {
	fmt.Printf("%s[%s] [INFO] %s%s\n", colorBlue, timestamp(), msg, colorReset)
}

func Infof(format string, args ...interface{}) {
	fmt.Printf("%s[%s] [INFO] %s%s\n", colorBlue, timestamp(), fmt.Sprintf(format, args...), colorReset)
}

func Debug(msg string) {
	if EnableDebug {
		fmt.Printf("%s[%s] [DEBUG] %s%s\n", colorGreen, timestamp(), msg, colorReset)
	}
}

func Debugf(format string, args ...interface{}) {
	if EnableDebug {
		fmt.Printf("%s[%s] [DEBUG] %s%s\n", colorGreen, timestamp(), fmt.Sprintf(format, args...), colorReset)
	}
}

func Warn(msg string) {
	fmt.Printf("%s[%s] [WARN] %s%s\n", colorYellow, timestamp(), msg, colorReset)
}

func Warnf(format string, args ...interface{}) {
	fmt.Printf("%s[%s] [WARN] %s%s\n", colorYellow, timestamp(), fmt.Sprintf(format, args...), colorReset)
}

func Error(msg string) {
	fmt.Printf("%s[%s] [ERROR] %s%s\n", colorRed, timestamp(), msg, colorReset)
}

func Errorf(format string, args ...interface{}) {
	fmt.Printf("%s[%s] [ERROR] %s%s\n", colorRed, timestamp(), fmt.Sprintf(format, args...), colorReset)
}
