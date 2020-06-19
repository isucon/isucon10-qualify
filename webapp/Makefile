SHELL=/bin/bash

export URL="http://localhost:1323/api/estate/range"
export NOW:=`date '+%Y_%m_%d_%T'`
export PPROF=cpu.pprof
export FILETIME=`tail -n 1 ./logs/record_list.log`
spec/measurement:
	echo ${NOW}
	docker-compose down -v
	docker-compose up --build -d
	wayt http -u ${URL}
	curl -s http://localhost:1323/debug/pprof/profile?seconds=80 > ./logs/apps/cpu.pprof
	docker-compose down
	mv ./logs/apps/cpu.pprof ./logs/apps/${NOW}${PPROF}
	echo ${NOW} >> ./logs/record_list.log
spec/pprof:
	go tool pprof ./logs/apps/${FILETIME}cpu.pprof

spec/mysql:
	querydigest -f ./logs/mysql/mysql-slow.log

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
