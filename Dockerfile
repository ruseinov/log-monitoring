FROM golang:1.11-stretch as builder
ADD . /go/src/github.com/ruseinov/log-monitoring
WORKDIR /go/src/github.com/ruseinov/log-monitoring/cmd/monitor/
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .
RUN touch /tmp/access.log

FROM scratch
WORKDIR /app
COPY --from=builder /go/src/github.com/ruseinov/log-monitoring/cmd/monitor/app .
COPY --from=builder /tmp/access.log /tmp/access.log
ENTRYPOINT ["./app"]

