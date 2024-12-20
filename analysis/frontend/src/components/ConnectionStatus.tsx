import React, { useEffect, useRef } from 'react';
import { APIService } from '../services/api';

interface ConnectionStatusProps {
  isConnected: boolean;
}

function ConnectionStatus({ isConnected }: ConnectionStatusProps) {
  const prevConnectedRef = useRef(isConnected);

  useEffect(() => {
    if (prevConnectedRef.current !== isConnected) {
      prevConnectedRef.current = isConnected;
    }
  }, [isConnected]);

  if (!localStorage.getItem('token')) return null;

  return (
    <div 
      className={`fixed top-4 right-4 flex items-center gap-2 rounded-full px-4 py-2 shadow-lg transition-all duration-300 ${
        isConnected ? 'bg-green-50 border-green-500' : 'bg-red-50 border-red-500'
      } border`}
    >
      <div className={`w-3 h-3 rounded-full ${
        isConnected ? 'bg-green-500 animate-pulse' : 'bg-red-500'
      }`} />
      <span className={`text-sm font-medium ${
        isConnected ? 'text-green-700' : 'text-red-700'
      }`}>
        {isConnected ? 'Đã kết nối thành công' : 'Đang kết nối lại...'}
      </span>
    </div>
  );
}

export default React.memo(ConnectionStatus);
