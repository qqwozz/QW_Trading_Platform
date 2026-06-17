import { useState } from 'react';
import { Link } from 'react-router-dom';

export default function LoginPage({ onLogin }: { onLogin: (email: string, password: string) => Promise<void> }) {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);
    try {
      await onLogin(email, password);
    } catch (err: any) {
      setError(err.message || 'Login failed');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center" style={{ background: '#0b0e11' }}>
      <div className="w-full max-w-md p-8 rounded-xl" style={{ background: '#1e2329' }}>
        <div className="text-center mb-8">
          <div className="text-2xl font-bold mb-1" style={{ color: '#fcd535' }}>QW Trading</div>
          <div className="text-sm" style={{ color: '#848e9c' }}>Sign in to your account</div>
        </div>
        <form onSubmit={handleSubmit}>
          {error && (
            <div className="mb-4 p-3 rounded-lg text-sm" style={{ background: 'rgba(246,70,93,0.1)', color: '#f6465d' }}>
              {error}
            </div>
          )}
          <div className="mb-4">
            <label className="block text-sm mb-1.5" style={{ color: '#848e9c' }}>Email</label>
            <input
              type="email"
              value={email}
              onChange={e => setEmail(e.target.value)}
              className="w-full px-4 py-3 rounded-lg text-sm outline-none"
              style={{ background: '#0b0e11', border: '1px solid #2b3139', color: '#eaecef' }}
              required
            />
          </div>
          <div className="mb-6">
            <label className="block text-sm mb-1.5" style={{ color: '#848e9c' }}>Password</label>
            <input
              type="password"
              value={password}
              onChange={e => setPassword(e.target.value)}
              className="w-full px-4 py-3 rounded-lg text-sm outline-none"
              style={{ background: '#0b0e11', border: '1px solid #2b3139', color: '#eaecef' }}
              required
            />
          </div>
          <button
            type="submit"
            disabled={loading}
            className="w-full py-3 rounded-lg font-semibold text-sm transition-opacity"
            style={{ background: '#fcd535', color: '#0b0e11', opacity: loading ? 0.7 : 1 }}
          >
            {loading ? 'Signing in...' : 'Sign In'}
          </button>
        </form>
        <div className="text-center mt-6 text-sm" style={{ color: '#848e9c' }}>
          Don't have an account?{' '}
          <Link to="/register" style={{ color: '#fcd535' }}>Register</Link>
        </div>
      </div>
    </div>
  );
}
