FROM golang
ADD . /go/src/github.com/tobyjsullivan/ocs-orders
RUN  go install github.com/tobyjsullivan/ocs-orders
CMD /go/bin/ocs-orders
