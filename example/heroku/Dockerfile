ARG COMMAND
FROM golang:1.17.6-alpine
WORKDIR /heroku 
COPY . .
# always get latest version of library, do not use this in production :)
RUN go get github.com/B3rs/gork@main
WORKDIR /heroku/client
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o client -ldflags="-w -s" .
WORKDIR /heroku/worker
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o worker -ldflags="-w -s" .

CMD exec $COMMAND