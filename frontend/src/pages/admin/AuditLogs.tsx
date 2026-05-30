import { useState } from 'react'
import { Shield, ChevronDown, ChevronUp, X } from 'lucide-react'
import { useAuditLogs, useAuditStats, useAuditLog } from '../../lib/hooks'

function AuditDetailPanel({ id, onClose }: { id: string; onClose: () => void }) {
  const { data: log, isLoading } = useAuditLog(id)
  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-8">
        <div className="h-6 w-6 animate-spin rounded-full border-2 border-gray-300 border-t-gray-900" />
      </div>
    )
  }
  if (!log) return null
  return (
    <div className="rounded-xl border border-gray-200 bg-white p-6 shadow-sm">
      <div className="mb-4 flex items-center justify-between">
        <h2 className="text-lg font-semibold text-gray-900">Audit Log Detail</h2>
        <button onClick={onClose} className="text-gray-400 hover:text-gray-600">
          <X className="h-5 w-5" />
        </button>
      </div>
      <div className="grid grid-cols-2 gap-3 text-sm">
        <div><span className="text-gray-500">User:</span> <span className="font-medium">{log.user_email}</span></div>
        <div><span className="text-gray-500">Action:</span> <span className="font-medium font-mono">{log.action}</span></div>
        <div><span className="text-gray-500">Resource:</span> <span className="font-medium">{log.resource_type}</span></div>
        <div><span className="text-gray-500">Resource ID:</span> <span className="font-medium font-mono text-xs">{log.resource_id}</span></div>
        <div><span className="text-gray-500">IP:</span> <span className="font-medium font-mono">{log.ip_address}</span></div>
        <div><span className="text-gray-500">Date:</span> <span className="font-medium">{new Date(log.created_at).toLocaleString()}</span></div>
      </div>
      {log.changes && Object.keys(log.changes).length > 0 && (
        <div className="mt-4">
          <p className="text-sm font-semibold text-gray-700 mb-2">Changes</p>
          <pre className="rounded-lg bg-gray-50 p-3 text-xs text-gray-700 overflow-auto max-h-40">
            {JSON.stringify(log.changes, null, 2)}
          </pre>
        </div>
      )}
    </div>
  )
}

export default function AuditLogs() {
  const [filters, setFilters] = useState({
    user_id: '',
    action: '',
    resource_type: '',
    from: '',
    to: '',
  })
  const [expandedId, setExpandedId] = useState<string | null>(null)

  const cleanFilters = Object.fromEntries(
    Object.entries(filters).filter(([, v]) => v !== '')
  ) as Record<string, string>

  const { data, isLoading } = useAuditLogs(cleanFilters)
  const { data: stats } = useAuditStats()
  const logs = data?.items ?? []

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-2xl font-bold text-gray-900">Audit Logs</h1>
        <p className="mt-1 text-sm text-gray-500">Track all admin actions</p>
      </div>

      {/* Stats Cards */}
      {stats && (
        <div className="grid grid-cols-2 gap-4 sm:grid-cols-4">
          <div className="rounded-xl border border-gray-200 bg-white p-4 shadow-sm">
            <p className="text-xs text-gray-500">Total Actions</p>
            <p className="mt-1 text-2xl font-bold text-gray-900">{stats.total_actions}</p>
          </div>
          <div className="rounded-xl border border-gray-200 bg-white p-4 shadow-sm">
            <p className="text-xs text-gray-500">Recent (24h)</p>
            <p className="mt-1 text-2xl font-bold text-indigo-600">{stats.recent_activity}</p>
          </div>
          {stats.actions_by_user.slice(0, 2).map((u) => (
            <div key={u.user_email} className="rounded-xl border border-gray-200 bg-white p-4 shadow-sm">
              <p className="text-xs text-gray-500 truncate">{u.user_email}</p>
              <p className="mt-1 text-2xl font-bold text-gray-900">{u.count}</p>
            </div>
          ))}
        </div>
      )}

      {/* Filters */}
      <div className="grid grid-cols-2 gap-3 sm:grid-cols-5">
        <input
          type="text"
          placeholder="Filter by user ID..."
          value={filters.user_id}
          onChange={(e) => setFilters({ ...filters, user_id: e.target.value })}
          className="rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-gray-400 focus:outline-none"
        />
        <input
          type="text"
          placeholder="Filter by action..."
          value={filters.action}
          onChange={(e) => setFilters({ ...filters, action: e.target.value })}
          className="rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-gray-400 focus:outline-none"
        />
        <input
          type="text"
          placeholder="Resource type..."
          value={filters.resource_type}
          onChange={(e) => setFilters({ ...filters, resource_type: e.target.value })}
          className="rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-gray-400 focus:outline-none"
        />
        <input
          type="date"
          value={filters.from}
          onChange={(e) => setFilters({ ...filters, from: e.target.value })}
          className="rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-gray-400 focus:outline-none"
        />
        <input
          type="date"
          value={filters.to}
          onChange={(e) => setFilters({ ...filters, to: e.target.value })}
          className="rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-gray-400 focus:outline-none"
        />
      </div>

      {/* Expanded Detail */}
      {expandedId && (
        <AuditDetailPanel id={expandedId} onClose={() => setExpandedId(null)} />
      )}

      {/* Logs Table */}
      <div className="rounded-xl border border-gray-200 bg-white shadow-sm">
        {isLoading ? (
          <div className="flex items-center justify-center py-12">
            <div className="h-6 w-6 animate-spin rounded-full border-2 border-gray-300 border-t-gray-900" />
          </div>
        ) : logs.length === 0 ? (
          <div className="flex flex-col items-center justify-center py-16 text-gray-400">
            <Shield className="mb-3 h-10 w-10" />
            <p className="text-sm font-medium">No audit logs found</p>
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead>
                <tr className="border-b border-gray-100 text-left text-sm text-gray-500">
                  <th className="px-6 py-3 font-medium">User</th>
                  <th className="px-6 py-3 font-medium">Action</th>
                  <th className="px-6 py-3 font-medium">Resource</th>
                  <th className="px-6 py-3 font-medium">IP</th>
                  <th className="px-6 py-3 font-medium">Date</th>
                  <th className="px-6 py-3 font-medium text-right">Details</th>
                </tr>
              </thead>
              <tbody>
                {logs.map((log) => (
                  <tr key={log.id} className="border-b border-gray-50 last:border-0">
                    <td className="px-6 py-4 text-sm text-gray-900">{log.user_email}</td>
                    <td className="px-6 py-4">
                      <span className="rounded-full bg-gray-100 px-2.5 py-0.5 text-xs font-mono text-gray-700">
                        {log.action}
                      </span>
                    </td>
                    <td className="px-6 py-4 text-sm text-gray-600">{log.resource_type}</td>
                    <td className="px-6 py-4 text-sm font-mono text-gray-500">{log.ip_address}</td>
                    <td className="px-6 py-4 text-sm text-gray-500">
                      {new Date(log.created_at).toLocaleString()}
                    </td>
                    <td className="px-6 py-4 text-right">
                      <button
                        onClick={() => setExpandedId(expandedId === log.id ? null : log.id)}
                        className="rounded-md p-1.5 text-gray-400 hover:bg-gray-100 hover:text-gray-600 transition-colors"
                      >
                        {expandedId === log.id ? <ChevronUp className="h-4 w-4" /> : <ChevronDown className="h-4 w-4" />}
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </div>
  )
}
