cover:
	go test -p 1 ./... -coverpkg=./internal/... -coverprofile=coverage.out
	go tool cover -html=coverage.out
