import React, { useState, useCallback, useEffect, useRef } from "react";
import { motion, useMotionValue } from "framer-motion";
import { getStatistics } from '../services/api';
import "../styles/Statistics.css";
import "../styles/shared.css";
import { useResize } from '../hooks/useResize';

interface StatisticData {
  timeRange: string;
  words: {
    word: string;
    wrongCount: number;
  }[];
}

const Statistics: React.FC = () => {
  const { isResizing, startResize, size, setSize } = useResize({
    minWidth: 200,
    minHeight: 200
  });

  const [isDraggable, setIsDraggable] = useState(true);
  const [statisticData, setStatisticData] = useState<StatisticData[]>([]);
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
  const handleGetStatistics = useCallback(async () => {
    try {
      setLoading(true);
      const response = await getStatistics();
      if (response.data) {
        setStatisticData(response.data);
      }
    } catch (error) {
      console.error('Lỗi khi tải thống kê:', error);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    handleGetStatistics();
  }, [handleGetStatistics]);
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
      
      <div className="history-header">
        <span>Thống kê học tập</span>
      </div>

      <div ref={contentRef} className={`content-wrapper ${isCollapsed ? 'collapsed' : ''}`}>
        <div className="history-section">
          {loading ? (
            <div className="loading">Đang tải...</div>
          ) : (
            statisticData.map((item, index) => (
              <motion.div
                key={index} 
                className="history-item"
                initial={{ opacity: 0, scale: 0.8 }}
                animate={{ opacity: 1, scale: 1 }}
                transition={{
                  duration: 0.5,
                  type: "spring",
                  stiffness: 100,
                }}
              >
                {/* Nội dung thống kê */}
              </motion.div>
            ))
          )}
        </div>

        {/* Thêm nút vào đây, nằm dưới history-section */}
        <button
          className="statistics-button"
          onClick={handleGetStatistics}
          disabled={loading}
        >
          {loading ? "Đang tải..." : "Xem thống kê"}
        </button>
      </div>
    </motion.div>
  );
};

export default Statistics;
