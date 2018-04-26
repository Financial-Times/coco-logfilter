
[![Circle CI](https://circleci.com/gh/Financial-Times/coco-logfilter/tree/master.png?style=shield)](https://circleci.com/gh/Financial-Times/coco-logfilter/tree/master) [![Go Report Card](https://goreportcard.com/badge/github.com/Financial-Times/coco-logfilter)](https://goreportcard.com/report/github.com/Financial-Times/coco-logfilter) [![Coverage Status](https://coveralls.io/repos/github/Financial-Times/coco-logfilter/badge.svg)](https://coveralls.io/github/Financial-Times/coco-logfilter)
# coco-logfilter
Simple pre-aggregation logfilter to operate on the json output of journald.

## Installation
Download the project:

`go get github.com/Financial-Times/coco-logfilter/logfilter`

## Building

### Install dep

#### Mac
`brew install dep` / `brew upgrade dep`

#### Other Platforms
`curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh`

#### Run dep ensure
`dep ensure`

### Build
```
CGO_ENABLED=0 go build -a -installsuffix cgo -o coco-logfilter

docker build -t coco/coco-logfilter .
```

## Example use
```
echo '{"MESSAGE":"127.0.0.1 - - [21/Apr/2015:12:15:34 +0000] \"GET /eom-file/all/zzzz HTTP/1.1\" 200 12345 919 919"}' | ./coco-logfilter
```
