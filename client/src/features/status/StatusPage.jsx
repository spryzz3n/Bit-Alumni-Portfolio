import { useEffect, useState } from 'react'
import { apiGet } from '../../lib/api'

// state is one of: { kind: 'loading' } | { kind: 'ok', data } | { kind: 'error', message }
function useApi(path) {
  const [state, setState] = useState({ kind: 'loading' })
  useEffect(() => {
    let cancelled = false
    apiGet(path)
      .then((data) => !cancelled && setState({ kind: 'ok', data }))
      .catch((e) => !cancelled && setState({ kind: 'error', message: e.message }))
    return () => {
      cancelled = true
    }
  }, [path])
  return state
}

function Card({ title, state }) {
  return (
    <section className="rounded-lg border border-slate-200 bg-white p-4 shadow-sm">
      <h2 className="text-sm font-medium text-slate-500">{title}</h2>
      {state.kind === 'loading' && <p className="mt-2 text-slate-500">Loading…</p>}
      {state.kind === 'error' && <p className="mt-2 text-red-600">Error: {state.message}</p>}
      {state.kind === 'ok' && (
        <pre className="mt-2 overflow-x-auto text-sm text-emerald-700">
          {JSON.stringify(state.data, null, 2)}
        </pre>
      )}
    </section>
  )
}

export default function StatusPage() {
  const health = useApi('/health')
  const db = useApi('/db-check')
  return (
    <main className="mx-auto max-w-xl space-y-4 p-6">
      <h1 className="text-2xl font-semibold text-slate-900">Alumni Portal</h1>
      <p className="text-slate-600">Milestone 1: React → Go → MySQL connectivity check.</p>
      <Card title="API health" state={health} />
      <Card title="Database check" state={db} />
    </main>
  )
}
