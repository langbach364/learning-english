import React, { useState, useEffect } from 'react';
import { api } from '../services/api';

interface ScheduleData {
  // Định nghĩa kiểu dữ liệu cho lịch học
  id: number;
  date: string;
  words: string[];
}

const Schedule = () => {
  const [schedule, setSchedule] = useState<ScheduleData[]>([]);

  const fetchSchedule = async () => {
    try {
      const response = await api.post('/create_schedule');
      setSchedule(response.data);
    } catch (error) {
      console.error('Lỗi khi tải lịch học:', error);
    }
  };

  useEffect(() => {
    fetchSchedule();
  }, []);

  return (
    <div>
      <h1 className="text-2xl font-bold">Lịch học</h1>
      {/* Thêm nội dung lịch học */}
    </div>
  );
};

export default Schedule;