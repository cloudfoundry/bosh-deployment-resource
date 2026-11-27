FROM --platform=$BUILDPLATFORM golang:1.27-alpine AS builder
ARG TARGETOS
ARG TARGETARCH
WORKDIR /src
COPY . .
RUN for cmd in check in out; do \
      CGO_ENABLED=0 GOOS="$TARGETOS" GOARCH="$TARGETARCH" go build -trimpath -o "/out/$cmd" "./cmd/$cmd"; \
    done

FROM alpine:3.24 AS resource
LABEL org.opencontainers.image.source="https://github.com/cloudfoundry/bosh-deployment-resource"
RUN apk add --no-cache bash ca-certificates git openssh-client tar
COPY --from=builder /out/ /opt/resource/
