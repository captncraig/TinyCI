from golang:1.4

ADD . /go/src/github.com/captncraig/tinyci
RUN go get github.com/captncraig/tinyci
RUN go install github.com/captncraig/tinyci

ENTRYPOINT /go/bin/tinyci
ENV TINYCI-SCRIPT-DIR /tinyci
VOLUME /tinyci
EXPOSE 4567