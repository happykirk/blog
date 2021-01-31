FROM golang:latest

WORKDIR $GOPATH/src/github.com/happykirk/blog
COPY . $GOPATH/src/github.com/happykirk/blog
RUN go build .

EXPOSE 8000
ENTRYPOINT ["./blog"]
