import { describe, expect, it } from 'vitest';
import { buildScenarioPreviewUrl, isScenarioPathname } from './preview';

describe('preview pathname', () => {
  it('accepts exact pathname and encodes the test token', () => {
    expect(isScenarioPathname('/add-item')).toBe(true);
    expect(
      buildScenarioPreviewUrl('/add-item', 'token+/=', 'https://preview.example').toString(),
    ).toBe('https://preview.example/add-item?token=token%2B%2F%3D');
  });

  it.each([
    'https://example.com/add-item',
    '//example.com/add-item',
    '/add-item?draft=1',
    '/add-item#step-1',
    '\\add-item',
    '/catalog/../checkout',
    '/add item',
  ])('rejects non-exact pathname %s', (value) => {
    expect(isScenarioPathname(value)).toBe(false);
    expect(() => buildScenarioPreviewUrl(value, 'token', 'https://preview.example')).toThrow(TypeError);
  });
});
