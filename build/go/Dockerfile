#go version
FROM golang:1.23.3

WORKDIR /app
COPY .env /app
RUN mkdir -p /app/src
WORKDIR /app/src
RUN go install github.com/air-verse/air@latest

# airを実行
CMD ["air", "-c", ".air.toml"]
