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

              {activeSection === 'Overview' ? (
                <div className="grid gap-6 xl:grid-cols-2">
                  <div className="rounded-[24px] border border-[#948979]/20 bg-[#222831] p-6">
                    <div className="flex items-center justify-between">
                      <h3 className="text-lg font-semibold text-[#DFD0B8]">Current Activity</h3>
                      <span className="rounded-full border border-emerald-400/20 bg-emerald-500/10 px-3 py-1 text-sm text-emerald-300">
                        Online
                      </span>
                    </div>
                    <p className="mt-4 text-sm leading-7 text-[#DFD0B8]/80">
                      The agent is actively reviewing recent context and preparing smart responses for the next step.
                    </p>
                  </div>

                  <div className="rounded-[24px] border border-[#948979]/20 bg-[#222831] p-6">
                    <h3 className="text-lg font-semibold text-[#DFD0B8]">Statistics</h3>
                    <div className="mt-4 grid gap-4 sm:grid-cols-2">
                      <div>
                        <p className="text-3xl font-semibold text-[#948979]">24</p>
                        <p className="text-sm text-[#DFD0B8]/70">Tasks completed</p>
                      </div>
                      <div>
                        <p className="text-3xl font-semibold text-[#948979]">92%</p>
                        <p className="text-sm text-[#DFD0B8]/70">Response quality</p>
                      </div>
                    </div>
                  </div>

                  <div className="rounded-[24px] border border-[#948979]/20 bg-[#222831] p-6">
                    <h3 className="text-lg font-semibold text-[#DFD0B8]">Recent Feed</h3>
                    <ul className="mt-4 space-y-3 text-sm text-[#DFD0B8]/80">
                      <li>• New prompt captured from the latest session</li>
                      <li>• Context summary refreshed successfully</li>
                      <li>• Topic cluster updated with 3 new signals</li>
                    </ul>
                  </div>

                  <div className="rounded-[24px] border border-[#948979]/20 bg-[#222831] p-6">
                    <h3 className="text-lg font-semibold text-[#DFD0B8]">Topic Queue</h3>
                    <ul className="mt-4 space-y-3 text-sm text-[#DFD0B8]/80">
                      <li>• Product strategy</li>
                      <li>• Audience segmentation</li>
                      <li>• Messaging refinement</li>
                    </ul>
                  </div>

                  <div className="rounded-[24px] border border-[#948979]/20 bg-[#222831] p-6 xl:col-span-2">
                    <h3 className="text-lg font-semibold text-[#DFD0B8]">Memory Summary</h3>
                    <p className="mt-4 text-sm leading-7 text-[#DFD0B8]/80">
                      The agent remembers recent preferences, workspace goals, and conversation priorities to keep future recommendations aligned.
                    </p>
                  </div>
                </div>
              ) : activeSection === 'Feed' ? (
                <div className="space-y-4">
                  {[
                    {
                      content: 'Drafted a new positioning angle for the launch narrative based on the latest audience signals.',
                      time: '2m ago',
                      rationale: 'This aligns with the current campaign objectives and strengthens the core message.',
                      sources: ['Customer interviews', 'Support transcripts', 'Recent product notes'],
                    },
                    {
                      content: 'Flagged a shift in interest toward automation workflows and summarized the top concerns.',
                      time: '18m ago',
                      rationale: 'The pattern suggests the audience is prioritizing efficiency over feature breadth.',
                      sources: ['Usage analytics', 'Feedback survey', 'Market watchlist'],
                    },
                    {
                      content: 'Prepared a follow-up content angle that connects product updates with customer value stories.',
                      time: '45m ago',
                      rationale: 'This helps bridge product launches with stronger narrative continuity.',
                      sources: ['Content calendar', 'Recent launches', 'Persona notes'],
                    },
                  ].map((post, index) => (
                    <div key={index} className="rounded-[24px] border border-[#948979]/20 bg-[#222831] p-6">
                      <div className="flex items-start justify-between gap-4">
                        <div>
                          <p className="text-base leading-7 text-[#DFD0B8]">{post.content}</p>
                          <p className="mt-3 text-sm text-[#948979]">{post.time}</p>
                        </div>
                        <span className="rounded-full border border-[#948979]/20 bg-[#393E46] px-3 py-1 text-xs uppercase tracking-[0.25em] text-[#DFD0B8]/70">
                          Feed
                        </span>
                      </div>

                      <div className="mt-4">
                        <p className="text-sm font-semibold text-[#DFD0B8]">Rationale</p>
                        <p className="mt-1 text-sm leading-7 text-[#DFD0B8]/80">{post.rationale}</p>
                      </div>

                      <div className="mt-4">
                        <p className="text-sm font-semibold text-[#DFD0B8]">Sources</p>
                        <div className="mt-2 flex flex-wrap gap-2">
                          {post.sources.map((source) => (
                            <span key={source} className="rounded-full border border-[#948979]/20 bg-[#393E46] px-3 py-1 text-xs text-[#DFD0B8]/80">
                              {source}
                            </span>
                          ))}
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              ) : activeSection === 'Topics' ? (
                <div className="grid gap-6 lg:grid-cols-2">
                  <div className="rounded-[24px] border border-[#948979]/20 bg-[#222831] p-6">
                    <div className="flex items-center justify-between">
                      <h3 className="text-lg font-semibold text-[#DFD0B8]">Accepted Topics</h3>
                      <span className="rounded-full border border-emerald-400/20 bg-emerald-500/10 px-3 py-1 text-sm text-emerald-300">
                        4 selected
                      </span>
                    </div>
                    <ul className="mt-4 space-y-3 text-sm text-[#DFD0B8]/80">
                      <li className="rounded-2xl border border-[#948979]/20 bg-[#393E46] px-4 py-3">Product strategy for launch planning</li>
                      <li className="rounded-2xl border border-[#948979]/20 bg-[#393E46] px-4 py-3">Customer pain points and messaging</li>
                      <li className="rounded-2xl border border-[#948979]/20 bg-[#393E46] px-4 py-3">Automation opportunities in workflows</li>
                    </ul>
                  </div>

                  <div className="rounded-[24px] border border-[#948979]/20 bg-[#222831] p-6">
                    <div className="flex items-center justify-between">
                      <h3 className="text-lg font-semibold text-[#DFD0B8]">Rejected Topics</h3>
                      <span className="rounded-full border border-rose-400/20 bg-rose-500/10 px-3 py-1 text-sm text-rose-300">
                        2 skipped
                      </span>
                    </div>
                    <ul className="mt-4 space-y-3 text-sm text-[#DFD0B8]/80">
                      <li className="rounded-2xl border border-[#948979]/20 bg-[#393E46] px-4 py-3">
                        <div className="font-medium text-[#DFD0B8]">Broad trend speculation</div>
                        <p className="mt-2 text-xs leading-6 text-[#DFD0B8]/70">Reason for rejection: Too vague for a focused agent workflow.</p>
                      </li>
                      <li className="rounded-2xl border border-[#948979]/20 bg-[#393E46] px-4 py-3">
                        <div className="font-medium text-[#DFD0B8]">Unrelated competitor noise</div>
                        <p className="mt-2 text-xs leading-6 text-[#DFD0B8]/70">Reason for rejection: Low signal value and limited actionability.</p>
                      </li>
                    </ul>
                  </div>
                </div>
              ) : activeSection === 'Memory' ? (
                <div className="grid gap-6 lg:grid-cols-2">
                  <div className="rounded-[24px] border border-[#948979]/20 bg-[#222831] p-6">
                    <h3 className="text-lg font-semibold text-[#DFD0B8]">Interests</h3>
                    <ul className="mt-4 space-y-3 text-sm text-[#DFD0B8]/80">
                      <li className="rounded-2xl border border-[#948979]/20 bg-[#393E46] px-4 py-3">Customer journey optimization</li>
                      <li className="rounded-2xl border border-[#948979]/20 bg-[#393E46] px-4 py-3">Narrative clarity and positioning</li>
                      <li className="rounded-2xl border border-[#948979]/20 bg-[#393E46] px-4 py-3">B2B content strategy</li>
                    </ul>
                  </div>

                  <div className="rounded-[24px] border border-[#948979]/20 bg-[#222831] p-6">
                    <h3 className="text-lg font-semibold text-[#DFD0B8]">Recent Topics</h3>
                    <ul className="mt-4 space-y-3 text-sm text-[#DFD0B8]/80">
                      <li className="rounded-2xl border border-[#948979]/20 bg-[#393E46] px-4 py-3">Launch sequencing and narrative pacing</li>
                      <li className="rounded-2xl border border-[#948979]/20 bg-[#393E46] px-4 py-3">Customer objections and objections handling</li>
                      <li className="rounded-2xl border border-[#948979]/20 bg-[#393E46] px-4 py-3">Automation and scaling communications</li>
                    </ul>
                  </div>
                </div>
              ) : activeSection === 'Analytics' ? (
                <div className="grid gap-6 lg:grid-cols-2">
                  <div className="rounded-[24px] border border-[#948979]/20 bg-[#222831] p-6">
                    <p className="text-sm uppercase tracking-[0.3em] text-[#948979]">Total Activity</p>
                    <p className="mt-4 text-4xl font-semibold text-[#DFD0B8]">1,248</p>
                    <p className="mt-2 text-sm text-[#DFD0B8]/70">Mock backend metric for the current agent lifecycle.</p>
                  </div>

                  <div className="rounded-[24px] border border-[#948979]/20 bg-[#222831] p-6">
                    <p className="text-sm uppercase tracking-[0.3em] text-[#948979]">Engagement</p>
                    <p className="mt-4 text-4xl font-semibold text-[#DFD0B8]">87%</p>
                    <p className="mt-2 text-sm text-[#DFD0B8]/70">Average interaction quality based on mock responses.</p>
                  </div>

                  <div className="rounded-[24px] border border-[#948979]/20 bg-[#222831] p-6">
                    <p className="text-sm uppercase tracking-[0.3em] text-[#948979]">Response Speed</p>
                    <p className="mt-4 text-4xl font-semibold text-[#DFD0B8]">1.8s</p>
                    <p className="mt-2 text-sm text-[#DFD0B8]/70">Average completion time for queued actions.</p>
                  </div>

                  <div className="rounded-[24px] border border-[#948979]/20 bg-[#222831] p-6">
                    <p className="text-sm uppercase tracking-[0.3em] text-[#948979]">Coverage</p>
                    <p className="mt-4 text-4xl font-semibold text-[#DFD0B8]">94%</p>
                    <p className="mt-2 text-sm text-[#DFD0B8]/70">Topic coverage across the active knowledge set.</p>
                  </div>
                </div>
              ) : (
                <p className="max-w-2xl text-base leading-8 text-[#DFD0B8]/80">{activeContent.description}</p>
              )}
            </div>
          </main>
        </div>
      </div>
    </div>
  );
}
