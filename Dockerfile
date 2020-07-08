FROM golang AS builder

WORKDIR /resqu
ADD . .
RUN go build -ldflags "-s -w" -o resqu


FROM debian
WORKDIR /resqu
COPY --from=builder /resqu/resqu .

CMD [ "/resqu/resqu" ]