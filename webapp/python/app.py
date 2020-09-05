from os import getenv
import subprocess
import flask
from mysql.connector.pooling import MySQLConnectionPool
import humps

LIMIT = 20
NAZOTTE_LIMIT = 50

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
    raise NotImplementedError()  # TODO


@app.route("/api/chair/search/condition", methods=["GET"])
def get_chair_search_condition():
    raise NotImplementedError()  # TODO


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
    return {}


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
