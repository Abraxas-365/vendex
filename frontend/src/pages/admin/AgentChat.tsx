import { useState, useRef, useEffect } from 'react'
import { Send, Bot, User, Loader2, Wrench } from 'lucide-react'
import type { AgentMessage } from '../../types'

// Placeholder messages for demo
const initialMessages: AgentMessage[] = [
  {
    role: 'assistant',
    content:
      'Hello! I\'m the Hada Commerce assistant. I can help you create product pages, write descriptions, generate promo campaigns, and more. What would you like to do?',
    timestamp: new Date(Date.now() - 60000).toISOString(),
  },
]

type AgentStatus = 'idle' | 'thinking' | 'tool_use'

export default function AgentChat() {
  const [messages, setMessages] = useState<AgentMessage[]>(initialMessages)
  const [input, setInput] = useState('')
  const [status, setStatus] = useState<AgentStatus>('idle')
  const [toolName, setToolName] = useState('')
  const messagesEndRef = useRef<HTMLDivElement>(null)
  const textareaRef = useRef<HTMLTextAreaElement>(null)

  // Scroll to bottom when messages change
  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [messages, status])

  // Auto-resize textarea
  useEffect(() => {
    if (textareaRef.current) {
      textareaRef.current.style.height = 'auto'
      textareaRef.current.style.height = `${Math.min(textareaRef.current.scrollHeight, 160)}px`
    }
  }, [input])

  function formatTime(dateStr: string): string {
    return new Date(dateStr).toLocaleTimeString('en-US', {
      hour: '2-digit',
      minute: '2-digit',
    })
  }

  async function handleSend() {
    const trimmed = input.trim()
    if (!trimmed || status !== 'idle') return

    const userMsg: AgentMessage = {
      role: 'user',
      content: trimmed,
      timestamp: new Date().toISOString(),
    }

    setMessages((prev) => [...prev, userMsg])
    setInput('')

    // Simulate agent response
    setStatus('thinking')

    // Simulate tool use after a brief delay
    setTimeout(() => {
      setStatus('tool_use')
      setToolName('search_products')
    }, 1500)

    // Simulate final response
    setTimeout(() => {
      const assistantMsg: AgentMessage = {
        role: 'assistant',
        content: `I've processed your request: "${trimmed}". This is a placeholder response — the actual WebSocket connection to the agent backend is a TODO. Once connected, I'll be able to create pages, manage products, and run campaigns for you.`,
        timestamp: new Date().toISOString(),
      }
      setMessages((prev) => [...prev, assistantMsg])
      setStatus('idle')
      setToolName('')
    }, 3000)
  }

  function handleKeyDown(e: React.KeyboardEvent) {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault()
      handleSend()
    }
  }

  return (
    <div className="flex h-[calc(100vh-8rem)] flex-col">
      {/* Header */}
      <div className="border-b border-gray-200 pb-4">
        <div className="flex items-center gap-3">
          <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-gray-900 text-white">
            <Bot className="h-5 w-5" />
          </div>
          <div>
            <h1 className="text-lg font-bold text-gray-900">AI Agent</h1>
            <p className="text-xs text-gray-500">
              {status === 'idle'
                ? 'Ready to help'
                : status === 'thinking'
                  ? 'Thinking...'
                  : `Using tool: ${toolName}`}
            </p>
          </div>
        </div>
      </div>

      {/* Messages */}
      <div className="flex-1 overflow-y-auto py-6">
        <div className="mx-auto max-w-3xl space-y-6">
          {messages.map((msg, idx) => (
            <div
              key={idx}
              className={`flex gap-3 ${msg.role === 'user' ? 'flex-row-reverse' : ''}`}
            >
              {/* Avatar */}
              <div
                className={`flex h-8 w-8 shrink-0 items-center justify-center rounded-lg ${
                  msg.role === 'assistant'
                    ? 'bg-gray-100 text-gray-600'
                    : 'bg-blue-600 text-white'
                }`}
              >
                {msg.role === 'assistant' ? (
                  <Bot className="h-4 w-4" />
                ) : (
                  <User className="h-4 w-4" />
                )}
              </div>

              {/* Message bubble */}
              <div
                className={`max-w-[75%] rounded-2xl px-4 py-3 ${
                  msg.role === 'assistant'
                    ? 'bg-gray-100 text-gray-800'
                    : 'bg-blue-600 text-white'
                }`}
              >
                <p className="whitespace-pre-wrap text-sm leading-relaxed">{msg.content}</p>
                <p
                  className={`mt-1.5 text-[10px] ${
                    msg.role === 'assistant' ? 'text-gray-400' : 'text-blue-200'
                  }`}
                >
                  {formatTime(msg.timestamp)}
                </p>
              </div>
            </div>
          ))}

          {/* Thinking / Tool use indicator */}
          {status !== 'idle' && (
            <div className="flex gap-3">
              <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-gray-100 text-gray-600">
                <Bot className="h-4 w-4" />
              </div>
              <div className="flex items-center gap-2 rounded-2xl bg-gray-100 px-4 py-3">
                {status === 'thinking' ? (
                  <>
                    <Loader2 className="h-4 w-4 animate-spin text-gray-500" />
                    <span className="text-sm text-gray-500">Thinking...</span>
                  </>
                ) : (
                  <>
                    <Wrench className="h-4 w-4 animate-pulse text-amber-500" />
                    <span className="text-sm text-gray-500">
                      Using tool:{' '}
                      <span className="font-mono font-medium text-amber-600">{toolName}</span>
                    </span>
                  </>
                )}
              </div>
            </div>
          )}

          <div ref={messagesEndRef} />
        </div>
      </div>

      {/* Input area */}
      <div className="border-t border-gray-200 pt-4">
        <div className="mx-auto max-w-3xl">
          <div className="flex items-end gap-3 rounded-xl border border-gray-200 bg-white p-3 shadow-sm focus-within:border-gray-400 focus-within:ring-1 focus-within:ring-gray-400 transition-shadow">
            <textarea
              ref={textareaRef}
              value={input}
              onChange={(e) => setInput(e.target.value)}
              onKeyDown={handleKeyDown}
              placeholder="Ask the agent to create a landing page, write product descriptions, or generate promo codes..."
              rows={1}
              className="max-h-40 min-h-[2.5rem] flex-1 resize-none text-sm text-gray-900 placeholder:text-gray-400 focus:outline-none"
            />
            <button
              onClick={handleSend}
              disabled={!input.trim() || status !== 'idle'}
              className="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-gray-900 text-white transition-colors hover:bg-gray-800 disabled:bg-gray-200 disabled:text-gray-400"
            >
              <Send className="h-4 w-4" />
            </button>
          </div>
          <p className="mt-2 text-center text-[10px] text-gray-400">
            Press Enter to send, Shift+Enter for new line. WebSocket connection is a TODO.
          </p>
        </div>
      </div>
    </div>
  )
}
