FROM ubuntu:12.04

# env vars
ENV HOME /home/deployer
ENV GOPATH /home/deployer/go
ENV GOOS linux
ENV GOARCH amd64
ENV CGO_ENABLED 0
ENV PATH /usr/local/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/games

RUN mkdir -p /home/deployer/gosrc
RUN echo 'GOPATH=/home/deployer/gosrc' >> ~/.bashrc

# apt
RUN echo "deb http://mirror.anl.gov/pub/ubuntu precise main universe" >> /etc/apt/sources.list
RUN apt-get update
RUN apt-get install -y build-essential mercurial git-core subversion wget default-jre

# go 1.2 tarball
RUN wget -qO- https://go.googlecode.com/files/go1.2.linux-amd64.tar.gz | tar -C /usr/local -xzf -

RUN mkdir -p /home/deployer/logs

WORKDIR /home/deployer/gosrc
RUN mkdir -p /home/deployer/gosrc/src/github.com/curt-labs/GoAPI
ADD . /home/deployer/gosrc/src/github.com/curt-labs/GoAPI
WORKDIR /home/deployer/gosrc/src/github.com/curt-labs/GoAPI
RUN export GOPATH=/home/deployer/gosrc && go get
RUN export GOPATH=/home/deployer/gosrc && go build -o GoAPI ./index.go

EXPOSE 443