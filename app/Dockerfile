# syntax=docker/dockerfile:1.16
ARG BINARY
FROM gcr.io/distroless/static:nonroot

COPY ${BINARY}/denon-to-mqtt-gw /denon

ENTRYPOINT ["/denon", "/var/lib/denon-to-mqtt-gw/config.json"]
