import React from 'react';

interface LearnButtonProps {
  onClick: () => void;
  loading: boolean;
}

const LearnButton: React.FC<LearnButtonProps> = ({ onClick, loading }) => {
  return (
    <button 
      onClick={onClick}
      disabled={loading}
      className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 disabled:opacity-50"
    >
      {loading ? 'Đang tải...' : 'Học từ mới'}
    </button>
  );
};

export default LearnButton;
