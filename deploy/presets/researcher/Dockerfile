# vendex-preset-researcher — Competitor Research preset for vendex agent.
# Extends base with Python for data analysis and Playwright for web scraping.

FROM vendex-preset-base:latest

# Install Python for data analysis scripts
RUN apk add --no-cache python3 py3-pip

# Copy the preset frontend (competitor dashboard UI)
COPY frontend/dist /frontend/dist

# Reports workspace
RUN mkdir -p /workspace/reports /workspace/data
