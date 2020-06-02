# coding:utf-8
import random
import time
import string
import os
from faker import Faker
fake = Faker("ja_JP")
Faker.seed(19700101)
random.seed(19700101)
base_dir = os.path.dirname(__file__)
DESCRIPTION_LINES_FILE = os.path.join(base_dir, "description_chair.txt")
OUTPUT_FILE = os.path.join(base_dir, "db/2_DummyChairData.sql")
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

def generate_chair_dummy_data():
    thumbnail = "/images/chair/{}.jpg".format(fake.sha256(raw_output=False))
    name= "".join([
        fake.word(ext_word_list=CHAIR_NAME_PREFIX_LIST),
        fake.word(ext_word_list=CHAIR_PROPERTY_LIST),
        fake.word(ext_word_list=CHAIR_NAME_LIST)
    ])
    price = random.randint(1000, 20000)
    height = random.randint(CHAIR_MIN_CENTIMETER, CHAIR_MAX_CENTIMETER)
    width = random.randint(CHAIR_MIN_CENTIMETER, CHAIR_MAX_CENTIMETER)
    depth = random.randint(CHAIR_MIN_CENTIMETER, CHAIR_MAX_CENTIMETER)
    color = fake.word(ext_word_list=CHAIR_COLOR_LIST)
    view_count = random.randint(3000, 1000000)
    stock = random.randint(1, 10)
    description = random.choice(desc_lines)
    features_length = random.randint(0, len(CHAIR_FEATURE_LIST) - 1)
    features = ",".join(fake.words(nb=features_length, ext_word_list=CHAIR_FEATURE_LIST, unique=True))
    kind = fake.word(ext_word_list=CHAIR_KIND_LIST)
    return f"('{thumbnail}', '{name}', '{price}', '{height}', '{width}', '{depth}', '{view_count}', '{stock}', '{color}', '{description}', '{features}', '{kind}')"

if __name__ == "__main__":
    with open(DESCRIPTION_LINES_FILE, mode="r") as description_lines:
        desc_lines = description_lines.readlines()

    with open(OUTPUT_FILE, mode="w") as sqlfile:
        sqlfile.write(sqlCommands)
        if RECORD_COUNT % BULK_INSERT_COUNT != 0:
            raise Exception("The results of RECORD_COUNT and BULK_INSERT_COUNT need to be a divisible number. RECORD_COUNT = {}, BULK_INSERT_COUNT = {}".format(RECORD_COUNT, BULK_INSERT_COUNT))

        for _ in range(RECORD_COUNT // BULK_INSERT_COUNT):
            bulk_list = [generate_chair_dummy_data() for i in range(BULK_INSERT_COUNT)]
            sqlCommand = f"""
            insert into chair
                (thumbnail, name, price, height, width, depth, view_count, stock, color, description, features, kind)
                values {", ".join(bulk_list)};
            """
            sqlfile.write(sqlCommand)
