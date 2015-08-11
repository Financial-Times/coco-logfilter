FROM gliderlabs/alpine:3.2
ADD logfilter/coco-logfilter /coco-logfilter
CMD /coco-logfilter -environment=$ENV
