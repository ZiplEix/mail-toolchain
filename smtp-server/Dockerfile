FROM golang:1.23

WORKDIR /app

COPY ./shared ./shared

WORKDIR /app/smtp-server

COPY smtp-server/go.mod smtp-server/go.sum ./
# RUN go mod edit -replace github.com/ZiplEix/mail-toolchain/shared=./shared
RUN go mod download -x

COPY ./smtp-server .

RUN go build -o main .

EXPOSE 2525

CMD ["./main"]
