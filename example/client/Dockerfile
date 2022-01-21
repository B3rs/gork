FROM golang:1.17.6-alpine
WORKDIR /app 
COPY . .
WORKDIR /app/example/client
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o client -ldflags="-w -s" .
CMD ["./client"]