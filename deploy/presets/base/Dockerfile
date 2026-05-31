# vendex-preset-base — Base image for all vendex agent preset containers.
# Provides: Node.js runtime, lightweight HTTP server (sirv), workspace structure.
# Preset images extend this with their own frontend build and workspace content.

FROM node:22-alpine

# Install lightweight static file server + common utilities
RUN apk add --no-cache curl jq && \
    npm install -g sirv-cli@3

# Create workspace structure
RUN mkdir -p /workspace /frontend/dist /tools

# Default workspace volume mount point
VOLUME /workspace

# Expose frontend port (preset UI) and workspace file server port
EXPOSE 8080 9091

# Health check on frontend port
HEALTHCHECK --interval=10s --timeout=3s --start-period=5s \
    CMD curl -f http://localhost:8080/ || exit 1

# Default: serve the frontend on :8080 and workspace files on :9091
# Preset images override CMD or ENTRYPOINT as needed.
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
