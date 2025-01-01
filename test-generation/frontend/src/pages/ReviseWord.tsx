import React, { useState, useEffect } from 'react';
import { api } from '../services/api';

interface ReviseWordData {
  // Định nghĩa kiểu dữ liệu cho từ ôn tập
  id: number;
  word: string;
  meaning: string;
}

const ReviseWord = () => {
  const [reviseWords, setReviseWords] = useState<ReviseWordData[]>([]);

  const fetchReviseWords = async () => {
    try {
      const response = await api.post('/revise_word');
      setReviseWords(response.data);
    } catch (error) {
      console.error('Lỗi khi tải từ ôn tập:', error);
    }
  };

  useEffect(() => {
    fetchReviseWords();
  }, []);

  return (
    <div>
      <h1 className="text-2xl font-bold">Ôn tập từ vựng</h1>
      {/* Thêm nội dung ôn tập */}
    </div>
  );
};

export default ReviseWord;