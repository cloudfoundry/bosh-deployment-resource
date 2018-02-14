FROM concourse/buildroot:git
MAINTAINER https://github.com/cloudfoundry/bosh-deployment-resource

RUN curl -s -L https://github.com/git-lfs/git-lfs/releases/download/v2.3.4/git-lfs-linux-amd64-2.3.4.tar.gz | tar xz && \
    ./git-lfs-*/install.sh && rm -r git-lfs-*

ADD check /opt/resource/check
ADD in /opt/resource/in
ADD out /opt/resource/out

RUN chmod +x /opt/resource/*
