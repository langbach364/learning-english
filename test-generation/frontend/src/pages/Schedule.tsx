import React, { useState, useCallback } from "react";
import { motion, useMotionValue } from "framer-motion";
import gsap from "gsap";
import { createSchedule } from "../services/api";
import "../styles/Schedule.css";

import { Word } from '../types';
import { useWords } from "../context/WordContext";

interface ScheduleData {
  date: string;
  words: Word[];
}

const Schedule: React.FC = () => {
  const { learnedWords } = useWords();
  const [scheduleData, setScheduleData] = useState<ScheduleData | null>(null);
  const [loading, setLoading] = useState(false);
  
  const x = useMotionValue(0);
  const y = useMotionValue(0);

  const formatDate = (isoDate: string) => {
    const date = new Date(isoDate);
    const day = date.getDate().toString().padStart(2, '0');
    const month = (date.getMonth() + 1).toString().padStart(2, '0');
    const year = date.getFullYear();
    return `${day}-${month}-${year}`;
  };

  const handleCreateSchedule = useCallback(async () => {
    try {
      setLoading(true);
      if (learnedWords.length === 0) {
        console.log('Chưa có từ nào được học');
        alert('Vui lòng học một số từ mới trước!');
        return;
      }
    
      console.log('Các từ đã học và số lần sai:', learnedWords);
      const response = await createSchedule(learnedWords);
      console.log('Lịch học được tạo:', response.data);
    
      if (response.data) {
        const currentDate = new Date().toISOString().split('T')[0];
        const formattedDate = formatDate(currentDate);
      
        const newSchedule = {
          date: formattedDate, // Sử dụng date đã format
          words: learnedWords
        };
        console.log('Lịch học mới:', newSchedule);
        setScheduleData(newSchedule);
      }
    } catch (error) {
      console.error('Lỗi tạo lịch:', error);
    } finally {
      setLoading(false);
    }
  }, [learnedWords]);  return (
    <div className="container mx-auto p-4">
      <motion.div 
        className="schedule-content"
        drag
        dragMomentum={false}
        dragElastic={0}
        style={{ x, y }}
      >
        <motion.span className="schedule-header">
          <span>Lịch học từ vựng</span>
        </motion.span>

        <div className="schedule-container">
          {scheduleData && (
            <motion.div
              className="schedule-item"
              initial={{ opacity: 0, scale: 0.8 }}
              animate={{ opacity: 1, scale: 1 }}
              transition={{
                duration: 0.5,
                type: "spring",
                stiffness: 100,
              }}
            >
              <div className="date-card">
                <strong className="date-text">{scheduleData.date}</strong>
                <ul className="word-list">
                  {scheduleData.words.map((word, idx) => (
                    <li key={idx} className="word-item">
                      {word.word} {word.wrongCount > 0 ? `(${word.wrongCount})` : ''}
                    </li>
                  ))}
                </ul>
              </div>
            </motion.div>
          )}
        </div>

        <button
          className="schedule-button"
          onClick={handleCreateSchedule}
          disabled={loading}
        >
          {loading ? "Đang tạo lịch..." : "Tạo lịch học"}
        </button>
      </motion.div>
    </div>
  );
};
export default Schedule;