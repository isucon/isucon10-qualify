const express = require("express");
const mysql = require("mysql");
const path = require("path");
const cp = require("child_process");
const promisify = require("util").promisify;
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
      chairs
    });
  } catch (e) {
    next(e);
  } finally {
    await connection.destroy();
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
    res.json(chair);
  } catch (e) {
    await connection.rollback();
    next(e);
  } finally {
    await connection.destroy();
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
      res.status(400).send("Not Found");
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
    await connection.destroy();
  }
});

//e.GET("/api/estate/:id", getEstateDetail)
//	e.GET("/api/estate/search", searchEstates)
//	e.POST("/api/estate/req_doc/:id", postEstateRequestDocument)
//	e.POST("/api/estate/nazotte", searchEstateNazotte)
//	e.GET("/api/estate/search/condition", getEstateSearchCondition)

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
    await connection.destroy();
  }
});

app.listen(PORT, () => {
  console.log(`Listening ${PORT}`);
});
