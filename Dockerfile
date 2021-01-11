FROM golang:1.9.4 AS build
RUN apt-get update 
RUN apt-get install -y xz-utils zip rsync
RUN go get github.com/golang/lint/golint

ENV PATH /go/bin:$PATH

WORKDIR /go/src

COPY . /go/src/
RUN CGO_ENABLED=0 go build -ldflags="-w -s -X main.VERSION=1 -linkmode external -extldflags -static" -o bin/hello-world


FROM alpine:3.11.6 AS production
RUN apk add --update bash libressl curl fping libcap && rm -rf /var/cache/apk/* && mkdir -p /opt 
COPY --from=build /go/src/bin/hello-world /opt/hello-world/
COPY img/* /opt/hello-world/

WORKDIR /opt/hello-world
ENTRYPOINT ["/opt/hello-world/hello-world"]
