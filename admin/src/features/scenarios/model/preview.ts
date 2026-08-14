const SCENARIO_PATHNAME_ORIGIN = 'https://scenario.local';

export function isScenarioPathname(value: string): boolean {
  if (!value || value.startsWith('//') || value.includes('\\') || /\s/.test(value)) return false;

  try {
    const url = new URL(value, SCENARIO_PATHNAME_ORIGIN);
    return (
      value.startsWith('/') &&
      url.origin === SCENARIO_PATHNAME_ORIGIN &&
      url.pathname === value &&
      url.search === '' &&
      url.hash === '' &&
      url.href === `${SCENARIO_PATHNAME_ORIGIN}${value}`
    );
  } catch {
    return false;
  }
}

export function buildScenarioPreviewUrl(
  pagePattern: string,
  token: string,
  base = process.env.NEXT_PUBLIC_PREVIEW_URL || 'http://localhost:3000',
): URL {
  if (!isScenarioPathname(pagePattern)) {
    throw new TypeError('Для проверки нужен путь вида /catalog без домена, параметров после ? и #');
  }
  const url = new URL(pagePattern, base);
  url.searchParams.set('token', token);
  return url;
}
