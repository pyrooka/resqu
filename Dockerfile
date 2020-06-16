FROM golang:buster

WORKDIR /resqu

COPY db ./db
COPY config.go go.mod go.sum main.go ./

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["resqu"]