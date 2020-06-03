SHELL=/bin/bash
run/go:
	docker-compose -f docker-compose.yaml up --build
run/test/with-makedata:
	docker-compose -f docker-compose-test.yaml down -v
	if [ -e ./mysql/db/1_DummyEstateData.sql ]; then rm ./mysql/db/1_DummyEstateData.sql; fi
	if [ -e ./mysql/db/2_DummyChairData.sql ];then rm ./mysql/db/2_DummyChairData.sql; fi
	python3 ./mysql/make_chair_data_SQL.py
	python3 ./mysql/make_estate_data_SQL.py
	docker-compose -f docker-compose-test.yaml up --build
run/test:
	docker-compose -f docker-compose-test.yaml down -v
	docker-compose -f docker-compose-test.yaml up --build
# run/node:
# 	docker-compose -f docker-compose-node.yaml up
# run/java:
# 	docker-compose -f docker-compose-java.yaml up

# .PHONY: run/node run/java
