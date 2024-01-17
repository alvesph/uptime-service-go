FROM golang:1.21-alpine AS base

ENV DEBIAN_FRONTEND=noninterative

ARG UID=1000
ARG DIR=/app

FROM base AS current
WORKDIR ${DIR}

COPY . ./
RUN go mod download
RUN go build -o /main

CMD ["/main"]
