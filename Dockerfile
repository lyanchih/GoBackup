FROM ubuntu:14.04
MAINTAINER Lyan Hung <lyanchih@gmail.com>

RUN apt-get -qq update; apt-get install -qy golang postgresql git Mercurial

WORKDIR /src

COPY ./src /src

VOLUME ["/src/config"]

EXPOSE 80

ENV GOPATH /src

RUN go get code.google.com/p/google-api-go-client/drive/v2 github.com/golang/oauth2 github.com/stacktic/dropbox && go build -o bin/gobackup *.go
