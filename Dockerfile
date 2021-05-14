FROM golang:1.9.2 
ADD . /go/src/myapp
WORKDIR /go/src/myapp
RUN go get myapp
RUN go install
ENTRYPOINT ["/go/bin/myapp"]