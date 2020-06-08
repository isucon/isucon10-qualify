# coding:utf-8
import random
import time
import string
import json
import os
from faker import Faker
fake = Faker("ja_JP")
Faker.seed(19700101)
random.seed(19700101)

DESCRIPTION_LINES_FILE = "./description_chair.txt"
OUTPUT_SQL_FILE = "./result/2_DummyChairData.sql"
OUTPUT_TXT_FILE = "./result/chair_json.txt"
CHAIR_IMAGE_ORIGIN_DIR = "./origin/chair"
CHAIR_IMAGE_PUBLIC_DIR = "../webapp/frontend/public/images/chair"
RECORD_COUNT = 10 ** 4
BULK_INSERT_COUNT = 500
CHAIR_MIN_CENTIMETER = 30
CHAIR_MAX_CENTIMETER = 200
sqlCommands = "use isuumo;\n"

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

CHAIR_NAME_LIST= [
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
    "高さ調節可",
    "フットレスト"
]

CHAIR_KIND_LIST = [
    "ゲーミングチェア",
    "座椅子",
    "エルゴノミクス",
    "ハンモック"
]

CHAIR_IMAGE_LIST = [
    os.path.join(CHAIR_IMAGE_ORIGIN_DIR, "1501E5C34A2B8EE645480ED1CC6442CD5929FE7616E20513574628096163DF0C.jpg"),
    os.path.join(CHAIR_IMAGE_ORIGIN_DIR, "3E880A828B1DBFACB42209724583B56EF28466E45E2BF3704475EA02B19BDBFC.jpg"),
    os.path.join(CHAIR_IMAGE_ORIGIN_DIR, "9120C2E3CAF5CD376C1B14899C2FD31438A839D1F6B6F8A52091392E0B9168FC.jpg")
]

CHAIR_IMAGE_DATA = [read_src_file_data(image) for image in CHAIR_IMAGE_LIST]

def read_src_file_data(file_path):
    with open(file_path, mode='rb') as img:
        return img.read()

def generate_chair_dummy_data(chair_id):
    features_length = random.randint(0, len(CHAIR_FEATURE_LIST) - 1)

    new_chair_image_hash = fake.sha256(raw_output=False)
    new_chair_image_path = os.path.join(CHAIR_IMAGE_PUBLIC_DIR, "{}.jpg".format(new_chair_image_hash))
    src_image_data = random.choice(CHAIR_IMAGE_DATA)

    with open(new_chair_image_path, mode='wb') as dst_image_data:
        dst_image_data.write(src_image_data + new_chair_image_hash.encode('utf-8'))

    return {
        "id": chair_id,
        "thumbnail": "/images/chair/{}.jpg".format(new_chair_image_path),
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
        "view_count": random.randint(3000, 1000000),
        "stock": random.randint(1, 10),
        "description": random.choice(desc_lines).strip(),
        "features": ",".join(fake.words(nb=features_length, ext_word_list=CHAIR_FEATURE_LIST, unique=True)),
        "kind": fake.word(ext_word_list=CHAIR_KIND_LIST)
    }

if __name__ == "__main__":
    for image in CHAIR_IMAGE_LIST:
        filename, _ = os.path.splitext(os.path.basename(image))
        with open(os.path.join(CHAIR_IMAGE_PUBLIC_DIR, "{}.jpg".format(filename)), mode='wb') as dst_image_data:
            dst_image_data.write(read_src_file_data(image) + filename.encode('utf-8'))

    with open(DESCRIPTION_LINES_FILE, mode="r") as description_lines:
        desc_lines = description_lines.readlines()

    with open(OUTPUT_SQL_FILE, mode='w') as sqlfile, open(OUTPUT_TXT_FILE, mode='w') as txtfile:
        sqlfile.write(sqlCommands)
        if RECORD_COUNT % BULK_INSERT_COUNT != 0:
            raise Exception("The results of RECORD_COUNT and BULK_INSERT_COUNT need to be a divisible number. RECORD_COUNT = {}, BULK_INSERT_COUNT = {}".format(RECORD_COUNT, BULK_INSERT_COUNT))

        chair_id = 1
        for _ in range(RECORD_COUNT // BULK_INSERT_COUNT):
            bulk_list = [generate_chair_dummy_data(chair_id + i) for i in range(BULK_INSERT_COUNT)]
            chair_id += BULK_INSERT_COUNT
            sqlCommand = f"""INSERT INTO chair (id, thumbnail, name, price, height, width, depth, view_count, stock, color, description, features, kind) VALUES {", ".join(map(lambda chair: f"('{chair['id']}', '{chair['thumbnail']}', '{chair['name']}', '{chair['price']}', '{chair['height']}', '{chair['width']}', '{chair['depth']}', '{chair['view_count']}', '{chair['stock']}', '{chair['color']}', '{chair['description']}', '{chair['features']}', '{chair['kind']}')", bulk_list))};"""
            sqlfile.write(sqlCommand)

            for chair in bulk_list:
                json_string = json.dumps({
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
                txtfile.write(json_string + "\n")
