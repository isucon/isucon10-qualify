from os import getenv
import flask
from mysql.connector.pooling import MySQLConnectionPool

app = flask.Flask(__name__)

cnxpool = MySQLConnectionPool(
    user=getenv("MYSQL_USER"),
    password=getenv("MYSQL_PASSWORD"),
    host=getenv("MYSQL_HOST"),
    database=getenv("MYSQL_DATABASE"),
)


def select_query(query, dictionary=True):
    cnx = cnxpool.get_connection()
    try:
        cur = cnx.cursor(dictionary=dictionary)
        cur.execute(query)
        return cur.fetchall()
    finally:
        cnx.close()


if __name__ == "__main__":
    app.run(port=getenv("SERVER_PORT", 1323), debug=True, threaded=True)
