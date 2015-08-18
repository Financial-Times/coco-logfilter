FROM gliderlabs/alpine:3.2

ADD . /logfilter
RUN apk --update add go git\
  && export GOPATH=/.gopath \
  && go get github.com/Financial-Times/coco-logfilter \
  && cd logfilter/logfilter \
  && go build \
  && mv logfilter /coco-logfilter \
  && apk del go git \
  && rm -rf $GOPATH /var/cache/apk/*

CMD /coco-logfilter -environment=$ENV
