import type {
  CheckResponse,
  ErrorResponse,
  GenerateResponse,
  VoucherFormData,
} from '../types'

const API_BASE = import.meta.env.VITE_API_URL ?? ''

async function parseJson<T>(response: Response): Promise<T> {
  return response.json() as Promise<T>
}

export function formatDisplayDate(date: Date): string {
  const day = String(date.getDate()).padStart(2, '0')
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const year = date.getFullYear()
  return `${day}-${month}-${year}`
}

export function parseDisplayDate(value: string): Date | null {
  if (!/^\d{2}-\d{2}-\d{4}$/.test(value)) {
    return null
  }

  const [day, month, year] = value.split('-').map(Number)
  const date = new Date(year, month - 1, day)

  if (
    date.getFullYear() !== year ||
    date.getMonth() !== month - 1 ||
    date.getDate() !== day
  ) {
    return null
  }

  return date
}

// UI pakai DD-MM-YYYY, API minta YYYY-MM-DD
export function toApiDate(displayDate: string): string {
  const [day, month, year] = displayDate.split('-')
  return `${year}-${month}-${day}`
}

export async function checkVoucher(
  flightNumber: string,
  date: string,
): Promise<CheckResponse> {
  const response = await fetch(`${API_BASE}/api/check`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Accept: 'application/json',
    },
    body: JSON.stringify({ flightNumber, date }),
  })

  const data = await parseJson<CheckResponse & ErrorResponse>(response)

  if (!response.ok) {
    throw new Error(firstError(data) ?? 'Gagal memeriksa voucher.')
  }

  return data
}

export async function generateVoucher(
  form: VoucherFormData,
): Promise<GenerateResponse> {
  const response = await fetch(`${API_BASE}/api/generate`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Accept: 'application/json',
    },
    body: JSON.stringify({
      name: form.name,
      id: form.id,
      flightNumber: form.flightNumber,
      date: toApiDate(form.date),
      aircraft: form.aircraft,
    }),
  })

  const data = await parseJson<GenerateResponse & ErrorResponse>(response)

  if (!response.ok) {
    throw new Error(firstError(data) ?? 'Gagal membuat voucher.')
  }

  return data
}

function firstError(data: ErrorResponse): string | undefined {
  if (data.message) {
    return data.message
  }

  if (!data.errors) {
    return undefined
  }

  const messages = Object.values(data.errors)
  return messages[0]?.[0]
}
