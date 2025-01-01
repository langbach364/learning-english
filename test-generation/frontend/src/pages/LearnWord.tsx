import React, { useState, useCallback } from "react";
import { motion, useMotionValue } from "framer-motion";
import gsap from "gsap";
import { learnWord } from "../services/api";
import "../styles/LearnWord.css";

interface Word {
  id: number;
  word: string;
  meaning: string;
}

const LearnWord: React.FC = () => {
  const [words, setWords] = useState<Word[]>([]);
  const [loading, setLoading] = useState(false);
  
  const x = useMotionValue(0);
  const y = useMotionValue(0);

  const handleLearnClick = useCallback(async () => {
    try {
      setLoading(true);
      const data = await learnWord(5);
      setWords(data as Word[]);

      gsap.from(".vocabulary-card", {
        scale: 0.9,
        opacity: 0,
        duration: 0.5,
        ease: "back.out",
      });

      gsap.from(".vocabulary-item", {
        y: 50,
        opacity: 0,
        stagger: 0.2,
        duration: 0.8,
        ease: "elastic.out(1, 0.5)",
      });
    } catch (error) {
      console.error("Error:", error);
    } finally {
      setLoading(false);
    }
  }, []);

  return (
    <motion.main 
      className="draggable-container"
      drag
      dragMomentum={false}
      style={{ x, y }}
    >
      <div className="container mx-auto p-4">
        <div className="vocabulary-card">
          <div className="vocabulary-content">
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
                    <strong className="word-text">{word.word}</strong>
                    <p className="word-meaning">{word.meaning}</p>
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
          </div>
        </div>
      </div>
    </motion.main>
  );
};

export default LearnWord;