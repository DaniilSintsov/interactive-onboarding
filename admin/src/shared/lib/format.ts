import type { Dayjs } from 'dayjs';

export function formatDate(value: string): string {
  return new Intl.DateTimeFormat('ru-RU', { dateStyle: 'medium' }).format(new Date(value));
}

export function formatPercent(value: number): string {
  return new Intl.NumberFormat('ru-RU', { style: 'percent', maximumFractionDigits: 1 }).format(value);
}

export function periodQuery(range: [Dayjs, Dayjs] | null): string {
  if (!range) return '';
  const query = new URLSearchParams({
    from: range[0].startOf('day').toISOString(),
    to: range[1].add(1, 'day').startOf('day').toISOString(),
  });
  return `?${query}`;
}
