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

ENV REDIS_PASSWORD eC0mm3rc3
ENV REDIS_CLIENT_ADDRESS 172.17.42.1:6379
ENV REDIS_MASTER_ADDRESS 162.222.179.124:6379

# Database Connection String Variables
ENV DATABASE_HOST 173.255.114.206:3306
ENV DATABASE_PROTOCOL tcp
ENV DATABASE_USERNAME curtDuser2
ENV DATABASE_PASSWORD eC0mm3rc3
ENV CURT_DEV_NAME CurtDev
ENV PCDB_NAME pcdb
ENV VCDB_NAME vcdb
ENV SITE_MONITOR_NAME SiteMonitor
ENV ADMIN_NAME admin
ENV FTP_HOST 216.17.90.82
ENV FTP_USERNAME eCommerceFTP
ENV FTP_PASSWORD 3GaJPaAZ

# SMTP Server Variables
ENV EMAIL_SERVER smtp.gmail.com
ENV EMAIL_ADDRESS no-reply@curtmfg.com
ENV EMAIL_USERNAME no-reply@curtmfg.com
ENV EMAIL_PASSWORD eC0mm3rc3
ENV EMAIL_SSL true
ENV EMAIL_PORT 587

RUN mkdir -p /home/deployer/logs

WORKDIR /home/deployer/gosrc
RUN mkdir -p /home/deployer/gosrc/src/github.com/curt-labs/GoAPI
ADD . /home/deployer/gosrc/src/github.com/curt-labs/GoAPI
WORKDIR /home/deployer/gosrc/src/github.com/curt-labs/GoAPI
RUN export GOPATH=/home/deployer/gosrc && go get
RUN export GOPATH=/home/deployer/gosrc && go build -o GoAPI ./index.go

EXPOSE 443