import React from "react";
import styled from "@emotion/styled";
import { DictionaryWordType, DictionaryExample } from "../types/dictionary";
import {
  slideDown,
  fadeIn,
  bounce,
  swing,
  pulse,
  slideLeft,
} from "../styles/animation";

interface DictionaryDisplayProps {
  content: DictionaryWordType;
  structure: "WordClass" | "Sentence";
}

const AnimatedContainer = styled.div`
  animation: ${fadeIn} 0.5s ease-out;
`;

const WordTypeCard = styled.div`
  animation: ${slideDown} 0.6s ease-out;
  transition: all 0.3s ease;
  &:hover {
    transform: translateY(-4px);
    box-shadow: 0 8px 16px rgba(0, 0, 0, 0.1);
  }
`;

const DefinitionBlock = styled.div`
  animation: ${slideLeft} 0.5s ease-out;
  transition: all 0.3s ease;
  &:hover {
    transform: scale(1.01);
  }
`;

const ExampleContainer = styled.div`
  animation: ${fadeIn} 0.7s ease-out;
  &:hover {
    animation: ${swing} 1s ease-in-out;
  }
`;

const DictionaryDisplay: React.FC<DictionaryDisplayProps> = ({
  content,
  structure,
}) => {
  const cleanKey = (text: string): string => {
    const cleanedText = text
      .replace(/^[*:+\s]+/, "")
      .replace(/^Ví dụ/i, "")
      .trim();
    return (
      cleanedText.charAt(0).toUpperCase() + cleanedText.slice(1).toLowerCase()
    );
  };

  const cleanValue = (text: string): string => {
    const cleanedText = text
      .replace(/^[*:+\s]+/, '')
      .replace(/^Ví dụ:?\s*/i, '')
      .replace(/\((?:\d+|vi|en)\)/gi, '')  // Chỉ xóa (số), (vi), (en)
      .trim();
  
    return cleanedText.charAt(0).match(/[A-Za-z]/) 
      ? cleanedText.charAt(0).toUpperCase() + cleanedText.slice(1).toLowerCase()
      : cleanedText;
  };
  
  const renderWords = (text: string): JSX.Element[] => {
    return text.split(" ").map((word, index) => (
      <span
        key={index}
        className="selectable-text mx-1 hover:bg-blue-100 hover-float"
      >
        {word}
      </span>
    ));
  };

  const renderDefinitionBlock = (
    definition: string,
    examples: DictionaryExample
  ) => {
    return (
      <DefinitionBlock className="ml-6 mb-4">
        <div className="text-indigo-800 font-medium mb-2">
          {renderWords(definition)}
        </div>
        {examples.EN && (
          <div className="ml-4 space-y-2">
            {examples.EN.map((en, idx) => (
              <div
                key={idx}
                className="bg-blue-50 p-3 rounded flex items-start gap-8"
              >
                <div className="text-blue-700 flex-1 min-w-0">
                  {renderWords(cleanValue(en))}
                </div>
                {examples.VI && examples.VI[idx] && (
                  <div className="text-gray-600 flex-1 min-w-0">
                    {renderWords(cleanValue(examples.VI[idx]))}
                  </div>
                )}
              </div>
            ))}
          </div>
        )}
      </DefinitionBlock>
    );
  };

  const renderContent = () => {
    if (structure !== "WordClass") return null;

    return Object.entries(content).map(([wordType, definitions], index) => (
      <WordTypeCard
        key={index}
        className="bg-white rounded-xl shadow-md border border-blue-200 hover:shadow-lg transition-shadow mb-6"
      >
        <div className="p-6 bg-gradient-to-b from-white to-blue-50">
          <h3 className="font-bold text-blue-900 text-lg mb-4">
            {cleanKey(wordType)}
          </h3>
          <div className="space-y-4">
            {Object.entries(definitions).map(([def, examples], idx) => (
              <div key={idx}>
                {renderDefinitionBlock(cleanValue(def), examples)}
              </div>
            ))}
          </div>
        </div>
      </WordTypeCard>
    ));
  };

  return (
    <AnimatedContainer className="space-y-6">
      {renderContent()}
    </AnimatedContainer>
  );
};

export default DictionaryDisplay;
