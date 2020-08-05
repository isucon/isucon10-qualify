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

DESCRIPTION_LINES_FILE = "./description.txt"
OUTPUT_SQL_FILE = "./result/1_DummyEstateData.sql"
OUTPUT_TXT_FILE = "./result/estate_json.txt"
OUTPUT_FIXTURE_FILE = "./result/estate_condition.json"
ESTATE_IMAGE_ORIGIN_DIR = "./origin/estate"
ESTATE_IMAGE_PUBLIC_DIR = "../webapp/frontend/public/images/estate"
ESTATE_DUMMY_IMAGE_NUM = 1000
RECORD_COUNT = 10 ** 4
BULK_INSERT_COUNT = 500
DOOR_MIN_CENTIMETER = 30
DOOR_MAX_CENTIMETER = 200
DOOR_HEIGHT_RANGE_SEPARATORS = [80, 110, 150]
DOOR_WIDTH_RANGE_SEPARATORS = [80, 110, 150]
RENT_RANGE_SEPARATORS = [50000, 100000, 150000]
MIN_VIEW_COUNT = 3000
MAX_VIEW_COUNT = 1000000

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
]

ESTATE_FEATURE_FOR_VERIFY = "デザイナーズ物件"

ESTATE_IMAGE_HASH_LIST = [fake.sha256(
    raw_output=False) for _ in range(ESTATE_DUMMY_IMAGE_NUM)]


def generate_ranges_from_separator(separators):
    before = -1
    ranges = []

    for i, separator in enumerate(separators + [-1]):
        ranges.append({
            "id": i,
            "min": before,
            "max": separator
        })
        before = separator

    return ranges


def read_src_file_data(file_path):
    with open(file_path, mode='rb') as img:
        return img.read()


