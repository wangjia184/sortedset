.PHONY: test
test:
	go test -v -coverprofile .coverage.txt ./...
	go tool cover -func .coverage.txt