.PHONY: help install dependencies clean

help:
	@cat $(firstword $(MAKEFILE_LIST))

install: \
	dependencies \
	build

build: \
	vendor

dependencies:
	type php > /dev/null
	type composer > /dev/null

vendor:
	composer install

clean:
	rm -rf vendor
