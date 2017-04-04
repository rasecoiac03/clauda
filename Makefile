HOME?=$$(HOME)
TESTS?=$$(go list ./... | egrep -v "vendor")

build:
	CGO_ENABLED=0 go build -v -a --installsuffix cgo --ldflags="-s" -o clauda

test:
	go test -v $(TESTS)

test-cover-html:
	./testing/coverage --html

test-coveralls:
	./testing/coverage --coveralls

run:
	go run main.go

install:
	CGO_ENABLED=0 go install -v -a --installsuffix cgo --ldflags="-s"

.PHONY: build test test-cover-html test-coveralls run install
