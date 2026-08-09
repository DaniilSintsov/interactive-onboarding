export type OnboardingCatalogMetadata = {
  page_paths: string[];
};

const metadataPrefix = 'onboarding-catalog:';

export function parseOnboardingCatalogMetadata(description: string): OnboardingCatalogMetadata | null {
  if (!description.startsWith(metadataPrefix)) return null;
  try {
    const metadata = JSON.parse(description.slice(metadataPrefix.length)) as Partial<OnboardingCatalogMetadata>;
    if (
      !Array.isArray(metadata.page_paths) ||
      metadata.page_paths.some((pagePath) => typeof pagePath !== 'string' || !pagePath.trim().startsWith('/'))
    ) {
      return null;
    }
    return { page_paths: [...new Set(metadata.page_paths.map((pagePath) => pagePath.trim()).filter(Boolean))] };
  } catch {
    return null;
  }
}
