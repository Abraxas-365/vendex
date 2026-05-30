import { useState } from 'react'
import { Users } from 'lucide-react'
import type { SocialProvider } from '../../types'
import { useSocialAccounts } from '../../lib/hooks'

const providerColors: Record<SocialProvider, string> = {
  google: 'bg-red-100 text-red-800',
  facebook: 'bg-blue-100 text-blue-800',
}

export default function SocialAccounts() {
  const [providerFilter, setProviderFilter] = useState<SocialProvider | ''>('')

  const { data, isLoading } = useSocialAccounts(
    providerFilter ? { provider: providerFilter } : undefined
  )
  const accounts = data?.items ?? []

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-2xl font-bold text-gray-900">Social Login Accounts</h1>
        <p className="mt-1 text-sm text-gray-500">Linked social provider accounts across customers</p>
      </div>

      {/* Filter */}
      <div className="flex items-center gap-3">
        <label className="text-sm font-medium text-gray-700">Provider:</label>
        <select
          value={providerFilter}
          onChange={(e) => setProviderFilter(e.target.value as SocialProvider | '')}
          className="rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-gray-400 focus:outline-none"
        >
          <option value="">All</option>
          <option value="google">Google</option>
          <option value="facebook">Facebook</option>
        </select>
      </div>

      {/* Stats */}
      {data && (
        <div className="flex items-center gap-2 text-sm text-gray-500">
          <span>Showing {accounts.length} of {data.total} accounts</span>
        </div>
      )}

      {/* Table */}
      <div className="rounded-xl border border-gray-200 bg-white shadow-sm">
        {isLoading ? (
          <div className="flex items-center justify-center py-12">
            <div className="h-6 w-6 animate-spin rounded-full border-2 border-gray-300 border-t-gray-900" />
          </div>
        ) : accounts.length === 0 ? (
          <div className="flex flex-col items-center justify-center py-16 text-gray-400">
            <Users className="mb-3 h-10 w-10" />
            <p className="text-sm font-medium">No social accounts found</p>
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead>
                <tr className="border-b border-gray-100 text-left text-sm text-gray-500">
                  <th className="px-6 py-3 font-medium">Customer</th>
                  <th className="px-6 py-3 font-medium">Provider</th>
                  <th className="px-6 py-3 font-medium">Provider User ID</th>
                  <th className="px-6 py-3 font-medium">Linked</th>
                </tr>
              </thead>
              <tbody>
                {accounts.map((acc) => (
                  <tr key={acc.id} className="border-b border-gray-50 last:border-0">
                    <td className="px-6 py-4">
                      <div className="text-sm font-medium text-gray-900">{acc.customer_name}</div>
                      <div className="text-xs text-gray-500">{acc.customer_id}</div>
                    </td>
                    <td className="px-6 py-4">
                      <span className={`inline-flex rounded-full px-2.5 py-0.5 text-xs font-medium ${providerColors[acc.provider]}`}>
                        {acc.provider}
                      </span>
                    </td>
                    <td className="px-6 py-4 text-sm font-mono text-gray-500">{acc.provider_user_id}</td>
                    <td className="px-6 py-4 text-sm text-gray-500">
                      {new Date(acc.created_at).toLocaleDateString()}
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
