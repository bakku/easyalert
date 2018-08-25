FROM golang:1.10-stretch

WORKDIR /go/src/github.com/bakku/easyalert
COPY . .

RUN go get -v github.com/bakku/gom/cmd/gom

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
RUN dep ensure

RUN make go_build

CMD [ "cmd/easyalert/easyalert" ]
