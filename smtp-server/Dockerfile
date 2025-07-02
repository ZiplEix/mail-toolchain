FROM golang:1.23

WORKDIR /app

RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY go.mod go.sum ./

RUN --mount=type=cache,target=/go/pkg/mod/ \
--mount=type=bind,source=go.sum,target=go.sum \
--mount=type=bind,source=go.mod,target=go.mod \
go mod download -x

COPY . .

RUN go build -o main .

EXPOSE 2525

CMD ["./main"]
