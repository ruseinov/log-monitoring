deps:
	rm -rf vendor
	dep ensure --vendor-only

docker: deps
	docker build -t ruseinov/log-monitoring:latest .

build: deps
	CGO_ENABLED=0 go build -a -installsuffix cgo -o log-monitor `pwd`/cmd/monitor
	mv `pwd`/log-monitor $(GOPATH)/bin/