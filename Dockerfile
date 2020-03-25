FROM golang:alpine

WORKDIR /go/src/app

COPY . .

RUN go build main.go
CMD ./main