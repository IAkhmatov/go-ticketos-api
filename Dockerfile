FROM docker.io/library/golang:1.21-alpine AS builder

RUN apk add build-base

WORKDIR /app
COPY go.mod ./
RUN go mod download

COPY cmd cmd
COPY internal internal
COPY Makefile .

RUN make build

FROM scratch

WORKDIR /
COPY --from=builder app/bin/api /api

EXPOSE 8080
ENTRYPOINT ["/api"]