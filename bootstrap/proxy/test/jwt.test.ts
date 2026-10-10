import { describe, expect, test } from "bun:test";
import { createJwtVerifier, InvalidTokenError } from "../src/auth/jwt.ts";
import { randomSecret, signToken } from "./helper.ts";

const kid = "test-1";

describe("JWT の検査", () => {
  const secret = randomSecret();
  const verify = createJwtVerifier({ kid, secret }, () => 1_000);

  const cases: { name: string; token: () => Promise<string>; ok: boolean }[] = [
    { name: "正しい署名と期限", token: () => signToken(secret, kid, { exp: 2_000 }), ok: true },
    { name: "署名が別の鍵", token: () => signToken(randomSecret(), kid, { exp: 2_000 }), ok: false },
    { name: "kid が違う", token: () => signToken(secret, "other", { exp: 2_000 }), ok: false },
    { name: "alg が none", token: () => signToken(secret, kid, { exp: 2_000 }, { alg: "none" }), ok: false },
    { name: "alg が HS512 の表示", token: () => signToken(secret, kid, { exp: 2_000 }, { alg: "HS512" }), ok: false },
    { name: "期限切れ", token: () => signToken(secret, kid, { exp: 500 }), ok: false },
    { name: "exp が無い", token: () => signToken(secret, kid, { exp: undefined }), ok: false },
    { name: "3 つに分かれていない", token: async () => "abc.def", ok: false },
    { name: "base64url でない文字", token: async () => "a!b.c.d", ok: false },
  ];

  for (const tt of cases) {
    test(tt.name, async () => {
      const token = await tt.token();
      if (tt.ok) {
        const claims = await verify(token);
        expect(claims.sid).toBe("s-1");
      } else {
        await expect(verify(token)).rejects.toBeInstanceOf(InvalidTokenError);
      }
    });
  }
});
