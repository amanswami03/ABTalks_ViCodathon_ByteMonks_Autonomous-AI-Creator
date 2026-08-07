import { useState } from 'react';
import { useParams } from 'react-router-dom';

const navItems = ['Overview', 'Feed', 'Topics', 'Memory', 'Analytics', 'Settings'];

const sectionContent = {
  Overview: {
    title: 'Overview',
    description: 'A high-level snapshot of your agent’s current activity and status.',
  },
  Feed: {
    title: 'Feed',
    description: 'Recent updates, prompts, and interactions will appear here.',
  },
  Topics: {
    title: 'Topics',
    description: 'Curated themes and conversation clusters will be managed here.',
  },
  Memory: {
    title: 'Memory',
    description: 'Persistent context and learned details will be shown here.',
  },
  Analytics: {
    title: 'Analytics',
    description: 'Performance trends and engagement metrics will be visualized here.',
  },
  Settings: {
    title: 'Settings',
    description: 'Agent preferences, permissions, and workspace configuration will live here.',
  },
};

export default function DashboardLayout() {
  const { agentId } = useParams();
  const [activeSection, setActiveSection] = useState('Overview');

  const activeContent = sectionContent[activeSection];

  return (
    <div className="min-h-screen bg-[#222831] text-[#DFD0B8]">
      <div className="mx-auto flex min-h-screen max-w-7xl flex-col lg:flex-row">
        <aside className="w-full border-b border-[#948979]/20 bg-[#393E46] p-6 lg:w-72 lg:border-b-0 lg:border-r">
          <div className="mb-8">
            <p className="text-sm font-semibold uppercase tracking-[0.3em] text-[#948979]">ABTalks</p>
            <h2 className="mt-2 text-xl font-semibold">Agent Console</h2>
          </div>

          <nav className="space-y-2">
            {navItems.map((item) => (
              <button
                key={item}
                type="button"
                onClick={() => setActiveSection(item)}
                className={`flex w-full items-center rounded-2xl px-4 py-3 text-left text-sm font-medium transition ${
                  activeSection === item
                    ? 'bg-[#948979] text-[#222831]'
                    : 'text-[#DFD0B8]/80 hover:bg-[#222831] hover:text-[#DFD0B8]'
                }`}
              >
                {item}
              </button>
            ))}
          </nav>
        </aside>

        <div className="flex-1">
          <header className="border-b border-[#948979]/20 bg-[#393E46]/80 px-6 py-4 backdrop-blur">
            <div className="flex flex-wrap items-center justify-between gap-4">
              <div className="flex items-center gap-4">
                <div className="flex h-11 w-11 items-center justify-center rounded-2xl bg-[#948979] text-lg font-semibold text-[#222831]">
                  A
                </div>
                <div>
                  <p className="text-sm text-[#DFD0B8]/70">Agent Workspace</p>
                  <h1 className="text-xl font-semibold">{agentId || 'Agent'}</h1>
                </div>
              </div>

              <div className="flex flex-wrap items-center gap-3">
                <div className="rounded-full border border-[#948979]/30 bg-[#948979]/10 px-4 py-2 text-sm text-[#DFD0B8]">
                  Active Status
                </div>
                <div className="rounded-full border border-emerald-400/20 bg-emerald-500/10 px-4 py-2 text-sm text-emerald-300">
                  Online
                </div>
              </div>
            </div>
          </header>

          <main className="p-6">
            <div className="rounded-[28px] border border-[#948979]/20 bg-[#393E46] p-8 shadow-2xl shadow-black/20">
              <div className="mb-6 flex items-center justify-between">
                <div>
                  <p className="text-sm uppercase tracking-[0.3em] text-[#948979]">{activeContent.title}</p>
                  <h2 className="mt-2 text-2xl font-semibold text-[#DFD0B8]">{activeContent.title}</h2>
                </div>
                <div className="rounded-full border border-[#948979]/20 bg-[#222831] px-4 py-2 text-sm text-[#DFD0B8]/80">
                  {agentId || 'demo-agent'}
                </div>
              </div>

              <p className="max-w-2xl text-base leading-8 text-[#DFD0B8]/80">{activeContent.description}</p>
            </div>
          </main>
        </div>
      </div>
    </div>
  );
}
