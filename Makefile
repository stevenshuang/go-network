.PHONY: clean build

EXEC_NAME := main

clean:
	rm -rf $(EXEC_NAME)

build: clean
	@go build -a -o $(EXEC_NAME)

