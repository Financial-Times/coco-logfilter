FROM golang:1.8-alpine

ENV PROJECT=coco-logfilter
COPY . /${PROJECT}-sources/

RUN apk --no-cache --virtual .build-dependencies add git \
  && ORG_PATH="github.com/Financial-Times" \
  && REPO_PATH="${ORG_PATH}/${PROJECT}" \
  && mkdir -p $GOPATH/src/${ORG_PATH} \
  # Linking the project sources in the GOPATH folder
  && ln -s /${PROJECT}-sources $GOPATH/src/${REPO_PATH} \
  && cd $GOPATH/src/${REPO_PATH} \
  && go get ./... \
  && cd logfilter \
  && go build \
  && ls -lah \
  && mv logfilter /${PROJECT} \
  && apk del .build-dependencies \
  && rm -rf $GOPATH /var/cache/apk/*

WORKDIR /

CMD [ "/coco-logfilter" ]
