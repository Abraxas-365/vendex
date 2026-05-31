#!/bin/sh
set -e

echo "[vendex-preset] Starting preset container..."
echo "[vendex-preset] Frontend port: 8080"
echo "[vendex-preset] Workspace file server port: 9091"

# Serve frontend UI on :8080 (from /frontend/dist)
if [ -d "/frontend/dist" ] && [ "$(ls -A /frontend/dist 2>/dev/null)" ]; then
    sirv /frontend/dist --port 8080 --cors --single &
    echo "[vendex-preset] Frontend serving from /frontend/dist"
else
    # Fallback: serve a placeholder page
    mkdir -p /frontend/dist
    cat > /frontend/dist/index.html <<'EOF'
<!DOCTYPE html>
<html><head><title>Preset Loading</title></head>
<body style="display:flex;align-items:center;justify-content:center;height:100vh;font-family:system-ui;">
<p>Preset frontend not configured. Override /frontend/dist in your preset image.</p>
</body></html>
EOF
    sirv /frontend/dist --port 8080 --cors --single &
    echo "[vendex-preset] Frontend serving placeholder"
fi

# Serve workspace files on :9091 (browseable, for the agent to read/write)
sirv /workspace --port 9091 --cors &
echo "[vendex-preset] Workspace file server ready"

# Keep container alive
wait
