FROM golang:1-buster

WORKDIR /resqu
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["resqu"]