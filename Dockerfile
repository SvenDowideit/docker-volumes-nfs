FROM golang:1.5

#TODO: really should extract this into its own onbuild Dockerfile

# Use out version of go-wrapper to add vendoring
COPY go-wrapper /usr/local/bin/

COPY . /go/src/app
RUN go-wrapper download
RUN go-wrapper install
