FROM golang:alpine

RUN mkdir /app && cd /app
WORKDIR /app
COPY . /app

RUN go install github.com/julesguesnon/gomon
RUN export PATH=$PATH:$GOPATH/bin
CMD ["gomon","server.go"]