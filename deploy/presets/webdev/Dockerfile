# vendex-preset-webdev — Web Designer preset for vendex agent.
# Extends vendex-preset-base with a live preview frontend and starter templates.

FROM vendex-preset-base:latest

# Copy the preset frontend (built React/Vite app showing live preview + file tree)
COPY frontend/dist /frontend/dist

# Copy starter workspace templates (HTML/CSS examples the agent can reference)
COPY workspace /workspace

# The base entrypoint handles serving frontend on :8080 and workspace on :9091
