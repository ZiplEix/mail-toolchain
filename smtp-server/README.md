# SMTP Server

This is a minimalist SMTP server written in Go, designed to receive and store incoming emails.

# Features

- Handles basic SMTP commands (HELO, MAIL FROM, RCPT TO, DATA, QUIT)
- STARTTLS support
- Structured logging
- Stores incoming emails into PostgreSQL
- Supports multiple recipients

## Usage

1. Start the server:

```bash
go run main.go
```

2. Connect with a mail client (Thunderbird, Outlook) or via Telnet/OpenSSL:

```bash
telnet localhost 2525
```

## Environment

Set your PostgreSQL connection string in `../.env`: