import type { JwtKey } from "./auth/jwt.ts";

export type Config = {
  jwtKey: JwtKey;
  upstream: string;
};

// readConfig は環境変数（Bun の process.env、Workers の env）から設定を読む。
// エラーには値を出さず、キーの名前だけを出す（docs/conventions.md の 11 節と同じ）。
export function readConfig(env: Record<string, string | undefined>): Config {
  const secret = env.FUNADANSU_JWT_KEY;
  const kid = env.FUNADANSU_JWT_KID;
  const upstream = env.FUNADANSU_UPSTREAM_URL;
  for (const [key, value] of [
    ["FUNADANSU_JWT_KEY", secret],
    ["FUNADANSU_JWT_KID", kid],
    ["FUNADANSU_UPSTREAM_URL", upstream],
  ] as const) {
    if (!value) {
      throw new Error(`${key} が必要です`);
    }
  }
  return {
    jwtKey: { kid: kid!, secret: secret! },
    upstream: upstream!,
  };
}
