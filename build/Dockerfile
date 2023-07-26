FROM quay.io/centos/centos:stream9

RUN dnf install -y epel-release
RUN dnf install -y golang

RUN mkdir -p /go/src/github.com/cpustat
WORKDIR /go/src/github.com/cpustat
ADD . /go/src/github.com/cpustat/

