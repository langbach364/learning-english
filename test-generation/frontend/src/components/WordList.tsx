import React from 'react';

interface Word {
  id: number;
  word: string;
  meaning: string;
}

interface WordListProps {
  words: Word[];
}

const WordList: React.FC<WordListProps> = ({ words }) => {
  return (
    <div className="mt-4 space-y-2">
      {words.map((word) => (
        <div key={word.id} className="p-3 border rounded">
          <h3 className="font-bold">{word.word}</h3>
          <p>{word.meaning}</p>
        </div>
      ))}
    </div>
  );
};

export default WordList;
