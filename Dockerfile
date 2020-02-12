FROM golang AS build_base

WORKDIR /app

COPY . /app
RUN go mod download
RUN go build .

FROM alpine

WORKDIR /app
COPY --from=build_base /app /app

EXPOSE 3000
CMD ["/app/main"]
