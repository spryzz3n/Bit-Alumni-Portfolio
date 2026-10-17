const API_URL = import.meta.env.VITE_API_URL ?? 'http://localhost:8080/api/v1'

export async function apiGet(path) {
  const res = await fetch(`${API_URL}${path}`, { headers: { Accept: 'application/json' } })
  const body = await res.json().catch(() => null)
  if (!res.ok) {
    throw new Error(body?.error?.message ?? `Request failed (${res.status})`)
  }
  return body
}
