FROM golang as builder

WORKDIR /app 
COPY . .
WORKDIR /app/example/client
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o client -ldflags="-w -s" .
CMD ["./client"]