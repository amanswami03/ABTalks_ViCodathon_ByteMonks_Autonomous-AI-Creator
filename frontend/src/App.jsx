import { Routes, Route, useNavigate } from 'react-router-dom';
import Landing from './pages/Landing';
import DashboardLayout from './layouts/DashboardLayout';

function DashboardRedirect() {
  const navigate = useNavigate();
  const agentId = localStorage.getItem('agentId');

  // Redirect immediately on mount.
  if (agentId) {
    navigate(`/dashboard/${agentId}`);
    return null;
  }
  navigate('/');
  return null;
}

export default function App() {
  return (
    <Routes>
      <Route path="/" element={<Landing />} />
      <Route path="/create" element={<Landing />} />
      <Route path="/dashboard" element={<DashboardRedirect />} />
      <Route path="/dashboard/:agentId" element={<DashboardLayout />} />
      <Route path="/dashboard/:agentId/feed" element={<DashboardLayout />} />
      <Route path="/dashboard/:agentId/topics" element={<DashboardLayout />} />
      <Route path="/dashboard/:agentId/memory" element={<DashboardLayout />} />
      <Route path="/dashboard/:agentId/analytics" element={<DashboardLayout />} />
      <Route path="/dashboard/:agentId/settings" element={<DashboardLayout />} />
    </Routes>
  );
}
