FROM golang:1.10-stretch

WORKDIR /go/src/github.com/bakku/easyalert
COPY . .

RUN go get -v github.com/bakku/gom/cmd/gom

RUN dep ensure

RUN make build

CMD [ "cmd/easyalert/easyalert" ]
