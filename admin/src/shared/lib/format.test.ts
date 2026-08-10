import dayjs from 'dayjs';
import { describe, expect, it } from 'vitest';
import { periodQuery } from './format';

describe('periodQuery', () => {
  it('uses an inclusive first day and exclusive day after the range', () => {
    const query = new URLSearchParams(periodQuery([dayjs('2026-08-01'), dayjs('2026-08-06')]).slice(1));

    expect(dayjs(query.get('from')).format('YYYY-MM-DD')).toBe('2026-08-01');
    expect(dayjs(query.get('to')).format('YYYY-MM-DD')).toBe('2026-08-07');
  });

  it('omits both parameters for all time', () => {
    expect(periodQuery(null)).toBe('');
  });
});
