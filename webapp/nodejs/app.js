const express = require("express");
const mysql = require("mysql");
const path = require("path");
const cp = require("child_process");
const util = require("util");
const camelcaseKeys = require("camelcase-keys");
const promisify = util.promisify;
const exec = promisify(cp.exec);
const chairSearchCondition = require("../fixture/chair_condition.json");
const estateSearchCondition = require("../fixture/estate_condition.json");

const PORT = process.env.PORT ?? 1323;
const LIMIT = 20;
const NAZOTTE_LIMIT = 50;
const dbinfo = {
  host: process.env.MYSQL_HOST ?? "127.0.0.1",
  port: process.env.MYSQL_PORT ?? 3306,
  user: process.env.MYSQL_USER ?? "isucon",
  password: process.env.MYSQL_PASS ?? "isucon",
  database: process.env.MYSQL_DBNAME ?? "isuumo",
};

const app = express();
const db = mysql.createPool(dbinfo);
app.set("db", db);

app.use(express.json());
app.post("/initialize", async (req, res, next) => {
  try {
    const dbdir = path.resolve("..", "mysql", "db");
    const dbfiles = ["0_Schema.sql", "1_DummyEstateData.sql", "2_DummyChairData.sql"];
    const execfiles = dbfiles.map((file) => path.join(dbdir, file));
    for (const execfile of execfiles) {
      await exec(`mysql -h ${dbinfo.host} -u ${dbinfo.user} -p${dbinfo.password} -P ${dbinfo.port} ${dbinfo.database} < ${execfile}`);
    }
    res.json({
      language: "nodejs",
    });
  } catch (e) {
    next(e);
  }
});

app.get("/api/chair/search", async (req, res, next) => {
  const searchQueries = [];
  const queryParams = [];
  const {priceRangeId, heightRangeId, widthRangeId, depthRangeId, color, features, page, perPage, } = req.query;

  if (priceRangeId != null) {
    const chairPrice = chairSearchCondition["price"].ranges[priceRangeId];
    if (chairPrice == null) {
      res.status(400).send("priceRangeID invalid");
      return;
    }

    if (chairPrice.min !== -1) {
      searchQueries.push("price >= ? ");
      queryParams.push(chairPrice.min);
    }

    if (chairPrice.max !== -1) {
      searchQueries.push("price < ? ");
      queryParams.push(chairPrice.max);
    }
  }

  if (heightRangeId != null) {
    const chairHeight = chairSearchCondition["height"].ranges[heightRangeId];
    if (chairHeight == null) {
      res.status(400).send("heightRangeId invalid");
      return;
    }

    if (chairHeight.min !== -1) {
      searchQueries.push("height >= ? ");
      queryParams.push(chairHeight.min);
    }

    if (chairHeight.max !== -1) {
      searchQueries.push("height < ? ");
      queryParams.push(chairHeight.max);
    }
  }

  if (widthRangeId != null) {
    const chairWidth = chairSearchCondition["width"].ranges[widthRangeId];
    if (chairWidth == null) {
      res.status(400).send("widthRangeId invalid");
      return;
    }

    if (chairWidth.min !== -1) {
      searchQueries.push("width >= ? ");
      queryParams.push(chairWidth.min);
    }

    if (chairWidth.max !== -1) {
      searchQueries.push("width < ? ");
      queryParams.push(chairWidth.max);
    }
  }

  if (depthRangeId != null) {
    const chairDepth = chairSearchCondition["depth"].ranges[depthRangeId];
    if (chairDepth == null) {
      res.status(400).send("depthRangeId invalid");
      return;
    }

    if (chairDepth.min !== -1) {
      searchQueries.push("depth >= ? ");
      queryParams.push(chairDepth.min);
    }

    if (chairDepth.max !== -1) {
      searchQueries.push("depth < ? ");
      queryParams.push(chairDepth.max);
    }
  }

  if (color != null) {
    searchQueries.push("color = ? ");
    queryParams.push(color);
  }

  if (features != null) {
    const featureConditions = features.split(",");
    for (const featureCondition of featureConditions) {
      searchQueries.push("features LIKE CONCAT('%', ?, '%')");
      queryParams.push(featureCondition);
    }
  }

  if (searchQueries.length === 0) {
    res.status(400).send("Search condition not found");
    return;
  }

  searchQueries.push("stock > 0");

	if (!page || page != +page) {
    res.status(400).send(`page condition invalid ${page}`);
    return;
	}

	if (!perPage || perPage != +perPage) {
    res.status(400).send("perPage condition invalid");
    return;
  }

  const pageNum = parseInt(page, 10);
  const perPageNum = parseInt(perPage, 10);
  
  const sqlprefix = "SELECT * FROM chair WHERE ";
  const searchCondition = searchQueries.join(" AND ");
  const limitOffset = " ORDER BY view_count DESC, id ASC LIMIT ? OFFSET ?";
  const countprefix = "SELECT COUNT(*) as count FROM chair WHERE ";

  const getConnection = promisify(db.getConnection.bind(db));
  const connection = await getConnection();
  const query = promisify(connection.query.bind(connection));
  try {
    const [{count}] = await query(`${countprefix}${searchCondition}`, queryParams);
    queryParams.push(perPageNum, perPageNum * pageNum);
    const chairs = await query(`${sqlprefix}${searchCondition}${limitOffset}`, queryParams);
    res.json({
      count,
      chairs: camelizeKeys(chairs),
    });
  } catch (e) {
    next(e);
  } finally {
    await connection.release();
  }
});

