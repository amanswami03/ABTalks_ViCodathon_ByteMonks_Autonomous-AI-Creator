import axios from 'axios';

const api = axios.create({
  baseURL:
    import.meta.env.VITE_API_BASE_URL ||
    'https://abtalks-vicodathon-bytemonks-autonomous.onrender.com/api',
  headers: {
    'Content-Type': 'application/json',
  },
});

export const initializeAgent = async ({ name, domain }) => {
  const response = await api.post('/agent/init', {
    persona: {
      name,
      domain,
    },
  });

  return response.data;
};

export const getAgentFeed = async (agentId) => {
  const response = await api.get('/agent/feed', {
    params: { agentId },
  });

  return response.data;
};

export const getAgentActivity = async (agentId) => {
  const response = await api.get(`/agent/${agentId}/activity`);
  return response.data;
};

export const getAgentLogs = async (agentId) => {
  const response = await api.get(`/agent/${agentId}/logs`);
  return response.data;
};

export const getAgentDetails = async (agentId) => {
  const response = await api.get(`/agent/${agentId}`);
  return response.data;
};

export default api;
