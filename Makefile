VERSION=1.0.0
BINARY_PATH=./bin
CSV_PARSER_BINARY_NAME=$(BINARY_PATH)/csv-parser.bin

clean:
	@ rm -rf bin
	@ rm -f cover.out.tmp

build:
	@ echo " ---         BUILDING CSV PARSER     --- "
	@ $(MAKE) clean
	@ go build -ldflags "-s -w -X main.version=$(VERSION)" -o $(CSV_PARSER_BINARY_NAME) cmd/main.go
	@ echo " ---     BUILD FINISHED      --- "

test:
	@ go test -v -race ./... -coverprofile=cover.out.tmp

cover-html:
	@ go tool cover -html=cover.out.tmp
