FROM golang:1.16-alpine

RUN mkdir -p /opt/src
WORKDIR /opt/src
COPY . .

RUN apk --no-cache add gcc libc-dev

RUN go build -a -trimpath -ldflags "-s -w -extldflags '-static'" -installsuffix cgo -tags netgo -o /bin/action .
RUN rm  -rf /opt/src;

ENTRYPOINT ["/bin/action"]
