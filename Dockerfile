FROM alpine:3.9

RUN apk update \
        && apk upgrade \
        && apk add --no-cache \
        ca-certificates \
        && update-ca-certificates 2>/dev/null || true

COPY ./bin/ ./
COPY ./config/*.yml ./config/
