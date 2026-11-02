ARG go_version=1
FROM golang:${go_version} AS build
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -o /opt/resource/check ./cmd/check \
    && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -o /opt/resource/in ./cmd/in \
    && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -o /opt/resource/out ./cmd/out

FROM alpine:3
RUN apk add --no-cache bash ca-certificates curl git openssh-client tar gzip
COPY --from=build /opt/resource/ /opt/resource/