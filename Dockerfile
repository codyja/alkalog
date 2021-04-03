FROM golang:1.16.3
WORKDIR /go/src/alkalog
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["alkalog"]