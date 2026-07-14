import { useState, type FormEvent } from 'react'
import DatePicker from 'react-datepicker'
import 'react-datepicker/dist/react-datepicker.css'
import {
  checkVoucher,
  formatDisplayDate,
  generateVoucher,
  parseDisplayDate,
  toApiDate,
} from '../api/voucher'
import { AIRCRAFT_OPTIONS, type VoucherFormData } from '../types'

const emptyForm: VoucherFormData = {
  name: '',
  id: '',
  flightNumber: '',
  date: '',
  aircraft: '',
}

type FieldErrors = Partial<Record<keyof VoucherFormData, string>>

type Status =
  | { kind: 'idle' }
  | { kind: 'loading' }
  | { kind: 'success'; seats: string[] }
  | { kind: 'error'; message: string }

const fieldClass =
  'w-full rounded-lg border border-line bg-sand px-3.5 py-2.5 text-ink outline-none transition placeholder:text-mist/60 focus:border-sky focus:ring-2 focus:ring-sky/25'

const fieldErrorClass =
  'w-full rounded-lg border border-danger bg-sand px-3.5 py-2.5 text-ink outline-none transition placeholder:text-mist/60 focus:border-danger focus:ring-2 focus:ring-danger/20'

const labelClass = 'mb-1.5 block text-sm font-medium text-ocean'

function validateForm(form: VoucherFormData): FieldErrors {
  const errors: FieldErrors = {}

  if (!form.name.trim()) {
    errors.name = 'Crew Name wajib diisi.'
  }

  if (!form.id.trim()) {
    errors.id = 'Crew ID wajib diisi.'
  }

  if (!form.flightNumber.trim()) {
    errors.flightNumber = 'Flight Number wajib diisi.'
  }

  if (!form.date) {
    errors.date = 'Flight Date wajib diisi.'
  } else if (!parseDisplayDate(form.date)) {
    errors.date = 'Format tanggal harus DD-MM-YYYY.'
  }

  if (!form.aircraft) {
    errors.aircraft = 'Aircraft Type wajib dipilih.'
  }

  return errors
}

