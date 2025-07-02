#!/bin/bash

HOST="localhost"
PORT="2525"

echo "==> Testing SMTP server on $HOST:$PORT"

{
  sleep 1
  echo "EHLO localhost"
  sleep 0.5
  echo "MAIL FROM:<alice@example.com>"
  sleep 0.5
  echo "RCPT TO:<bob@example.com>"
  sleep 0.5
  echo "DATA"
  sleep 0.5
  echo "salut"
  echo "."
  sleep 0.5
  echo "QUIT"
} | nc "$HOST" "$PORT"
