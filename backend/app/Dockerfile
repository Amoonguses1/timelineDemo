FROM golang:1.22

WORKDIR /go/src/app
COPY . .
RUN go mod download && go mod verify
RUN go build -o /usr/local/bin/server ./cmd
ENTRYPOINT ["/usr/local/bin/server"]
