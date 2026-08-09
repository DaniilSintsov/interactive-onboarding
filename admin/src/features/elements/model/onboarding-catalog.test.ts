import { describe, expect, it } from 'vitest';
import { parseOnboardingCatalogMetadata } from './onboarding-catalog';

describe('parseOnboardingCatalogMetadata', () => {
  it('reads normalized page paths from CI metadata', () => {
    expect(
      parseOnboardingCatalogMetadata(
        'onboarding-catalog:{"page_paths":[" /add-item/details ","/add-item/details"]}',
      ),
    ).toEqual({ page_paths: ['/add-item/details'] });
  });

  it('ignores regular descriptions and invalid metadata', () => {
    expect(parseOnboardingCatalogMetadata('Обычное описание')).toBeNull();
    expect(parseOnboardingCatalogMetadata('onboarding-catalog:{bad-json}')).toBeNull();
    expect(parseOnboardingCatalogMetadata('onboarding-catalog:{"page_paths":["details"]}')).toBeNull();
  });
});
