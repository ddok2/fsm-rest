FROM golang:1.16-alpine

WORKDIR /go/src/blockchain.automation

COPY . .

RUN CGO_ENABLED=0 go build main.go

CMD ["./main"]
