import React, { useEffect, useState } from "react";

interface FileEntry {
  name: string;
  path: string;
  isDir: boolean;
}

// The workspace file server runs on :9091 inside the container.
// In production, this is proxied through the vendex backend.
const WORKSPACE_URL = window.location.port === "5173"
  ? "http://localhost:9091"  // dev mode
  : `${window.location.origin}/workspace`;

export function App() {
  const [files, setFiles] = useState<FileEntry[]>([]);
  const [selectedFile, setSelectedFile] = useState<string | null>(null);
  const [previewKey, setPreviewKey] = useState(0);

  // Poll workspace files
  useEffect(() => {
    const fetchFiles = async () => {
      try {
        const res = await fetch(`${WORKSPACE_URL}/`);
        if (res.ok) {
          const html = await res.text();
          // Parse directory listing from sirv
          const entries = parseDirectoryListing(html);
          setFiles(entries);
        }
      } catch {
        // Container may not be ready yet
      }
    };

    fetchFiles();
    const interval = setInterval(fetchFiles, 3000);
    return () => clearInterval(interval);
  }, []);

  // Listen for refresh messages from parent (vendex admin)
  useEffect(() => {
    const handler = (event: MessageEvent) => {
      if (event.data?.type === "refresh-preview") {
        setPreviewKey((k) => k + 1);
      }
      if (event.data?.type === "select-file") {
        setSelectedFile(event.data.path);
      }
    };
    window.addEventListener("message", handler);
    return () => window.removeEventListener("message", handler);
  }, []);

  const previewUrl = selectedFile
    ? `${WORKSPACE_URL}/${selectedFile}`
    : null;

  return (
    <div style={styles.container}>
      {/* File Tree Panel */}
      <aside style={styles.sidebar}>
        <h3 style={styles.sidebarTitle}>Workspace Files</h3>
        <div style={styles.fileList}>
          {files.length === 0 && (
            <p style={styles.empty}>No files yet. The agent will create them.</p>
          )}
          {files.map((f) => (
            <button
              key={f.path}
              onClick={() => setSelectedFile(f.path)}
              style={{
                ...styles.fileItem,
                ...(selectedFile === f.path ? styles.fileItemActive : {}),
              }}
            >
              <span>{f.isDir ? "\ud83d\udcc1" : "\ud83d\udcc4"}</span>
              <span>{f.name}</span>
            </button>
          ))}
        </div>
        <button
          onClick={() => setPreviewKey((k) => k + 1)}
          style={styles.refreshBtn}
        >
          Refresh Preview
        </button>
      </aside>

      {/* Live Preview Panel */}
      <main style={styles.main}>
        {previewUrl ? (
          <iframe
            key={previewKey}
            src={previewUrl}
            style={styles.iframe}
            title="Live Preview"
          />
        ) : (
          <div style={styles.placeholder}>
            <h2>Live Preview</h2>
            <p>Select a file from the sidebar or wait for the agent to create pages.</p>
          </div>
        )}
      </main>
    </div>
  );
}

function parseDirectoryListing(html: string): FileEntry[] {
  // sirv serves a basic HTML directory listing with links
  const entries: FileEntry[] = [];
  const linkRegex = /<a[^>]*href="([^"]+)"[^>]*>([^<]+)<\/a>/g;
  let match;
  while ((match = linkRegex.exec(html)) !== null) {
    const href = match[1];
    const name = match[2];
    if (name === "../" || href === "../") continue;
    entries.push({
      name: name.replace(/\/$/, ""),
      path: href.replace(/^\//, "").replace(/\/$/, ""),
      isDir: href.endsWith("/"),
    });
  }
  return entries;
}

const styles: Record<string, React.CSSProperties> = {
  container: {
    display: "flex",
    height: "100vh",
    fontFamily: "system-ui, -apple-system, sans-serif",
    background: "#0f0f0f",
    color: "#e0e0e0",
  },
  sidebar: {
    width: 260,
    borderRight: "1px solid #2a2a2a",
    display: "flex",
    flexDirection: "column",
    padding: 16,
    background: "#161616",
  },
  sidebarTitle: {
    margin: "0 0 16px 0",
    fontSize: 14,
    fontWeight: 600,
    color: "#999",
    textTransform: "uppercase" as const,
    letterSpacing: "0.05em",
  },
  fileList: {
    flex: 1,
    overflow: "auto",
  },
  fileItem: {
    display: "flex",
    alignItems: "center",
    gap: 8,
    width: "100%",
    padding: "8px 12px",
    border: "none",
    borderRadius: 6,
    background: "transparent",
    color: "#e0e0e0",
    cursor: "pointer",
    fontSize: 13,
    textAlign: "left" as const,
  },
  fileItemActive: {
    background: "#2a2a2a",
    color: "#fff",
  },
  refreshBtn: {
    marginTop: 12,
    padding: "10px 16px",
    border: "1px solid #333",
    borderRadius: 6,
    background: "#1a1a1a",
    color: "#e0e0e0",
    cursor: "pointer",
    fontSize: 13,
  },
  main: {
    flex: 1,
    display: "flex",
    alignItems: "center",
    justifyContent: "center",
  },
  iframe: {
    width: "100%",
    height: "100%",
    border: "none",
    background: "#fff",
  },
  placeholder: {
    textAlign: "center" as const,
    color: "#666",
  },
  empty: {
    fontSize: 13,
    color: "#666",
    padding: "8px 0",
  },
};