app.get("/api/chair/search/condition", (req, res, next) => {
  res.json(chairSearchCondition);
});

app.get("/api/chair/:id", async (req, res, next) => {
  const getConnection = promisify(db.getConnection.bind(db));
  const connection = await getConnection();
  const query = promisify(connection.query.bind(connection));
  try {
    const id = req.params.id;
    const [chair] = await query("SELECT * FROM chair WHERE id = ?", [id]);
    if (chair == null) {
      res.status(404).send("Not Found");
      return;
    }
    await connection.beginTransaction();
    await query("UPDATE chair SET view_count = ? WHERE id = ?", [chair.view_count+1, id]);
    await connection.commit();
    res.json(camelizeKeys(chair));
  } catch (e) {
    await connection.rollback();
    next(e);
  } finally {
    await connection.release();
  }
});

app.post("/api/chair/buy/:id", async (req, res, next) => {
  const getConnection = promisify(db.getConnection.bind(db));
  const connection = await getConnection();
  const query = promisify(connection.query.bind(connection));
  try {
    const id = req.params.id;
    const [chair] = await query("SELECT * FROM chair WHERE id = ? AND stock > 0", [id]);
    if (chair == null) {
      res.status(404).send("Not Found");
      return;
    }
    await connection.beginTransaction();
    await query("UPDATE chair SET stock = ? WHERE id = ?", [chair.stock-1, id]);
    await connection.commit();
    res.json({ ok: true });
  } catch (e) {
    await connection.rollback();
    next(e);
  } finally {
    await connection.release();
  }
});

