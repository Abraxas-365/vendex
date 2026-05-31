import { useState, useRef, useEffect, useCallback } from 'react'
import { useNavigate } from '@tanstack/react-router'
import {
  Bot,
  Send,
  Loader2,
  Wrench,
  StopCircle,
  ChevronLeft,
  AlertTriangle,
  User,
  RefreshCw,
} from 'lucide-react'
import { useAgentSession, useStopAgentSession, useMarketplacePresets, useSessionHistory } from '../../lib/hooks'
import * as api from '../../lib/api'
import type { AgentEvent, SessionStatus } from '../../types'

// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------

interface ChatMsg {
  role: 'user' | 'assistant'
  content: string
  timestamp: string
  toolName?: string
}

type AgentStatus = 'idle' | 'thinking' | 'tool_use'

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

function statusBadgeClass(status: SessionStatus): string {
  switch (status) {
    case 'creating':
      return 'bg-yellow-100 text-yellow-700'
    case 'running':
      return 'bg-green-100 text-green-700'
    case 'stopped':
      return 'bg-slate-100 text-slate-600'
    case 'failed':
      return 'bg-red-100 text-red-700'
  }
}

function formatTime(dateStr: string): string {
  return new Date(dateStr).toLocaleTimeString('en-US', {
    hour: '2-digit',
    minute: '2-digit',
  })
}

// ---------------------------------------------------------------------------
// Chat panel
// ---------------------------------------------------------------------------

interface ChatPanelProps {
  sessionId: string
  presetId: string
}

