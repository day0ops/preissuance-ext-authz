FROM alpine:3.20
RUN apk add --no-cache ca-certificates

ARG TARGETARCH
COPY bin/preissuance-ext-authz-linux-${TARGETARCH} /usr/local/bin/preissuance-ext-authz

ENTRYPOINT ["/usr/local/bin/preissuance-ext-authz"]
