FROM golang:1.8.1

ENV GOPATH /go
RUN go get -u github.com/kardianos/govendor
ADD . /go/src/cubiclerebels.com/vault2secrets/
WORKDIR /go/src/cubiclerebels.com/vault2secrets/
RUN govendor fetch +all
RUN go build
ENTRYPOINT ./vault2secrets
