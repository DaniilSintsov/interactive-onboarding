export function buildScenarioPreviewUrl(
  pagePattern: string,
  token: string,
  base = process.env.NEXT_PUBLIC_PREVIEW_URL || 'http://localhost:3000',
): URL {
  const url = new URL(pagePattern, base);
  url.searchParams.set('token', token);
  return url;
}
