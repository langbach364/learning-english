import React, { useState, useCallback, useEffect, useRef } from "react";
import { motion, useMotionValue } from "framer-motion";
import { learnWord } from "../services/api";
import "../styles/LearnWord.css";
import "../styles/shared.css";
import { useWords } from "../context/WordContext";
import { Word } from '../types';
import { useResize } from '../hooks/useResize';
import WordList from "../components/WordList";

const LearnWord: React.FC = () => {
  const { isResizing, startResize, size, setSize } = useResize({
    minWidth: 150,
    minHeight: 150
  });

  const [isDraggable, setIsDraggable] = useState(true);
  const [isCollapsed, setIsCollapsed] = useState(false);
  const contentRef = useRef<HTMLDivElement>(null);
  const { setLearnedWords } = useWords();
  const [words, setWords] = useState<Word[]>([]);
  const [loading, setLoading] = useState(false);
  
  const x = useMotionValue(0);
  const y = useMotionValue(0);

  useEffect(() => {
    setIsDraggable(!isResizing);
  }, [isResizing]);

  const handleLearnClick = useCallback(async () => {
    try {
      setLoading(true);
      const data = await learnWord(5);
      const wordsWithCount = (data as Word[]).map(word => ({
        ...word,
        wrongCount: 0
      }));
      setWords(wordsWithCount);
      setLearnedWords(wordsWithCount);
    } catch (error) {
      console.error("Error:", error);
    } finally {
      setLoading(false);
    }
  }, [setLearnedWords]);

  const toggleCollapse = () => {
    setIsCollapsed(prev => !prev);
    if (!isCollapsed) {
      setSize(current => ({...current, height: 60}));
    } else {
      const contentHeight = contentRef.current?.scrollHeight || 400;
      setSize(current => ({...current, height: Math.max(contentHeight + 80, 400)}));
    }
  };
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
      {/* Các nút resize vẫn giữ nguyên vị trí và chức năng */}
      <div 
        className="resize-handle resize-handle-bottom-right"
        onMouseDown={(e) => startResize(e, 'bottom-right')}
      />
      <div 
        className="resize-handle resize-handle-bottom-left"
        onMouseDown={(e) => startResize(e, 'bottom-left')}
      />
      
      {/* Nội dung khác */}
      <button className="collapse-button" onClick={toggleCollapse}>
        {isCollapsed ? '▼' : '▲'}
      </button>
      
      <div className="vocabulary-header">
        <span>Học từ mới</span>
      </div>

      <div ref={contentRef} className={`content-wrapper ${isCollapsed ? 'collapsed' : ''}`}>
        <div className="words-section">
          <WordList 
            words={words} 
            containerWidth={size.width - 40}
          />
        </div>

        <button
          className="learn-button"
          onClick={handleLearnClick}
          disabled={loading}
        >
          {loading ? "Đang tải..." : "Học từ mới"}
        </button>
      </div>
    </motion.div>
  );
}
export default LearnWord;
