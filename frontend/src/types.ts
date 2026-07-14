export type AircraftType = 'ATR' | 'Airbus 320' | 'Boeing 737 Max'

export type VoucherFormData = {
  name: string
  id: string
  flightNumber: string
  date: string
  aircraft: AircraftType | ''
}

export type CheckResponse = {
  exists: boolean
}

export type GenerateResponse = {
  success: boolean
  seats: string[]
}

export type ErrorResponse = {
  success?: boolean
  message?: string
  errors?: Record<string, string[]>
}

export const AIRCRAFT_OPTIONS: AircraftType[] = [
  'ATR',
  'Airbus 320',
  'Boeing 737 Max',
]
