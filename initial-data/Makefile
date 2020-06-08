all: chair_data estate_data

refresh_chair_data:
	find ../webapp/frontend/public/images/chair -name *.jpg | xargs rm

refresh_estate_data:
	find ../webapp/frontend/public/images/estate -name *.jpg | xargs rm

chair_data: make_chair_data.py
	python3 make_chair_data.py
	cp result/2_DummyChairData.sql ../webapp/mysql/db/2_DummyChairData.sql

estate_data: make_estate_data.py
	python3 make_estate_data.py
	cp result/1_DummyEstateData.sql ../webapp/mysql/db/1_DummyEstateData.sql

