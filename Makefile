##
# Snippetbox
#
# @file
# @version 0.1


.DEFAULT_GOAL := build

.PHONY:vet build

vet:
	go vet ./...

build: vet
	go build

run: vet
	go run .

clean:
	go clean -x

# end
