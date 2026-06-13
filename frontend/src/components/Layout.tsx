import { NavLink, Outlet } from 'react-router-dom'
import { useAuth } from '../hooks/useAuth'

export default function Layout() {
  const { user, logout } = useAuth()

  return (
    <div className="layout">
      <aside className="sidebar">
        <div className="sidebar-logo">
          <span className="logo-icon">QW</span>
          <span className="logo-text">Trading</span>
        </div>
        <nav className="sidebar-nav">
          <NavLink to="/trade" className={({ isActive }) => `nav-link${isActive ? ' active' : ''}`}>
            <span className="nav-icon">&#9650;&#9660;</span>
            Trade
          </NavLink>
          <NavLink to="/portfolio" className={({ isActive }) => `nav-link${isActive ? ' active' : ''}`}>
            <span className="nav-icon">&#9632;</span>
            Portfolio
          </NavLink>
          <NavLink to="/history" className={({ isActive }) => `nav-link${isActive ? ' active' : ''}`}>
            <span className="nav-icon">&#9776;</span>
            History
          </NavLink>
        </nav>
        <div className="sidebar-footer">
          <div className="user-info">
            <span className="user-avatar">{user?.username?.[0]?.toUpperCase() ?? '?'}</span>
            <span className="user-name">{user?.username ?? 'Guest'}</span>
          </div>
          <button className="btn-logout" onClick={logout}>Logout</button>
        </div>
      </aside>
      <main className="main-content">
        <Outlet />
      </main>
    </div>
  )
}
