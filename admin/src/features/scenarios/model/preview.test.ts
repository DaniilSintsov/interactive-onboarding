import { describe, expect, it } from 'vitest';
import { buildScenarioPreviewUrl } from './preview';

describe('buildScenarioPreviewUrl', () => {
  it('keeps scenario path and query while encoding the test token', () => {
    expect(
      buildScenarioPreviewUrl('/add-item?draft=1', 'token+/=', 'https://preview.example').toString(),
    ).toBe('https://preview.example/add-item?draft=1&token=token%2B%2F%3D');
  });
});
