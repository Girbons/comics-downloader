FROM golang:latest

RUN apt-get update -qq \
    && apt-get install -y -q --no-install-recommends \
        libgl1-mesa-dev \
        xorg-dev \
        gosu \
        curl \
    && apt-get -qy autoremove \
    && apt-get clean \
&& rm -r /var/lib/apt/lists/*;


RUN curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(go env GOPATH)/bin v1.16.0
RUN go get github.com/mattn/goveralls
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
