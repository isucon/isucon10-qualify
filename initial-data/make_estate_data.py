# coding:utf-8
import random
import time
import string
import json
import os
import shutil
from faker import Faker
fake = Faker('ja_JP')
Faker.seed(19700101)
random.seed(19700101)

DESCRIPTION_LINES_FILE = "./description_estate.txt"
OUTPUT_SQL_FILE = "./result/1_DummyEstateData.sql"
OUTPUT_TXT_FILE = "./result/estate_json.txt"
ESTATE_IMAGE_FILE_PATH = "../webapp/frontend/public/images/estate"
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

ESTATE_IMAGE_LIST = [
    os.path.join(ESTATE_IMAGE_FILE_PATH, "1501E5C34A2B8EE645480ED1CC6442CD5929FE7616E20513574628096163DF0C.jpg"),
    os.path.join(ESTATE_IMAGE_FILE_PATH, "3E880A828B1DBFACB42209724583B56EF28466E45E2BF3704475EA02B19BDBFC.jpg"),
    os.path.join(ESTATE_IMAGE_FILE_PATH, "9120C2E3CAF5CD376C1B14899C2FD31438A839D1F6B6F8A52091392E0B9168FC.jpg")
]

def generate_estate_dummy_data():
    latlng = fake.local_latlng(country_code='JP', coords_only=True)
    feature_length = random.randint(0, len(ESTATE_FEATURE_LIST) - 1)

    new_estate_image_name = os.path.join(ESTATE_IMAGE_FILE_PATH, "{}.jpg".format(fake.sha256(raw_output=False)))
    src_estate_image_filename = fake.word(ext_word_list=ESTATE_IMAGE_LIST)

    shutil.copy(src_estate_image_filename, new_estate_image_name)

    return {
        "thumbnail": '/images/estate/{}.jpg'.format(new_estate_image_name),
        "name": fake.word(ext_word_list=BUILDING_NAME_LIST).format(name=fake.last_name()),
        "latitude": float(latlng[0]),
        "longitude": float(latlng[1]),
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
                json_string = json.dumps({
                    "thumbnail": estate["thumbnail"],
                    "name": estate["name"],
                    "latitude": estate["latitude"],
                    "longitude": estate["longitude"],
                    "address": estate["address"],
                    "rent": estate["rent"],
                    "doorHeight": estate["door_height"],
                    "doorWidth": estate["door_width"],
                    "viewCount": estate["view_count"],
                    "description": estate["description"],
                    "features": estate["features"]
                }, ensure_ascii=False)
                txtfile.write(json_string + "\n")
