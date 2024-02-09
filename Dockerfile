FROM golang:1.21-alpine AS base

WORKDIR /app

RUN mkdir -p /app/.cache && chown -R nobody:nogroup /app/.cache

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /main

FROM alpine:latest AS current

USER nobody:nogroup

COPY --from=base /main /main

WORKDIR /app
CMD ["/main"]

