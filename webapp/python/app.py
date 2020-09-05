from os import getenv
import json
import subprocess
import flask
from werkzeug.exceptions import BadRequest, NotFound
from mysql.connector.pooling import MySQLConnectionPool
from humps import camelize

LIMIT = 20
NAZOTTE_LIMIT = 50

chair_search_condition = json.load(open("../fixture/chair_condition.json", "r"))
estate_search_condition = json.load(open("../fixture/estate_condition.json", "r"))

app = flask.Flask(__name__)

cnxpool = MySQLConnectionPool(
    user=getenv("MYSQL_USER"),
    password=getenv("MYSQL_PASSWORD"),
    host=getenv("MYSQL_HOST"),
    database=getenv("MYSQL_DATABASE"),
)


def select_all(query, *args, dictionary=True):
    cnx = cnxpool.get_connection()
    try:
        cur = cnx.cursor(dictionary=dictionary)
        cur.execute(query, *args)
        return cur.fetchall()
    finally:
        cnx.close()


def select_row(*args, **kwargs):
    rows = select_all(*args, **kwargs)
    return rows[0] if len(rows) > 0 else None


@app.route("/initialize", methods=["POST"])
def post_initialize():
    subprocess.call("../mysql/db/init.sh")
    return {"language": "python"}


@app.route("/api/estate/low_priced", methods=["GET"])
def get_estate_low_priced():
    rows = select_all("SELECT * FROM estate ORDER BY rent ASC, id ASC LIMIT %s", (LIMIT,))
    return {"estates": camelize(rows)}


@app.route("/api/chair/low_priced", methods=["GET"])
def get_chair_low_priced():
    rows = select_all("SELECT * FROM chair WHERE stock > 0 ORDER BY price ASC, id ASC LIMIT %s", (LIMIT,))
    return {"chairs": camelize(rows)}


@app.route("/api/chair/search", methods=["GET"])
def get_chair_search():
    args = flask.request.args

    conditions = []
    params = []

    if args.get("priceRangeId"):
        for _range in chair_search_condition["price"]["ranges"]:
            if _range["id"] == int(args.get("priceRangeId")):
                price = _range
                break
        else:
            raise BadRequest("priceRangeID invalid")
        if price["min"] != -1:
            conditions.append("price >= %s")
            params.append(price["min"])
        if price["min"] != -1:
            conditions.append("price < %s")
            params.append(price["max"])

    if args.get("heightRangeId"):
        for _range in chair_search_condition["height"]["ranges"]:
            if _range["id"] == int(args.get("heightRangeId")):
                height = _range
                break
        else:
            raise BadRequest("heightRangeId invalid")
        if height["min"] != -1:
            conditions.append("height >= %s")
            params.append(height["min"])
        if height["min"] != -1:
            conditions.append("height < %s")
            params.append(height["max"])

    if args.get("widthRangeId"):
        for _range in chair_search_condition["width"]["ranges"]:
            if _range["id"] == int(args.get("widthRangeId")):
                width = _range
                break
        else:
            raise BadRequest("widthRangeId invalid")
        if width["min"] != -1:
            conditions.append("width >= %s")
            params.append(width["min"])
        if width["min"] != -1:
            conditions.append("width < %s")
            params.append(width["max"])

    if args.get("depthRangeId"):
        for _range in chair_search_condition["depth"]["ranges"]:
            if _range["id"] == int(args.get("depthRangeId")):
                depth = _range
                break
        else:
            raise BadRequest("depthRangeId invalid")
        if depth["min"] != -1:
            conditions.append("depth >= %s")
            params.append(depth["min"])
        if depth["min"] != -1:
            conditions.append("depth < %s")
            params.append(depth["max"])

    if args.get("kind"):
        conditions.append("kind = %s")
        params.append(args.get("kind"))

    if args.get("color"):
        conditions.append("color = %s")
        params.append(args.get("color"))

    if args.get("features"):
        for feature_confition in args.get("features").split(","):
            conditions.append("features LIKE CONCAT('%', %s, '%')")
            params.append(feature_confition)

    if len(conditions) == 0:
        raise BadRequest("Search condition not found")

    conditions.append("stock > 0")

    try:
        page = int(args.get("page"))
    except (TypeError, ValueError):
        raise BadRequest("Invalid format page parameter")

    try:
        per_page = int(args.get("perPage"))
    except (TypeError, ValueError):
        raise BadRequest("Invalid format perPage parameter")

    search_condition = " AND ".join(conditions)

    query = f"SELECT COUNT(*) as count FROM chair WHERE {search_condition}"
    count = select_row(query, params)["count"]

    query = f"SELECT * FROM chair WHERE {search_condition} ORDER BY popularity DESC, id ASC LIMIT %s OFFSET %s"
    chairs = select_all(query, params + [per_page, per_page * page])

    return {"count": count, "chairs": camelize(chairs)}


@app.route("/api/chair/search/condition", methods=["GET"])
def get_chair_search_condition():
    return chair_search_condition


@app.route("/api/chair/<int:chair_id>", methods=["GET"])
def get_chair(chair_id):
    rows = select_all("SELECT * FROM chair WHERE id = %s", (chair_id,))
    if len(rows) == 0 or rows[0]["stock"] <= 0:
        raise NotFound()
    return camelize(rows[0])


