FROM golang:1.12 AS build

WORKDIR /app
ADD . /app
RUN go mod download
RUN go build src/apis/main.go
CMD ["/app/main"]