import React, { createContext, useContext, useState } from 'react';
import { Word } from '../types';

interface WordContextType {
  learnedWords: Word[];
  setLearnedWords: (words: Word[]) => void;
}

const WordContext = createContext<WordContextType | undefined>(undefined);

export const WordProvider: React.FC<{children: React.ReactNode}> = ({ children }) => {
  const [learnedWords, setLearnedWords] = useState<Word[]>([]);

  return (
    <WordContext.Provider value={{ learnedWords, setLearnedWords }}>
      {children}
    </WordContext.Provider>
  );
};

export const useWords = () => {
  const context = useContext(WordContext);
  if (!context) {
    throw new Error('useWords must be used within a WordProvider');
  }
  return context;
};
