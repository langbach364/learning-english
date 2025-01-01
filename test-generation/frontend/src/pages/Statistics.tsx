import React, { useState, useEffect } from 'react';
import { api } from '../services/api';

interface StatisticsData {
  // Định nghĩa kiểu dữ liệu cho thống kê
  totalWords: number;
  learnedWords: number;
  progress: number;
}

const Statistics = () => {
  const [statistics, setStatistics] = useState<StatisticsData | null>(null);

  const fetchStatistics = async () => {
    try {
      const response = await api.post('/get_statistics');
      setStatistics(response.data);
    } catch (error) {
      console.error('Lỗi khi tải thống kê:', error);
    }
  };

  useEffect(() => {
    fetchStatistics();
  }, []);

  return (
    <div>
      <h1 className="text-2xl font-bold">Thống kê</h1>
      {/* Thêm nội dung thống kê */}
    </div>
  );
};

export default Statistics;