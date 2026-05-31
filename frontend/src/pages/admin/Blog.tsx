import { useState } from 'react'
import { Newspaper, Plus, Trash2, Edit2, Send, Archive } from 'lucide-react'
import type { BlogPost, BlogCategory } from '../../types'
import {
  useBlogPosts,
  useCreateBlogPost,
  useUpdateBlogPost,
  useDeleteBlogPost,
  usePublishBlogPost,
  useArchiveBlogPost,
  useBlogCategories,
  useCreateBlogCategory,
  useDeleteBlogCategory,
} from '../../lib/hooks'

const emptyPost: Partial<BlogPost> = {
  title: '',
  slug: '',
  content: '',
  excerpt: '',
  author: '',
  tags: [],
  seo_title: '',
  seo_description: '',
}

const statusColors: Record<string, string> = {
  draft: 'bg-yellow-100 text-yellow-700',
  published: 'bg-green-100 text-green-700',
  archived: 'bg-gray-100 text-gray-600',
}

type Tab = 'posts' | 'categories'

export default function Blog() {
  const [tab, setTab] = useState<Tab>('posts')
  const [page] = useState(1)
  const [showPostDialog, setShowPostDialog] = useState(false)
  const [editingPost, setEditingPost] = useState<BlogPost | null>(null)
  const [postForm, setPostForm] = useState<Partial<BlogPost>>(emptyPost)
  const [tagsInput, setTagsInput] = useState('')

  const [showCatDialog, setShowCatDialog] = useState(false)
  const [catForm, setCatForm] = useState<Partial<BlogCategory>>({ name: '', slug: '', description: '' })

  const { data, isLoading } = useBlogPosts({ page, page_size: 20 })
  const posts = data?.items ?? []

  const { data: categories, isLoading: catsLoading } = useBlogCategories()

  const createPost = useCreateBlogPost()
  const updatePost = useUpdateBlogPost()
  const deletePost = useDeleteBlogPost()
  const publishPost = usePublishBlogPost()
  const archivePost = useArchiveBlogPost()

  const createCat = useCreateBlogCategory()
  const deleteCat = useDeleteBlogCategory()

  function openCreatePost() {
    setEditingPost(null)
    setPostForm(emptyPost)
    setTagsInput('')
    setShowPostDialog(true)
  }

  function openEditPost(p: BlogPost) {
    setEditingPost(p)
    setPostForm({
      title: p.title,
      slug: p.slug,
      content: p.content,
      excerpt: p.excerpt,
      author: p.author,
      seo_title: p.seo_title,
      seo_description: p.seo_description,
    })
    setTagsInput((p.tags ?? []).join(', '))
    setShowPostDialog(true)
  }

  function handlePostSubmit(e: React.FormEvent) {
    e.preventDefault()
    const tags = tagsInput
      .split(',')
      .map((t) => t.trim())
      .filter(Boolean)
    const payload = { ...postForm, tags }

    if (editingPost) {
      updatePost.mutate(
        { id: editingPost.id, data: payload },
        { onSuccess: () => setShowPostDialog(false) },
      )
    } else {
      createPost.mutate(payload, { onSuccess: () => setShowPostDialog(false) })
    }
  }

  function handleCatSubmit(e: React.FormEvent) {
    e.preventDefault()
    createCat.mutate(catForm, {
      onSuccess: () => {
        setShowCatDialog(false)
        setCatForm({ name: '', slug: '', description: '' })
      },
    })
  }

  const postPending = createPost.isPending || updatePost.isPending

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <Newspaper className="h-6 w-6 text-indigo-600" />
          <div>
            <h1 className="text-2xl font-bold text-gray-900">Blog</h1>
            <p className="mt-0.5 text-sm text-gray-500">Manage blog posts and categories</p>
          </div>
        </div>
        <button
          onClick={tab === 'posts' ? openCreatePost : () => setShowCatDialog(true)}
          className="inline-flex items-center gap-2 rounded-lg bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-700 transition-colors"
        >
          <Plus className="h-4 w-4" />
          {tab === 'posts' ? 'New Post' : 'New Category'}
        </button>
      </div>

      {/* Tabs */}
      <div className="flex gap-1 border-b border-gray-200">
        {(['posts', 'categories'] as const).map((t) => (
          <button
            key={t}
            onClick={() => setTab(t)}
            className={`px-4 py-2 text-sm font-medium border-b-2 transition-colors capitalize ${
              tab === t
                ? 'border-indigo-600 text-indigo-600'
                : 'border-transparent text-gray-500 hover:text-gray-700'
            }`}
          >
            {t}
          </button>
        ))}
      </div>

      {tab === 'posts' && (
        <div className="rounded-xl border border-gray-200 bg-white shadow-sm overflow-hidden">
          {isLoading ? (
            <div className="flex items-center justify-center py-12">
              <div className="h-6 w-6 animate-spin rounded-full border-2 border-gray-300 border-t-indigo-600" />
            </div>
          ) : posts.length === 0 ? (
            <div className="flex flex-col items-center justify-center py-16 text-gray-400">
              <Newspaper className="mb-3 h-10 w-10" />
              <p className="text-sm font-medium">No blog posts yet</p>
            </div>
          ) : (
            <table className="w-full text-sm">
              <thead className="border-b border-gray-100 bg-gray-50">
                <tr>
                  <th className="px-6 py-3 text-left font-medium text-gray-500">Title</th>
                  <th className="px-6 py-3 text-left font-medium text-gray-500">Author</th>
                  <th className="px-6 py-3 text-left font-medium text-gray-500">Status</th>
                  <th className="px-6 py-3 text-left font-medium text-gray-500">Published</th>
                  <th className="px-6 py-3 text-right font-medium text-gray-500">Actions</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-50">
                {posts.map((p) => (
                  <tr key={p.id} className="hover:bg-gray-50/50 transition-colors">
                    <td className="px-6 py-4">
                      <p className="font-medium text-gray-900">{p.title}</p>
                      <p className="text-xs text-gray-400 mt-0.5">{p.slug}</p>
                    </td>
                    <td className="px-6 py-4 text-gray-600">{p.author}</td>
                    <td className="px-6 py-4">
                      <span className={`inline-flex rounded-full px-2 py-0.5 text-xs font-medium ${statusColors[p.status] ?? 'bg-gray-100 text-gray-600'}`}>
                        {p.status}
                      </span>
                    </td>
                    <td className="px-6 py-4 text-gray-500 text-xs">
                      {p.published_at ? new Date(p.published_at).toLocaleDateString() : '—'}
                    </td>
                    <td className="px-6 py-4">
                      <div className="flex items-center justify-end gap-1">
                        {p.status === 'draft' && (
                          <button
                            onClick={() => publishPost.mutate(p.id)}
                            className="rounded-md p-1.5 text-gray-400 hover:bg-green-50 hover:text-green-600 transition-colors"
                            title="Publish"
                          >
                            <Send className="h-4 w-4" />
                          </button>
                        )}
                        {p.status === 'published' && (
                          <button
                            onClick={() => archivePost.mutate(p.id)}
                            className="rounded-md p-1.5 text-gray-400 hover:bg-orange-50 hover:text-orange-600 transition-colors"
                            title="Archive"
                          >
                            <Archive className="h-4 w-4" />
                          </button>
                        )}
                        <button
                          onClick={() => openEditPost(p)}
                          className="rounded-md p-1.5 text-gray-400 hover:bg-blue-50 hover:text-blue-600 transition-colors"
                          title="Edit"
                        >
                          <Edit2 className="h-4 w-4" />
                        </button>
                        <button
                          onClick={() => deletePost.mutate(p.id)}
                          className="rounded-md p-1.5 text-gray-400 hover:bg-red-50 hover:text-red-600 transition-colors"
                          title="Delete"
                        >
                          <Trash2 className="h-4 w-4" />
                        </button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </div>
      )}

      {tab === 'categories' && (
        <div className="rounded-xl border border-gray-200 bg-white shadow-sm overflow-hidden">
          {catsLoading ? (
            <div className="flex items-center justify-center py-12">
              <div className="h-6 w-6 animate-spin rounded-full border-2 border-gray-300 border-t-indigo-600" />
            </div>
          ) : (categories ?? []).length === 0 ? (
            <div className="flex flex-col items-center justify-center py-16 text-gray-400">
              <p className="text-sm font-medium">No categories yet</p>
            </div>
          ) : (
            <table className="w-full text-sm">
              <thead className="border-b border-gray-100 bg-gray-50">
                <tr>
                  <th className="px-6 py-3 text-left font-medium text-gray-500">Name</th>
                  <th className="px-6 py-3 text-left font-medium text-gray-500">Slug</th>
                  <th className="px-6 py-3 text-left font-medium text-gray-500">Description</th>
                  <th className="px-6 py-3 text-right font-medium text-gray-500">Actions</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-50">
                {(categories ?? []).map((c) => (
                  <tr key={c.id} className="hover:bg-gray-50/50 transition-colors">
                    <td className="px-6 py-4 font-medium text-gray-900">{c.name}</td>
                    <td className="px-6 py-4 text-gray-600">{c.slug}</td>
                    <td className="px-6 py-4 text-gray-500">{c.description || '—'}</td>
                    <td className="px-6 py-4 text-right">
                      <button
                        onClick={() => deleteCat.mutate(c.id)}
                        className="rounded-md p-1.5 text-gray-400 hover:bg-red-50 hover:text-red-600 transition-colors"
                        title="Delete"
                      >
                        <Trash2 className="h-4 w-4" />
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </div>
      )}

      {/* Post Dialog */}
      {showPostDialog && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4">
          <div className="w-full max-w-2xl rounded-2xl bg-white shadow-2xl max-h-[90vh] overflow-y-auto">
            <div className="border-b border-gray-100 px-6 py-4">
              <h2 className="text-lg font-semibold text-gray-900">
                {editingPost ? 'Edit Post' : 'New Post'}
              </h2>
            </div>
            <form onSubmit={handlePostSubmit} className="space-y-4 px-6 py-5">
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-xs font-medium text-gray-700 mb-1">Title *</label>
                  <input
                    required
                    value={postForm.title ?? ''}
                    onChange={(e) => setPostForm((f) => ({ ...f, title: e.target.value }))}
                    className="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
                  />
                </div>
                <div>
                  <label className="block text-xs font-medium text-gray-700 mb-1">Slug *</label>
                  <input
                    required
                    value={postForm.slug ?? ''}
                    onChange={(e) => setPostForm((f) => ({ ...f, slug: e.target.value }))}
                    className="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
                  />
                </div>
              </div>
              <div>
                <label className="block text-xs font-medium text-gray-700 mb-1">Author</label>
                <input
                  value={postForm.author ?? ''}
                  onChange={(e) => setPostForm((f) => ({ ...f, author: e.target.value }))}
                  className="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
                />
              </div>
              <div>
                <label className="block text-xs font-medium text-gray-700 mb-1">Excerpt</label>
                <textarea
                  rows={2}
                  value={postForm.excerpt ?? ''}
                  onChange={(e) => setPostForm((f) => ({ ...f, excerpt: e.target.value }))}
                  className="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
                />
              </div>
              <div>
                <label className="block text-xs font-medium text-gray-700 mb-1">Content</label>
                <textarea
                  rows={6}
                  value={postForm.content ?? ''}
                  onChange={(e) => setPostForm((f) => ({ ...f, content: e.target.value }))}
                  className="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500 font-mono"
                />
              </div>
              <div>
                <label className="block text-xs font-medium text-gray-700 mb-1">Tags (comma separated)</label>
                <input
                  value={tagsInput}
                  onChange={(e) => setTagsInput(e.target.value)}
                  placeholder="tag1, tag2, tag3"
                  className="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
                />
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-xs font-medium text-gray-700 mb-1">SEO Title</label>
                  <input
                    value={postForm.seo_title ?? ''}
                    onChange={(e) => setPostForm((f) => ({ ...f, seo_title: e.target.value }))}
                    className="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
                  />
                </div>
                <div>
                  <label className="block text-xs font-medium text-gray-700 mb-1">SEO Description</label>
                  <input
                    value={postForm.seo_description ?? ''}
                    onChange={(e) => setPostForm((f) => ({ ...f, seo_description: e.target.value }))}
                    className="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
                  />
                </div>
              </div>
              {(createPost.error ?? updatePost.error) && (
                <p className="text-sm text-red-600">{(createPost.error ?? updatePost.error)?.message}</p>
              )}
              <div className="flex justify-end gap-3 pt-2">
                <button
                  type="button"
                  onClick={() => setShowPostDialog(false)}
                  className="rounded-lg border border-gray-200 px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50 transition-colors"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  disabled={postPending}
                  className="rounded-lg bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-700 disabled:opacity-50 transition-colors"
                >
                  {postPending ? 'Saving…' : editingPost ? 'Save Changes' : 'Create'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Category Dialog */}
      {showCatDialog && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4">
          <div className="w-full max-w-md rounded-2xl bg-white shadow-2xl">
            <div className="border-b border-gray-100 px-6 py-4">
              <h2 className="text-lg font-semibold text-gray-900">New Category</h2>
            </div>
            <form onSubmit={handleCatSubmit} className="space-y-4 px-6 py-5">
              <div>
                <label className="block text-xs font-medium text-gray-700 mb-1">Name *</label>
                <input
                  required
                  value={catForm.name ?? ''}
                  onChange={(e) => setCatForm((f) => ({ ...f, name: e.target.value }))}
                  className="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
                />
              </div>
              <div>
                <label className="block text-xs font-medium text-gray-700 mb-1">Slug *</label>
                <input
                  required
                  value={catForm.slug ?? ''}
                  onChange={(e) => setCatForm((f) => ({ ...f, slug: e.target.value }))}
                  className="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
                />
              </div>
              <div>
                <label className="block text-xs font-medium text-gray-700 mb-1">Description</label>
                <textarea
                  rows={2}
                  value={catForm.description ?? ''}
                  onChange={(e) => setCatForm((f) => ({ ...f, description: e.target.value }))}
                  className="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
                />
              </div>
              {createCat.error && (
                <p className="text-sm text-red-600">{createCat.error.message}</p>
              )}
              <div className="flex justify-end gap-3 pt-2">
                <button
                  type="button"
                  onClick={() => setShowCatDialog(false)}
                  className="rounded-lg border border-gray-200 px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50 transition-colors"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  disabled={createCat.isPending}
                  className="rounded-lg bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-700 disabled:opacity-50 transition-colors"
                >
                  {createCat.isPending ? 'Creating…' : 'Create'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  )
}
