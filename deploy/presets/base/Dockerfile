# Stage 1: Build the live-preview frontend
FROM node:22-alpine AS frontend-build
WORKDIR /app
COPY frontend/package.json ./
RUN npm install
COPY frontend/ ./
RUN npx vite build

# Stage 2: Runtime
FROM node:22-alpine

RUN apk add --no-cache curl jq chromium

# sirv for static file serving
RUN npm install -g sirv-cli@3

# Set up workspace
RUN mkdir -p /workspace/pages /workspace/assets
COPY workspace/ /workspace/

# Copy built frontend
COPY --from=frontend-build /app/dist /frontend/dist

VOLUME /workspace
EXPOSE 8080 9091

HEALTHCHECK --interval=10s --timeout=3s --start-period=5s \
  CMD curl -f http://localhost:8080/ || exit 1

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]
