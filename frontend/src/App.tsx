import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { useState, useEffect } from 'react';
import { api, setToken, clearToken, getToken } from './lib/api';
import LoginPage from './pages/LoginPage';
import RegisterPage from './pages/RegisterPage';
import DashboardPage from './pages/DashboardPage';
import TradePage from './pages/TradePage';
import PortfolioPage from './pages/PortfolioPage';
import HistoryPage from './pages/HistoryPage';
import Layout from './components/Layout';

function App() {
  const [user, setUser] = useState<{ email: string; username: string } | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (getToken()) {
      api.auth.me()
        .then(res => setUser(res.data))
        .catch(() => { clearToken(); })
        .finally(() => setLoading(false));
    } else {
      setLoading(false);
    }
  }, []);

  const handleLogin = async (email: string, password: string) => {
    const res = await api.auth.login(email, password);
    setToken(res.access_token);
    const me = await api.auth.me();
    setUser(me.data);
  };

  const handleRegister = async (email: string, username: string, password: string) => {
    await api.auth.register(email, username, password);
    await handleLogin(email, password);
  };

  const handleLogout = () => {
    clearToken();
    setUser(null);
  };

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center" style={{ background: '#0b0e11' }}>
        <div className="text-lg" style={{ color: '#848e9c' }}>Loading...</div>
      </div>
    );
  }

  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={!user ? <LoginPage onLogin={handleLogin} /> : <Navigate to="/" />} />
        <Route path="/register" element={!user ? <RegisterPage onRegister={handleRegister} /> : <Navigate to="/" />} />
        <Route path="/" element={user ? <Layout user={user} onLogout={handleLogout} /> : <Navigate to="/login" />}>
          <Route index element={<DashboardPage />} />
          <Route path="trade" element={<TradePage />} />
          <Route path="portfolio" element={<PortfolioPage />} />
          <Route path="history" element={<HistoryPage />} />
        </Route>
      </Routes>
    </BrowserRouter>
  );
}

export default App;
