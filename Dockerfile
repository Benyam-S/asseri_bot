FROM golang:1.15.7-alpine3.13

ADD . /asseri
WORKDIR /asseri/servers/botclient

RUN apk add git
RUN go mod download
RUN go build -o main .

EXPOSE 443

CMD  ["/asseri/servers/botclient/main"]