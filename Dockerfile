FROM golang:alpine

COPY . /concourse/pullrequest-resource

ENV GOPATH /concourse/pullrequest-resource
ENV PATH ${GOPATH}/bin:${PATH}

WORKDIR /concourse/pullrequest-resource

RUN go build -o /assets/out src/pullrequest/cmd/out/main.go
RUN go build -o /assets/in src/pullrequest/cmd/in/main.go
RUN go build -o /assets/check src/pullrequest/cmd/check/main.go
