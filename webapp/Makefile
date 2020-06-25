SHELL=/bin/bash

include ../bench/Makefile
export URL="http://localhost:1323/api/estate/range"
export NOW:=`date '+%Y_%m_%d_%H_%M'`
export PPROF=cpu.pprof
export MYSQL=mysql-slow.log
export NGINX=access.log
export FILETIME=`tail -n 1 ${PWD}/logs/record_list.log`

.DEFAULT_GOAL := help


spec/measurement: ## to measure performance, docker-compose up all and log-rotate
	echo ${NOW}
	docker-compose down -v
	docker-compose up -d --build mysql nginx api-server frontend
	wayt http -u ${URL}
	cd $(shell dirname ${PWD})/bench && make
	curl -s http://localhost:1323/debug/pprof/profile?seconds=80 > ./logs/apps/${PPROF} &
	sleep 5
	$(shell dirname ${PWD})/bench/bench -target-url "http://127.0.0.1:8080" -data-dir $(shell dirname ${PWD})/initial-data
	mv ${PWD}/logs/apps/${PPROF} ${PWD}/logs/apps/${NOW}${PPROF}
	mv ${PWD}/logs/mysql/${MYSQL} ${PWD}/logs/mysql/${NOW}${MYSQL}
	mv ${PWD}/logs/nginx/${NGINX} ${PWD}/logs/nginx/${NOW}${NGINX}
	echo ${NOW} >> ${PWD}/logs/record_list.log
	docker-compose down -v

spec/pprof: ## analyze pprof
	go tool pprof ./logs/apps/${FILETIME}${PPROF}

spec/mysql: ## analyze slow-query with querydigest
	querydigest -f ${PWD}/logs/mysql/${FILETIME}${MYSQL}

spec/accesslog: ## analyze accesslog with kataribe (kataribe.toml is in the same dir of this Makefile)
	cat ${PWD}/logs/nginx/${FILETIME}${NGINX}| kataribe

dev/api-server: ## docker-compose up with minimal api-server
	docker-compose down -v
	docker-compose -f docker-compose.yaml up --build mysql api-server

.PHONY: dev/test-with-makedata
dev/test-with-makedata: ## make new initial data and do test with docker-compose-test.yaml
	docker-compose -f docker-compose-test.yaml down -v
	if [ -e ./mysql/db/1_DummyEstateData.sql ]; then rm ./mysql/db/1_DummyEstateData.sql; fi
	if [ -e ./mysql/db/2_DummyChairData.sql ];then rm ./mysql/db/2_DummyChairData.sql; fi
	python3 ./mysql/make_chair_data_SQL.py
	python3 ./mysql/make_estate_data_SQL.py
	docker-compose -f docker-compose-test.yaml up --build

.PHONY: dev/test
dev/test: ## do test of api-server
	docker-compose -f docker-compose-test.yaml down -v
	docker-compose -f docker-compose-test.yaml up --build

# run/node:
# 	docker-compose -f docker-compose-node.yaml up
# run/java:
# 	docker-compose -f docker-compose-java.yaml up

# .PHONY: run/node run/java
#
#
.PHONY: help
help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+/?[a-zA-Z_-]*:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
