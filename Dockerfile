FROM golang:alpine

RUN apk add openssh git

RUN mkdir /app && cd /app
WORKDIR /app
COPY . /app

RUN export PATH=$PATH:$GOPATH/bin
RUN go get -u github.com/julesguesnon/gomon
RUN go get -u github.com/gorilla/mux

EXPOSE 3030
CMD ["gomon","server.go"]