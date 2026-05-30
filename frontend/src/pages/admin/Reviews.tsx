import { useState } from 'react'
import { Star, X, MessageSquare } from 'lucide-react'
import type { ReviewStatus } from '../../types'
import { useReviews, useApproveReview, useRejectReview, useRespondToReview } from '../../lib/hooks'

const statusColors: Record<ReviewStatus, string> = {
  pending: 'bg-yellow-100 text-yellow-800',
  approved: 'bg-green-100 text-green-800',
  rejected: 'bg-red-100 text-red-800',
}

export default function Reviews() {
  const [statusFilter, setStatusFilter] = useState<ReviewStatus | ''>('')
  const [respondingId, setRespondingId] = useState<string | null>(null)
  const [responseText, setResponseText] = useState('')

  const { data, isLoading } = useReviews(statusFilter ? { status: statusFilter } : undefined)
  const approveReview = useApproveReview()
  const rejectReview = useRejectReview()
  const respondToReview = useRespondToReview()

  const reviews = data?.items ?? []

  function handleRespond(e: React.FormEvent) {
    e.preventDefault()
    if (!respondingId) return
    respondToReview.mutate({ id: respondingId, response: responseText })
    setRespondingId(null)
    setResponseText('')
  }

  function renderStars(rating: number) {
    return (
      <div className="flex items-center gap-0.5">
        {[1, 2, 3, 4, 5].map((n) => (
          <Star
            key={n}
            className={`h-3.5 w-3.5 ${n <= rating ? 'fill-yellow-400 text-yellow-400' : 'text-gray-300'}`}
          />
        ))}
      </div>
    )
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Reviews</h1>
          <p className="mt-1 text-sm text-gray-500">Moderate customer product reviews</p>
        </div>
      </div>

      {/* Filter */}
      <div className="flex items-center gap-3">
        <label className="text-sm font-medium text-gray-700">Status:</label>
        <select
          value={statusFilter}
          onChange={(e) => setStatusFilter(e.target.value as ReviewStatus | '')}
          className="rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-gray-400 focus:outline-none"
        >
          <option value="">All</option>
          <option value="pending">Pending</option>
          <option value="approved">Approved</option>
          <option value="rejected">Rejected</option>
        </select>
      </div>

      {/* Reviews Table */}
      <div className="rounded-xl border border-gray-200 bg-white shadow-sm">
        {isLoading ? (
          <div className="flex items-center justify-center py-12">
            <div className="h-6 w-6 animate-spin rounded-full border-2 border-gray-300 border-t-gray-900" />
          </div>
        ) : reviews.length === 0 ? (
          <div className="flex flex-col items-center justify-center py-16 text-gray-400">
            <Star className="mb-3 h-10 w-10" />
            <p className="text-sm font-medium">No reviews found</p>
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead>
                <tr className="border-b border-gray-100 text-left text-sm text-gray-500">
                  <th className="px-6 py-3 font-medium">Customer</th>
                  <th className="px-6 py-3 font-medium">Product</th>
                  <th className="px-6 py-3 font-medium">Rating</th>
                  <th className="px-6 py-3 font-medium">Review</th>
                  <th className="px-6 py-3 font-medium">Status</th>
                  <th className="px-6 py-3 font-medium">Date</th>
                  <th className="px-6 py-3 font-medium text-right">Actions</th>
                </tr>
              </thead>
              <tbody>
                {reviews.map((review) => (
                  <tr key={review.id} className="border-b border-gray-50 last:border-0">
                    <td className="px-6 py-4">
                      <div className="text-sm font-medium text-gray-900">{review.customer_name}</div>
                    </td>
                    <td className="px-6 py-4 text-sm text-gray-600">{review.product_name}</td>
                    <td className="px-6 py-4">{renderStars(review.rating)}</td>
                    <td className="px-6 py-4">
                      <div className="text-sm font-medium text-gray-900">{review.title}</div>
                      <div className="mt-0.5 max-w-xs text-xs text-gray-500 line-clamp-2">{review.body}</div>
                      {review.admin_response && (
                        <div className="mt-1 max-w-xs text-xs text-indigo-600 italic line-clamp-1">
                          Response: {review.admin_response}
                        </div>
                      )}
                    </td>
                    <td className="px-6 py-4">
                      <span className={`inline-flex rounded-full px-2.5 py-0.5 text-xs font-medium ${statusColors[review.status]}`}>
                        {review.status}
                      </span>
                    </td>
                    <td className="px-6 py-4 text-sm text-gray-500">
                      {new Date(review.created_at).toLocaleDateString()}
                    </td>
                    <td className="px-6 py-4 text-right">
                      <div className="flex items-center justify-end gap-2">
                        {review.status === 'pending' && (
                          <>
                            <button
                              onClick={() => approveReview.mutate(review.id)}
                              className="rounded-lg border border-green-200 px-2.5 py-1 text-xs font-medium text-green-700 hover:bg-green-50 transition-colors"
                            >
                              Approve
                            </button>
                            <button
                              onClick={() => rejectReview.mutate(review.id)}
                              className="rounded-lg border border-red-200 px-2.5 py-1 text-xs font-medium text-red-700 hover:bg-red-50 transition-colors"
                            >
                              Reject
                            </button>
                          </>
                        )}
                        <button
                          onClick={() => { setRespondingId(review.id); setResponseText(review.admin_response ?? '') }}
                          className="rounded-md p-1.5 text-gray-400 hover:bg-gray-100 hover:text-gray-600 transition-colors"
                          title="Respond"
                        >
                          <MessageSquare className="h-4 w-4" />
                        </button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>

      {/* Response Dialog */}
      {respondingId && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40">
          <div className="w-full max-w-md rounded-xl bg-white p-6 shadow-xl">
            <div className="mb-4 flex items-center justify-between">
              <h2 className="text-lg font-semibold text-gray-900">Admin Response</h2>
              <button onClick={() => setRespondingId(null)} className="text-gray-400 hover:text-gray-600">
                <X className="h-5 w-5" />
              </button>
            </div>
            <form onSubmit={handleRespond} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700">Response</label>
                <textarea
                  rows={4}
                  required
                  value={responseText}
                  onChange={(e) => setResponseText(e.target.value)}
                  placeholder="Write your response to this review..."
                  className="mt-1 w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-gray-400 focus:outline-none"
                />
              </div>
              <div className="flex gap-3">
                <button
                  type="submit"
                  className="flex-1 rounded-lg bg-gray-900 py-2.5 text-sm font-medium text-white hover:bg-gray-800 transition-colors"
                >
                  Save Response
                </button>
                <button
                  type="button"
                  onClick={() => setRespondingId(null)}
                  className="flex-1 rounded-lg border border-gray-200 bg-white py-2.5 text-sm font-medium text-gray-700 hover:bg-gray-50 transition-colors"
                >
                  Cancel
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  )
}
