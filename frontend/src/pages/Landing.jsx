import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { initializeAgent } from '../services/api';

const domains = ['AI Security', 'Machine Learning', 'AI Ethics', 'Developer Advocacy', 'Open Source', 'Product Strategy'];

export default function Landing() {
  const [personaName, setPersonaName] = useState('');
  const [domain, setDomain] = useState(domains[0]);
  const [isInitializing, setIsInitializing] = useState(false);
  const [errorMessage, setErrorMessage] = useState('');
  const [customTopic, setCustomTopic] = useState('');
  const [customTopicError, setCustomTopicError] = useState('');
  const [isInitializingTopicAgent, setIsInitializingTopicAgent] = useState(false);
  const navigate = useNavigate();

  // If an agentId already exists in localStorage, go straight to the dashboard.
  useEffect(() => {
    const existing = localStorage.getItem('agentId');
    if (existing) {
      navigate(`/dashboard/${existing}`);
      return;
    }
    navigate('/create');
  }, [navigate]);

  const handleSubmit = async (event) => {
    event.preventDefault();

    if (!personaName.trim()) {
      return;
    }

    setIsInitializing(true);
    setErrorMessage('');

    try {
      const response = await initializeAgent({
        name: personaName.trim(),
        domain,
      });

      navigate(`/dashboard/${response.agentId}`);
    } catch (error) {
      setErrorMessage('Unable to initialize the agent right now. Please try again.');
      setIsInitializing(false);
    }
  };

  const handleCustomTopicSubmit = async (event) => {
    event.preventDefault();
    if (!customTopic.trim()) {
      setCustomTopicError('Please enter a topic or question.');
      return;
    }

    setIsInitializingTopicAgent(true);
    setCustomTopicError('');

    try {
      const response = await initializeAgent({
        name: personaName.trim() || 'Autonomous Persona',
        domain,
        topic: customTopic.trim(),
      });

      navigate(`/dashboard/${response.agentId}`);
    } catch (error) {
      setCustomTopicError('Unable to initialize the agent for that topic. Please try again.');
    } finally {
      setIsInitializingTopicAgent(false);
    }
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

          {errorMessage ? (
            <p className="text-sm text-rose-300">{errorMessage}</p>
          ) : null}
        </form>

        <div className="mt-10 rounded-3xl border border-[#948979]/20 bg-[#222831] p-6">
          <h2 className="text-2xl font-semibold text-[#DFD0B8]">Launch a topic-driven agent</h2>
          <p className="mt-2 text-sm leading-7 text-[#DFD0B8]/70">Enter any topic or question and start the agent so it can publish through its normal, high-quality workflow.</p>

          <form onSubmit={handleCustomTopicSubmit} className="mt-6 space-y-4">
            <textarea
              value={customTopic}
              onChange={(event) => setCustomTopic(event.target.value)}
              rows={4}
              placeholder="Type any topic or question here..."
              className="w-full rounded-2xl border border-[#948979]/20 bg-[#393E46] px-4 py-3 text-[#DFD0B8] outline-none transition focus:border-[#948979]"
            />
            <div className="flex flex-wrap items-center gap-3">
              <button
                type="submit"
                disabled={isInitializingTopicAgent}
                className="rounded-full bg-[#948979] px-5 py-3 text-sm font-semibold text-[#222831] transition hover:opacity-90 disabled:cursor-not-allowed disabled:opacity-70"
              >
                {isInitializingTopicAgent ? 'Initializing...' : 'Start Agent'}
              </button>
              {customTopicError ? <p className="text-sm text-rose-300">{customTopicError}</p> : null}
            </div>
          </form>
        </div>
      </section>
    </main>
  );
}
