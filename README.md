# coco-logfilter
simple pre-aggregation logfilter to operate on the json output of journald

##Building
```
cd logfilter
CGO_ENABLED=0 go build -a -installsuffix cgo -o coco-logfilter .
cd ..

docker build -t coco/coco-logfilter .
```

##Installation
```
go get github.com/Financial-Times/coco-logfilter/logfilter

```

##Example use
```
echo '{"MESSAGE":"127.0.0.1 - - [21/Apr/2015:12:15:34 +0000] \"GET /eom-file/all/zzzz HTTP/1.1\" 200 12345 919 919"}' | logfilter
```
