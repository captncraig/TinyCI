from golang:1.4

ADD . /go/src/tinyci
WORKDIR /go/src/tinyci

RUN go get -v
RUN go install

EXPOSE 4567

VOLUME /scripts
ENV TINYCI-SCRIPT-DIR /scripts

ENV DOCKER_HOST unix:///tmp/docker.sock
ENTRYPOINT /go/bin/tinyci
