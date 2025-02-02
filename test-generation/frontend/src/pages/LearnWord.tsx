import React, { useState, useCallback } from "react";
import { motion, useMotionValue } from "framer-motion";
import gsap from "gsap";
import { learnWord } from "../services/api";
import "../styles/LearnWord.css";
import { useWords } from "../context/WordContext";
import { Word } from '../types';

const LearnWord: React.FC = () => {
  const { setLearnedWords } = useWords();
  const [words, setWords] = useState<Word[]>([]);
  const [loading, setLoading] = useState(false);
  
  const x = useMotionValue(0);
  const y = useMotionValue(0);

  const handleWrongAnswer = (wordId: number) => {
    setWords(prevWords => {
      const updatedWords = prevWords.map(word => {
        if (word.id === wordId) {
          const newCount = (word.wrongCount || 0) + 1;
          console.log(`Từ "${word.word}" đã sai ${newCount} lần`);
          return { ...word, wrongCount: newCount };
        }
        return word;
      });
      console.log('Danh sách từ và số lần sai:', updatedWords);
      return updatedWords;
    });
  };

  const handleLearnClick = useCallback(async () => {
    try {
      setLoading(true);
      const data = await learnWord(5);
      const wordsWithCount = (data as Word[]).map(word => ({
        ...word,
        wrongCount: 0
      }));
      console.log('Từ mới được tải:', wordsWithCount);
      setWords(wordsWithCount);
      setLearnedWords(wordsWithCount);
    } catch (error) {
      console.error("Error:", error);
    } finally {
      setLoading(false);
    }
  }, [setLearnedWords]);
  return (
    <div className="container mx-auto p-4">
      <motion.div 
        className="vocabulary-content"
        drag
        dragMomentum={false}
        dragElastic={0}
        style={{ x, y }}
      >
        <motion.span className="vocabulary-header">
          <span>Học từ mới</span>
        </motion.span>

        <div className="words-container">
          {words.map((word, index) => (
            <motion.div
              key={word.id || index}
              className="vocabulary-item"
              initial={{ opacity: 0, scale: 0.8 }}
              animate={{ opacity: 1, scale: 1 }}
              transition={{
                delay: index * 0.1,
                duration: 0.5,
                type: "spring",
                stiffness: 100,
              }}
            >
              <div className="word-card">
                <div className="word-info">
                  <strong className="word-text">{word.word}</strong>
                  <span className="wrong-count">
                    {word.wrongCount > 0 ? `(${word.wrongCount})` : ''}
                  </span>
                </div>
              </div>
            </motion.div>
          ))}
        </div>

        <button
          className="learn-button"
          onClick={handleLearnClick}
          disabled={loading}
        >
          {loading ? "Đang tải..." : "Học từ mới"}
        </button>
      </motion.div>
    </div>
  );
};
export default LearnWord;