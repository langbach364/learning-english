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
  background: white;
  padding: 1.5rem;
  border-radius: 0.75rem;
  margin-bottom: 1.5rem;

  @media (max-width: 480px) {
    padding: 0.75rem;
    margin-bottom: 1rem;
  }

  &:hover {
    transform: translateY(-4px);
    box-shadow: 0 8px 16px rgba(0, 0, 0, 0.1);
  }
`;

const DefinitionBlock = styled.div`
  animation: ${slideLeft} 0.5s ease-out;
  transition: all 0.3s ease;
  margin-left: 1.5rem;
  margin-bottom: 1rem;

  @media (max-width: 480px) {
    margin-left: 0.75rem;
    margin-bottom: 0.75rem;
  }

  &:hover {
    transform: scale(1.01);
  }
`;

const ExampleContainer = styled.div`
  animation: ${fadeIn} 0.7s ease-out;
  margin-left: 1rem;
  margin-top: 0.5rem;

  @media (max-width: 480px) {
    margin-left: 0.5rem;
    margin-top: 0.25rem;
  }

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
      .replace(/\((?:\d+|vi|en)\)/gi, '')
      .trim();
  
    return cleanedText.charAt(0).match(/[A-Za-z]/) 
      ? cleanedText.charAt(0).toUpperCase() + cleanedText.slice(1).toLowerCase()
      : cleanedText;
  };
  
  const renderWords = (text: string): JSX.Element[] => {
    return text.split(" ").map((word, index) => (
      <span
        key={index}
        className="selectable-text mx-0.5 md:mx-1 hover:bg-blue-100 hover-float text-sm md:text-base"
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
      <DefinitionBlock>
        <div className="text-indigo-800 font-medium mb-1 md:mb-2 text-sm md:text-base">
          {renderWords(definition)}
        </div>
        {examples.EN && (
          <div className="space-y-1 md:space-y-2">
            {examples.EN.map((en, idx) => (
              <div
                key={idx}
                className="bg-blue-50 p-2 md:p-3 rounded flex flex-col md:flex-row items-start gap-2 md:gap-8"
              >
                <div className="text-blue-700 flex-1 min-w-0 text-sm md:text-base">
                  {renderWords(cleanValue(en))}
                </div>
                {examples.VI && examples.VI[idx] && (
                  <div className="text-gray-600 flex-1 min-w-0 text-sm md:text-base">
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
      <WordTypeCard key={index}>
        <div className="bg-gradient-to-b from-white to-blue-50">
          <h3 className="font-bold text-blue-900 text-base md:text-lg mb-2 md:mb-4">
            {cleanKey(wordType)}
          </h3>
          <div className="space-y-2 md:space-y-4">
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
    <AnimatedContainer className="space-y-4 md:space-y-6">
      {renderContent()}
    </AnimatedContainer>
  );
};

export default DictionaryDisplay;
