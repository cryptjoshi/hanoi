FROM golang:1.23.0-alpine3.20 as builder

ENV GO111MODULE=on
ENV GOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0
ADD . /src
RUN cd /src && go build -o popmio -ldflags "-w -s" main.go

FROM alpine:latest

RUN apk update \
        && apk upgrade \
        && apk add --no-cache \
        ca-certificates \
        && update-ca-certificates 2>/dev/null || true

WORKDIR /app
COPY --from=builder /src/js /app/js
COPY --from=builder /src/css /app/css
COPY --from=builder /src/popmio /app/
CMD /app/popmio