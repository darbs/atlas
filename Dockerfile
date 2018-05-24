#FROM golang:latest
FROM alpine:latest

RUN apk --update-cache --allow-untrusted \
        --repository http://dl-4.alpinelinux.org/alpine/edge/community \
        --arch=x86_64 add \
    git \
    go \
    bash \
    curl \
    && rm -rf /var/cache/apk/*

ENV GOROOT /usr/lib/go
ENV GOPATH /go
ENV PATH /go/bin:$PATH
ENV PROJPATH ${GOPATH}/src/github.com/darbs/atlas

RUN mkdir -p ${PROJPATH} ${GOPATH}/bin

RUN curl https://glide.sh/get | sh

COPY . ${PROJPATH}
WORKDIR ${PROJPATH}
#RUN ls -a

RUN glide install
RUN go build -o main .

EXPOSE 80

CMD ["/atlas/main"]