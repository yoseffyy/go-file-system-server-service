FROM golang:alpine3.11

WORKDIR /go/src/app

COPY . .

RUN go build main.go
CMD ./main