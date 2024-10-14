/** @type {import('next').NextConfig} */
 
const nextConfig = {
    async rewrites() {
        return [
          {
            source: '/:lng/dashboard/agents/:prefix/action',
            destination: '/[lng]/dashboard/agents/[prefix]/action',
          },
        ]
      },
};

export default nextConfig;
