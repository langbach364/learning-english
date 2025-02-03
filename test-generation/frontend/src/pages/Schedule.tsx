import React, { useState, useCallback, useEffect, useRef } from "react";
import { motion, useMotionValue } from "framer-motion";
import { createSchedule } from "../services/api";
import "../styles/Schedule.css";
import "../styles/shared.css";
import { Word } from '../types';
import { useWords } from "../context/WordContext";
import { useResize } from '../hooks/useResize';

interface ScheduleData {
  date: string;
  words: Word[];
}

const formatDate = (isoDate: string) => {
  const date = new Date(isoDate);
  return `${date.getDate().toString().padStart(2, '0')}-${(date.getMonth() + 1).toString().padStart(2, '0')}-${date.getFullYear()}`;
};

const Schedule: React.FC = () => {
  const { isResizing, startResize, size, setSize } = useResize({
    minWidth: 200,
    minHeight: 200
  });

  const [isDraggable, setIsDraggable] = useState(true);
  const { learnedWords } = useWords();
  const [scheduleData, setScheduleData] = useState<ScheduleData | null>(null);
  const [loading, setLoading] = useState(false);
  const [isCollapsed, setIsCollapsed] = useState(false);
  const contentRef = useRef<HTMLDivElement>(null);
  
  const x = useMotionValue(0);
  const y = useMotionValue(0);

  useEffect(() => {
    setIsDraggable(!isResizing);
  }, [isResizing]);

  const toggleCollapse = () => {
    setIsCollapsed(prev => !prev);
    if (!isCollapsed) {
      setSize(current => ({...current, height: 60}));
    } else {
      const contentHeight = contentRef.current?.scrollHeight || 400;
      setSize(current => ({...current, height: Math.max(contentHeight + 80, 400)}));
    }
  };

  const handleCreateSchedule = useCallback(async () => {
    if (learnedWords.length === 0) {
      alert('Vui lòng học một số từ mới trước!');
      return;
    }

    try {
      setLoading(true);
      const response = await createSchedule(learnedWords);
      
      if (response) {
        const currentDate = new Date().toISOString().split('T')[0];
        setScheduleData({
          date: formatDate(currentDate),
          words: learnedWords
        });
      }
    } catch (error) {
      console.error('Lỗi tạo lịch:', error);
    } finally {
      setLoading(false);
    }
  }, [learnedWords]);

  return (
    <motion.div 
      className="collapsible-content"
      style={{ 
        x, y,
        width: size.width,
        height: size.height,
        fontSize: `${Math.max(10, size.width * 0.02)}px`,
        overflow: isCollapsed ? 'hidden' : 'visible',
      }}
      drag={isDraggable}
      dragMomentum={false}
      dragElastic={0}
    >
      <button className="collapse-button" onClick={toggleCollapse}>
        {isCollapsed ? '▼' : '▲'}
      </button>

      <div className="resize-handle resize-handle-bottom-right"
           onMouseDown={(e) => startResize(e, 'bottom-right')} />
      <div className="resize-handle resize-handle-bottom-left"
           onMouseDown={(e) => startResize(e, 'bottom-left')} />
      
      <div className="schedule-header">
        <span>Lịch học từ vựng</span>
      </div>

      <div ref={contentRef} className={`content-wrapper ${isCollapsed ? 'collapsed' : ''}`}>
        <div className="schedule-section">
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
                <div className="date-header">
                  <div className="date-badge">
                    <span className="date-day">{scheduleData.date.split('-')[0]}</span>
                    <span className="date-month-year">
                      {scheduleData.date.split('-')[1]}-{scheduleData.date.split('-')[2]}
                    </span>
                  </div>
                </div>
                <ul className="word-list">
                  {scheduleData.words.map((word, idx) => (
                    <li key={idx} className="word-item">
                      <span className="word-text">{word.word}</span>
                      <span className={`wrong-count-badge ${word.wrongCount === 0 ? 'success' : 'error'}`}>
                        {word.wrongCount === 0 ? 'Chưa sai lần nào' : `Sai: ${word.wrongCount} lần`}
                      </span>
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
      </div>
    </motion.div>
  );
};

export default Schedule;
