FROM golang:alpine

COPY . /concourse/pullrequest

ENV GOPATH /concourse/pullrequest-resource
ENV PATH ${GOPATH}/bin:${PATH}

RUN go build -o /assets/out /concourse/pullrequest-resource/cmd/out
RUN go build -o /assets/in /concourse/pullrequest-resource/cmd/in
RUN go build -o /assets/check /concourse/pullrequest-resource/cmd/check
