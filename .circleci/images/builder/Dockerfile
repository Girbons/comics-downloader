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

RUN go get github.com/mattn/goveralls
