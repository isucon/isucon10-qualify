from os import getenv
import json
import subprocess
import flask
from mysql.connector.pooling import MySQLConnectionPool
import humps

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


def select_query(query, *args, dictionary=True):
    cnx = cnxpool.get_connection()
    try:
        cur = cnx.cursor(dictionary=dictionary)
        cur.execute(query, *args)
        return cur.fetchall()
    finally:
        cnx.close()


@app.route("/initialize", methods=["POST"])
def post_initialize():
    subprocess.call("../mysql/db/init.sh")
    return {"language": "python"}


@app.route("/api/estate/low_priced", methods=["GET"])
def get_estate_low_priced():
    rows = select_query("SELECT * FROM estate ORDER BY rent ASC, id ASC LIMIT %s", (LIMIT,))
    return {"estates": humps.camelize(rows)}


@app.route("/api/chair/low_priced", methods=["GET"])
def get_chair_low_priced():
    rows = select_query("SELECT * FROM chair WHERE stock > 0 ORDER BY price ASC, id ASC LIMIT %s", (LIMIT,))
    return {"chairs": humps.camelize(rows)}


@app.route("/api/chair/search", methods=["GET"])
def get_chair_search():
    args = flask.request.args

    search_queries = []
    query_params = []

    if args.get("priceRangeId"):
        chair_price = chair_search_condition["price"]["ranges"][args.get("priceRangeId")]
        if chair_price is None:
            return abort(400, "priceRangeID invalid")
        if chair_price["min"] != -1:
            search_queries.append("price >= %s")
            query_params.append(chair_price["min"])
        if chair_price["min"] != -1:
            search_queries.append("price < %s")
            query_params.append(chair_price["max"])

    if args.get("heightRangeId"):
        chair_height = chair_search_condition["height"][args.get("heightRangeId")]
        if chair_height is None:
            return abort(400, "heightRangeId invalid")
        if chair_height["min"] != -1:
            search_queries.append("height >= %s")
            query_params.append(chair_height["min"])
        if chair_height["min"] != -1:
            search_queries.append("height < %s")
            query_params.append(chair_height["max"])

    if args.get("widthRangeId"):
        chair_width = chair_search_condition["width"][args.get("widthRangeId")]
        if chair_width is None:
            return abort(400, "widthRangeId invalid")
        if chair_width["min"] != -1:
            search_queries.append("width >= %s")
            query_params.append(chair_width["min"])
        if chair_width["min"] != -1:
            search_queries.append("width < %s")
            query_params.append(chair_width["max"])

    if args.get("depthRangeId"):
        chair_width = chair_search_condition["depth"][args.get("depthRangeId")]
        if chair_depth is None:
            return abort(400, "depthRangeId invalid")
        if chair_depth["min"] != -1:
            search_queries.append("depth >= %s")
            query_params.append(chair_depth["min"])
        if chair_depth["min"] != -1:
            search_queries.append("depth < %s")
            query_params.append(chair_depth["max"])

    if args.get("kind"):
        search_queries.append("kind = %s")
        query_params.append(args.get("kind"))

    if args.get("color"):
        search_queries.append("color = %s")
        query_params.append(args.get("color"))

    if args.get("features"):
        for feature_confition in args.get("features").split(","):
            search_queries.append("features LIKE CONCAT('%', %s, '%')")
            query_params.append(feature_confition)

    if len(search_queries) == 0:
        return abort(400, "Search condition not found")

    search_queries.append("stock > 0")

    try:
        page = int(args.get("page"))
    except (TypeError, ValueError):
        return abort(400, "Invalid format page parameter")

    try:
        per_page = int(args.get("perPage"))
    except (TypeError, ValueError):
        return abort(400, "Invalid format perPage parameter")

    search_condition = " AND ".join(search_queries)

    query = f"SELECT COUNT(*) as count FROM chair WHERE {search_condition}"
    count = select_query(query, query_params)[0]["count"]

    query = f"SELECT * FROM chair WHERE {search_condition} ORDER BY popularity DESC, id ASC LIMIT %s OFFSET %s"
    chairs = select_query(query, query_params + [per_page, per_page * page])

    return {"count": count, "chairs": humps.camelize(chairs)}


@app.route("/api/chair/search/condition", methods=["GET"])
def get_chair_search_condition():
    return chair_search_condition


@app.route("/api/chair/<int:chair_id>", methods=["GET"])
def get_chair(chair_id):
    raise NotImplementedError()  # TODO


@app.route("/api/chair/buy/<int:chair_id>", methods=["POST"])
def post_chair_buy(chair_id):
    raise NotImplementedError()  # TODO


@app.route("/api/estate/search", methods=["GET"])
def get_estate_search():
    raise NotImplementedError()  # TODO


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


@app.route("/api/recommended_estate/<int:estate_id>", methods=["GET"])
def get_recommended_estate(estate_id):
    raise NotImplementedError()  # TODO


@app.route("/api/chair", methods=["POST"])
def post_chair():
    raise NotImplementedError()  # TODO


@app.route("/api/estate", methods=["POST"])
def post_estate():
    raise NotImplementedError()  # TODO


if __name__ == "__main__":
    app.run(host="0.0.0.0", port=getenv("SERVER_PORT", 1323), debug=True, threaded=True)
