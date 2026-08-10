import { NextRequest } from "next/server";

export function GET(request: NextRequest) {
  const query = request.nextUrl.searchParams.get("q")?.trim().toLowerCase() ?? "";
  const items = query.includes("книг")
    ? ["Электронные книги", "Книги", "Аксессуары для чтения"]
    : ["Электронные книги", "Коллекционирование", "Другое"];

  return Response.json({ items });
}