function ChatPanel({ sessionId, presetId }: ChatPanelProps) {
  const [messages, setMessages] = useState<ChatMsg[]>([])
  const [input, setInput] = useState('')
  const [status, setStatus] = useState<AgentStatus>('idle')
  const [toolName, setToolName] = useState('')
  const messagesEndRef = useRef<HTMLDivElement>(null)
  const textareaRef = useRef<HTMLTextAreaElement>(null)
  const abortRef = useRef<AbortController | null>(null)

  const { data: history } = useSessionHistory(sessionId)

  // Seed from chat history on load
  useEffect(() => {
    if (history && history.length > 0) {
      setMessages(
        history
          .filter((m) => m.role === 'user' || m.role === 'assistant')
          .map((m) => ({
            role: m.role as 'user' | 'assistant',
            content: m.content,
            timestamp: m.created_at,
          })),
      )
    }
  }, [history])

  // Scroll to bottom
  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [messages, status])

  // Auto-resize textarea
  useEffect(() => {
    if (textareaRef.current) {
      textareaRef.current.style.height = 'auto'
      textareaRef.current.style.height = `${Math.min(textareaRef.current.scrollHeight, 120)}px`
    }
  }, [input])

  const handleSend = useCallback(async () => {
    const trimmed = input.trim()
    if (!trimmed || status !== 'idle') return

    const userMsg: ChatMsg = {
      role: 'user',
      content: trimmed,
      timestamp: new Date().toISOString(),
    }

    setMessages((prev) => [...prev, userMsg])
    setInput('')
    setStatus('thinking')

    // Abort any previous stream
    abortRef.current?.abort()
    const controller = new AbortController()
    abortRef.current = controller

    let assistantText = ''

    try {
      const token = api.getAccessToken()
      const baseUrl = import.meta.env.VITE_API_BASE_URL ?? '/api/v1'

      const res = await fetch(`${baseUrl}/agent/chat`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          ...(token ? { Authorization: `Bearer ${token}` } : {}),
        },
        body: JSON.stringify({
          message: trimmed,
          session_id: sessionId,
          preset_id: presetId,
        }),
        signal: controller.signal,
      })

      if (!res.ok || !res.body) {
        throw new Error(`HTTP ${res.status}`)
      }

      const reader = res.body.getReader()
      const decoder = new TextDecoder()
      let buffer = ''

      // Placeholder assistant message while streaming
      const assistantMsg: ChatMsg = {
        role: 'assistant',
        content: '',
        timestamp: new Date().toISOString(),
      }
      setMessages((prev) => [...prev, assistantMsg])

      while (true) {
        const { done, value } = await reader.read()
        if (done) break

        buffer += decoder.decode(value, { stream: true })
        const lines = buffer.split('\n')
        buffer = lines.pop() ?? ''

        for (const line of lines) {
          if (!line.startsWith('data: ')) continue
          const raw = line.slice(6).trim()
          if (raw === '' || raw === '[DONE]') continue

          let event: AgentEvent
          try {
            event = JSON.parse(raw) as AgentEvent
          } catch {
            continue
          }

          switch (event.kind) {
            case 'text_delta':
              assistantText += event.text ?? ''
              setMessages((prev) => {
                const updated = [...prev]
                const last = updated[updated.length - 1]
                if (last && last.role === 'assistant') {
                  updated[updated.length - 1] = { ...last, content: assistantText }
                }
                return updated
              })
              setStatus('thinking')
              break

            case 'tool_start':
              setStatus('tool_use')
              setToolName(event.tool_name ?? '')
              break

            case 'tool_end':
              setStatus('thinking')
              setToolName('')
              break

            case 'turn_end':
              setStatus('idle')
              setToolName('')
              break

            case 'error':
              setStatus('idle')
              setMessages((prev) => [
                ...prev,
                {
                  role: 'assistant',
                  content: `Error: ${event.error ?? 'Unknown error'}`,
                  timestamp: new Date().toISOString(),
                },
              ])
              break
          }
        }
      }

      setStatus('idle')
      setToolName('')
    } catch (err: unknown) {
      if (err instanceof Error && err.name === 'AbortError') return
      setStatus('idle')
      setMessages((prev) => [
        ...prev,
        {
          role: 'assistant',
          content: err instanceof Error ? err.message : 'An error occurred.',
          timestamp: new Date().toISOString(),
        },
      ])
    }
  }, [input, status, sessionId, presetId])

  function handleKeyDown(e: React.KeyboardEvent) {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault()
      void handleSend()
    }
  }

  return (
    <div className="flex flex-col h-full">
      {/* Messages */}
      <div className="flex-1 overflow-y-auto p-4 space-y-4">
        {messages.length === 0 && (
          <div className="flex flex-col items-center justify-center h-full text-center py-12">
            <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-indigo-100 text-indigo-600 mb-3">
              <Bot size={22} />
            </div>
            <p className="text-sm font-medium text-slate-600">Agent ready</p>
            <p className="text-xs text-slate-400 mt-1">Send a message to get started</p>
          </div>
        )}

        {messages.map((msg, idx) => (
          <div
            key={idx}
            className={`flex gap-2 ${msg.role === 'user' ? 'flex-row-reverse' : ''}`}
          >
            {/* Avatar */}
            <div
              className={`flex h-7 w-7 shrink-0 items-center justify-center rounded-lg ${
                msg.role === 'assistant'
                  ? 'bg-slate-100 text-slate-600'
                  : 'bg-indigo-600 text-white'
              }`}
            >
              {msg.role === 'assistant' ? (
                <Bot size={14} />
              ) : (
                <User size={14} />
              )}
            </div>

            {/* Bubble */}
            <div
              className={`max-w-[82%] rounded-2xl px-3 py-2.5 ${
                msg.role === 'assistant'
                  ? 'bg-slate-100 text-slate-800'
                  : 'bg-indigo-600 text-white'
              }`}
            >
              <p className="whitespace-pre-wrap text-xs leading-relaxed">{msg.content}</p>
              <p
                className={`mt-1 text-[10px] ${
                  msg.role === 'assistant' ? 'text-slate-400' : 'text-indigo-200'
                }`}
              >
                {formatTime(msg.timestamp)}
              </p>
            </div>
          </div>
        ))}

        {/* Thinking / Tool use indicator */}
        {status !== 'idle' && (
          <div className="flex gap-2">
            <div className="flex h-7 w-7 shrink-0 items-center justify-center rounded-lg bg-slate-100 text-slate-600">
              <Bot size={14} />
            </div>
            <div className="flex items-center gap-2 rounded-2xl bg-slate-100 px-3 py-2.5">
              {status === 'thinking' ? (
                <>
                  <Loader2 size={13} className="animate-spin text-slate-500" />
                  <span className="text-xs text-slate-500">Thinking…</span>
                </>
              ) : (
                <>
                  <Wrench size={13} className="animate-pulse text-amber-500" />
                  <span className="text-xs text-slate-500">
                    Using:{' '}
                    <span className="font-mono font-medium text-amber-600">{toolName}</span>
                  </span>
                </>
              )}
            </div>
          </div>
        )}

        <div ref={messagesEndRef} />
      </div>

      {/* Input */}
      <div className="border-t border-slate-200 p-3">
        <div className="flex items-end gap-2 rounded-xl border border-slate-200 bg-white p-2.5 shadow-sm focus-within:border-indigo-300 focus-within:ring-1 focus-within:ring-indigo-200 transition-shadow">
          <textarea
            ref={textareaRef}
            value={input}
            onChange={(e) => setInput(e.target.value)}
            onKeyDown={handleKeyDown}
            placeholder="Message the agent…"
            rows={1}
            className="max-h-28 min-h-[2rem] flex-1 resize-none text-xs text-slate-900 placeholder:text-slate-400 focus:outline-none"
          />
          <button
            onClick={() => void handleSend()}
            disabled={!input.trim() || status !== 'idle'}
            className="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-indigo-600 text-white transition-colors hover:bg-indigo-700 disabled:bg-slate-200 disabled:text-slate-400"
          >
            <Send size={13} />
          </button>
        </div>
      </div>
    </div>
  )
}

