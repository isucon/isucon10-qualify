import { Application, Router } from "https://deno.land/x/oak/mod.ts";

const currentEnv = Deno.env.toObject();
const decoder = new TextDecoder();

const PORT = currentEnv.PORT ?? 1323;
const LIMIT = 20;
const NAZOTTE_LIMIT = 50;
const dbinfo = {
  hostname: currentEnv.MYSQL_HOST ?? "127.0.0.1",
  port: currentEnv.MYSQL_PORT ?? 3306,
  username: currentEnv.MYSQL_USER ?? "isucon",
  password: currentEnv.MYSQL_PASS ?? "isucon",
  db: currentEnv.MYSQL_DBNAME ?? "isuumo",
  poolSize: 10,
};
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

const app = new Application();
app.use(router.routes());
app.use(router.allowedMethods());

await app.listen({ port: +PORT });
