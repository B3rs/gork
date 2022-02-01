FROM golang:1.17.6-alpine
WORKDIR /app 
COPY . .
# always get latest version of library, do not use this in production :)
RUN go get github.com/B3rs/gork@main
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o worker -ldflags="-w -s" .
CMD ["./worker"]