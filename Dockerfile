FROM golang:1.4.2
RUN go get github.com/Financial-Times/coco-logfilter/logfilter
CMD $GOPATH/bin/logfilter