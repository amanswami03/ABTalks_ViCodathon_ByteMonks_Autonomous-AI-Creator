import { useEffect, useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';

export default function Home() {
  const [agentId, setAgentId] = useState(null);
  const navigate = useNavigate();

  useEffect(() => {
    const existing = localStorage.getItem('agentId');
    setAgentId(existing);
  }, []);

  const handleCreateNew = () => {
    localStorage.removeItem('agentId');
    navigate('/create');
  };

  return (
    <main className="flex min-h-screen items-center justify-center bg-[#222831] px-6 py-10 text-[#DFD0B8]">
      <section className="w-full max-w-3xl rounded-[32px] border border-[#948979]/20 bg-[#393E46] p-8 shadow-2xl shadow-black/20 sm:p-10 lg:p-12">
        <div className="max-w-2xl">
          <p className="mb-4 inline-flex rounded-full border border-[#948979]/30 bg-[#948979]/10 px-3 py-1 text-sm font-medium text-[#DFD0B8]">
            Autonomous AI Creator
          </p>
          <h1 className="text-4xl font-semibold tracking-tight text-[#DFD0B8] sm:text-5xl">
            Create or resume your AI workspace.
          </h1>
          <p className="mt-4 text-lg leading-8 text-[#DFD0B8]/80">
            Start by choosing a topic and persona, then let the agent publish insights automatically.
          </p>
        </div>

        <div className="mt-10 space-y-4">
          <button
            type="button"
            onClick={handleCreateNew}
            className="inline-flex items-center justify-center rounded-full bg-[#948979] px-6 py-3 font-semibold text-[#222831] transition hover:opacity-90"
          >
            Create new agent
          </button>

          {agentId && (
            <div className="rounded-2xl border border-[#948979]/20 bg-[#222831] p-5 text-sm text-[#DFD0B8]/80">
              <p className="font-semibold text-[#DFD0B8]">Existing agent available</p>
              <p className="mt-2">You can resume your current workspace, or create a new one to choose a different topic.</p>
              <div className="mt-4 flex flex-col gap-3 sm:flex-row">
                <Link
                  to={`/dashboard/${agentId}`}
                  className="rounded-full bg-[#948979] px-5 py-3 text-sm font-semibold text-[#222831]"
                >
                  Resume agent
                </Link>
                <button
                  type="button"
                  onClick={handleCreateNew}
                  className="rounded-full border border-[#948979]/20 px-5 py-3 text-sm font-semibold text-[#DFD0B8] hover:border-[#DFD0B8]/40"
                >
                  Choose new topic
                </button>
              </div>
            </div>
          )}
        </div>
      </section>
    </main>
  );
}
