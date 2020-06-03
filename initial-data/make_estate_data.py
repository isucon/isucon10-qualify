# coding:utf-8
import random
import time
import string
import json
from faker import Faker
fake = Faker('ja_JP')
Faker.seed(19700101)
random.seed(19700101)

DESCRIPTION_LINES_FILE = "./description_estate.txt"
OUTPUT_SQL_FILE = "./result/1_DummyEstateData.sql"
OUTPUT_TXT_FILE = "./result/estate_json.txt"
RECORD_COUNT = 10 ** 4
BULK_INSERT_COUNT = 500
DOOR_MIN_CENTIMETER = 30
DOOR_MAX_CENTIMETER = 200
sqlCommands = ""
sqlCommands += "use isuumo;\n"

BUILDING_NAME_LIST = [
    "{name}ISUビルディング",
    "ISUアパート {name}",
    "ISU{name}レジデンス",
    "ISUガーデン {name}",
    "{name} ISUマンション",
    "{name} ISUビル"
]

ESTATE_FEATURE_LIST = [
    "バストイレ別",
    "駅から徒歩5分",
    "ペット飼育可能",
    "デザイナーズ物件",
]

def generate_estate_dummy_data():
    latlng = fake.local_latlng(country_code='JP', coords_only=True)
    feature_length = random.randint(0, len(ESTATE_FEATURE_LIST) - 1)
    return {
        "thumbnail": '/images/estate/{}.jpg'.format(fake.sha256(raw_output=False)),
        "name": fake.word(ext_word_list=BUILDING_NAME_LIST).format(name=fake.last_name()),
        "latitude": latlng[0],
        "longitude": latlng[1],
        "address": fake.address(),
        "rent": random.randint(50000, 200000),
        "door_height": random.randint(DOOR_MIN_CENTIMETER, DOOR_MAX_CENTIMETER),
        "door_width": random.randint(DOOR_MIN_CENTIMETER, DOOR_MAX_CENTIMETER),
        "view_count": random.randint(3000, 1000000),
        "description": random.choice(desc_lines).strip(),
        "features": ','.join(fake.words(nb=feature_length, ext_word_list=ESTATE_FEATURE_LIST, unique=True))
    }

if __name__ == '__main__':
    with open(DESCRIPTION_LINES_FILE, mode='r') as description_lines:
        desc_lines = description_lines.readlines()

    with open(OUTPUT_SQL_FILE, mode='w') as sqlfile, open(OUTPUT_TXT_FILE, mode='w') as txtfile:
        sqlfile.write(sqlCommands)
        if RECORD_COUNT % BULK_INSERT_COUNT != 0:
            raise Exception("The results of RECORD_COUNT and BULK_INSERT_COUNT need to be a divisible number. RECORD_COUNT = {}, BULK_INSERT_COUNT = {}".format(RECORD_COUNT, BULK_INSERT_COUNT))

        for _ in range(RECORD_COUNT//BULK_INSERT_COUNT):
            bulk_list = [generate_estate_dummy_data() for i in range(BULK_INSERT_COUNT)]
            sqlCommand = f"""INSERT INTO estate (thumbnail, name, latitude, longitude, address, rent, door_height, door_width, view_count, description, features) VALUES {', '.join(map(lambda estate: f"('{estate['thumbnail']}', '{estate['name']}', '{estate['latitude']}' , '{estate['longitude']}', '{estate['address']}', '{estate['rent']}', '{estate['door_height']}', '{estate['door_width']}', '{estate['view_count']}', '{estate['description']}', '{estate['features']}')", bulk_list))};"""
            sqlfile.write(sqlCommand)

            for estate in bulk_list:
                txtfile.write(json.dumps(estate, ensure_ascii=False) + "\n")
