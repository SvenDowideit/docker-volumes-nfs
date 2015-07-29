FROM golang:1.5

#TODO: really should extract this into its own onbuild Dockerfile

# turn on golang experiment to add vendoring
# see ttps://medium.com/@freeformz/go-1-5-s-vendor-experiment-fd3e830f52c3
ENV GO15VENDOREXPERIMENT 1

COPY . /go/src/app
WORKDIR /go/src/app
RUN go-wrapper download
RUN go-wrapper install
RUN go build
ENTRYPOINT ["./app"]
