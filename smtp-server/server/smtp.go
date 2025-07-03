package server

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"net"
	"strings"

	"github.com/ZiplEix/mail-toolchain/shared/database"
	"github.com/ZiplEix/mail-toolchain/shared/logger"
	"github.com/ZiplEix/mail-toolchain/smtp-server/server/internal"
)

var tlsConfig *tls.Config

func LoadTLSConfig(certPath, keyPath string) error {
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return fmt.Errorf("failed to load TLS certificate and key: %v", err)
	}
	tlsConfig = &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	return nil
}

func LunchSMTPServer() {
	listener, err := net.Listen("tcp", ":2525")
	if err != nil {
		panic(fmt.Sprintf("Failed to start SMTP server: %v", err))
	}
	logger.Info("SMTP server listening on port 2525")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go HandleConnection(conn)
	}
}

func HandleConnection(conn net.Conn) {
	defer func() {
		if r := recover(); r != nil {
			logger.Event(conn, fmt.Sprintf("Recovered from panic: %v", r))
			conn.Close()
		}
	}()

	defer conn.Close()
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	sendLine := func(s string) {
		_, _ = writer.WriteString(s + "\r\n")
		writer.Flush()
	}

	sendError := func(code int, msg string) {
		resp := fmt.Sprintf("%d %s", code, msg)
		sendLine(resp)
		logger.Event(conn, "S: "+resp)
	}

	sendLine("220 localhost SMTP ready")
	logger.Event(conn, "Connection opened")

	var from string
	var toList []string
	var dataLines []string
	mode := "command"

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			logger.Event(conn, "Connection closed")
			break
		}
		line = strings.TrimRight(line, "\r\n")
		logger.Event(conn, fmt.Sprintf("C: %s", line))

		switch mode {
		case "command":
			switch {
			case strings.HasPrefix(line, "NOOP"):
				sendLine("250 OK")

			case strings.HasPrefix(line, "EHLO"):
				sendLine("250-localhost Hello")
				sendLine("250-PIPELINING")
				sendLine("250-SIZE 35882577")
				sendLine("250-STARTTLS")
				sendLine("250 HELP")

			case strings.HasPrefix(line, "HELO"):
				sendLine("250 localhost Hello")

			case strings.HasPrefix(line, "MAIL FROM:"):
				addr := extractEmailAddress(line, "MAIL FROM:")
				if !internal.IsValidEmail(addr) {
					sendError(550, "Invalid sender address")
					continue
				}
				from = addr
				sendLine("250 OK")

			case strings.HasPrefix(line, "RCPT TO:"):
				addr := extractEmailAddress(line, "RCPT TO:")
				if !internal.IsValidEmail(addr) {
					sendError(550, "Invalid recipient address")
					continue
				}
				toList = append(toList, addr)
				sendLine("250 OK")

			case strings.HasPrefix(line, "DATA"):
				if from == "" || len(toList) == 0 {
					sendError(503, "Bad sequence: MAIL FROM and RCPT TO required before DATA")
					continue
				}
				sendLine("354 End data with <CR><LF>.<CR><LF>")
				mode = "data"

			case strings.HasPrefix(line, "RSET"):
				from = ""
				toList = nil
				dataLines = nil
				sendLine("250 OK: Reset state")

			case strings.HasPrefix(line, "VRFY"):
				arg := strings.TrimSpace(strings.TrimPrefix(line, "VRFY"))
				if internal.IsValidEmail(arg) {
					sendLine("250 User exists")
				} else {
					sendError(550, "No such user")
				}

			case line == "QUIT":
				sendLine("221 Bye")
				logger.Event(conn, "Connection closed via QUIT")
				return

			case strings.HasPrefix(line, "STARTTLS"):
				if tlsConfig == nil {
					sendError(454, "TLS not available")
					continue
				}
				sendLine("220 Ready to start TLS")
				tlsConn := tls.Server(conn, tlsConfig)
				err := tlsConn.Handshake()
				if err != nil {
					logger.Event(conn, fmt.Sprintf("TLS handshake failed: %v", err))
					return
				}
				conn = tlsConn
				reader = bufio.NewReader(conn)
				writer = bufio.NewWriter(conn)
				logger.Event(conn, "TLS connection established")

			default:
				sendError(500, "Unrecognized command")
			}

		case "data":
			if line == "." {
				err := database.SaveMail(from, toList, dataLines)
				if err != nil {
					sendError(550, "Error saving message to database")
					logger.Event(conn, fmt.Sprintf("Error saving mail from %s to %v: %v", from, toList, err))
				} else {
					sendLine("250 OK: message accepted")
					logger.Event(conn, fmt.Sprintf("Mail saved from %s to %v", from, toList))
				}
				mode = "command"
				from = ""
				toList = nil
				dataLines = nil
			} else {
				dataLines = append(dataLines, line)
			}
		}
	}
}

func extractEmailAddress(line string, prefix string) string {
	rest := strings.TrimSpace(strings.TrimPrefix(line, prefix))

	start := strings.Index(rest, "<")
	end := strings.Index(rest, ">")

	if start == -1 || end == -1 || end <= start {
		return ""
	}

	return rest[start : end+1]
}
