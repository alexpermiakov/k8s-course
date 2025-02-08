FROM --platform=linux/amd64 golang:1.23.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main . 

FROM --platform=linux/amd64 alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .

RUN apk add --no-cache tzdata

EXPOSE 8080

CMD ["./main"]