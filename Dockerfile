FROM ubuntu:focal AS mirror

ENV TZ=UTC
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone \
    && echo 'deb mirror://mirrors.ubuntu.com/mirrors.txt focal main restricted universe multiverse\n\
deb mirror://mirrors.ubuntu.com/mirrors.txt focal-updates main restricted universe multiverse\n\
deb mirror://mirrors.ubuntu.com/mirrors.txt focal-security main restricted universe multiverse' > /etc/apt/sources.list \
    && apt update

FROM mirror AS dev-usql

RUN apt update && apt install -y curl build-essential git

ARG GO_VERSION=1.24.3
ARG GO_DISTRO=linux-amd64
RUN curl -O https://dl.google.com/go/go${GO_VERSION}.${GO_DISTRO}.tar.gz \
&& tar xzvf go${GO_VERSION}.${GO_DISTRO}.tar.gz -C /usr/local \
&& rm go${GO_VERSION}.${GO_DISTRO}.tar.gz
ENV PATH=/usr/local/go/bin:$PATH

WORKDIR /opt
COPY . .