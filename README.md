
[![Circle CI](https://circleci.com/gh/Financial-Times/coco-logfilter/tree/master.png?style=shield)](https://circleci.com/gh/Financial-Times/coco-logfilter/tree/master)[![Go Report Card](https://goreportcard.com/badge/github.com/Financial-Times/coco-logfilter)](https://goreportcard.com/report/github.com/Financial-Times/coco-logfilter) [![Coverage Status](https://coveralls.io/repos/github/Financial-Times/coco-logfilter/badge.svg)](https://coveralls.io/github/Financial-Times/coco-logfilter)
# coco-logfilter
Simple pre-aggregation logfilter to operate on the json output of journald.

## Building
```
cd logfilter
CGO_ENABLED=0 go build -a -installsuffix cgo -o coco-logfilter .
cd ..

docker build -t coco/coco-logfilter .
```

## Installation
Download the project:
```
go get github.com/Financial-Times/coco-logfilter/logfilter

```
Use govendor for dependencies:
```
go get github.com/kardianos/govendor
govendor sync
```

## Example use
```
echo '{"MESSAGE":"127.0.0.1 - - [21/Apr/2015:12:15:34 +0000] \"GET /eom-file/all/zzzz HTTP/1.1\" 200 12345 919 919"}' | logfilter
```