export function VoucherForm() {
  const [form, setForm] = useState<VoucherFormData>(emptyForm)
  const [fieldErrors, setFieldErrors] = useState<FieldErrors>({})
  const [status, setStatus] = useState<Status>({ kind: 'idle' })

  const selectedDate = parseDisplayDate(form.date)

  function update<K extends keyof VoucherFormData>(
    key: K,
    value: VoucherFormData[K],
  ) {
    setForm((prev) => ({ ...prev, [key]: value }))
    setFieldErrors((prev) => ({ ...prev, [key]: undefined }))
  }

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()

    const errors = validateForm(form)
    setFieldErrors(errors)

    if (Object.keys(errors).length > 0) {
      setStatus({ kind: 'idle' })
      return
    }

    setStatus({ kind: 'loading' })

    try {
      const apiDate = toApiDate(form.date)
      const check = await checkVoucher(form.flightNumber.trim(), apiDate)

      if (check.exists) {
        setStatus({
          kind: 'error',
          message: `Voucher untuk penerbangan ${form.flightNumber.trim()} pada tanggal ${form.date} sudah pernah digenerate.`,
        })
        return
      }

      const result = await generateVoucher({
        ...form,
        name: form.name.trim(),
        id: form.id.trim(),
        flightNumber: form.flightNumber.trim(),
      })

      setStatus({ kind: 'success', seats: result.seats })
    } catch (error) {
      setStatus({
        kind: 'error',
        message:
          error instanceof Error
            ? error.message
            : 'Terjadi kesalahan. Coba lagi.',
      })
    }
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-5" noValidate>
      <div className="grid gap-4 sm:grid-cols-2">
        <div>
          <label htmlFor="crewName" className={labelClass}>
            Crew Name
          </label>
          <input
            id="crewName"
            type="text"
            value={form.name}
            onChange={(e) => update('name', e.target.value)}
            className={fieldErrors.name ? fieldErrorClass : fieldClass}
            placeholder="Sarah"
            autoComplete="name"
            aria-invalid={Boolean(fieldErrors.name)}
          />
          {fieldErrors.name && (
            <p className="mt-1.5 text-sm text-danger">{fieldErrors.name}</p>
          )}
        </div>

        <div>
          <label htmlFor="crewId" className={labelClass}>
            Crew ID
          </label>
          <input
            id="crewId"
            type="text"
            value={form.id}
            onChange={(e) => update('id', e.target.value)}
            className={fieldErrors.id ? fieldErrorClass : fieldClass}
            placeholder="98123"
            aria-invalid={Boolean(fieldErrors.id)}
          />
          {fieldErrors.id && (
            <p className="mt-1.5 text-sm text-danger">{fieldErrors.id}</p>
          )}
        </div>

        <div>
          <label htmlFor="flightNumber" className={labelClass}>
            Flight Number
          </label>
          <input
            id="flightNumber"
            type="text"
            value={form.flightNumber}
            onChange={(e) => update('flightNumber', e.target.value.toUpperCase())}
            className={fieldErrors.flightNumber ? fieldErrorClass : fieldClass}
            placeholder="GA102"
            aria-invalid={Boolean(fieldErrors.flightNumber)}
          />
          {fieldErrors.flightNumber && (
            <p className="mt-1.5 text-sm text-danger">{fieldErrors.flightNumber}</p>
          )}
        </div>

        <div>
          <label htmlFor="flightDate" className={labelClass}>
            Flight Date
          </label>
          <DatePicker
            id="flightDate"
            selected={selectedDate}
            onChange={(date: Date | null) => {
              update('date', date ? formatDisplayDate(date) : '')
            }}
            dateFormat="dd-MM-yyyy"
            placeholderText="DD-MM-YYYY"
            className={fieldErrors.date ? fieldErrorClass : fieldClass}
            calendarClassName="garuda-datepicker"
            showMonthDropdown
            showYearDropdown
            dropdownMode="select"
          />
          {fieldErrors.date && (
            <p className="mt-1.5 text-sm text-danger">{fieldErrors.date}</p>
          )}
        </div>
      </div>

      <div>
        <label htmlFor="aircraft" className={labelClass}>
          Aircraft Type
        </label>
        <select
          id="aircraft"
          value={form.aircraft}
          onChange={(e) =>
            update('aircraft', e.target.value as VoucherFormData['aircraft'])
          }
          className={fieldErrors.aircraft ? fieldErrorClass : fieldClass}
          aria-invalid={Boolean(fieldErrors.aircraft)}
        >
          <option value="" disabled>
            Pilih tipe pesawat
          </option>
          {AIRCRAFT_OPTIONS.map((option) => (
            <option key={option} value={option}>
              {option}
            </option>
          ))}
        </select>
        {fieldErrors.aircraft && (
          <p className="mt-1.5 text-sm text-danger">{fieldErrors.aircraft}</p>
        )}
      </div>

      <button
        type="submit"
        disabled={status.kind === 'loading'}
        className="inline-flex w-full items-center justify-center gap-2 rounded-lg bg-sky px-4 py-3 text-sm font-semibold text-white transition hover:bg-sky-deep disabled:cursor-not-allowed disabled:opacity-60 sm:w-auto sm:min-w-52"
      >
        {status.kind === 'loading' ? 'Generating...' : 'Generate Vouchers'}
      </button>

      {status.kind === 'error' && (
        <div
          role="alert"
          className="rounded-lg border border-danger/40 bg-danger/10 px-4 py-3 text-sm text-danger"
        >
          {status.message}
        </div>
      )}

      {status.kind === 'success' && (
        <div
          role="status"
          className="rounded-lg border border-ok/35 bg-ok/10 px-4 py-4"
        >
          <p className="mb-3 text-sm font-medium text-ok">
            Voucher berhasil digenerate
          </p>
          <div className="flex flex-wrap gap-3">
            {status.seats.map((seat) => (
              <span
                key={seat}
                className="rounded-md border border-ocean/25 bg-sand px-4 py-2 font-mono text-lg font-semibold tracking-wide text-ocean"
              >
                {seat}
              </span>
            ))}
          </div>
        </div>
      )}
    </form>
  )
}
