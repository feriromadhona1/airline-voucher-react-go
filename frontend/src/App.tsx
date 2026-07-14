import { VoucherForm } from './components/VoucherForm'

export default function App() {
  return (
    <main className="mx-auto flex min-h-screen w-full max-w-3xl flex-col justify-center px-5 py-12 sm:px-8">
      <header className="mb-10">
        <p className="mb-2 text-xs font-semibold uppercase tracking-[0.22em] text-sky-deep">
          Crew Desk
        </p>
        <h1 className="text-4xl font-bold tracking-tight text-ocean sm:text-5xl">
          Airline Voucher
        </h1>
        <p className="mt-3 max-w-xl text-base leading-relaxed text-mist">
          Assign 3 random seat vouchers per flight.
        </p>
      </header>

      <section className="rounded-2xl border border-line bg-panel p-6 shadow-[0_18px_50px_rgba(10,77,140,0.12)] sm:p-8">
        <VoucherForm />
      </section>
    </main>
  )
}
