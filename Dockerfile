FROM golang:1.17.3
WORKDIR /go/src/app
COPY go/ .
RUN go install -v ./...

CMD ["gocurrency"]