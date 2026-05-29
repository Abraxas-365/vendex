import { useState, useRef } from 'react'
import { Upload, Image, FileIcon, X, Search } from 'lucide-react'
import type { Media as MediaType } from '../../types'
import { useMedia, useUploadMedia, useDeleteMedia } from '../../lib/hooks'

function formatFileSize(bytes: number): string {
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`
}

function formatDate(dateStr: string): string {
  return new Date(dateStr).toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
  })
}

function isImage(contentType: string): boolean {
  return contentType.startsWith('image/')
}

export default function Media() {
  const [searchQuery, setSearchQuery] = useState('')
  const [isDragging, setIsDragging] = useState(false)
  const fileInputRef = useRef<HTMLInputElement>(null)

  const { data, isLoading } = useMedia()
  const uploadMedia = useUploadMedia()
  const deleteMedia = useDeleteMedia()

  const items: MediaType[] = data?.items ?? []

  const filtered = items.filter(
    (m) =>
      searchQuery === '' ||
      m.filename.toLowerCase().includes(searchQuery.toLowerCase()),
  )

  function handleFiles(files: FileList | null) {
    if (!files) return
    for (let i = 0; i < files.length; i++) {
      const file = files[i]
      if (file) {
        uploadMedia.mutate({ file })
      }
    }
  }

  function handleDragOver(e: React.DragEvent) {
    e.preventDefault()
    setIsDragging(true)
  }

  function handleDragLeave(e: React.DragEvent) {
    e.preventDefault()
    setIsDragging(false)
  }

  function handleDrop(e: React.DragEvent) {
    e.preventDefault()
    setIsDragging(false)
    handleFiles(e.dataTransfer.files)
  }

  function handleDelete(item: MediaType) {
    if (confirm(`Delete "${item.filename}"?`)) {
      deleteMedia.mutate(item.id)
    }
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Media</h1>
          <p className="mt-1 text-sm text-gray-500">
            Upload and manage images and files
          </p>
        </div>
        <button
          onClick={() => fileInputRef.current?.click()}
          className="inline-flex items-center gap-2 rounded-lg bg-gray-900 px-4 py-2.5 text-sm font-medium text-white shadow-sm hover:bg-gray-800 transition-colors"
        >
          <Upload className="h-4 w-4" />
          Upload Files
        </button>
        <input
          ref={fileInputRef}
          type="file"
          multiple
          accept="image/*,application/pdf,.svg"
          onChange={(e) => handleFiles(e.target.files)}
          className="hidden"
        />
      </div>

      {/* Upload Dropzone */}
      <div
        onDragOver={handleDragOver}
        onDragLeave={handleDragLeave}
        onDrop={handleDrop}
        onClick={() => fileInputRef.current?.click()}
        className={`flex cursor-pointer flex-col items-center justify-center rounded-xl border-2 border-dashed p-8 transition-colors ${
          isDragging
            ? 'border-gray-900 bg-gray-50'
            : 'border-gray-200 bg-white hover:border-gray-400 hover:bg-gray-50'
        }`}
      >
        <Upload
          className={`mb-3 h-8 w-8 ${isDragging ? 'text-gray-900' : 'text-gray-400'}`}
        />
        <p className="text-sm font-medium text-gray-600">
          {isDragging ? 'Drop files here' : 'Drag and drop files here, or click to browse'}
        </p>
        <p className="mt-1 text-xs text-gray-400">PNG, JPG, SVG, PDF up to 10MB</p>
        {uploadMedia.isPending && (
          <div className="mt-3 flex items-center gap-2 text-sm text-gray-500">
            <div className="h-4 w-4 animate-spin rounded-full border-2 border-gray-300 border-t-gray-900" />
            Uploading...
          </div>
        )}
      </div>

      {/* Search */}
      <div className="relative">
        <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-400" />
        <input
          type="text"
          placeholder="Search files..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className="w-full rounded-lg border border-gray-200 bg-white py-2.5 pl-10 pr-4 text-sm text-gray-900 placeholder:text-gray-400 focus:border-gray-400 focus:outline-none focus:ring-1 focus:ring-gray-400"
        />
      </div>

      {/* Media Grid */}
      {isLoading ? (
        <div className="flex items-center justify-center py-12">
          <div className="h-6 w-6 animate-spin rounded-full border-2 border-gray-300 border-t-gray-900" />
        </div>
      ) : filtered.length === 0 ? (
        <div className="flex flex-col items-center justify-center rounded-xl border border-gray-200 bg-white py-16 text-gray-400">
          <Image className="mb-3 h-10 w-10" />
          <p className="text-sm font-medium">No media files found</p>
          <p className="mt-1 text-xs">Upload your first file to get started.</p>
        </div>
      ) : (
        <div className="grid grid-cols-2 gap-4 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5">
          {filtered.map((item) => (
            <div
              key={item.id}
              className="group relative overflow-hidden rounded-xl border border-gray-200 bg-white shadow-sm transition-shadow hover:shadow-md"
            >
              {/* Thumbnail */}
              <div className="aspect-square bg-gray-100">
                {isImage(item.content_type) ? (
                  <img
                    src={item.url}
                    alt={item.alt || item.filename}
                    className="h-full w-full object-cover"
                  />
                ) : (
                  <div className="flex h-full w-full items-center justify-center">
                    <FileIcon className="h-12 w-12 text-gray-300" />
                  </div>
                )}
              </div>

              {/* Delete overlay */}
              <div className="absolute inset-0 flex items-start justify-end p-2 opacity-0 group-hover:opacity-100 transition-opacity">
                <button
                  onClick={(e) => {
                    e.stopPropagation()
                    handleDelete(item)
                  }}
                  className="rounded-lg bg-white/90 p-1.5 text-gray-600 shadow-sm backdrop-blur-sm hover:bg-red-50 hover:text-red-600 transition-colors"
                  title="Delete"
                >
                  <X className="h-4 w-4" />
                </button>
              </div>

              {/* File info */}
              <div className="p-3">
                <p className="truncate text-sm font-medium text-gray-900" title={item.filename}>
                  {item.filename}
                </p>
                <div className="mt-1 flex items-center justify-between text-xs text-gray-400">
                  <span>{formatFileSize(item.size)}</span>
                  <span>{formatDate(item.created_at)}</span>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
