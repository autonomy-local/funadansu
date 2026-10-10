import { afterAll, beforeAll, describe, expect, test } from "bun:test";
import { createApp } from "../src/app.ts";
import { randomSecret, signToken } from "./helper.ts";

// 中継の先は、実際の HTTP サーバー（ローカルの一時ポート）。受け取った内容を返すだけの見本。
const upstream = Bun.serve({
  port: 0,
  hostname: "127.0.0.1",
  async fetch(req) {
    const url = new URL(req.url);
    return Response.json({
      method: req.method,
      path: url.pathname,
      query: url.search,
      authorization: req.headers.get("authorization"),
      body: req.method === "POST" ? await req.text() : null,
    }, { status: 201, headers: { "x-upstream": "yes" } });
  },
});
const upstreamOrigin = `http://127.0.0.1:${upstream.port}`;

const kid = "test-1";
const secret = randomSecret();
const app = createApp({
  jwtKey: { kid, secret },
  upstream: upstreamOrigin,
  now: () => 1_000,
});

afterAll(() => upstream.stop(true));

describe("入口の認証", () => {
  test("Authorization が無いと 401 で、中継しない", async () => {
    const res = await app.request("/bootstrap/operators/1");
    expect(res.status).toBe(401);
    expect(await res.json()).toEqual({ status: 401, message: "認証に失敗しました" });
  });

  test("署名が不正な JWT は 401", async () => {
    const token = await signToken(randomSecret(), kid, { exp: 2_000 });
    const res = await app.request("/bootstrap/operators/1", { headers: { Authorization: `Bearer ${token}` } });
    expect(res.status).toBe(401);
  });

  test("署名が正しい JWT は通る", async () => {
    const token = await signToken(secret, kid, { exp: 2_000 });
    const res = await app.request("/bootstrap/operators/1", { headers: { Authorization: `Bearer ${token}` } });
    expect(res.status).toBe(201);
  });
});

describe("Go への中継", () => {
  let token: string;
  beforeAll(async () => {
    token = await signToken(secret, kid, { exp: 2_000 });
  });

  test("GET のパス・クエリ・JWT が、そのまま Go に渡る", async () => {
    const res = await app.request("/bootstrap/operators/1?x=2", { headers: { Authorization: `Bearer ${token}` } });
    expect(res.status).toBe(201);
    expect(res.headers.get("x-upstream")).toBe("yes");
    expect(await res.json()).toEqual({
      method: "GET",
      path: "/bootstrap/operators/1",
      query: "?x=2",
      authorization: `Bearer ${token}`,
      body: null,
    });
  });

  test("POST の本文が、そのまま Go に渡る", async () => {
    const res = await app.request("/bootstrap/operators", {
      method: "POST",
      headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
      body: '{"id":1}',
    });
    expect((await res.json()).body).toBe('{"id":1}');
  });

  test("Go に届かないときは、想定外のエラーの形（503）で返す", async () => {
    const down = createApp({ jwtKey: { kid, secret }, upstream: "http://127.0.0.1:1", now: () => 1_000 });
    const res = await down.request("/bootstrap/operators/1", { headers: { Authorization: `Bearer ${token}` } });
    expect(res.status).toBe(503);
    expect(await res.json()).toEqual({ status: 503, message: "未定義のエラーが発生しました" });
  });
});
