FROM golang:alpine as builder

WORKDIR /app 
COPY . .
WORKDIR /app/example/worker
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o worker -ldflags="-w -s" .
CMD ["./worker"]