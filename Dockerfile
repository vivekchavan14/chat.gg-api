FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o chat.gg-api ./main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/chat.gg-api .

EXPOSE 8080

CMD ["./chat.gg-api"]
