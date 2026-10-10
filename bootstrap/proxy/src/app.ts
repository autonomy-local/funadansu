import { Hono } from "hono";
import { createJwtVerifier, InvalidTokenError, type JwtKey } from "./auth/jwt.ts";
import { unauthorized, unexpected } from "./error.ts";
import { type FetchLike, relay } from "./relay.ts";

export type Deps = {
  jwtKey: JwtKey;
  upstream: string; // Go のサービスのオリジン
  upstreamTimeoutMs?: number;
  fetch?: FetchLike;
  now?: () => number; // テストで時刻を固定するため
};

// createApp は、Workers と Bun の両方で同じアプリを作る。実行環境の違いは src/entry/ に閉じ込める。
export function createApp(deps: Deps) {
  const app = new Hono();
  const verify = createJwtVerifier(deps.jwtKey, deps.now);

  app.onError((err, c) => {
    console.error(JSON.stringify({
      time: new Date().toISOString(),
      level: "ERROR",
      msg: "unexpected error",
      kind: "system",
      err: err.message,
    }));
    return unexpected(c);
  });

  // 認証が要るのは、業務の経路だけ。JWT の検査は入口で行い、通ったものだけ Go へ渡す。
  app.use("/bootstrap/*", async (c, next) => {
    const auth = c.req.header("Authorization") ?? "";
    const match = /^Bearer (.+)$/.exec(auth);
    if (!match) {
      return unauthorized(c);
    }
    try {
      await verify(match[1]);
    } catch (err) {
      if (err instanceof InvalidTokenError) {
        return unauthorized(c);
      }
      throw err;
    }
    await next();
  });

  app.all(
    "/bootstrap/*",
    relay({
      upstream: deps.upstream,
      timeoutMs: deps.upstreamTimeoutMs ?? 10_000,
      // 呼び出し時に this が変わると Workers の fetch は失敗するため、包んでから渡す。
      fetch: deps.fetch ?? ((input, init) => globalThis.fetch(input, init)),
    }),
  );

  return app;
}
