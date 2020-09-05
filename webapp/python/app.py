from os import getenv
import flask

app = flask.Flask(__name__)

if __name__ == "__main__":
    app.run(port=getenv("SERVER_PORT", 1323), debug=True, threaded=True)
