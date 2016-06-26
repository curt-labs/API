FROM golang:1.6

RUN mkdir -p /home/deployer/gosrc/src/github.com/curt-labs/API
ADD . /home/deployer/gosrc/src/github.com/curt-labs/API
WORKDIR /home/deployer/gosrc/src/github.com/curt-labs/API
RUN export GOPATH=/home/deployer/gosrc && go get
RUN export GOPATH=/home/deployer/gosrc && go build -o API ./index.go

ENTRYPOINT /home/deployer/gosrc/src/github.com/curt-labs/API/API

EXPOSE 8080
