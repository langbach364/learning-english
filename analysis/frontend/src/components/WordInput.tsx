import React, { useState, useRef } from 'react';
import WordCollector from './WordCollector';
import { DictionaryWordType, GrammarDetail } from '../types/dictionary';

interface WordInputProps {
  word: string;
  setWord: (word: string) => void;
  onSubmit: (e: React.FormEvent) => void;
  isLoading: boolean;
  data: DictionaryWordType | GrammarDetail;
}

const WordInput: React.FC<WordInputProps> = ({ word, setWord, onSubmit, isLoading, data }) => {
  const [showAlert, setShowAlert] = useState(false);
  const [previousWord, setPreviousWord] = useState('');
  const isSearchSubmit = useRef(false);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    isSearchSubmit.current = true;
    
    if (previousWord && word.toLowerCase() === previousWord.toLowerCase()) {
      setShowAlert(true);
      setTimeout(() => setShowAlert(false), 3000);
      return;
    }

    setPreviousWord(word);
    onSubmit(e);
    isSearchSubmit.current = false;
  };

  return (
    <div className="w-full max-w-4xl mx-auto px-4">
      {showAlert && isSearchSubmit.current && previousWord && (
        <div className="mb-3 md:mb-4 p-3 md:p-4 rounded-lg bg-yellow-50 border border-yellow-200 text-yellow-700 text-sm md:text-base">
          Từ "{word}" đang được hiển thị định nghĩa
        </div>
      )}
      
      <form onSubmit={handleSubmit} className="mb-6 md:mb-8">
        <div className="flex flex-col md:flex-row items-center gap-2 md:gap-4">
          <input
            type="text"
            value={word}
            onChange={(e) => setWord(e.target.value)}
            className="w-full md:flex-1 px-3 md:px-4 py-2 md:py-2.5 text-sm md:text-base border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="Nhập từ cần tra cứu..."
            disabled={isLoading}
          />
          <div className="flex w-full md:w-auto gap-2">
            <button
              type="submit"
              disabled={isLoading}
              className="flex-1 md:flex-none px-4 md:px-6 py-2 md:py-2.5 text-sm md:text-base bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 transition-colors"
            >
              {isLoading ? 'Đang tìm...' : 'Tìm kiếm'}
            </button>
            <WordCollector data={data} />
          </div>
        </div>
      </form>
    </div>
  );
};

export default WordInput;