@app.route("/api/chair/buy/<int:chair_id>", methods=["POST"])
def post_chair_buy(chair_id):
    cnx = cnxpool.get_connection()
    try:
        cnx.start_transaction()
        cur = cnx.cursor(dictionary=True)
        cur.execute("SELECT * FROM chair WHERE id = %s AND stock > 0 FOR UPDATE", (chair_id,))
        chair = cur.fetchone()
        if chair is None:
            raise NotFound()
        cur.execute("UPDATE chair SET stock = %s WHERE id = %s", (chair["stock"] - 1, chair_id))
        cnx.commit()
        return {"ok": True}
    except Exception as e:
        cnx.rollback()
        raise e
    finally:
        cnx.close()


@app.route("/api/estate/search", methods=["GET"])
def get_estate_search():
    args = flask.request.args

    conditions = []
    params = []

    if args.get("doorHeightRangeId"):
        for _range in estate_search_condition["doorHeight"]["ranges"]:
            if _range["id"] == int(args.get("doorHeightRangeId")):
                door_height = _range
                break
        else:
            raise BadRequest("doorHeightRangeId invalid")
        if door_height["min"] != -1:
            conditions.append("door_height >= %s")
            params.append(door_height["min"])
        if door_height["min"] != -1:
            conditions.append("door_height < %s")
            params.append(door_height["max"])

    if args.get("doorWidthRangeId"):
        for _range in estate_search_condition["doorWidth"]["ranges"]:
            if _range["id"] == int(args.get("doorWidthRangeId")):
                door_width = _range
                break
        else:
            raise BadRequest("doorWidthRangeId invalid")
        if door_width["min"] != -1:
            conditions.append("door_width >= %s")
            params.append(door_width["min"])
        if door_width["min"] != -1:
            conditions.append("door_width < %s")
            params.append(door_width["max"])

    if args.get("rentRangeId"):
        for _range in estate_search_condition["rent"]["ranges"]:
            if _range["id"] == int(args.get("rentRangeId")):
                rent = _range
                break
        else:
            raise BadRequest("rentRangeId invalid")
        if rent["min"] != -1:
            conditions.append("rent >= %s")
            params.append(rent["min"])
        if rent["min"] != -1:
            conditions.append("rent < %s")
            params.append(rent["max"])

    if args.get("features"):
        for feature_confition in args.get("features").split(","):
            conditions.append("features LIKE CONCAT('%', %s, '%')")
            params.append(feature_confition)

    if len(conditions) == 0:
        raise BadRequest("Search condition not found")

    try:
        page = int(args.get("page"))
    except (TypeError, ValueError):
        raise BadRequest("Invalid format page parameter")

    try:
        per_page = int(args.get("perPage"))
    except (TypeError, ValueError):
        raise BadRequest("Invalid format perPage parameter")

    search_condition = " AND ".join(conditions)

    query = f"SELECT COUNT(*) as count FROM estate WHERE {search_condition}"
    count = select_row(query, params)["count"]

    query = f"SELECT * FROM estate WHERE {search_condition} ORDER BY popularity DESC, id ASC LIMIT %s OFFSET %s"
    chairs = select_all(query, params + [per_page, per_page * page])

    return {"count": count, "chairs": camelize(chairs)}


@app.route("/api/estate/search/condition", methods=["GET"])
def get_estate_search_condition():
    return estate_search_condition


@app.route("/api/estate/req_doc/<int:estate_id>", methods=["POST"])
def post_estate_req_doc(estate_id):
    raise NotImplementedError()  # TODO


@app.route("/api/estate/nazotte", methods=["POST"])
def post_estate_nazotte():
    raise NotImplementedError()  # TODO


@app.route("/api/estate/<int:estate_id>", methods=["GET"])
def get_estate(estate_id):
    raise NotImplementedError()  # TODO


# low_priced の方のエンドポイントはフロントエンドのバグ (URL 間違い) に対応するためのエイリアス (TODO フロントエンドが直ったら消す)
@app.route("/api/estate/low_priced/<int:chair_id>", methods=["GET"])
@app.route("/api/recommended_estate/<int:chair_id>", methods=["GET"])
def get_recommended_estate(chair_id):
    chair = select_row("SELECT * FROM chair WHERE id = %s", (chair_id,))
    if chair is None:
        raise BadRequest(f"Invalid format searchRecommendedEstateWithChair id : {chair_id}")
    w, h, d = chair["width"], chair["height"], chair["depth"]
    query = (
        "SELECT * FROM estate"
        " WHERE (door_width >= %s AND door_height >= %s)"
        "    OR (door_width >= %s AND door_height >= %s)"
        "    OR (door_width >= %s AND door_height >= %s)"
        "    OR (door_width >= %s AND door_height >= %s)"
        "    OR (door_width >= %s AND door_height >= %s)"
        "    OR (door_width >= %s AND door_height >= %s)"
        " ORDER BY popularity DESC, id ASC"
        " LIMIT %s"
    )
    estates = select_all(query, (w, h, w, d, h, w, h, d, d, w, d, h, LIMIT))
    return {"estates": camelize(estates)}


@app.route("/api/chair", methods=["POST"])
def post_chair():
    raise NotImplementedError()  # TODO


@app.route("/api/estate", methods=["POST"])
def post_estate():
    raise NotImplementedError()  # TODO


if __name__ == "__main__":
    app.run(host="0.0.0.0", port=getenv("SERVER_PORT", 1323), debug=True, threaded=True)
