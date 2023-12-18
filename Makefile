TEST_PKG=$$(go list ./... | grep -v -e /mocks/ -e /docs$$)
COV_PKG  = $$(echo $(TEST_PKG) | tr ' ' ',' )
COV_OUTPUT=cover.out
COV = $$(go tool cover -func $(COV_OUTPUT) | grep total | awk '{print substr($$3, 1, length($$3))}')

unittest:
	go test -short $(TEST_PKG)

test-coverage:
	echo ${COV_PKG}
	go clean -testcache
	go test -count=1 -coverpkg=$(COV_PKG) -coverprofile=$(COV_OUTPUT) $(TEST_PKG)
	echo test coverage: $(COV)

test-coverage-html:
	go test -coverprofile=$(COV_OUTPUT) $(TEST_PKG)
	go tool cover -html=$(COV_OUTPUT)

mock:
	mockery --dir=. --exclude compose --exclude mocks --replace-type cloud.google.com/go/internal/pubsub=cloud.google.com/go/pubsub