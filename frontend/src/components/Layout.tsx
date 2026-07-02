import { Outlet, Link, useLocation } from 'react-router-dom';
import { LayoutDashboard, TrendingUp, Briefcase, History, LogOut } from 'lucide-react';

const navItems = [
  { path: '/', label: 'Dashboard', icon: LayoutDashboard },
  { path: '/trade', label: 'Trade', icon: TrendingUp },
  { path: '/portfolio', label: 'Portfolio', icon: Briefcase },
  { path: '/history', label: 'History', icon: History },
];

export default function Layout({ user, onLogout }: { user: { email: string; username: string }; onLogout: () => void }) {
  const location = useLocation();

  return (
    <div className="flex h-screen">
      <aside className="w-60 flex flex-col" style={{ background: '#1e2329', borderRight: '1px solid #2b3139' }}>
        <div className="p-5 border-b" style={{ borderColor: '#2b3139' }}>
          <div className="text-xl font-bold tracking-wide" style={{ color: '#fcd535' }}>QW Trading</div>
          <div className="text-xs mt-1 tracking-wider" style={{ color: '#848e9c' }}>Exchange Platform</div>
        </div>
        <nav className="flex-1 p-3 mt-2">
          {navItems.map(item => {
            const Icon = item.icon;
            const active = location.pathname === item.path;
            return (
              <Link
                key={item.path}
                to={item.path}
                className="flex items-center gap-3 px-3 py-2.5 rounded-lg mb-1 text-sm transition-all duration-200"
                style={{
                  background: active ? 'linear-gradient(135deg, rgba(252,213,53,0.12), rgba(252,213,53,0.04))' : 'transparent',
                  color: active ? '#fcd535' : '#848e9c',
                  borderLeft: active ? '3px solid #fcd535' : '3px solid transparent',
                }}
                onMouseEnter={(e) => {
                  if (!active) {
                    e.currentTarget.style.background = 'rgba(255,255,255,0.04)';
                    e.currentTarget.style.color = '#eaecef';
                  }
                }}
                onMouseLeave={(e) => {
                  if (!active) {
                    e.currentTarget.style.background = 'transparent';
                    e.currentTarget.style.color = '#848e9c';
                  }
                }}
              >
                <Icon size={18} />
                {item.label}
              </Link>
            );
          })}
        </nav>
        <div className="p-3 border-t" style={{ borderColor: '#2b3139' }}>
          <div className="flex items-center gap-2.5 mb-3 px-2">
            <div
              className="w-9 h-9 rounded-full flex items-center justify-center text-sm font-bold"
              style={{
                background: 'linear-gradient(135deg, #fcd535, #f0b90b)',
                color: '#0b0e11',
              }}
            >
              {user.username[0].toUpperCase()}
            </div>
            <div className="flex-1 min-w-0">
              <div className="text-sm font-medium truncate">{user.username}</div>
              <div className="text-xs truncate" style={{ color: '#848e9c' }}>{user.email}</div>
            </div>
          </div>
          <button
            onClick={onLogout}
            className="flex items-center gap-2 w-full px-3 py-2 rounded-lg text-sm transition-all duration-200"
            style={{ color: '#848e9c' }}
            onMouseEnter={(e) => {
              e.currentTarget.style.background = 'rgba(246,70,93,0.1)';
              e.currentTarget.style.color = '#f6465d';
            }}
            onMouseLeave={(e) => {
              e.currentTarget.style.background = 'transparent';
              e.currentTarget.style.color = '#848e9c';
            }}
          >
            <LogOut size={16} />
            Logout
          </button>
        </div>
      </aside>
      <main className="flex-1 overflow-auto" style={{ background: '#0b0e11' }}>
        <Outlet />
      </main>
    </div>
  );
}
