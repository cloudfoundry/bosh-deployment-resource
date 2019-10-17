FROM ubuntu:16.04

# Packages
RUN DEBIAN_FRONTEND=noninteractive apt-get -y -qq update && apt-get -y -qq install \
  gcc \
  git-core \
  make \
  python-software-properties \
  software-properties-common \
  wget \
  curl

WORKDIR /tmp/docker-build

# Golang
ENV GO_VERSION=1.12.9
ENV GO_SHA2SUM=ac2a6efcc1f5ec8bdc0db0a988bb1d301d64b6d61b7e8d9e42f662fbb75a2b9b

RUN curl -LO https://storage.googleapis.com/golang/go${GO_VERSION}.linux-amd64.tar.gz && \
    echo "${GO_SHA2SUM}  go${GO_VERSION}.linux-amd64.tar.gz" > go_${GO_VERSION}_SHA2SUM && \
    shasum -a 256 -cw --status go_${GO_VERSION}_SHA2SUM
RUN tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz
ENV GOPATH /root/go
RUN mkdir -p /root/go/bin
ENV PATH $PATH:/usr/local/go/bin:$GOPATH/bin
RUN go get github.com/onsi/ginkgo
RUN go install github.com/onsi/ginkgo/...
RUN go get golang.org/x/lint/golint

# Google SDK
ENV GCLOUD_VERSION=257.0.0
ENV GCLOUD_SHA2SUM=2b9eb732206f9c171b6eb9a22083efb98006cc09f1e53f09468626160d3a7cf8

RUN wget https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-sdk-${GCLOUD_VERSION}-linux-x86_64.tar.gz \
    -O gcloud_${GCLOUD_VERSION}_linux_amd64.tar.gz && \
    echo "${GCLOUD_SHA2SUM}  gcloud_${GCLOUD_VERSION}_linux_amd64.tar.gz" > gcloud_${GCLOUD_VERSION}_SHA2SUM && \
    shasum -a 256 -cw --status gcloud_${GCLOUD_VERSION}_SHA2SUM && \
    tar xvf gcloud_${GCLOUD_VERSION}_linux_amd64.tar.gz && \
    mv google-cloud-sdk / && cd /google-cloud-sdk  && ./install.sh

ENV PATH=$PATH:/google-cloud-sdk/bin

# Cleanup
RUN rm -rf /tmp/docker-build
