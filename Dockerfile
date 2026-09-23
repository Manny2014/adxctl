# build stage
FROM golang:1.19-alpine AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o /adxctl .

# final stage
FROM alpine:latest

WORKDIR /

COPY --from=build /adxctl /adxctl

ENTRYPOINT ["/adxctl"]