// ---------------------------------------------------------------------------
// Main workspace view
// ---------------------------------------------------------------------------

interface WorkspaceViewProps {
  sessionId: string
}

export default function WorkspaceView({ sessionId }: WorkspaceViewProps) {
  const navigate = useNavigate()
  const { data: session, isLoading, error } = useAgentSession(sessionId)
  const { data: presetsPage } = useMarketplacePresets()
  const stopSession = useStopAgentSession()

  const presets = presetsPage?.items ?? []
  const preset = session ? presets.find((p) => p.id === session.preset_id) : undefined

  function handleStop() {
    if (!session) return
    stopSession.mutate(session.id)
  }

  function handleBack() {
    void navigate({ to: '/admin/workspaces' })
  }

  if (isLoading) {
    return (
      <div className="flex h-full items-center justify-center">
        <div className="flex flex-col items-center gap-3">
          <Loader2 size={28} className="animate-spin text-indigo-600" />
          <p className="text-sm text-slate-500">Loading workspace…</p>
        </div>
      </div>
    )
  }

  if (error || !session) {
    return (
      <div className="flex h-full items-center justify-center">
        <div className="flex flex-col items-center gap-3 text-center">
          <AlertTriangle size={28} className="text-red-400" />
          <p className="text-sm font-medium text-slate-700">Workspace not found</p>
          <p className="text-xs text-slate-400">{error?.message ?? 'Unknown error'}</p>
          <button
            onClick={handleBack}
            className="mt-2 text-xs font-medium text-indigo-600 hover:underline"
          >
            ← Back to workspaces
          </button>
        </div>
      </div>
    )
  }

  return (
    <div className="flex flex-col h-[calc(100vh-8rem)]">
      {/* Top bar */}
      <div className="flex h-12 shrink-0 items-center justify-between border-b border-slate-200 bg-white px-4">
        <div className="flex items-center gap-3">
          <button
            onClick={handleBack}
            className="flex h-7 w-7 items-center justify-center rounded-lg text-slate-400 hover:bg-slate-100 hover:text-slate-600 transition-colors"
          >
            <ChevronLeft size={16} />
          </button>
          <div className="flex items-center gap-2">
            <div className="flex h-7 w-7 items-center justify-center rounded-lg bg-indigo-100 text-indigo-600">
              <Bot size={14} />
            </div>
            <div>
              <span className="text-sm font-semibold text-slate-800">
                {session.name || 'Workspace'}
              </span>
              {preset && (
                <span className="ml-2 text-xs text-slate-400">{preset.name}</span>
              )}
            </div>
          </div>
          <span
            className={`inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium capitalize ${statusBadgeClass(session.status)}`}
          >
            {session.status}
          </span>
        </div>

        <div className="flex items-center gap-2">
          {/* Reload iframe */}
          {session.status === 'running' && session.frontend_url && (
            <button
              onClick={() => {
                const iframe = document.querySelector<HTMLIFrameElement>('#workspace-iframe')
                if (iframe) iframe.src = iframe.src
              }}
              className="flex h-7 w-7 items-center justify-center rounded-lg text-slate-400 hover:bg-slate-100 hover:text-slate-600 transition-colors"
              title="Reload workspace frontend"
            >
              <RefreshCw size={14} />
            </button>
          )}

          {(session.status === 'running' || session.status === 'creating') && (
            <button
              onClick={handleStop}
              disabled={stopSession.isPending}
              className="inline-flex items-center gap-1.5 rounded-lg border border-red-200 px-3 py-1.5 text-xs font-medium text-red-600 hover:bg-red-50 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              <StopCircle size={13} />
              {stopSession.isPending ? 'Stopping…' : 'Stop'}
            </button>
          )}
        </div>
      </div>

      {/* Main split layout */}
      <div className="flex flex-1 overflow-hidden">
        {/* Left: Preset frontend iframe (70%) */}
        <div className="w-[70%] shrink-0 border-r border-slate-200 bg-slate-900 overflow-hidden">
          {session.status === 'creating' ? (
            <div className="flex h-full flex-col items-center justify-center text-white gap-4">
              <Loader2 size={32} className="animate-spin text-indigo-400" />
              <div className="text-center">
                <p className="text-sm font-medium">Starting workspace…</p>
                <p className="text-xs text-slate-400 mt-1">This may take a few seconds</p>
              </div>
            </div>
          ) : session.status === 'stopped' || session.status === 'failed' ? (
            <div className="flex h-full flex-col items-center justify-center text-white gap-3">
              <StopCircle size={32} className="text-slate-500" />
              <p className="text-sm font-medium text-slate-300">
                {session.status === 'failed' ? 'Workspace failed to start' : 'Workspace stopped'}
              </p>
            </div>
          ) : session.frontend_url ? (
            <iframe
              id="workspace-iframe"
              src={session.frontend_url}
              className="w-full h-full border-0"
              title={`Workspace: ${session.name || session.id}`}
              sandbox="allow-scripts allow-same-origin allow-forms allow-popups allow-downloads"
            />
          ) : (
            <div className="flex h-full flex-col items-center justify-center text-white gap-3">
              <Bot size={32} className="text-slate-500" />
              <p className="text-sm font-medium text-slate-300">No frontend UI available</p>
              <p className="text-xs text-slate-500">This preset doesn't have a frontend interface</p>
            </div>
          )}
        </div>

        {/* Right: Chat panel (30%) */}
        <div className="flex-1 flex flex-col overflow-hidden bg-white">
          {/* Chat header */}
          <div className="flex h-10 shrink-0 items-center border-b border-slate-100 px-4">
            <Bot size={14} className="text-indigo-600 mr-2" />
            <span className="text-xs font-semibold text-slate-600">Agent Chat</span>
          </div>

          {session.status === 'running' ? (
            <ChatPanel sessionId={session.id} presetId={session.preset_id} />
          ) : session.status === 'creating' ? (
            <div className="flex flex-1 flex-col items-center justify-center text-center p-6">
              <Loader2 size={22} className="animate-spin text-indigo-400 mb-3" />
              <p className="text-xs text-slate-500">Waiting for workspace to start…</p>
            </div>
          ) : (
            <div className="flex flex-1 flex-col items-center justify-center text-center p-6">
              <StopCircle size={22} className="text-slate-300 mb-3" />
              <p className="text-xs text-slate-500">Workspace is not running</p>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
