import { useEffect, useState } from 'react';
import { Link, useLocation, useNavigate, useParams } from 'react-router-dom';
import { getAgentFeed, getAgentActivity, getAgentDetails, getAgentLogs, getAgentTopics } from '../services/api';

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

function formatPostTime(value) {
  if (!value) {
    return 'Live';
  }

  try {
    return new Date(value).toLocaleString();
  } catch {
    return value;
  }
}

function getSectionFromPath(pathname) {
  const segment = pathname.split('/').filter(Boolean).pop();

  if (segment === 'feed') return 'Feed';
  if (segment === 'topics') return 'Topics';
  if (segment === 'memory') return 'Memory';
  if (segment === 'analytics') return 'Analytics';
  if (segment === 'settings') return 'Settings';

  return 'Overview';
}

export default function DashboardLayout() {
  const { agentId } = useParams();
  const location = useLocation();
  const navigate = useNavigate();
  const [activeSection, setActiveSection] = useState('Overview');
  const [posts, setPosts] = useState([]);
  const [isLoadingFeed, setIsLoadingFeed] = useState(true);
  const [feedError, setFeedError] = useState('');
  const [activity, setActivity] = useState(null);
  const [agentDetails, setAgentDetails] = useState(null);
  const [isLoadingActivity, setIsLoadingActivity] = useState(false);
  const [activityError, setActivityError] = useState('');
  const [logs, setLogs] = useState([]);
  const [isLoadingLogs, setIsLoadingLogs] = useState(false);
  const [logsError, setLogsError] = useState('');
  const [topics, setTopics] = useState({ accepted: [], rejected: [] });
  const [isLoadingTopics, setIsLoadingTopics] = useState(false);
  const [topicsError, setTopicsError] = useState('');

  useEffect(() => {
    setActiveSection(getSectionFromPath(location.pathname));
  }, [location.pathname]);

  // Persist agentId when the dashboard is visited so the experience is remembered.
  useEffect(() => {
    if (agentId) {
      try {
        localStorage.setItem('agentId', agentId);
      } catch {
        // noop
      }
    }
  }, [agentId]);

  const activeContent = sectionContent[activeSection];
  const latestPost = posts[0];
  const overviewFeed = activeSection === 'Feed' ? posts : posts.slice(0, 3);
  const topicQueueFeed = posts;
  const sourceCount = posts.reduce((count, post) => count + (post.sources?.length || 0), 0);
  const responseQuality = posts.length > 0 ? `${Math.min(99, 70 + posts.length * 2)}%` : '—';

  const formattedLastPublished = activity?.lastPublishedAt ? formatPostTime(activity.lastPublishedAt) : 'Never';
  const formattedNextRun = activity?.nextRunAt ? formatPostTime(activity.nextRunAt) : 'Pending';

  useEffect(() => {
    if (!agentId) {
      setPosts([]);
      setIsLoadingFeed(false);
      return undefined;
    }

    let isMounted = true;

    const loadFeed = async () => {
      setIsLoadingFeed(true);

      try {
        const response = await getAgentFeed(agentId);
        if (isMounted) {
          setPosts(response.posts || []);
          setFeedError('');
        }
      } catch (error) {
        if (isMounted) {
          setFeedError('Unable to load the live feed from the backend.');
          setPosts([]);
        }
      } finally {
        if (isMounted) {
          setIsLoadingFeed(false);
        }
      }
    };

    const loadActivity = async () => {
      setIsLoadingActivity(true);

      try {
        const response = await getAgentActivity(agentId);
        if (isMounted) {
          setActivity(response);
          setActivityError('');
        }
      } catch (error) {
        if (isMounted) {
          setActivityError('Unable to load agent activity.');
          setActivity(null);
        }
      } finally {
        if (isMounted) {
          setIsLoadingActivity(false);
        }
      }
    };

    const loadDetails = async () => {
      try {
        const response = await getAgentDetails(agentId);
        if (isMounted) {
          setAgentDetails(response || null);
        }
      } catch (err) {
        if (isMounted) setAgentDetails(null);
      }
    };

    const loadTopics = async () => {
      setIsLoadingTopics(true);
      try {
        const response = await getAgentTopics(agentId);
        if (isMounted) {
          setTopics(response || { accepted: [], rejected: [] });
          setTopicsError('');
        }
      } catch (err) {
        if (isMounted) {
          setTopicsError('Unable to load topics from the backend.');
          setTopics({ accepted: [], rejected: [] });
        }
      } finally {
        if (isMounted) setIsLoadingTopics(false);
      }
    };

    const loadLogs = async () => {
      setIsLoadingLogs(true);

      try {
        const response = await getAgentLogs(agentId);
        if (isMounted) {
          setLogs(response || []);
          setLogsError('');
        }
      } catch (error) {
        if (isMounted) {
          setLogsError('Unable to load agent logs.');
          setLogs([]);
        }
      } finally {
        if (isMounted) {
          setIsLoadingLogs(false);
        }
      }
    };

    loadFeed();
    loadActivity();
    loadDetails();
    loadLogs();
    loadTopics();

    const refreshInterval = window.setInterval(() => {
      loadFeed();
      loadActivity();
      loadLogs();
      loadTopics();
    }, 15000);

    return () => {
      isMounted = false;
      window.clearInterval(refreshInterval);
    };
  }, [agentId]);

  return (
    <div className="min-h-screen bg-[#222831] text-[#DFD0B8]">
      <div className="mx-auto flex min-h-screen max-w-7xl flex-col lg:flex-row">
        <aside className="hidden lg:block w-full border-b border-[#948979]/20 bg-[#393E46] p-6 lg:w-72 lg:border-b-0 lg:border-r">
          <div className="mb-8">
            <p className="text-sm font-semibold uppercase tracking-[0.3em] text-[#948979]">ABTalks</p>
            <h2 className="mt-2 text-xl font-semibold">Agent Console</h2>
          </div>

          <nav className="space-y-2">
            {navItems.map((item) => (
              <button
                key={item}
                type="button"
                onClick={() => {
                  const routeKey = item === 'Overview' ? '' : item.toLowerCase();
                  const targetPath = routeKey
                    ? `/dashboard/${agentId}/${routeKey}`
                    : `/dashboard/${agentId}`;

                  navigate(targetPath);
                }}
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

        <div className="block w-full border-b border-[#948979]/10 bg-[#393E46]/80 p-4 lg:hidden">
          <div className="flex gap-2 overflow-x-auto pb-1">
            {navItems.map((item) => (
              <button
                key={item}
                type="button"
                onClick={() => {
                  const routeKey = item === 'Overview' ? '' : item.toLowerCase();
                  const targetPath = routeKey
                    ? `/dashboard/${agentId}/${routeKey}`
                    : `/dashboard/${agentId}`;

                  navigate(targetPath);
                }}
                className={`min-w-[120px] rounded-full border px-3 py-2 text-xs font-semibold transition ${
                  activeSection === item
                    ? 'border-[#948979] bg-[#948979] text-[#222831]'
                    : 'border-[#948979]/20 bg-[#222831] text-[#DFD0B8]/80 hover:border-[#DFD0B8]/40 hover:text-[#DFD0B8]'
                }`}
              >
                {item}
              </button>
            ))}
          </div>
        </div>

        <div className="flex-1">
          <header className="border-b border-[#948979]/20 bg-[#393E46]/80 px-6 py-4 backdrop-blur">
            <div className="flex flex-wrap items-center justify-between gap-4">
              <div className="flex items-center gap-4">
                <div className="flex h-11 w-11 items-center justify-center rounded-2xl bg-[#948979] text-lg font-semibold text-[#222831]">
                  {(agentDetails && agentDetails.name && agentDetails.name[0]) || (agentId ? agentId[0] : 'A')}
                </div>
                <div>
                  <p className="text-sm text-[#DFD0B8]/70">Agent Workspace</p>
                  <h1 className="text-xl font-semibold">{(agentDetails && agentDetails.name) || agentId || 'Agent'}</h1>
                  {agentDetails?.domain ? (
                    <p className="text-xs text-[#DFD0B8]/60">{agentDetails.domain}</p>
                  ) : null}
                </div>
              </div>

              <div className="flex flex-wrap items-center gap-3">
                <Link
                  to="/"
                  className="rounded-full border border-[#948979]/30 bg-[#948979]/10 px-4 py-2 text-sm font-semibold text-[#DFD0B8] transition hover:bg-[#948979]/20"
                >
                  Home
                </Link>
                <div className="rounded-full border border-[#948979]/30 bg-[#948979]/10 px-4 py-2 text-sm text-[#DFD0B8]">
                  Active Status
                </div>
                <div className="rounded-full border border-emerald-400/20 bg-emerald-500/10 px-4 py-2 text-sm text-[#34d399]">
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
                      <div>
                        <h3 className="text-lg font-semibold text-[#DFD0B8]">Current Activity</h3>
                        <p className="text-sm text-[#DFD0B8]/70">{activity?.currentTask || 'Reading live sources'}</p>
                      </div>
                      <span className="rounded-full border border-emerald-400/20 bg-emerald-500/10 px-3 py-1 text-sm text-emerald-300">
                        {activity?.status || 'Online'}
                      </span>
                    </div>
                    <p className="mt-4 text-sm leading-7 text-[#DFD0B8]/80">
                      {latestPost?.text || (isLoadingFeed ? 'Loading live activity...' : 'The agent is waiting for its first published post.')}
                    </p>
                    <div className="mt-4 grid gap-3 sm:grid-cols-2">
                      <div className="rounded-2xl bg-[#393E46] p-3 text-sm text-[#DFD0B8]/80">
                        <p className="font-semibold text-[#DFD0B8]">Last published</p>
                        <p className="mt-1">{formattedLastPublished}</p>
                      </div>
                      <div className="rounded-2xl bg-[#393E46] p-3 text-sm text-[#DFD0B8]/80">
                        <p className="font-semibold text-[#DFD0B8]">Next run</p>
                        <p className="mt-1">{formattedNextRun}</p>
                      </div>
                    </div>
                  </div>

                  <div className="rounded-[24px] border border-[#948979]/20 bg-[#222831] p-6">
                    <h3 className="text-lg font-semibold text-[#DFD0B8]">Statistics</h3>
                    <div className="mt-4 grid gap-4 sm:grid-cols-2">
                      <div>
                        <p className="text-3xl font-semibold text-[#948979]">{posts.length}</p>
                        <p className="text-sm text-[#DFD0B8]/70">Tasks completed</p>
                      </div>
                      <div>
                        <p className="text-3xl font-semibold text-[#948979]">{responseQuality}</p>
                        <p className="text-sm text-[#DFD0B8]/70">Response quality</p>
                      </div>
                    </div>
                  </div>

                  <div className="rounded-[24px] border border-[#948979]/20 bg-[#222831] p-6">
                    <h3 className="text-lg font-semibold text-[#DFD0B8]">Recent Feed</h3>
                    <ul className="mt-4 space-y-3 text-sm text-[#DFD0B8]/80">
                      {overviewFeed.length > 0 ? (
                        overviewFeed.map((post) => (
                          <li key={post.id || post.text} className="leading-7">
                            • {post.text}
                          </li>
                        ))
                      ) : (
                        <li>{isLoadingFeed ? 'Loading posts...' : 'No posts published yet.'}</li>
                      )}
                    </ul>
                  </div>

                  <div className="rounded-[24px] border border-[#948979]/20 bg-[#222831] p-6">
                    <h3 className="text-lg font-semibold text-[#DFD0B8]">Topic Queue</h3>
                    <ul className="mt-4 space-y-3 text-sm text-[#DFD0B8]/80">
                      {topicQueueFeed.length > 0 ? (
                        topicQueueFeed.map((post) => (
                          <li key={`${post.id || post.text}-topic`} className="leading-7">
                            • {post.text}
                          </li>
                        ))
                      ) : (
                        <li>{isLoadingFeed ? 'Loading topics...' : 'No topics queued yet.'}</li>
                      )}
                    </ul>
                  </div>

                  <div className="rounded-[24px] border border-[#948979]/20 bg-[#222831] p-6 xl:col-span-2">
                    <h3 className="text-lg font-semibold text-[#DFD0B8]">Memory Summary</h3>
                    <p className="mt-4 text-sm leading-7 text-[#DFD0B8]/80">
                      {latestPost?.rationale || 'The agent will publish fresh rationale and source context as new posts arrive.'}
                    </p>
                  </div>

                  <div className="rounded-[24px] border border-[#948979]/20 bg-[#222831] p-6 xl:col-span-2">
                    <div className="flex items-center justify-between">
                      <h3 className="text-lg font-semibold text-[#DFD0B8]">Recent Log Activity</h3>
                      <span className="text-sm text-[#DFD0B8]/70">{logs.length} entries</span>
                    </div>
                    <div className="mt-4 space-y-3 text-sm text-[#DFD0B8]/80">
                      {logsError ? (
                        <p className="text-rose-300">{logsError}</p>
                      ) : isLoadingLogs ? (
                        <p>Loading logs...</p>
                      ) : logs.length > 0 ? (
                        logs.slice(0, 5).map((entry) => (
                          <div key={entry.time} className="rounded-2xl border border-[#948979]/20 bg-[#393E46] p-4">
                            <p className="font-semibold text-[#DFD0B8]">{entry.action}</p>
                            <p className="mt-1 text-xs text-[#DFD0B8]/70">{new Date(entry.time).toLocaleString()}</p>
                            <p className="mt-2 leading-6">{entry.details}</p>
                          </div>
                        ))
                      ) : (
                        <p>No log activity available yet.</p>
                      )}
                    </div>
                  </div>
                </div>
              ) : activeSection === 'Feed' ? (
                <div className="space-y-4">
                  {feedError ? (
                    <p className="text-sm text-rose-300">{feedError}</p>
                  ) : null}

                  {posts.length > 0 ? (
                    posts.map((post) => (
                      <div key={post.id || post.text} className="rounded-[24px] border border-[#948979]/20 bg-[#222831] p-6">
                        <div className="flex items-start justify-between gap-4">
                          <div>
                            <p className="break-words text-base leading-7 text-[#DFD0B8]">{post.text}</p>
                            <p className="mt-3 text-sm text-[#948979]">{formatPostTime(post.createdAt)}</p>
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
                            {(post.sources || []).map((source) => (
                              <span key={source} className="rounded-full border border-[#948979]/20 bg-[#393E46] px-3 py-1 text-xs text-[#DFD0B8]/80">
                                {source}
                              </span>
                            ))}
                          </div>
                        </div>
                      </div>
                    ))
                  ) : (
                    <div className="rounded-[24px] border border-[#948979]/20 bg-[#222831] p-6">
                      <p className="text-base leading-7 text-[#DFD0B8]">
                        {isLoadingFeed ? 'Loading feed from the backend...' : 'No posts published yet.'}
                      </p>
                    </div>
                  )}
                </div>
              ) : activeSection === 'Topics' ? (
                <div className="grid gap-6 lg:grid-cols-2">
                  <div className="rounded-[24px] border border-[#948979]/20 bg-[#222831] p-6">
                    <div className="flex items-center justify-between">
                      <h3 className="text-lg font-semibold text-[#DFD0B8]">Accepted Topics</h3>
                      <span className="rounded-full border border-emerald-400/20 bg-emerald-500/10 px-3 py-1 text-sm text-emerald-300">
                        {topics.accepted && topics.accepted.length > 0 ? `${topics.accepted.length} selected` : '0 selected'}
                      </span>
                    </div>
                    <ul className="mt-4 space-y-3 text-sm text-[#DFD0B8]/80">
                      {isLoadingTopics ? (
                        <li>Loading topics...</li>
                      ) : topics.accepted && topics.accepted.length > 0 ? (
                        topics.accepted.map((post) => (
                          <li key={`${post.id || post.text}-accepted`} className="rounded-2xl border border-[#948979]/20 bg-[#393E46] px-4 py-3">
                            {post.text}
                          </li>
                        ))
                      ) : (
                        <li className="rounded-2xl border border-[#948979]/20 bg-[#393E46] px-4 py-3">No accepted topics yet.</li>
                      )}
                    </ul>
                  </div>

                  <div className="rounded-[24px] border border-[#948979]/20 bg-[#222831] p-6">
                    <div className="flex items-center justify-between">
                      <h3 className="text-lg font-semibold text-[#DFD0B8]">Rejected Topics</h3>
                      <span className="rounded-full border border-rose-400/20 bg-rose-500/10 px-3 py-1 text-sm text-rose-300">
                        {topics.rejected && topics.rejected.length > 0 ? `${topics.rejected.length} skipped` : '0 skipped'}
                      </span>
                    </div>
                    <ul className="mt-4 space-y-3 text-sm text-[#DFD0B8]/80">
                      {topicsError ? (
                        <li className="rounded-2xl border border-[#948979]/20 bg-[#393E46] px-4 py-3">
                          <div className="font-medium text-[#DFD0B8]">Live backend unavailable</div>
                          <p className="mt-2 text-xs leading-6 text-[#DFD0B8]/70">Reason for rejection: The backend could not be reached.</p>
                        </li>
                      ) : isLoadingTopics ? (
                        <li>Loading rejected topics...</li>
                      ) : topics.rejected && topics.rejected.length > 0 ? (
                        topics.rejected.map((t) => (
                          <li key={`${t.title}-${t.time}`} className="rounded-2xl border border-[#948979]/20 bg-[#393E46] px-4 py-3">
                            <div className="font-medium text-[#DFD0B8]">{t.title || 'Untitled'}</div>
                            <p className="mt-2 text-xs leading-6 text-[#DFD0B8]/70">Reason: {t.reason}</p>
                          </li>
                        ))
                      ) : (
                        <li className="rounded-2xl border border-[#948979]/20 bg-[#393E46] px-4 py-3">
                          <div className="font-medium text-[#DFD0B8]">No rejected topics yet</div>
                          <p className="mt-2 text-xs leading-6 text-[#DFD0B8]/70">The agent has not rejected any topics yet.</p>
                        </li>
                      )}
                    </ul>
                  </div>

                </div>
              ) : activeSection === 'Memory' ? (
                <div className="grid gap-6 lg:grid-cols-2">
                  <div className="rounded-[24px] border border-[#948979]/20 bg-[#222831] p-6">
                    <h3 className="text-lg font-semibold text-[#DFD0B8]">Interests</h3>
                    <ul className="mt-4 space-y-3 text-sm text-[#DFD0B8]/80">
                      {topicQueueFeed.length > 0 ? (
                        topicQueueFeed.map((post) => (
                          <li key={`${post.id || post.text}-memory`} className="rounded-2xl border border-[#948979]/20 bg-[#393E46] px-4 py-3">
                            {post.text}
                          </li>
                        ))
                      ) : (
                        <li className="rounded-2xl border border-[#948979]/20 bg-[#393E46] px-4 py-3">No memory nodes yet.</li>
                      )}
                    </ul>
                  </div>

                  <div className="rounded-[24px] border border-[#948979]/20 bg-[#222831] p-6">
                    <h3 className="text-lg font-semibold text-[#DFD0B8]">Recent Topics</h3>
                    <ul className="mt-4 space-y-3 text-sm text-[#DFD0B8]/80">
                      {topicQueueFeed.length > 0 ? (
                        topicQueueFeed.map((post) => (
                          <li key={`${post.id || post.text}-recent`} className="rounded-2xl border border-[#948979]/20 bg-[#393E46] px-4 py-3">
                            {post.text}
                          </li>
                        ))
                      ) : (
                        <li className="rounded-2xl border border-[#948979]/20 bg-[#393E46] px-4 py-3">No recent topics available.</li>
                      )}
                    </ul>
                  </div>
                </div>
              ) : activeSection === 'Analytics' ? (
                <div className="grid gap-6 lg:grid-cols-2">
                  <div className="rounded-[24px] border border-[#948979]/20 bg-[#222831] p-6">
                    <p className="text-sm uppercase tracking-[0.3em] text-[#948979]">Total Activity</p>
                    <p className="mt-4 text-4xl font-semibold text-[#DFD0B8]">{posts.length}</p>
                    <p className="mt-2 text-sm text-[#DFD0B8]/70">Live count of published posts returned by the backend feed endpoint.</p>
                  </div>

                  <div className="rounded-[24px] border border-[#948979]/20 bg-[#222831] p-6">
                    <p className="text-sm uppercase tracking-[0.3em] text-[#948979]">Engagement</p>
                    <p className="mt-4 text-4xl font-semibold text-[#DFD0B8]">{responseQuality}</p>
                    <p className="mt-2 text-sm text-[#DFD0B8]/70">Derived from the latest live feed response.</p>
                  </div>

                  <div className="rounded-[24px] border border-[#948979]/20 bg-[#222831] p-6">
                    <p className="text-sm uppercase tracking-[0.3em] text-[#948979]">Sources</p>
                    <p className="mt-4 text-4xl font-semibold text-[#DFD0B8]">{sourceCount}</p>
                    <p className="mt-2 text-sm text-[#DFD0B8]/70">Aggregated source references from the current feed data.</p>
                  </div>

                  <div className="rounded-[24px] border border-[#948979]/20 bg-[#222831] p-6">
                    <p className="text-sm uppercase tracking-[0.3em] text-[#948979]">Coverage</p>
                    <p className="mt-4 text-4xl font-semibold text-[#DFD0B8]">{posts.length > 0 ? 'Live' : 'Waiting'}</p>
                    <p className="mt-2 text-sm text-[#DFD0B8]/70">Coverage state updates as new posts are returned by the backend.</p>
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
