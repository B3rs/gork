FROM golang:1.17.6-alpine
WORKDIR /app 
COPY . .
WORKDIR /app/example/worker
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o worker -ldflags="-w -s" .
CMD ["./worker"]