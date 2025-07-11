package server

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"

	"github.com/ZiplEix/mail-toolchain/shared/database"
	"github.com/ZiplEix/mail-toolchain/shared/logger"
	"github.com/toorop/go-dkim"
)

func dataMode(session *Session, line string, dataLines []string) ([]string, error) {
	if line == "." {
		pemKey, err := EncodePrivateKeyToPEM(privateKey)
		if err != nil {
			session.SendError(550, "Error encoding private key")
			logger.Errorf("Error encoding private key: %v", err)
			return dataLines, err
		}

		signedDataLines, err := signWithDKIM(
			[]byte(strings.Join(dataLines, "\r\n")),
			extractDomain(session.From), "default",
			pemKey,
		)
		if err != nil {
			session.SendError(550, "Error signing message")
			logger.Event(session.Conn, fmt.Sprintf("Error signing mail from %s to %v: %v", session.From, session.ToList, err))
			return dataLines, err
		}

		err = database.SaveMail(session.From, session.ToList, BytesToLines(signedDataLines))
		if err != nil {
			session.SendError(550, "Error saving message to database")
			logger.Event(session.Conn, fmt.Sprintf("Error saving mail from %s to %v: %v", session.From, session.ToList, err))
		} else {
			session.SendLine("250 OK: message accepted")
			logger.Event(session.Conn, fmt.Sprintf("Mail saved from %s to %v", session.From, session.ToList))
		}
		session.Mode = "command"
		session.From = ""
		session.ToList = nil
		dataLines = nil
	} else {
		dataLines = append(dataLines, line)
	}
	return dataLines, nil
}

func BytesToLines(data []byte) []string {
	var lines []string
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

func signWithDKIM(message []byte, domain, selector string, privateKey []byte) ([]byte, error) {
	options := dkim.NewSigOptions()
	options.PrivateKey = privateKey
	options.Domain = "baptiste.zip"
	options.Selector = "default"
	options.Headers = []string{"from", "to", "subject", "date"}
	options.Canonicalization = "relaxed/relaxed"
	options.AddSignatureTimestamp = true

	email := make([]byte, len(message))
	copy(email, message)

	err := dkim.Sign(&email, options)
	if err != nil {
		return nil, err
	}
	return email, nil
}

func extractDomain(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) == 2 {
		return parts[1]
	}
	return ""
}
