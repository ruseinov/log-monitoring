deps:
	rm -rf vendor
	dep ensure --vendor-only
docker:
	docker build -t ruseinov/log-monitoring:latest .
