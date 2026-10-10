// テスト用の補助。鍵は実行のたびに乱数で作り、コードにも Git にも残さない。
export function randomSecret(): string {
  const bytes = crypto.getRandomValues(new Uint8Array(32));
  return Buffer.from(bytes).toString("base64url");
}

function base64url(input: string | Uint8Array): string {
  const bytes = typeof input === "string" ? new TextEncoder().encode(input) : input;
  return Buffer.from(bytes).toString("base64url");
}

// signToken は HS256 で署名したトークンを作る。header と payload は上書きできる（不正なトークンを作るため）。
export async function signToken(
  secret: string,
  kid: string,
  payload: Record<string, unknown> = {},
  header: Record<string, unknown> = {},
): Promise<string> {
  const h = base64url(JSON.stringify({ alg: "HS256", typ: "JWT", kid, ...header }));
  const p = base64url(JSON.stringify({ sid: "s-1", typ: "operator", exp: Date.now() / 1000 + 300, ...payload }));
  const key = await crypto.subtle.importKey(
    "raw",
    new TextEncoder().encode(secret),
    { name: "HMAC", hash: "SHA-256" },
    false,
    ["sign"],
  );
  const sig = await crypto.subtle.sign("HMAC", key, new TextEncoder().encode(`${h}.${p}`));
  return `${h}.${p}.${base64url(new Uint8Array(sig))}`;
}
