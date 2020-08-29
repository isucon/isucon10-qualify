.PHONY: help install dependencies run clean

PORT := 1323

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

run:
	php -S localhost:$(PORT) -t public public/index.php

vendor:
	composer install

clean:
	rm -rf vendor
