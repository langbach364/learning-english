import React from 'react';
import { motion } from 'framer-motion';

interface StatisticsCardProps {
  timeRange: string;
  words: {
    word: string;
    wrongCount: number;
  }[];
}

const StatisticsCard: React.FC<StatisticsCardProps> = ({ timeRange, words }) => {
  return (
    <motion.div
      className="statistics-item"
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
            <span className="date-day">{timeRange}</span>
          </div>
        </div>
        <ul className="word-list">
          {words.map((word, idx) => (
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
  );
};

export default StatisticsCard;
