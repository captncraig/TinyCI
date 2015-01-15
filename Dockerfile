from golang:1.4

ADD . /go/src/tinyci
WORKDIR /go/src/tinyci

RUN go get -v
RUN go install

EXPOSE 4567

VOLUME /scripts

ENTRYPOINT /go/bin/tinyci
