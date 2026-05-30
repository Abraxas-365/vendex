import { useState, useRef } from 'react'
import { Download, Upload, FileText, AlertCircle, CheckCircle, X } from 'lucide-react'
import type { ImportResult } from '../../types'
import { exportProducts, exportOrders, exportCustomers, importProducts } from '../../lib/api'

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

function triggerDownload(blob: Blob, filename: string) {
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = filename
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
  URL.revokeObjectURL(url)
}

// ---------------------------------------------------------------------------
// Export section
// ---------------------------------------------------------------------------

interface ExportCardProps {
  label: string
  description: string
  filename: string
  onExport: () => Promise<Blob>
}

function ExportCard({ label, description, filename, onExport }: ExportCardProps) {
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [success, setSuccess] = useState(false)

  async function handleExport() {
    setLoading(true)
    setError(null)
    setSuccess(false)
    try {
      const blob = await onExport()
      triggerDownload(blob, filename)
      setSuccess(true)
      setTimeout(() => setSuccess(false), 3000)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Export failed')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="rounded-xl border border-slate-200 bg-white p-5">
      <div className="flex items-start justify-between gap-4">
        <div className="flex items-start gap-3">
          <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-indigo-50">
            <FileText size={18} className="text-indigo-600" />
          </div>
          <div>
            <p className="text-sm font-semibold text-slate-800">{label}</p>
            <p className="mt-0.5 text-xs text-slate-500">{description}</p>
          </div>
        </div>
        <button
          onClick={() => void handleExport()}
          disabled={loading}
          className="flex shrink-0 items-center gap-2 rounded-lg bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-700 disabled:opacity-60"
        >
          <Download size={15} />
          {loading ? 'Exporting…' : 'Export CSV'}
        </button>
      </div>
      {error && (
        <div className="mt-3 flex items-center gap-2 rounded-lg bg-red-50 px-3 py-2 text-sm text-red-700">
          <AlertCircle size={15} className="shrink-0" />
          {error}
        </div>
      )}
      {success && (
        <div className="mt-3 flex items-center gap-2 rounded-lg bg-green-50 px-3 py-2 text-sm text-green-700">
          <CheckCircle size={15} className="shrink-0" />
          Download started successfully.
        </div>
      )}
    </div>
  )
}

// ---------------------------------------------------------------------------
// Import section
// ---------------------------------------------------------------------------

function ImportSection() {
  const fileInputRef = useRef<HTMLInputElement>(null)
  const [selectedFile, setSelectedFile] = useState<File | null>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [result, setResult] = useState<ImportResult | null>(null)

  function handleFileChange(e: React.ChangeEvent<HTMLInputElement>) {
    const file = e.target.files?.[0] ?? null
    setSelectedFile(file)
    setResult(null)
    setError(null)
  }

  async function handleImport() {
    if (!selectedFile) return
    setLoading(true)
    setError(null)
    setResult(null)
    try {
      const res = await importProducts(selectedFile)
      setResult(res)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Import failed')
    } finally {
      setLoading(false)
    }
  }

  function handleClear() {
    setSelectedFile(null)
    setResult(null)
    setError(null)
    if (fileInputRef.current) fileInputRef.current.value = ''
  }

  return (
    <div className="rounded-xl border border-slate-200 bg-white p-6 space-y-5">
      <div>
        <p className="text-sm font-semibold text-slate-800">Import Products</p>
        <p className="mt-0.5 text-xs text-slate-500">
          Upload a CSV file to bulk-import products. The first row must be a header row.
        </p>
      </div>

      {/* File picker */}
      <div className="flex items-center gap-3">
        <label className="flex cursor-pointer items-center gap-2 rounded-lg border border-dashed border-slate-300 bg-slate-50 px-4 py-3 text-sm text-slate-600 hover:bg-slate-100 transition-colors">
          <Upload size={16} className="text-slate-400" />
          <span>{selectedFile ? selectedFile.name : 'Choose CSV file…'}</span>
          <input
            ref={fileInputRef}
            type="file"
            accept=".csv"
            className="sr-only"
            onChange={handleFileChange}
          />
        </label>
        {selectedFile && (
          <button onClick={handleClear} className="text-slate-400 hover:text-slate-600">
            <X size={16} />
          </button>
        )}
      </div>

      {/* Import button */}
      <button
        onClick={() => void handleImport()}
        disabled={!selectedFile || loading}
        className="flex items-center gap-2 rounded-lg bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-700 disabled:opacity-50"
      >
        <Upload size={15} />
        {loading ? 'Importing…' : 'Import Products'}
      </button>

      {/* Error */}
      {error && (
        <div className="flex items-center gap-2 rounded-lg bg-red-50 px-3 py-2 text-sm text-red-700">
          <AlertCircle size={15} className="shrink-0" />
          {error}
        </div>
      )}

      {/* Results */}
      {result && (
        <div className="space-y-3">
          <div className="flex items-center gap-4 rounded-lg bg-green-50 px-4 py-3">
            <CheckCircle size={18} className="shrink-0 text-green-600" />
            <div className="text-sm">
              <span className="font-semibold text-green-800">{result.imported}</span>
              <span className="text-green-700"> of </span>
              <span className="font-semibold text-green-800">{result.total}</span>
              <span className="text-green-700"> rows imported successfully.</span>
            </div>
          </div>

          {result.errors.length > 0 && (
            <div className="space-y-2">
              <p className="text-xs font-semibold uppercase tracking-wide text-red-600">
                {result.errors.length} error{result.errors.length !== 1 ? 's' : ''}
              </p>
              <div className="overflow-hidden rounded-lg border border-red-200">
                <table className="w-full text-sm">
                  <thead className="bg-red-50 text-xs font-semibold uppercase tracking-wide text-red-600">
                    <tr>
                      <th className="px-4 py-2.5 text-left w-20">Row</th>
                      <th className="px-4 py-2.5 text-left">Error</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-red-100">
                    {result.errors.map((e, i) => (
                      <tr key={i} className="bg-white">
                        <td className="px-4 py-2.5 text-slate-500">{e.row}</td>
                        <td className="px-4 py-2.5 text-red-700">{e.error}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </div>
          )}
        </div>
      )}
    </div>
  )
}

// ---------------------------------------------------------------------------
// Main page
// ---------------------------------------------------------------------------

export default function ImportExport() {
  return (
    <div className="space-y-8">
      {/* Header */}
      <div>
        <h1 className="text-xl font-semibold text-slate-800">Import / Export</h1>
        <p className="mt-0.5 text-sm text-slate-500">
          Bulk-export store data to CSV or import products from a file.
        </p>
      </div>

      {/* Export section */}
      <section className="space-y-3">
        <h2 className="text-sm font-semibold uppercase tracking-widest text-slate-400">Export</h2>
        <div className="grid gap-3 sm:grid-cols-1 lg:grid-cols-3">
          <ExportCard
            label="Products"
            description="Export all products with pricing, stock, and metadata."
            filename="products-export.csv"
            onExport={exportProducts}
          />
          <ExportCard
            label="Orders"
            description="Export all orders with customer and line-item details."
            filename="orders-export.csv"
            onExport={exportOrders}
          />
          <ExportCard
            label="Customers"
            description="Export all customers with contact information."
            filename="customers-export.csv"
            onExport={exportCustomers}
          />
        </div>
      </section>

      {/* Import section */}
      <section className="space-y-3">
        <h2 className="text-sm font-semibold uppercase tracking-widest text-slate-400">Import</h2>
        <ImportSection />
      </section>
    </div>
  )
}
