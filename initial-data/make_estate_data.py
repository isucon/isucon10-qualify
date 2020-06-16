# coding:utf-8
import random
import time
import string
import json
import os
import glob
from faker import Faker
fake = Faker('ja_JP')
Faker.seed(19700101)
random.seed(19700101)

DESCRIPTION_LINES_FILE = "./description_estate.txt"
OUTPUT_SQL_FILE = "./result/1_DummyEstateData.sql"
OUTPUT_TXT_FILE = "./result/estate_json.txt"
ESTATE_IMAGE_ORIGIN_DIR = "./origin/estate"
ESTATE_IMAGE_PUBLIC_DIR = "../webapp/frontend/public/images/estate"
ESTATE_DUMMY_IMAGE_NUM = 1000
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

ESTATE_IMAGE_HASH_LIST = [fake.sha256(
    raw_output=False) for _ in range(ESTATE_DUMMY_IMAGE_NUM)]


def read_src_file_data(file_path):
    with open(file_path, mode='rb') as img:
        return img.read()


def generate_estate_dummy_data(estate_id):
    latlng = fake.local_latlng(country_code='JP', coords_only=True)
    feature_length = random.randint(0, len(ESTATE_FEATURE_LIST) - 1)
    image_hash = fake.word(ext_word_list=ESTATE_IMAGE_HASH_LIST)

    return {
        "id": estate_id,
        "thumbnail": f'/images/estate/{image_hash}.png',
        "name": fake.word(ext_word_list=BUILDING_NAME_LIST).format(name=fake.last_name()),
        "latitude": float(latlng[0]) + random.normalvariate(mu=0.0, sigma=0.3),
        "longitude": float(latlng[1]) + random.normalvariate(mu=0.0, sigma=0.3),
        "address": fake.address(),
        "rent": random.randint(30000, 200000),
        "door_height": random.randint(DOOR_MIN_CENTIMETER, DOOR_MAX_CENTIMETER),
        "door_width": random.randint(DOOR_MIN_CENTIMETER, DOOR_MAX_CENTIMETER),
        "view_count": random.randint(3000, 1000000),
        "description": random.choice(desc_lines).strip(),
        "features": ','.join(fake.words(nb=feature_length, ext_word_list=ESTATE_FEATURE_LIST, unique=True))
    }


if __name__ == '__main__':
    for i, random_hash in enumerate(ESTATE_IMAGE_HASH_LIST):
        image_data_list = [read_src_file_data(
            image) for image in glob.glob(os.path.join(ESTATE_IMAGE_ORIGIN_DIR, "*.png"))]
        with open(os.path.join(ESTATE_IMAGE_PUBLIC_DIR, f"{random_hash}.png"), mode='wb') as image_file:
            image_file.write(
                image_data_list[i % len(image_data_list)] + random_hash.encode('utf-8'))

    with open(DESCRIPTION_LINES_FILE, mode='r') as description_lines:
        desc_lines = description_lines.readlines()

    with open(OUTPUT_SQL_FILE, mode='w') as sqlfile, open(OUTPUT_TXT_FILE, mode='w') as txtfile:
        sqlfile.write(sqlCommands)
        if RECORD_COUNT % BULK_INSERT_COUNT != 0:
            raise Exception("The results of RECORD_COUNT and BULK_INSERT_COUNT need to be a divisible number. RECORD_COUNT = {}, BULK_INSERT_COUNT = {}".format(
                RECORD_COUNT, BULK_INSERT_COUNT))

        estate_id = 1
        for _ in range(RECORD_COUNT//BULK_INSERT_COUNT):
            bulk_list = [generate_estate_dummy_data(
                estate_id + i) for i in range(BULK_INSERT_COUNT)]
            estate_id += BULK_INSERT_COUNT
            sqlCommand = f"""INSERT INTO estate (id, thumbnail, name, latitude, longitude, address, rent, door_height, door_width, view_count, description, features) VALUES {', '.join(map(lambda estate: f"('{estate['id']}', '{estate['thumbnail']}', '{estate['name']}', '{estate['latitude']}' , '{estate['longitude']}', '{estate['address']}', '{estate['rent']}', '{estate['door_height']}', '{estate['door_width']}', '{estate['view_count']}', '{estate['description']}', '{estate['features']}')", bulk_list))};"""
            sqlfile.write(sqlCommand)

            for estate in bulk_list:
                json_string = json.dumps({
                    "id": estate["id"],
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
