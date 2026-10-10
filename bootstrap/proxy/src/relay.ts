import type { Context, Handler } from "hono";
import { unexpected } from "./error.ts";

// テストで差し替えられるように、fetch は関数として受け取る。
export type FetchLike = (input: URL | Request, init?: RequestInit) => Promise<Response>;

export type RelayOptions = {
  upstream: string; // Go のサービスのオリジン（例：http://127.0.0.1:8080）
  timeoutMs: number; // 中継の時限
  fetch: FetchLike;
};

// 接続ごとの決まりの見出し。中継では転送しない（RFC 9110 の hop-by-hop）。
const hopByHop = new Set([
  "connection",
  "keep-alive",
  "proxy-authenticate",
  "proxy-authorization",
  "te",
  "trailer",
  "transfer-encoding",
  "upgrade",
  "host",
]);

// relay は、受け取ったリクエストを Go のサービスへそのまま渡し、応答を返す。
// Authorization は落とさない。サービスでも JWT の署名を検査するため（ADR 0007）。
export function relay(options: RelayOptions): Handler {
  return async (c: Context) => {
    const incoming = new URL(c.req.url);
    const target = new URL(incoming.pathname + incoming.search, options.upstream);

    const headers = new Headers();
    c.req.raw.headers.forEach((value, name) => {
      if (!hopByHop.has(name.toLowerCase())) {
        headers.set(name, value);
      }
    });

    const hasBody = c.req.method !== "GET" && c.req.method !== "HEAD";
    try {
      const res = await options.fetch(target, {
        method: c.req.method,
        headers,
        body: hasBody ? c.req.raw.body : undefined,
        // Bun の fetch は、ストリームを本文に渡すとき duplex が必要。Workers は無視する。
        duplex: "half",
        redirect: "manual",
        signal: AbortSignal.timeout(options.timeoutMs),
      } as RequestInit);

      const out = new Headers();
      res.headers.forEach((value, name) => {
        // fetch が展開した後の本文に、元の長さや符号化の見出しは合わない。
        if (!hopByHop.has(name.toLowerCase()) && !["content-length", "content-encoding"].includes(name.toLowerCase())) {
          out.set(name, value);
        }
      });
      return new Response(res.body, { status: res.status, headers: out });
    } catch (err) {
      console.error(JSON.stringify({
        time: new Date().toISOString(),
        level: "ERROR",
        msg: "upstream failed",
        kind: "system",
        err: err instanceof Error ? err.message : String(err),
      }));
      return unexpected(c);
    }
  };
}
