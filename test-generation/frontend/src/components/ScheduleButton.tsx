import React from 'react';

interface ScheduleButtonProps {
  onClick: () => void;
  loading: boolean;
}

const ScheduleButton: React.FC<ScheduleButtonProps> = ({ onClick, loading }) => {
  return (
    <button 
      onClick={onClick}
      disabled={loading}
      className="px-4 py-2 bg-green-500 text-white rounded hover:bg-green-600 disabled:opacity-50"
    >
      {loading ? 'Đang tạo lịch...' : 'Tạo lịch học'}
    </button>
  );
};

export default ScheduleButton;
