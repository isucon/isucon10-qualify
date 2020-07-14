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
OUTPUT_CSV_FILE = "./result/estateData.csv"
OUTPUT_TXT_FILE = "./result/estate_json.txt"
ESTATE_IMAGE_ORIGIN_DIR = "./origin/estate"
ESTATE_IMAGE_PUBLIC_DIR = "../webapp/frontend/public/images/estate"
ESTATE_DUMMY_IMAGE_NUM = 1000
RECORD_COUNT = 10 ** 4
DOOR_MIN_CENTIMETER = 30
DOOR_MAX_CENTIMETER = 200
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


def read_src_file_data(file_path):
    with open(file_path, mode='rb') as img:
        return img.read()


def dump_estate_to_csv_str(estate):
    # id, thumbnail, name, latitude, longitude, address, rent, door_height, door_width, view_count, description, features
    return f"{estate['id']},\"{estate['thumbnail']}\",\"{estate['name']}\",{estate['latitude']},{estate['longitude']},\"{estate['address']}\",{estate['rent']},{estate['door_height']},{estate['door_width']},{estate['view_count']},\"{estate['description']}\",\"{estate['features']}\""


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

    with open(OUTPUT_CSV_FILE, mode='w', encoding='utf-8') as csvfile, open(OUTPUT_TXT_FILE, mode='w', encoding='utf-8') as txtfile:
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

        estates = ESTATES_FOR_VERIFY + \
            [generate_estate_dummy_data(len(ESTATES_FOR_VERIFY) + i + 1)
             for i in range(RECORD_COUNT)]

        csvfile.write(
            "\n".join([dump_estate_to_csv_str(estate) for estate in estates]))

        txtfile.write("\n".join([dump_estate_to_json_str(estate)
                                 for estate in estates]) + "\n")
