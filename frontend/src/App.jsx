import { Routes, Route } from 'react-router-dom';
import Landing from './pages/Landing';
import DashboardLayout from './layouts/DashboardLayout';

export default function App() {
  return (
    <Routes>
      <Route path="/" element={<Landing />} />
      <Route path="/create" element={<Landing />} />
      <Route path="/dashboard/:agentId" element={<DashboardLayout />} />
      <Route path="/dashboard/:agentId/feed" element={<DashboardLayout />} />
      <Route path="/dashboard/:agentId/topics" element={<DashboardLayout />} />
      <Route path="/dashboard/:agentId/memory" element={<DashboardLayout />} />
      <Route path="/dashboard/:agentId/analytics" element={<DashboardLayout />} />
      <Route path="/dashboard/:agentId/settings" element={<DashboardLayout />} />
    </Routes>
  );
}
