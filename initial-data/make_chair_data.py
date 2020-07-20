# coding:utf-8
import random
import time
import string
import json
import os
import glob
from faker import Faker
fake = Faker("ja_JP")
Faker.seed(19700101)
random.seed(19700101)

DESCRIPTION_LINES_FILE = "./description_chair.txt"
OUTPUT_CSV_FILE = "./result/chairData.csv"
OUTPUT_TXT_FILE = "./result/chair_json.txt"
CHAIR_IMAGE_ORIGIN_DIR = "./origin/chair"
CHAIR_IMAGE_PUBLIC_DIR = "../webapp/frontend/public/images/chair"
CSV_COLUMNS_ORDER = "id,thumbnail,name,price,height,width,depth,view_count,stock,color,description,features,kind\n"
CHAIR_DUMMY_IMAGE_NUM = 1000
RECORD_COUNT = 10 ** 4
CHAIR_MIN_CENTIMETER = 30
CHAIR_MAX_CENTIMETER = 200
MIN_VIEW_COUNT = 3000
MAX_VIEW_COUNT = 1000000

CHAIR_COLOR_LIST = [
    "黒",
    "白",
    "赤",
    "青",
    "緑",
    "黄",
    "紫",
    "ピンク",
    "オレンジ",
    "水色",
    "ネイビー",
    "ベージュ"
]

CHAIR_NAME_PREFIX_LIST = [
    "ふわふわ",
    "エルゴノミクス",
    "こだわりの逸品",
    "[期間限定]",
    "[残りわずか]",
    "オフィス",
    "[30%OFF]",
    "【お買い得】",
    "シンプル",
    "大人気！",
    "【伝説の一品】",
    "【本格仕様】"
]

CHAIR_PROPERTY_LIST = [
    "社長の",
    "俺の",
    "回転式",
    "ありきたりの"
    "すごい",
    "ボロい",
    "普通の",
    "アンティークな",
    "パイプ",
    "モダンな",
    "金の",
    "子供用",
    "アウトドア",
]

CHAIR_NAME_LIST = [
    "イス",
    "チェア",
    "フロアチェア",
    "ソファー",
    "ゲーミングチェア",
    "座椅子",
    "ハンモック",
    "オフィスチェア",
    "ダイニングチェア",
    "パイプイス",
    "椅子"
]

CHAIR_FEATURE_LIST = [
    "折りたたみ可",
    "肘掛け",
    "キャスター",
    "リクライニング",
    "高さ調節可"
]

CHAIR_FEATURE_FOR_VERIFY = "フットレスト"

CHAIR_KIND_LIST = [
    "ゲーミングチェア",
    "座椅子",
    "エルゴノミクス",
    "ハンモック"
]

CHAIR_IMAGE_HASH_LIST = [fake.sha256(
    raw_output=False) for _ in range(CHAIR_DUMMY_IMAGE_NUM)]


def read_src_file_data(file_path):
    with open(file_path, mode='rb') as img:
        return img.read()


def dump_chair_to_csv_str(chair):
    # id, thumbnail, name, price, height, width, depth, view_count, stock, color, description, features, kind
    return f"{chair['id']},\"{chair['thumbnail']}\",\"{chair['name']}\",{chair['price']},{chair['height']},{chair['width']},{chair['depth']},{chair['view_count']},{chair['stock']},\"{chair['color']}\",\"{chair['description']}\",\"{chair['features']}\",\"{chair['kind']}\""


def dump_chair_to_json_str(chair):
    return json.dumps({
        "id": chair["id"],
        "thumbnail": chair["thumbnail"],
        "name": chair["name"],
        "price": chair["price"],
        "height": chair["height"],
        "width": chair["width"],
        "depth": chair["depth"],
        "color": chair["color"],
        "view_count": chair["view_count"],
        "stock": chair["stock"],
        "description": chair["description"],
        "features": chair["features"],
        "kind": chair["kind"]
    }, ensure_ascii=False)


