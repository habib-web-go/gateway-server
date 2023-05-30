# Set the default value for tags of docker images
ARG GO_IMAGE_TAG=latest

# Specify version of golang inside .env file
FROM golang:${GO_IMAGE_TAG} AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o ./server .

CMD ["./server"]