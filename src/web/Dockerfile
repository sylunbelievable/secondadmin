FROM oven/bun:1.3-alpine AS build
WORKDIR /app
COPY package.json bun.lock ./
RUN bun install --frozen-lockfile
COPY . .
ARG VITE_API_URL=
ENV VITE_API_URL=$VITE_API_URL
RUN bun run build

FROM oven/bun:1.3-alpine AS runtime-deps
WORKDIR /app
COPY package.json bun.lock ./
RUN bun install --production --frozen-lockfile

FROM oven/bun:1.3-alpine
WORKDIR /app
ENV NODE_ENV=production PORT=3000
COPY --from=build --chown=bun:bun /app/dist ./dist
COPY --from=runtime-deps --chown=bun:bun /app/node_modules ./node_modules
USER bun
EXPOSE 3000
HEALTHCHECK --interval=10s --timeout=3s --retries=5 CMD ["bun", "-e", "const r=await fetch('http://127.0.0.1:3000/login');process.exit(r.ok?0:1)"]
CMD ["bun", "dist/server/server.js"]
