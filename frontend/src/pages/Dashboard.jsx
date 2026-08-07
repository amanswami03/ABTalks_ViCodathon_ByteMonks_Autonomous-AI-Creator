import { useParams, Link } from 'react-router-dom';

export default function Dashboard() {
  const { agentId } = useParams();

  return (
    <main className="flex min-h-screen items-center justify-center bg-[#222831] px-6 text-[#DFD0B8]">
      <section className="w-full max-w-2xl rounded-[28px] border border-[#948979]/20 bg-[#393E46] p-10 shadow-2xl shadow-black/20">
        <p className="mb-4 inline-flex rounded-full border border-[#948979]/30 bg-[#948979]/10 px-3 py-1 text-sm font-medium text-[#DFD0B8]">
          Dashboard Route
        </p>
        <h1 className="text-3xl font-semibold tracking-tight sm:text-4xl">
          Agent Dashboard
        </h1>
        <p className="mt-4 text-lg text-[#DFD0B8]/80">
          Viewing agent: <span className="font-semibold text-[#948979]">{agentId}</span>
        </p>
        <div className="mt-8">
          <Link
            to="/"
            className="rounded-full border border-[#948979]/30 px-5 py-3 font-semibold text-[#DFD0B8] transition hover:border-[#DFD0B8]/40"
          >
            Back to Landing
          </Link>
        </div>
      </section>
    </main>
  );
}
