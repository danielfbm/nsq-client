FROM index.alauda.cn/alaudaorg/gobuild:1.8-alpine

COPY . /go/src/github.com/danielfbm/nsq-client

RUN  go build -ldflags "-w -s" -v -o /go/bin/nsq-client github.com/danielfbm/nsq-client \
   && chmod +x /go/bin/nsq-client

CMD ["nsq-client produce"]