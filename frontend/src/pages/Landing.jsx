import { useState } from 'react';
import { useNavigate } from 'react-router-dom';

const domains = ['Product Strategy', 'Marketing', 'Research', 'Operations'];

export default function Landing() {
  const [personaName, setPersonaName] = useState('');
  const [domain, setDomain] = useState(domains[0]);
  const [isInitializing, setIsInitializing] = useState(false);
  const navigate = useNavigate();

  const handleSubmit = (event) => {
    event.preventDefault();

    if (!personaName.trim()) {
      return;
    }

    setIsInitializing(true);

    window.setTimeout(() => {
      const slug = personaName.trim().toLowerCase().replace(/[^a-z0-9]+/g, '-');
      navigate(`/dashboard/${slug}`);
    }, 600);
  };

  return (
    <main className="flex min-h-screen items-center justify-center bg-[#222831] px-6 py-10 text-[#DFD0B8]">
      <section className="w-full max-w-3xl rounded-[32px] border border-[#948979]/20 bg-[#393E46] p-8 shadow-2xl shadow-black/20 sm:p-10 lg:p-12">
        <div className="max-w-2xl">
          <p className="mb-4 inline-flex rounded-full border border-[#948979]/30 bg-[#948979]/10 px-3 py-1 text-sm font-medium text-[#DFD0B8]">
            Autonomous AI Creator
          </p>
          <h1 className="text-4xl font-semibold tracking-tight text-[#DFD0B8] sm:text-5xl">
            Launch your next intelligent agent.
          </h1>
          <p className="mt-4 text-lg leading-8 text-[#DFD0B8]/80">
            Define a persona, select a domain, and initialize an agent experience tailored to your workflow.
          </p>
        </div>

        <form onSubmit={handleSubmit} className="mt-10 space-y-6">
          <div className="space-y-2">
            <label htmlFor="persona-name" className="text-sm font-medium uppercase tracking-[0.2em] text-[#DFD0B8]/70">
              Persona Name
            </label>
            <input
              id="persona-name"
              type="text"
              value={personaName}
              onChange={(event) => setPersonaName(event.target.value)}
              placeholder="e.g. Maya, the Growth Strategist"
              className="w-full rounded-2xl border border-[#948979]/20 bg-[#222831] px-4 py-3 text-[#DFD0B8] outline-none transition focus:border-[#948979]"
            />
          </div>

          <div className="space-y-2">
            <label htmlFor="ai-domain" className="text-sm font-medium uppercase tracking-[0.2em] text-[#DFD0B8]/70">
              AI Domain
            </label>
            <select
              id="ai-domain"
              value={domain}
              onChange={(event) => setDomain(event.target.value)}
              className="w-full rounded-2xl border border-[#948979]/20 bg-[#222831] px-4 py-3 text-[#DFD0B8] outline-none transition focus:border-[#948979]"
            >
              {domains.map((item) => (
                <option key={item} value={item}>
                  {item}
                </option>
              ))}
            </select>
          </div>

          <button
            type="submit"
            disabled={isInitializing}
            className="rounded-full bg-[#948979] px-6 py-3 font-semibold text-[#222831] transition hover:opacity-90 disabled:cursor-not-allowed disabled:opacity-70"
          >
            {isInitializing ? 'Initializing...' : 'Initialize Agent'}
          </button>
        </form>
      </section>
    </main>
  );
}
