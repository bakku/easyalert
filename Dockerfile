FROM golang:1.12.9-alpine3.9

WORKDIR /go/src/github.com/bakku/easyalert

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh make build-base && \
    go get -v github.com/bakku/gom/cmd/gom

COPY . .

RUN GO111MODULE=on make go_build

CMD [ "build/easyalert" ]
