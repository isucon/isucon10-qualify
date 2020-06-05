SHELL=/bin/bash

export API_SERVER_HOST=
export API_SERVER_PORT=


all: build export deploy

install: package.json
	npm ci

build: src
	npm run build

export: .next
	npm run export

deploy: out
	npm run deploy

