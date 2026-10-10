// Cloudflare Workers の入口。配線はここだけで行う（env は wrangler の vars と secrets）。
import { createApp } from "../app.ts";
import { readConfig } from "../config.ts";

export default {
  fetch(request: Request, env: Record<string, string | undefined>): Response | Promise<Response> {
    const config = readConfig(env);
    const app = createApp({ jwtKey: config.jwtKey, upstream: config.upstream });
    return app.fetch(request);
  },
};
