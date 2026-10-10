// コンテナ（Bun）の入口。配線はここだけで行う。
import { createApp } from "../app.ts";
import { readConfig } from "../config.ts";

const config = readConfig(process.env);
const app = createApp({ jwtKey: config.jwtKey, upstream: config.upstream });

// 既定は localhost。公開するときは、前段の経路（リバースプロキシなど）で制限する。
const addr = process.env.FUNADANSU_PROXY_ADDR ?? "127.0.0.1:8787";
const [hostname, port] = addr.split(/:(?=\d+$)/);

Bun.serve({ hostname, port: Number(port), fetch: app.fetch });
console.log(JSON.stringify({ time: new Date().toISOString(), level: "INFO", msg: "listening", kind: "system", addr }));
