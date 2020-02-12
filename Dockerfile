FROM golang AS build_base

WORKDIR /app

COPY . /app
RUN go mod download
RUN go build .

FROM alpine

WORKDIR /app
COPY --from=build_base /app /app

EXPOSE 8080
CMD ["/app/main"]
