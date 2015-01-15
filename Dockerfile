from golang:1.4

ADD . /go/tinyci
WORKDIR /go/tinyci

RUN go get -v
RUN go install

EXPOSE 4567

VOLUME /scripts

ENTRYPOINT /go/bin/tinyci