app.get("/api/estate/search", async (req, res, next) => {
  const searchQueries = [];
  const queryParams = [];
  const {doorHeightRangeId, doorWidthRangeId, rentRangeId, features, page, perPage, } = req.query;

  if (doorHeightRangeId != null) {
    const doorHeight = estateSearchCondition["doorHeight"].ranges[doorHeightRangeId];
    if (doorHeight == null) {
      res.status(400).send("doorHeightRangeId invalid");
      return;
    }

    if (doorHeight.min !== -1) {
      searchQueries.push("door_height >= ? ");
      queryParams.push(doorHeight.min);
    }

    if (doorHeight.max !== -1) {
      searchQueries.push("door_height < ? ");
      queryParams.push(doorHeight.max);
    }
  }

  if (doorWidthRangeId != null) {
    const doorWidth = estateSearchCondition["doorWidth"].ranges[doorWidthRangeId];
    if (doorWidth == null) {
      res.status(400).send("doorWidthRangeId invalid");
      return;
    }

    if (doorWidth.min !== -1) {
      searchQueries.push("door_width >= ? ");
      queryParams.push(doorWidth.min);
    }

    if (doorWidth.max !== -1) {
      searchQueries.push("door_width < ? ");
      queryParams.push(doorWidth.max);
    }
  }

  if (rentRangeId != null) {
    const rent = estateSearchCondition["rent"].ranges[rentRangeId];
    if (rent == null) {
      res.status(400).send("rentRangeId invalid");
      return;
    }

    if (rent.min !== -1) {
      searchQueries.push("rent >= ? ");
      queryParams.push(rent.min);
    }

    if (rent.max !== -1) {
      searchQueries.push("rent < ? ");
      queryParams.push(rent.max);
    }
  }

  if (features != null) {
    const featureConditions = features.split(",");
    for (const featureCondition of featureConditions) {
      searchQueries.push("features LIKE CONCAT('%', ?, '%')");
      queryParams.push(featureCondition);
    }
  }

  if (searchQueries.length === 0) {
    res.status(400).send("Search condition not found");
    return;
  }

	if (!page || page != +page) {
    res.status(400).send(`page condition invalid ${page}`);
    return;
	}

	if (!perPage || perPage != +perPage) {
    res.status(400).send("perPage condition invalid");
    return;
  }

  const pageNum = parseInt(page, 10);
  const perPageNum = parseInt(perPage, 10);
  
  const sqlprefix = "SELECT * FROM estate WHERE ";
  const searchCondition = searchQueries.join(" AND ");
  const limitOffset = " ORDER BY view_count DESC, id ASC LIMIT ? OFFSET ?";
  const countprefix = "SELECT COUNT(*) as count FROM estate WHERE ";

  const getConnection = promisify(db.getConnection.bind(db));
  const connection = await getConnection();
  const query = promisify(connection.query.bind(connection));
  try {
    const [{count}] = await query(`${countprefix}${searchCondition}`, queryParams);
    queryParams.push(perPageNum, perPageNum * pageNum);
    const estates = await query(`${sqlprefix}${searchCondition}${limitOffset}`, queryParams);
    res.json({
      count,
      estates: camelcaseKeys(estates),
    });
  } catch (e) {
    next(e);
  } finally {
    await connection.release();
  }
});

app.get("/api/estate/search/condition", (req, res, next) => {
  res.json(estateSearchCondition);
});

app.post("/api/estate/req_doc/:id", async (req, res, next) => {
  const id = req.params.id;
  const getConnection = promisify(db.getConnection.bind(db));
  const connection = await getConnection();
  const query = promisify(connection.query.bind(connection));
  try {
    const id = req.params.id;
    const [estate] = await query("SELECT * FROM estate WHERE id = ?", [id]);
    if (estate == null) {
      res.status(404).send("Not Found");
      return;
    }
    res.json({ ok: true });
  } catch (e) {
    next(e);
  } finally {
    await connection.release();
  }
});

