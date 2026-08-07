import { Outlet, Link, useParams } from 'react-router-dom';

const navItems = ['Overview', 'Agents', 'Insights', 'Settings'];

export default function DashboardLayout() {
  const { agentId } = useParams();

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
              <Link
                key={item}
                to="#"
                className="flex items-center rounded-2xl px-4 py-3 text-sm font-medium text-[#DFD0B8]/80 transition hover:bg-[#222831] hover:text-[#DFD0B8]"
              >
                {item}
              </Link>
            ))}
          </nav>
        </aside>

        <div className="flex-1">
          <header className="border-b border-[#948979]/20 bg-[#393E46]/80 px-6 py-4 backdrop-blur">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-[#DFD0B8]/70">Dashboard</p>
                <h1 className="text-xl font-semibold">{agentId || 'Agent'}</h1>
              </div>
              <div className="rounded-full border border-[#948979]/30 bg-[#948979]/10 px-4 py-2 text-sm text-[#DFD0B8]">
                Live Workspace
              </div>
            </div>
          </header>

          <main className="p-6">
            <div className="rounded-[28px] border border-[#948979]/20 bg-[#393E46] p-8 shadow-2xl shadow-black/20">
              <Outlet />
            </div>
          </main>
        </div>
      </div>
    </div>
  );
}
