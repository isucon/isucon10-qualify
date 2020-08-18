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

DESCRIPTION_LINES_FILE = "./description.txt"
OUTPUT_SQL_FILE = "./result/2_DummyChairData.sql"
OUTPUT_TXT_FILE = "./result/chair_json.txt"
OUTPUT_FIXTURE_FILE = "./result/chair_condition.json"
CHAIR_IMAGE_ORIGIN_DIR = "./origin/chair"
CHAIR_IMAGE_PUBLIC_DIR = "../webapp/frontend/public/images/chair"
CHAIR_DUMMY_IMAGE_NUM = 1000
RECORD_COUNT = (10 ** 4) * 3
BULK_INSERT_COUNT = 500
CHAIR_MIN_CENTIMETER = 30
CHAIR_MAX_CENTIMETER = 200
MIN_VIEW_COUNT = 3000
MAX_VIEW_COUNT = 1000000
HEIGHT_RANGE_SEPARATORS = [80, 110, 150]
WIDTH_RANGE_SEPARATORS = [80, 110, 150]
DEPTH_RANGE_SEPARATORS = [80, 110, 150]
PRICE_RANGE_SEPARATORS = [3000, 6000, 9000, 12000, 15000]

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
    "ヘッドレスト付き",
    "肘掛け付き",
    "キャスター付き",
    "アーム高さ調節可能",
    "リクライニング可能",
    "高さ調節可能",
    "通気性抜群",
    "メタルフレーム",
    "低反発",
    "木製",
    "背もたれつき",
    "回転可能",
    "レザー製",
    "昇降式",
    "デザイナーズ",
    "金属製",
    "プラスチック製",
    "法事用",
    "和風",
    "中華風",
    "西洋風",
    "イタリア製",
    "国産",
    "背もたれなし",
    "ラテン風",
    "布貼地",
    "スチール製",
    "メッシュ貼地",
    "オフィス用",
    "料理店用",
    "自宅用",
    "キャンプ用",
    "クッション性抜群",
    "モーター付き",
    "ベッド一体型",
    "ディスプレイ配置可能",
    "ミニ机付き",
    "スピーカー付属",
    "中国製",
    "アンティーク",
    "折りたたみ可能",
    "重さ500g以内",
    "24回払い無金利",
    "現代的デザイン",
    "近代的なデザイン",
    "ルネサンス的なデザイン",
    "アームなし",
    "オーダーメイド可能",
    "ポリカーボネート製",
]

CHAIR_FEATURE_FOR_VERIFY = "フットレスト付き"

CHAIR_KIND_LIST = [
    "ゲーミングチェア",
    "座椅子",
    "エルゴノミクス",
    "ハンモック"
]

CHAIR_IMAGE_HASH_LIST = [fake.sha256(
    raw_output=False) for _ in range(CHAIR_DUMMY_IMAGE_NUM)]


def generate_ranges_from_SEPARATORS(SEPARATORSs):
    before = -1
    ranges = []

    for i, SEPARATORS in enumerate(SEPARATORSs + [-1]):
        ranges.append({
            "id": i,
            "min": before,
            "max": SEPARATORS
        })
        before = SEPARATORS

    return ranges


def read_src_file_data(file_path):
    with open(file_path, mode='rb') as img:
        return img.read()


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
        "popularity": chair["popularity"],
        "stock": chair["stock"],
        "description": chair["description"],
        "features": chair["features"],
        "kind": chair["kind"]
    }, ensure_ascii=False)


