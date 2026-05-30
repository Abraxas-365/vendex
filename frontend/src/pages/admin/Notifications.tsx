import { useState } from 'react'
import { Bell, CheckCheck, Trash2, Check } from 'lucide-react'
import type { NotificationSeverity } from '../../types'
import {
  useNotifications,
  useUnreadCount,
  useMarkNotificationRead,
  useMarkAllNotificationsRead,
  useDeleteNotification,
} from '../../lib/hooks'

const typeColors: Record<NotificationSeverity, string> = {
  info: 'bg-blue-100 text-blue-800',
  warning: 'bg-yellow-100 text-yellow-800',
  error: 'bg-red-100 text-red-800',
  success: 'bg-green-100 text-green-800',
}

export default function Notifications() {
  const [readFilter, setReadFilter] = useState<'all' | 'unread' | 'read'>('all')

  const queryParams: { read?: boolean } = {}
  if (readFilter === 'unread') queryParams.read = false
  if (readFilter === 'read') queryParams.read = true

  const { data, isLoading } = useNotifications(queryParams)
  const { data: unreadData } = useUnreadCount()
  const notifications = data?.items ?? []

  const markRead = useMarkNotificationRead()
  const markAllRead = useMarkAllNotificationsRead()
  const deleteNotif = useDeleteNotification()

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div>
            <h1 className="text-2xl font-bold text-gray-900">Notifications</h1>
            <p className="mt-1 text-sm text-gray-500">In-app admin notifications</p>
          </div>
          {unreadData && unreadData.count > 0 && (
            <span className="inline-flex items-center rounded-full bg-red-100 px-2.5 py-0.5 text-xs font-semibold text-red-800">
              {unreadData.count} unread
            </span>
          )}
        </div>
        <button
          onClick={() => markAllRead.mutate()}
          className="inline-flex items-center gap-2 rounded-lg border border-gray-200 bg-white px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50 transition-colors"
        >
          <CheckCheck className="h-4 w-4" />
          Mark All Read
        </button>
      </div>

      {/* Filter tabs */}
      <div className="flex gap-1 border-b border-gray-200">
        {(['all', 'unread', 'read'] as const).map((f) => (
          <button
            key={f}
            onClick={() => setReadFilter(f)}
            className={`px-4 py-2 text-sm font-medium border-b-2 transition-colors capitalize ${
              readFilter === f
                ? 'border-gray-900 text-gray-900'
                : 'border-transparent text-gray-500 hover:text-gray-700'
            }`}
          >
            {f}
          </button>
        ))}
      </div>

      {/* Notifications List */}
      <div className="rounded-xl border border-gray-200 bg-white shadow-sm">
        {isLoading ? (
          <div className="flex items-center justify-center py-12">
            <div className="h-6 w-6 animate-spin rounded-full border-2 border-gray-300 border-t-gray-900" />
          </div>
        ) : notifications.length === 0 ? (
          <div className="flex flex-col items-center justify-center py-16 text-gray-400">
            <Bell className="mb-3 h-10 w-10" />
            <p className="text-sm font-medium">No notifications found</p>
          </div>
        ) : (
          <div className="divide-y divide-gray-50">
            {notifications.map((n) => (
              <div
                key={n.id}
                className={`flex items-start gap-4 px-6 py-4 transition-colors ${!n.read ? 'bg-blue-50/30' : ''}`}
              >
                <div className="flex-1 min-w-0">
                  <div className="flex items-center gap-2">
                    <span className={`inline-flex rounded-full px-2 py-0.5 text-xs font-medium ${typeColors[n.severity]}`}>
                      {n.severity}
                    </span>
                    {!n.read && (
                      <span className="h-2 w-2 rounded-full bg-blue-500" />
                    )}
                  </div>
                  <p className="mt-1 text-sm font-medium text-gray-900">{n.title}</p>
                  <p className="mt-0.5 text-sm text-gray-600">{n.body}</p>
                  {n.resource_type && (
                    <p className="mt-1 text-xs text-gray-400">{n.resource_type}: {n.resource_id}</p>
                  )}
                  <p className="mt-1 text-xs text-gray-400">{new Date(n.created_at).toLocaleString()}</p>
                </div>
                <div className="flex items-center gap-2 shrink-0">
                  {!n.read && (
                    <button
                      onClick={() => markRead.mutate(n.id)}
                      className="rounded-md p-1.5 text-gray-400 hover:bg-green-50 hover:text-green-600 transition-colors"
                      title="Mark as read"
                    >
                      <Check className="h-4 w-4" />
                    </button>
                  )}
                  <button
                    onClick={() => deleteNotif.mutate(n.id)}
                    className="rounded-md p-1.5 text-gray-400 hover:bg-red-50 hover:text-red-600 transition-colors"
                    title="Delete"
                  >
                    <Trash2 className="h-4 w-4" />
                  </button>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}
