.PHONY: run-server lint-go 

run-server:
	cd server && go run main.go

lint-go:
	cd server && golangci-lint run