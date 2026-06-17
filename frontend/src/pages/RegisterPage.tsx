import { useState } from 'react';
import { Link } from 'react-router-dom';

export default function RegisterPage({ onRegister }: { onRegister: (email: string, username: string, password: string) => Promise<void> }) {
  const [email, setEmail] = useState('');
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);
    try {
      await onRegister(email, username, password);
    } catch (err: any) {
      setError(err.message || 'Registration failed');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center" style={{ background: '#0b0e11' }}>
      <div className="w-full max-w-md p-8 rounded-xl" style={{ background: '#1e2329' }}>
        <div className="text-center mb-8">
          <div className="text-2xl font-bold mb-1" style={{ color: '#fcd535' }}>QW Trading</div>
          <div className="text-sm" style={{ color: '#848e9c' }}>Create your account</div>
        </div>
        <form onSubmit={handleSubmit}>
          {error && (
            <div className="mb-4 p-3 rounded-lg text-sm" style={{ background: 'rgba(246,70,93,0.1)', color: '#f6465d' }}>
              {error}
            </div>
          )}
          <div className="mb-4">
            <label className="block text-sm mb-1.5" style={{ color: '#848e9c' }}>Email</label>
            <input type="email" value={email} onChange={e => setEmail(e.target.value)}
              className="w-full px-4 py-3 rounded-lg text-sm outline-none"
              style={{ background: '#0b0e11', border: '1px solid #2b3139', color: '#eaecef' }} required />
          </div>
          <div className="mb-4">
            <label className="block text-sm mb-1.5" style={{ color: '#848e9c' }}>Username</label>
            <input type="text" value={username} onChange={e => setUsername(e.target.value)}
              className="w-full px-4 py-3 rounded-lg text-sm outline-none"
              style={{ background: '#0b0e11', border: '1px solid #2b3139', color: '#eaecef' }} required />
          </div>
          <div className="mb-6">
            <label className="block text-sm mb-1.5" style={{ color: '#848e9c' }}>Password</label>
            <input type="password" value={password} onChange={e => setPassword(e.target.value)}
              className="w-full px-4 py-3 rounded-lg text-sm outline-none"
              style={{ background: '#0b0e11', border: '1px solid #2b3139', color: '#eaecef' }} required />
          </div>
          <button type="submit" disabled={loading}
            className="w-full py-3 rounded-lg font-semibold text-sm"
            style={{ background: '#fcd535', color: '#0b0e11', opacity: loading ? 0.7 : 1 }}>
            {loading ? 'Creating account...' : 'Create Account'}
          </button>
        </form>
        <div className="text-center mt-6 text-sm" style={{ color: '#848e9c' }}>
          Already have an account? <Link to="/login" style={{ color: '#fcd535' }}>Sign In</Link>
        </div>
      </div>
    </div>
  );
}
