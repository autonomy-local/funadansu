// JWT の検査の見本。署名方式は検査する側で固定し、トークンのヘッダーの alg を信じない（ADR 0007）。
// Web Crypto だけを使うので、Cloudflare Workers と Bun の両方で動く。

export type JwtKey = {
  kid: string; // 鍵の識別子。ヘッダーの kid と一致しなければ拒否する
  secret: string; // 検証に使う鍵（この見本ではダミーの HMAC の鍵）
};

export type Claims = Record<string, unknown> & { exp: number };

export class InvalidTokenError extends Error {
  constructor(message: string) {
    super(message);
    this.name = "InvalidTokenError";
  }
}

const ALG = "HS256";

export function createJwtVerifier(key: JwtKey, now: () => number = () => Date.now() / 1000) {
  // 鍵の取り込みは初回の検証で一度だけ行う。
  let cryptoKey: Promise<CryptoKey> | undefined;
  const importKey = () => {
    cryptoKey ??= crypto.subtle.importKey(
      "raw",
      new TextEncoder().encode(key.secret),
      { name: "HMAC", hash: "SHA-256" },
      false,
      ["verify"],
    );
    return cryptoKey;
  };

  return async function verify(token: string): Promise<Claims> {
    const parts = token.split(".");
    if (parts.length !== 3) {
      throw new InvalidTokenError("形式が不正です");
    }
    const [headerPart, payloadPart, signaturePart] = parts;

    const header = parseJson(decodeBase64Url(headerPart));
    // alg は固定の値と一致するかだけを見る。他の方式への切り替えには従わない。
    if (header?.alg !== ALG) {
      throw new InvalidTokenError("alg が許可されていません");
    }
    if (header.kid !== key.kid) {
      throw new InvalidTokenError("kid が一致しません");
    }

    const signed = new TextEncoder().encode(`${headerPart}.${payloadPart}`);
    const valid = await crypto.subtle.verify(
      "HMAC",
      await importKey(),
      decodeBase64UrlBytes(signaturePart),
      signed,
    );
    if (!valid) {
      throw new InvalidTokenError("署名が不正です");
    }

    const payload = parseJson(decodeBase64Url(payloadPart));
    if (typeof payload?.exp !== "number" || payload.exp <= now()) {
      throw new InvalidTokenError("期限が切れているか、exp がありません");
    }
    return payload as Claims;
  };
}

function decodeBase64UrlBytes(part: string): Uint8Array<ArrayBuffer> {
  if (!/^[A-Za-z0-9_-]*$/.test(part)) {
    throw new InvalidTokenError("base64url が不正です");
  }
  const padded = part.replace(/-/g, "+").replace(/_/g, "/").padEnd(Math.ceil(part.length / 4) * 4, "=");
  try {
    return Uint8Array.from(atob(padded), (c) => c.charCodeAt(0));
  } catch {
    throw new InvalidTokenError("base64url が不正です");
  }
}

function decodeBase64Url(part: string): string {
  return new TextDecoder().decode(decodeBase64UrlBytes(part));
}

function parseJson(text: string): Record<string, unknown> {
  try {
    const value: unknown = JSON.parse(text);
    if (typeof value === "object" && value !== null && !Array.isArray(value)) {
      return value as Record<string, unknown>;
    }
  } catch {
    // 下で拒否する
  }
  throw new InvalidTokenError("JSON が不正です");
}