app.post("/api/estate/nazotte", async (req, res, next) => {
  const coordinates = req.body.coordinates;
  const longitudes = coordinates.map((c) => c.longitude);
  const latitudes = coordinates.map((c) => c.latitude);
  const boundingbox = {
    topleft: {
      longitude: Math.min(...longitudes),
      latitude: Math.min(...latitudes),
    },
    bottomright: {
      longitude: Math.max(...longitudes),
      latitude: Math.max(...latitudes),
    },
  };

  const getConnection = promisify(db.getConnection.bind(db));
  const connection = await getConnection();
  const query = promisify(connection.query.bind(connection));
  try {
    const estates = await query("SELECT * FROM estate WHERE latitude <= ? AND latitude >= ? AND longitude <= ? AND longitude >= ? ORDER BY view_count DESC, id ASC", [
      boundingbox.bottomright.latitude, boundingbox.topleft.latitude, boundingbox.bottomright.longitude, boundingbox.topleft.longitude,
    ]);

    const estatesInPolygon = [];
    for (const estate of estates) {
      const point = util.format("'POINT(%f %f)'", estate.latitude, estate.longitude);
      const sql = "SELECT * FROM estate WHERE id = ? AND ST_Contains(ST_PolygonFromText(%s), ST_GeomFromText(%s))";
      const coordinatesToText = util.format("'POLYGON((%s))'", coordinates.map((coordinate) => util.format("%f %f", coordinate.latitude, coordinate.longitude)).join(","));
      const sqlstr = util.format(sql, coordinatesToText, point)
      const [e] = await query(sqlstr, [estate.id]);
      if (e && Object.keys(e).length > 0) {
        estatesInPolygon.push(e);
      }
    }

    const results = {
      estates: [],
    };
    let i = 0;
    for (const estate of estatesInPolygon) {
      if (i >= NAZOTTE_LIMIT) {
        break;
      }
      results.estates.push(camelcaseKeys(estate));
      i++;
    }
    results.count = results.estates.length;
    res.json(results);
  } catch (e) {
    next(e);
  } finally {
    await connection.release();
  }
});

app.get("/api/estate/:id", async (req, res, next) => {
  const getConnection = promisify(db.getConnection.bind(db));
  const connection = await getConnection();
  const query = promisify(connection.query.bind(connection));
  try {
    const id = req.params.id;
    const [estate] = await query("SELECT * FROM estate WHERE id = ?", [id]);
    if (estate == null) {
      res.status(404).send("Not Found");
      return;
    }
    await connection.beginTransaction();
    await query("UPDATE estate SET view_count = ? WHERE id = ?", [estate.view_count+1, id]);
    await connection.commit();
    res.json(estate);
  } catch (e) {
    await connection.rollback();
    next(e);
  } finally {
    await connection.release();
  }
});

app.get("/api/recommended_estate", async (req, res, next) => {
  const getConnection = promisify(db.getConnection.bind(db));
  const connection = await getConnection();
  const query = promisify(connection.query.bind(connection));
  try {
    const es = await query("SELECT * FROM estate ORDER BY view_count DESC, id ASC LIMIT ?", [LIMIT]);
    const estates = es.map((estate) => camelcaseKeys(estate)); 
    res.json({estates});
  } catch (e) {
    next(e);
  } finally {
    await connection.release();
  } 
});

app.get("/api/recommended_estate/:id", async (req, res, next) => {
  const id = req.params.id;
  const getConnection = promisify(db.getConnection.bind(db));
  const connection = await getConnection();
  const query = promisify(connection.query.bind(connection));
  try {
    const [chair] = await query("SELECT * FROM chair WHERE id = ?", [id]);
    const w = chair.width;
    const h = chair.height;
    const d = chair.depth;
    const es = await query("SELECT * FROM estate where (door_width >= ? AND door_height>= ?) OR (door_width >= ? AND door_height>= ?) OR (door_width >= ? AND door_height>=?) OR (door_width >= ? AND door_height>=?) OR (door_width >= ? AND door_height>=?) OR (door_width >= ? AND door_height>=?) ORDER BY view_count DESC, id ASC LIMIT ?", [
      w, h, w, d, h, w, h, d, d, w, d, h, LIMIT
    ]);
    const estates = es.map((estate) => camelcaseKeys(estate)); 
    res.json({ estates });
  } catch (e) {
    next(e);
  } finally {
    await connection.release();
  } 
});

app.get("/api/recommended_chair", async (req, res, next) => {
  const getConnection = promisify(db.getConnection.bind(db));
  const connection = await getConnection();
  const query = promisify(connection.query.bind(connection));
  try {
    const cs = await query("SELECT * FROM chair WHERE stock > 0 ORDER BY view_count DESC, id ASC LIMIT ?", [LIMIT]);
    const chairs = cs.map((chair) => camelcaseKeys(chair)); 
    res.json({ chairs });
  } catch (e) {
    next(e);
  } finally {
    await connection.release();
  } 
});

app.listen(PORT, () => {
  console.log(`Listening ${PORT}`);
});