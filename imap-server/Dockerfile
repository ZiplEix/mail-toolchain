FROM golang:1.23

WORKDIR /app

COPY ./shared ./shared

WORKDIR /app/imap-server

COPY imap-server/go.mod imap-server/go.sum ./
# RUN go mod edit -replace github.com/ZiplEix/mail-toolchain/shared=./shared
RUN go mod download -x

COPY ./imap-server .

RUN go build -o main .

EXPOSE 2525

CMD ["./main"]
