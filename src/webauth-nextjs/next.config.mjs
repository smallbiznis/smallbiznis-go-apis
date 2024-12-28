/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  output: "export",
  compress: true,
  experimental: {
    missingSuspenseWithCSRBailout: false
  }
};

export default nextConfig;
