# Development Stage
FROM node:18-alpine AS development

ENV TZ=Asia/Bangkok

WORKDIR /app

# Update npm to the latest version and install pnpm and yarn globally
RUN npm install -g npm@10.9.0
RUN npm install -g pnpm

COPY pnpm-lock.yaml ./
COPY package.json ./

# Use pnpm to install dependencies
RUN pnpm install 
#--frozen-lockfile
 

COPY . .

EXPOSE 5003

# Use yarn to run the development server
CMD ["yarn", "run", "dev"]

# Builder Stage
FROM node:18-alpine AS builder

ENV TZ=Asia/Bangkok

WORKDIR /app

# Update npm to the latest version and install pnpm globally
RUN npm install -g npm@10.9.0  
RUN npm install -g pnpm
COPY pnpm-lock.yaml ./
COPY package.json ./

# Use pnpm to install dependencies
RUN pnpm install --frozen-lockfile

COPY . .

RUN pnpm run build

# Production Stage
FROM node:18-alpine AS production

ENV TZ=Asia/Bangkok

WORKDIR /app

# Copy the built artifacts from the builder stage
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/package.json ./package.json
COPY --from=builder /app/public ./public

# Set the environment variables (if needed)
ENV NODE_ENV=production

RUN addgroup -g 1001 -S nodejs
RUN adduser -S nextjs -u 1001
USER nextjs

EXPOSE 5003

CMD ["yarn", "start"]
