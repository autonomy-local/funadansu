import type { Context } from "hono";

// エラーの形は docs/conventions.md の 8 節に合わせる（{"status": 数値, "message": 文}）。
export function errorBody(status: number, message: string) {
  return { status, message };
}

export function unauthorized(c: Context) {
  return c.json(errorBody(401, "認証に失敗しました"), 401);
}

export function unexpected(c: Context) {
  return c.json(errorBody(503, "未定義のエラーが発生しました"), 503);
}