def dump_estate_to_json_str(estate):
    return json.dumps({
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


def generate_estate_dummy_data(estate_id, wrap={}):
    latlng = fake.local_latlng(country_code='JP', coords_only=True)
    feature_length = random.randint(0, len(ESTATE_FEATURE_LIST) - 1)
    image_hash = fake.word(ext_word_list=ESTATE_IMAGE_HASH_LIST)

    estate = {
        "id": estate_id,
        "thumbnail": f'/images/estate/{image_hash}.png',
        "name": fake.word(ext_word_list=BUILDING_NAME_LIST).format(name=fake.last_name()),
        "latitude": float(latlng[0]) + random.normalvariate(mu=0.0, sigma=0.3),
        "longitude": float(latlng[1]) + random.normalvariate(mu=0.0, sigma=0.3),
        "address": fake.address(),
        "rent": random.randint(30000, 200000),
        "door_height": random.randint(DOOR_MIN_CENTIMETER, DOOR_MAX_CENTIMETER),
        "door_width": random.randint(DOOR_MIN_CENTIMETER, DOOR_MAX_CENTIMETER),
        "view_count": random.randint(MIN_VIEW_COUNT, MAX_VIEW_COUNT),
        "description": random.choice(desc_lines).strip(),
        "features": ','.join(fake.words(nb=feature_length, ext_word_list=ESTATE_FEATURE_LIST, unique=True))
    }
    return dict(estate, **wrap)


if __name__ == '__main__':

    for i, random_hash in enumerate(ESTATE_IMAGE_HASH_LIST):
        image_data_list = [read_src_file_data(
            image) for image in glob.glob(os.path.join(ESTATE_IMAGE_ORIGIN_DIR, "*.png"))]
        with open(os.path.join(ESTATE_IMAGE_PUBLIC_DIR, f"{random_hash}.png"), mode='wb') as image_file:
            image_file.write(
                image_data_list[i % len(image_data_list)] + random_hash.encode('utf-8'))

    with open(DESCRIPTION_LINES_FILE, mode='r', encoding='utf-8') as description_lines:
        desc_lines = description_lines.readlines()

    with open(OUTPUT_SQL_FILE, mode='w', encoding='utf-8') as sqlfile, open(OUTPUT_TXT_FILE, mode='w', encoding='utf-8') as txtfile:
        if RECORD_COUNT % BULK_INSERT_COUNT != 0:
            raise Exception("The results of RECORD_COUNT and BULK_INSERT_COUNT need to be a divisible number. RECORD_COUNT = {}, BULK_INSERT_COUNT = {}".format(
                RECORD_COUNT, BULK_INSERT_COUNT))

        estate_id = 1

        ESTATES_FOR_VERIFY = [
            # 2回閲覧された後の検索で、順番が前に行くことを検証するためのデータ (2位 → 1位)
            generate_estate_dummy_data(1, {
                "features": ESTATE_FEATURE_FOR_VERIFY,
                "view_count": (MAX_VIEW_COUNT + MIN_VIEW_COUNT) // 2
            }),
            # 2回閲覧された後の検索で、順番が前に行くことを検証するためのデータ (1位 → 2位)
            generate_estate_dummy_data(2, {
                "features": ESTATE_FEATURE_FOR_VERIFY,
                "view_count": (MAX_VIEW_COUNT + MIN_VIEW_COUNT) // 2 + 1
            })
        ]

        sqlCommand = f"""INSERT INTO isuumo.estate (id, thumbnail, name, latitude, longitude, address, rent, door_height, door_width, view_count, description, features) VALUES {', '.join(map(lambda estate: f"('{estate['id']}', '{estate['thumbnail']}', '{estate['name']}', '{estate['latitude']}' , '{estate['longitude']}', '{estate['address']}', '{estate['rent']}', '{estate['door_height']}', '{estate['door_width']}', '{estate['view_count']}', '{estate['description']}', '{estate['features']}')", ESTATES_FOR_VERIFY))};"""
        sqlfile.write(sqlCommand)
        txtfile.write("\n".join([dump_estate_to_json_str(estate)
                                 for estate in ESTATES_FOR_VERIFY]) + "\n")

        estate_id += len(ESTATES_FOR_VERIFY)

        for _ in range(RECORD_COUNT//BULK_INSERT_COUNT):
            bulk_list = [generate_estate_dummy_data(
                estate_id + i) for i in range(BULK_INSERT_COUNT)]
            estate_id += BULK_INSERT_COUNT
            sqlCommand = f"""INSERT INTO isuumo.estate (id, thumbnail, name, latitude, longitude, address, rent, door_height, door_width, view_count, description, features) VALUES {', '.join(map(lambda estate: f"('{estate['id']}', '{estate['thumbnail']}', '{estate['name']}', '{estate['latitude']}' , '{estate['longitude']}', '{estate['address']}', '{estate['rent']}', '{estate['door_height']}', '{estate['door_width']}', '{estate['view_count']}', '{estate['description']}', '{estate['features']}')", bulk_list))};"""
            sqlfile.write(sqlCommand)
            txtfile.write("\n".join([dump_estate_to_json_str(estate)
                                     for estate in bulk_list]) + "\n")

    with open(OUTPUT_FIXTURE_FILE, mode='w', encoding='utf-8') as fixture_file:
        fixture_file.write(json.dumps({
            "doorWidth": {
                "prefix": "",
                "suffix": "cm",
                "ranges": generate_ranges_from_separator(DOOR_WIDTH_RANGE_SEPARATORS)
            },
            "doorHeight": {
                "prefix": "",
                "suffix": "cm",
                "ranges": generate_ranges_from_separator(DOOR_HEIGHT_RANGE_SEPARATORS)
            },
            "rent": {
                "prefix": "",
                "suffix": "円",
                "ranges": generate_ranges_from_separator(RENT_RANGE_SEPARATORS)
            },
            "feature": {
                "list": ESTATE_FEATURE_LIST + [ESTATE_FEATURE_FOR_VERIFY]
            }
        }, ensure_ascii=False, indent=2))
