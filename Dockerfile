FROM golang:alpine

COPY . /concourse/pullrequest-resource

ENV GOPATH /concourse/pullrequest-resource
ENV PATH ${GOPATH}/bin:${PATH}

# RUN go build -o /assets/out /concourse/pullrequest-resource/src/pullrequest/cmd/out
# RUN go build -o /assets/in /concourse/pullrequest-resource/src/pullrequest/cmd/in
# RUN go build -o /assets/check /concourse/pullrequest-resource/cmd/src/pullrequest/check
