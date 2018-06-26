FROM golang:alpine

RUN apk update && \
    apk add git

COPY . /concourse/pullrequest-resource

ENV GOPATH /concourse/pullrequest-resource
ENV PATH ${GOPATH}/bin:${PATH}

WORKDIR /concourse/pullrequest-resource

RUN go build -o /opt/resource/out src/pullrequest/cmd/out/main.go
RUN go build -o /opt/resource/in src/pullrequest/cmd/in/main.go
RUN go build -o /opt/resource/check src/pullrequest/cmd/check/main.go