def generate_chair_dummy_data(chair_id, wrap={}):
    features_length = random.randint(0, min(3, len(CHAIR_FEATURE_LIST)))
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
        "popularity": random.randint(MIN_VIEW_COUNT, MAX_VIEW_COUNT),
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

    with open(OUTPUT_SQL_FILE, mode='w', encoding='utf-8') as sqlfile, open(OUTPUT_TXT_FILE, mode='w', encoding='utf-8') as txtfile:
        if RECORD_COUNT % BULK_INSERT_COUNT != 0:
            raise Exception("The results of RECORD_COUNT and BULK_INSERT_COUNT need to be a divisible number. RECORD_COUNT = {}, BULK_INSERT_COUNT = {}".format(
                RECORD_COUNT, BULK_INSERT_COUNT))

        chair_id = 1

        CHAIRS_FOR_VERIFY = [
            # 購入された際に在庫が減ることを検証するためのデータ
            generate_chair_dummy_data(1, {
                "features": CHAIR_FEATURE_FOR_VERIFY,
                "stock": 1,
                "popularity": MIN_VIEW_COUNT
            }),
            # 2回閲覧された後の検索で、順番が前に行くことを検証するためのデータ (2位 → 1位)
            generate_chair_dummy_data(2, {
                "features": CHAIR_FEATURE_FOR_VERIFY,
                "popularity": (MAX_VIEW_COUNT + MIN_VIEW_COUNT) // 2
            }),
            # 2回閲覧された後の検索で、順番が前に行くことを検証するためのデータ (1位 → 2位)
            generate_chair_dummy_data(3, {
                "features": CHAIR_FEATURE_FOR_VERIFY,
                "popularity": (MAX_VIEW_COUNT + MIN_VIEW_COUNT) // 2 + 1
            })
        ]

        sqlCommand = f"""INSERT INTO isuumo.chair (id, thumbnail, name, price, height, width, depth, popularity, stock, color, description, features, kind) VALUES {", ".join(map(lambda chair: f"('{chair['id']}', '{chair['thumbnail']}', '{chair['name']}', '{chair['price']}', '{chair['height']}', '{chair['width']}', '{chair['depth']}', '{chair['popularity']}', '{chair['stock']}', '{chair['color']}', '{chair['description']}', '{chair['features']}', '{chair['kind']}')", CHAIRS_FOR_VERIFY))};"""
        sqlfile.write(sqlCommand)
        txtfile.write(
            "\n".join([dump_chair_to_json_str(chair) for chair in CHAIRS_FOR_VERIFY]) + "\n")

        chair_id += len(CHAIRS_FOR_VERIFY)

        for _ in range(RECORD_COUNT // BULK_INSERT_COUNT):
            bulk_list = [generate_chair_dummy_data(
                chair_id + i) for i in range(BULK_INSERT_COUNT)]
            chair_id += BULK_INSERT_COUNT
            sqlCommand = f"""INSERT INTO isuumo.chair (id, thumbnail, name, price, height, width, depth, popularity, stock, color, description, features, kind) VALUES {", ".join(map(lambda chair: f"('{chair['id']}', '{chair['thumbnail']}', '{chair['name']}', '{chair['price']}', '{chair['height']}', '{chair['width']}', '{chair['depth']}', '{chair['popularity']}', '{chair['stock']}', '{chair['color']}', '{chair['description']}', '{chair['features']}', '{chair['kind']}')", bulk_list))};"""
            sqlfile.write(sqlCommand)

            txtfile.write(
                "\n".join([dump_chair_to_json_str(chair) for chair in bulk_list]) + "\n")

    with open(OUTPUT_FIXTURE_FILE, mode='w', encoding='utf-8') as fixture_file:
        fixture_file.write(json.dumps({
            "height": {
                "prefix": "",
                "suffix": "cm",
                "ranges": generate_ranges_from_SEPARATORS(HEIGHT_RANGE_SEPARATORS)
            },
            "width": {
                "prefix": "",
                "suffix": "cm",
                "ranges": generate_ranges_from_SEPARATORS(WIDTH_RANGE_SEPARATORS)
            },
            "depth": {
                "prefix": "",
                "suffix": "cm",
                "ranges": generate_ranges_from_SEPARATORS(DEPTH_RANGE_SEPARATORS)
            },
            "price": {
                "prefix": "",
                "suffix": "円",
                "ranges": generate_ranges_from_SEPARATORS(PRICE_RANGE_SEPARATORS)
            },
            "color": {
                "list": CHAIR_COLOR_LIST
            },
            "feature": {
                "list": CHAIR_FEATURE_LIST + [CHAIR_FEATURE_FOR_VERIFY]
            },
            "kind": {
                "list": CHAIR_KIND_LIST
            }
        }, ensure_ascii=False, indent=2))