def generate_chair_dummy_data(chair_id, wrap={}):
    features_length = random.randint(0, len(CHAIR_FEATURE_LIST) - 1)
    image_hash = fake.word(ext_word_list=CHAIR_IMAGE_HASH_LIST)

    chair = {
        "id": chair_id,
        "thumbnail": f'/images/chair/{image_hash}.png',
        "name": "".join([
            fake.word(ext_word_list=CHAIR_NAME_PREFIX_LIST),
            fake.word(ext_word_list=CHAIR_PROPERTY_LIST),
            fake.word(ext_word_list=CHAIR_NAME_LIST)
        ]),
        "price": random.randint(1000, 20000),
        "height": random.randint(CHAIR_MIN_CENTIMETER, CHAIR_MAX_CENTIMETER),
        "width": random.randint(CHAIR_MIN_CENTIMETER, CHAIR_MAX_CENTIMETER),
        "depth": random.randint(CHAIR_MIN_CENTIMETER, CHAIR_MAX_CENTIMETER),
        "color": fake.word(ext_word_list=CHAIR_COLOR_LIST),
        "view_count": random.randint(MIN_VIEW_COUNT, MAX_VIEW_COUNT),
        "stock": random.randint(1, 10),
        "description": random.choice(desc_lines).strip(),
        "features": ",".join(fake.words(nb=features_length, ext_word_list=CHAIR_FEATURE_LIST, unique=True)),
        "kind": fake.word(ext_word_list=CHAIR_KIND_LIST)
    }

    return dict(chair, **wrap)


if __name__ == "__main__":
    for i, random_hash in enumerate(CHAIR_IMAGE_HASH_LIST):
        image_data_list = [read_src_file_data(
            image) for image in glob.glob(os.path.join(CHAIR_IMAGE_ORIGIN_DIR, "*.png"))]
        with open(os.path.join(CHAIR_IMAGE_PUBLIC_DIR, f"{random_hash}.png"), mode='wb') as image_file:
            image_file.write(
                image_data_list[i % len(image_data_list)] + random_hash.encode('utf-8'))

    with open(DESCRIPTION_LINES_FILE, mode='r', encoding='utf-8') as description_lines:
        desc_lines = description_lines.readlines()

    with open(OUTPUT_CSV_FILE, mode='w', encoding='utf-8') as csvfile, open(OUTPUT_TXT_FILE, mode='w', encoding='utf-8') as txtfile:
        CHAIRS_FOR_VERIFY = [
            # 購入された際に在庫が減ることを検証するためのデータ
            generate_chair_dummy_data(1, {
                "features": CHAIR_FEATURE_FOR_VERIFY,
                "stock": 1,
                "view_count": MIN_VIEW_COUNT
            }),
            # 2回閲覧された後の検索で、順番が前に行くことを検証するためのデータ (2位 → 1位)
            generate_chair_dummy_data(2, {
                "features": CHAIR_FEATURE_FOR_VERIFY,
                "view_count": (MAX_VIEW_COUNT + MIN_VIEW_COUNT) // 2
            }),
            # 2回閲覧された後の検索で、順番が前に行くことを検証するためのデータ (1位 → 2位)
            generate_chair_dummy_data(3, {
                "features": CHAIR_FEATURE_FOR_VERIFY,
                "view_count": (MAX_VIEW_COUNT + MIN_VIEW_COUNT) // 2 + 1
            })
        ]

        chairs = CHAIRS_FOR_VERIFY + \
            [generate_chair_dummy_data(len(CHAIRS_FOR_VERIFY) + i + 1)
             for i in range(RECORD_COUNT)]

        csvfile.write(
            CSV_COLUMNS_ORDER + "\n".join([dump_chair_to_csv_str(chair) for chair in chairs]))

        txtfile.write("\n".join([dump_chair_to_json_str(chair)
                                 for chair in chairs]) + "\n")
