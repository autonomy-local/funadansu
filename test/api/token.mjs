// テスト用の JWT（HS256）を作って標準出力に出します。鍵と kid は環境変数で受け取ります。
// 鍵は run.sh が実行のたびに乱数で作るため、コードにも Git にも残りません（bootstrap/proxy/test/helper.ts と同じ考え方）。
import { createHmac } from "node:crypto";

const key = process.env.FUNADANSU_JWT_KEY;
const kid = process.env.FUNADANSU_JWT_KID;
if (!key || !kid) {
  console.error("FUNADANSU_JWT_KEY と FUNADANSU_JWT_KID が必要です");
  process.exit(1);
}

const encode = (value) => Buffer.from(JSON.stringify(value)).toString("base64url");
const header = encode({ alg: "HS256", typ: "JWT", kid });
const payload = encode({ sid: "runn-sample", typ: "operator", exp: Math.floor(Date.now() / 1000) + 600 });
const signature = createHmac("sha256", key).update(`${header}.${payload}`).digest("base64url");

process.stdout.write(`${header}.${payload}.${signature}`);
