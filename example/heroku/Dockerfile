FROM golang:1.17.6-alpine
ARG COMMAND
WORKDIR /heroku 
COPY . .
# always get latest version of library, do not use this in production :)
RUN go get github.com/B3rs/gork@main
WORKDIR /heroku/$COMMAND
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $COMMAND -ldflags="-w -s" .

CMD ["./${$COMMAND}"]