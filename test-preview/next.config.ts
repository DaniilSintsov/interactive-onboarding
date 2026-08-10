import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  agentRules: false,
  reactStrictMode: true,
  transpilePackages: ["@interactive-onboarding/sdk"],
  experimental: { optimizePackageImports: ["antd"] },
};

export default nextConfig;
