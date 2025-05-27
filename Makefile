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
	go build ./cmd/web/

run: vet
	go build ./cmd/web/ && ./web

clean:
	go clean -x

# end
