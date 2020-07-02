SHELL=/bin/bash

MAKEFILE_DIR:=$(dir $(abspath $(lastword $(MAKEFILE_LIST))))
export NOW:=`date '+%Y_%m_%d_%H_%M'`

.DEFAULT_GOAL := help

ci/test: ## make new initial data and do test with docker-compose.yaml
	docker-compose -f ${MAKEFILE_DIR}webapp/docker-compose-test.yaml down -v
	if [ -e ${MAKEFILE_DIR}webapp/mysql/db/1_DummyEstateData.sql ]; then rm ${MAKEFILE_DIR}/webapp/mysql/db/1_DummyEstateData.sql; fi
	if [ -e ${MAKEFILE_DIR}webapp/mysql/db/2_DummyChairData.sql ];then rm ${MAKEFILE_DIR}/webapp/mysql/db/2_DummyChairData.sql; fi
	cd ${MAKEFILE_DIR}initial-data && python3 make_chair_data.py
	mv ${MAKEFILE_DIR}initial-data/result/2_DummyChairData.sql ${MAKEFILE_DIR}/webapp/mysql/db
	cd ${MAKEFILE_DIR}initial-data && python3 make_estate_data.py
	mv ${MAKEFILE_DIR}initial-data/result/1_DummyEstateData.sql ${MAKEFILE_DIR}/webapp/mysql/db
	docker-compose -f ${MAKEFILE_DIR}webapp/docker-compose-test.yaml up --exit-code-from test-server --build mysql api-server test-server 

.PHONY: help
help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+/?[a-zA-Z_-]*:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
