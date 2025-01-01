import React from 'react';
import { Outlet } from 'react-router-dom';

const Layout: React.FC = () => {
  return (
    <main className="container mx-auto p-4">
      <Outlet />
    </main>
  );
};

export default Layout;