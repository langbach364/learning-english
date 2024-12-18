import React, { useEffect } from 'react';

interface ConnectionStatusProps {
  isConnected: boolean;
}

function ConnectionStatus({ isConnected }: ConnectionStatusProps) {
  useEffect(() => {
    const message = isConnected 
      ? 'Kết nối WebSocket thành công!' 
      : 'Mất kết nối WebSocket, đang thử kết nối lại...';
    console.log(message);
  }, [isConnected]);

  return (
    <div 
      className={`fixed top-4 right-4 flex items-center gap-2 rounded-full px-4 py-2 shadow-lg transition-all duration-300 ${
        isConnected ? 'bg-green-50 border-green-500' : 'bg-red-50 border-red-500'
      } border`}
    >
      <div className={`w-3 h-3 rounded-full animate-pulse ${
        isConnected ? 'bg-green-500' : 'bg-red-500'
      }`} />
      <span className={`text-sm font-medium ${
        isConnected ? 'text-green-700' : 'text-red-700'
      }`}>
        {isConnected ? 'Đã kết nối thành công' : 'Mất kết nối'}
      </span>
    </div>
  );
}

export default ConnectionStatus;
