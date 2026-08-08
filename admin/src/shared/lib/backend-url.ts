export function buildBackendUrl(path: string[], search = '', base = process.env.ADMIN_API_URL): URL {
  const apiBase = new URL(base || 'http://localhost:8080');
  if (apiBase.protocol !== 'http:' && apiBase.protocol !== 'https:') {
    throw new Error('ADMIN_API_URL must use http or https');
  }
  apiBase.pathname = `/api/v1/${path.map(encodeURIComponent).join('/')}`;
  apiBase.search = search;
  return apiBase;
}
