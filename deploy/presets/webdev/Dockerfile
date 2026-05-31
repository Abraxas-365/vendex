# hada-preset-webdev — Web Designer preset for hada-commerce agent.
# Extends hada-preset-base with a live preview frontend and starter templates.

FROM hada-preset-base:latest

# Copy the preset frontend (built React/Vite app showing live preview + file tree)
COPY frontend/dist /frontend/dist

# Copy starter workspace templates (HTML/CSS examples the agent can reference)
COPY workspace /workspace

# The base entrypoint handles serving frontend on :8080 and workspace on :9091
