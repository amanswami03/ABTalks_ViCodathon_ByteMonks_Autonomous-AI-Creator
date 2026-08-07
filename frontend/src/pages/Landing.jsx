import { Link } from 'react-router-dom';

export default function Landing() {
  return (
    <main className="flex min-h-screen items-center justify-center bg-[#222831] px-6 text-[#DFD0B8]">
      <section className="w-full max-w-2xl rounded-[28px] border border-[#948979]/20 bg-[#393E46] p-10 shadow-2xl shadow-black/20">
        <p className="mb-4 inline-flex rounded-full border border-[#948979]/30 bg-[#948979]/10 px-3 py-1 text-sm font-medium text-[#DFD0B8]">
          ABTalks Frontend
        </p>
        <h1 className="text-4xl font-semibold tracking-tight sm:text-5xl">
          Build your autonomous AI experience.
        </h1>
        <p className="mt-4 text-lg text-[#DFD0B8]/80">
          A modern, minimalist foundation for your next product experience.
        </p>
        <div className="mt-8 flex flex-wrap gap-3">
          <Link
            to="/dashboard/demo-agent"
            className="rounded-full bg-[#948979] px-5 py-3 font-semibold text-[#222831] transition hover:opacity-90"
          >
            Open Dashboard
          </Link>
          <a
            href="https://reactrouter.com"
            target="_blank"
            rel="noreferrer"
            className="rounded-full border border-[#948979]/30 px-5 py-3 font-semibold text-[#DFD0B8] transition hover:border-[#DFD0B8]/40"
          >
            Learn Routing
          </a>
        </div>
      </section>
    </main>
  );
}
