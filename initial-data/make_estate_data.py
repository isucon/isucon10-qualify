# coding:utf-8
import random
import time
import string
import json
import os
from faker import Faker
fake = Faker('ja_JP')
Faker.seed(19700101)
random.seed(19700101)

DESCRIPTION_LINES_FILE = "./description_estate.txt"
OUTPUT_SQL_FILE = "./result/1_DummyEstateData.sql"
OUTPUT_TXT_FILE = "./result/estate_json.txt"
ESTATE_IMAGE_ORIGIN_DIR = "./origin/estate"
ESTATE_IMAGE_PUBLIC_DIR = "../webapp/frontend/public/images/estate"
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
    os.path.join(ESTATE_IMAGE_ORIGIN_DIR, "1501E5C34A2B8EE645480ED1CC6442CD5929FE7616E20513574628096163DF0C.jpg"),
    os.path.join(ESTATE_IMAGE_ORIGIN_DIR, "3E880A828B1DBFACB42209724583B56EF28466E45E2BF3704475EA02B19BDBFC.jpg"),
    os.path.join(ESTATE_IMAGE_ORIGIN_DIR, "9120C2E3CAF5CD376C1B14899C2FD31438A839D1F6B6F8A52091392E0B9168FC.jpg")
]

def read_src_file_data(file_path):
    with open(file_path, mode='rb') as img:
        return img.read()

ESTATE_IMAGE_DATA = [read_src_file_data(image) for image in ESTATE_IMAGE_LIST]

def generate_estate_dummy_data(estate_id):
    latlng = fake.local_latlng(country_code='JP', coords_only=True)
    feature_length = random.randint(0, len(ESTATE_FEATURE_LIST) - 1)

    new_estate_image_hash = fake.sha256(raw_output=False)
    new_estate_image_path = os.path.join(ESTATE_IMAGE_PUBLIC_DIR, "{}.jpg".format(new_estate_image_hash))
    src_image_data = random.choice(ESTATE_IMAGE_DATA)

    with open(new_estate_image_path, mode='wb') as dst_image_data:
        dst_image_data.write(src_image_data + new_estate_image_hash.encode('utf-8'))

    src_estate_image_filename = fake.word(ext_word_list=ESTATE_IMAGE_LIST)

    return {
        "id": estate_id,
        "thumbnail": '/images/estate/{}'.format(new_estate_image_path),
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
    for image in ESTATE_IMAGE_LIST:
        filename, _ = os.path.splitext(os.path.basename(image))
        with open(os.path.join(ESTATE_IMAGE_PUBLIC_DIR, "{}.jpg".format(filename)), mode='wb') as dst_image_data:
            dst_image_data.write(read_src_file_data(image) + filename.encode('utf-8'))

    with open(DESCRIPTION_LINES_FILE, mode='r') as description_lines:
        desc_lines = description_lines.readlines()

    with open(OUTPUT_SQL_FILE, mode='w') as sqlfile, open(OUTPUT_TXT_FILE, mode='w') as txtfile:
        sqlfile.write(sqlCommands)
        if RECORD_COUNT % BULK_INSERT_COUNT != 0:
            raise Exception("The results of RECORD_COUNT and BULK_INSERT_COUNT need to be a divisible number. RECORD_COUNT = {}, BULK_INSERT_COUNT = {}".format(RECORD_COUNT, BULK_INSERT_COUNT))

        estate_id = 1
        for _ in range(RECORD_COUNT//BULK_INSERT_COUNT):
            bulk_list = [generate_estate_dummy_data(estate_id + i) for i in range(BULK_INSERT_COUNT)]
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
