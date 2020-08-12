import { Application, Router, helpers } from "https://deno.land/x/oak/mod.ts";
import { Client } from "https://deno.land/x/mysql/mod.ts";

const currentEnv = Deno.env.toObject();
const decoder = new TextDecoder();

const PORT = currentEnv.PORT ?? 1323;
const LIMIT = 20;
const NAZOTTE_LIMIT = 50;
const dbinfo = {
  hostname: currentEnv.MYSQL_HOST ?? "127.0.0.1",
  port: parseInt(currentEnv.MYSQL_PORT ?? 3306),
  username: currentEnv.MYSQL_USER ?? "isucon",
  password: currentEnv.MYSQL_PASS ?? "isucon",
  db: currentEnv.MYSQL_DBNAME ?? "isuumo",
  poolSize: 10,
};

const chairSearchConditionJson = await Deno.readFile(
  "../fixture/chair_condition.json",
);
const chairSearchCondition = JSON.parse(
  decoder.decode(chairSearchConditionJson),
);

const estateSearchConditionJson = await Deno.readFile(
  "../fixture/estate_condition.json",
);
const estateSearchCondition = JSON.parse(
  decoder.decode(estateSearchConditionJson),
);

const db = await new Client().connect(dbinfo);
const router = new Router();

router.post("/initialize", async (ctx) => {
  const dbdir = "../mysql/db";
  const dbfiles = [
    "0_Schema.sql",
    "1_DummyEstateData.sql",
    "2_DummyChairData.sql",
  ];
  const execfiles = dbfiles.map((file) => `${dbdir}/${file}`);
  for (const execfile of execfiles) {
    const p = Deno.run({
      cmd: [
        "mysql",
        "-h",
        `${dbinfo.hostname}`,
        "-u",
        `${dbinfo.username}`,
        `-p${dbinfo.password}`,
        "-P",
        `${dbinfo.port}`,
        `${dbinfo.db}`,
      ],
      stdin: "piped",
      stdout: "piped",
    });
    const content = await Deno.readFile(execfile);
    let bytes = 0;
    while (bytes < content.byteLength) {
      let b = content.slice(bytes);
      bytes += await p.stdin?.write(b);
    }
    p.stdin?.close();
    const status = await p.status();
    if (!status.success) {
      const output = await p.output();
      throw new Error("Deno run is failed " + output);
    }
  }
  ctx.response.body = {
    language: "deno",
  };
});

router.get("/api/chair/search", async (ctx, next) => {
  const searchQueries = [];
  const queryParams = [];
  const {
    priceRangeId,
    heightRangeId,
    widthRangeId,
    depthRangeId,
    color,
    features,
    page,
    perPage,
  } = helpers.getQuery(ctx);

  if (priceRangeId != null) {
    const chairPrice = chairSearchCondition["price"].ranges[priceRangeId];
    if (chairPrice == null) {
      ctx.response.status = 400;
      ctx.response.body = "priceRangeID invalid";
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
      ctx.response.status = 400;
      ctx.response.body = "heightRangeId invalid";
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
      ctx.response.status = 400;
      ctx.response.body = "widthRangeId invalid";
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
      ctx.response.status = 400;
      ctx.response.body = "depthRangeId invalid";
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
    ctx.response.status = 400;
    ctx.response.body = "Search condition not found";
    return;
  }

  searchQueries.push("stock > 0");
  const pageNum = parseInt(page, 10);
  const perPageNum = parseInt(perPage, 10);

  if (!page || Number.isNaN(pageNum)) {
    ctx.response.status = 400;
    ctx.response.body = `page condition invalid ${page}`;
    return;
  }

  if (!perPage || Number.isNaN(perPageNum)) {
    ctx.response.status = 400;
    ctx.response.body = `perPage condition invalid ${perPage}`;
    return;
  }


  const sqlprefix = "SELECT * FROM chair WHERE ";
  const searchCondition = searchQueries.join(" AND ");
  const limitOffset = " ORDER BY view_count DESC, id ASC LIMIT ? OFFSET ?";
  const countprefix = "SELECT COUNT(*) as count FROM chair WHERE ";

  try {
    const [{ count }] = await db.query(
      `${countprefix}${searchCondition}`,
      queryParams,
    );
    queryParams.push(perPageNum, perPageNum * pageNum);
    const chairs = await db.query(
      `${sqlprefix}${searchCondition}${limitOffset}`,
      queryParams,
    );
    ctx.response.body = {
      count,
      chairs: chairs,
    };
  } catch (e) {
    ctx.response.status = 500;
    ctx.response.body = e.toString();
  }
});

router.get("/api/chair/search/condition", (ctx) => {
  ctx.response.body = chairSearchCondition;
});

router.get("/api/chair/:id", async (ctx) => {
  try {
    const id = ctx.params.id;
    const [chair] = await db.query("SELECT * FROM chair WHERE id = ?", [id]);
    if (chair == null || chair.stock <= 0) {
      ctx.response.status = 404;
      ctx.response.body = "Not Found";
      return;
    }
    await db.transaction(async (conn) => {
      await conn.execute("UPDATE chair SET view_count = ? WHERE id = ?", [chair.view_count+1, id]);
      await conn.execute("COMMIT");
    });
    ctx.response.body = chair;
  } catch (e) {
    ctx.response.status = 500;
    ctx.response.body = e.toString();
  }
});

router.post("/api/chair/buy/:id", async (ctx) => {
  try {
    const id = ctx.params.id;
    const [chair] = await db.query("SELECT * FROM chair WHERE id = ? AND stock > 0", [id]);
    if (chair == null) {
      ctx.response.status = 404;
      ctx.response.body = "Not Found";
      return;
    }

    await db.transaction(async (conn) => {
      await conn.execute("UPDATE chair SET stock = ? WHERE id = ?", [chair.stock-1, id]);
      await conn.execute("COMMIT");
    });
    
    ctx.response.body = { ok: true };
  } catch (e) {
    ctx.response.status = 500;
    ctx.response.body = e.toString(); 
  }
});
const app = new Application();
app.use(router.routes());
app.use(router.allowedMethods());

await app.listen({ port: +PORT });
