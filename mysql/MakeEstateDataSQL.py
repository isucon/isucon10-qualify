# coding:utf-8
import random
import time
import string
from faker import Faker
fake = Faker('ja_JP')
Faker.seed(19700101)
random.seed(19700101)

DESCRIPTION_LINES_FILE = "./description.txt"
OUTPUT_FILE = "./db/1_DummyData.sql"
RECORD_COUNT = 10 ** 4
BULK_INSERT_COUNT = 500
DOOR_MIN_CENTIMETER = 30
DOOR_MAX_CENTIMETER = 200
sqlCommands = ""
sqlCommands += "create database isuumo;\nuse isuumo;\n"

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

estate_table = """
create table estate (
    thumbnails varchar(256),
    name varchar(64),
    coordinate geometry not null,
    address varchar(128),
    rent integer,
    door_height integer,
    door_width integer,
    view_count integer default 0,
    description text,
    feature varchar(256)
)ENGINE=MyISAM;
"""

def generate_estate_dummy_data():
    thumbnails = ','.join(['{}.jpg'.format(fake.sha256(raw_output=False)) for i in range(3)])
    name= fake.word(ext_word_list=BUILDING_NAME_LIST).format(name=fake.last_name())
    #designer_id random int
    latitude, longitude = fake.local_latlng(country_code='JP', coords_only=True)
    address = fake.address()
    rent = fake.pyint(min_value=5, max_value=20) * 10000
    door_height = random.randint(DOOR_MIN_CENTIMETER, DOOR_MAX_CENTIMETER)
    door_width = random.randint(DOOR_MIN_CENTIMETER, DOOR_MAX_CENTIMETER)
    view_count = random.randint(3000, 1000000)
    description = random.choice(desc_lines)
    feature_length = random.randint(0, len(ESTATE_FEATURE_LIST) - 1)
    feature = ','.join(fake.words(nb=feature_length, ext_word_list=ESTATE_FEATURE_LIST, unique=True))
    return f"('{thumbnails}', '{name}', ST_GeomFromText('POINT({latitude} {longitude})', 6668), '{address}', '{rent}', '{door_height}', '{door_width}', '{view_count}', '{description}', '{feature}')"

if __name__ == '__main__':
    with open(DESCRIPTION_LINES_FILE, mode='r') as description_lines:
        desc_lines = description_lines.readlines()

    with open(OUTPUT_FILE, mode='w') as sqlfile:
        sqlfile.write(sqlCommands)
        sqlfile.write(estate_table)
        if RECORD_COUNT % BULK_INSERT_COUNT != 0:
            raise Exception("The results of RECORD_COUNT and BULK_INSERT_COUNT need to be a divisible number. RECORD_COUNT = {}, BULK_INSERT_COUNT = {}".format(RECORD_COUNT, BULK_INSERT_COUNT))

        for _ in range(RECORD_COUNT//BULK_INSERT_COUNT):
            bulk_list = [generate_estate_dummy_data() for i in range(BULK_INSERT_COUNT)]
            sqlCommand = f"""
            insert into estate
                (thumbnails, name, coordinate, address, rent, door_height, door_width, view_count, description, feature)
                values {', '.join(bulk_list)};
            """
            sqlfile.write(sqlCommand)

        sqlfile.write("create spatial index coord_index on estate (coordinate);")
