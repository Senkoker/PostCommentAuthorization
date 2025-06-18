FROM golang:1.23.4-alpine AS builder

WORKDIR /usr/local/src

RUN apk --no-cache add gcc gettext musl-dev

COPY ./go.mod ./

COPY ./go.sum ./

RUN go mod download

COPY ./ ./

RUN go build -o ./bin/app/cmd ./cmd/web/main.go

FROM alpine

WORKDIR /service

COPY --from=builder /usr/local/src/bin/app/cmd ./

CMD ["./cmd"]