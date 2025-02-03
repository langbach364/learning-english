import React from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { Word } from '../types';

interface WordListProps {
  words: Word[];
  containerWidth: number;
}

const WordList: React.FC<WordListProps> = ({ words, containerWidth }) => {
  const columns = Math.max(1, Math.floor(containerWidth / 250));
  const baseFontSize = Math.max(14, containerWidth * 0.02);

  const cardVariants = {
    hidden: { 
      opacity: 0,
      scale: 0.8 
    },
    visible: { 
      opacity: 1,
      scale: 1,
      transition: {
        duration: 0.3
      }
    },
    exit: {
      opacity: 0,
      scale: 0.5,
      transition: {
        duration: 0.2
      }
    }
  };

  return (
    <div className="words-container">
      <div 
        className="words-grid"
        style={{
          gridTemplateColumns: `repeat(${columns}, 1fr)`,
          fontSize: `${baseFontSize}px`
        }}
      >
        <AnimatePresence>
          {words.map((word, index) => (
            <motion.div
              key={`${word.id}-${index}`}
              className="word-card"
              variants={cardVariants}
              initial="hidden"
              animate="visible"
              exit="exit"
              layout
              style={{
                background: '#ffffff',
                borderRadius: '8px',
                padding: '1rem',
                boxShadow: '0 2px 4px rgba(0,0,0,0.1)',
                display: 'flex',
                flexDirection: 'column',
                gap: '0.5rem'
              }}
            >
              <h3 style={{ 
                margin: 0,
                fontSize: `${baseFontSize * 1.2}px`,
                fontWeight: 600,
                color: '#2d3748'
              }}>
                {word.word}
              </h3>
              
              {word.wrongCount > 0 && (
                <motion.span 
                  className="wrong-count"
                  initial={{ opacity: 0 }}
                  animate={{ opacity: 1 }}
                  style={{
                    fontSize: `${baseFontSize * 0.9}px`,
                    color: '#e53e3e',
                    padding: '0.25rem 0.5rem',
                    borderRadius: '4px',
                    background: '#fff5f5',
                    alignSelf: 'flex-start'
                  }}
                >
                  Sai: {word.wrongCount} lần
                </motion.span>
              )}
            </motion.div>
          ))}
        </AnimatePresence>
      </div>
    </div>
  );
};

export default WordList;