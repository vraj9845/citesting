.PHONY: setup 
setup:
	docker-compose up -d
	
.PHONY: test
test:
	go test -v ./...
