import React, { useState } from 'react';
import styled from '@emotion/styled';
import { css } from '@emotion/react';
import { GrammarDetail } from '../types/dictionary';
import { slideDown, fadeIn, scaleUp } from '../styles/animation';

interface GrammarDisplayProps {
  content: {
    detail: GrammarDetail;
    structure: "WordClass" | "Sentence";
  };
}

const KeyContainer = styled.div`
  animation: ${fadeIn} 0.5s ease-out;
  background: white;
  padding: 1.5rem;
  border-radius: 0.75rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  border: 1px solid #e2e8f0;
  margin-bottom: 1.5rem;

  @media (max-width: 480px) {
    padding: 0.75rem;
    margin-bottom: 0.75rem;
  }
`;

const KeyButton = styled.button<{ isSelected: boolean }>`
  transition: all 0.3s ease;
  padding: 0.75rem;
  border-radius: 0.5rem;
  text-align: left;
  width: 100%;
  background: ${props => props.isSelected ? '#3B82F6' : '#EFF6FF'};
  color: ${props => props.isSelected ? 'white' : '#1D4ED8'};
  font-size: var(--font-size-base);

  @media (max-width: 480px) {
    padding: 0.5rem;
    font-size: var(--font-size-sm);
  }

  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  }

  &:active {
    transform: translateY(0);
  }
`;

const ContentCard = styled.div`
  animation: ${slideDown} 0.4s ease-out, ${scaleUp} 0.3s ease-out;
  background: white;
  padding: 1.5rem;
  border-radius: 0.75rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  border: 1px solid #e2e8f0;
  margin-bottom: 1rem;
  transition: transform 0.3s ease;

  @media (max-width: 480px) {
    padding: 0.75rem;
    margin-bottom: 0.75rem;
  }

  &:hover {
    transform: translateY(-4px);
  }
`;

const ContentItem = styled.div<{ index: number }>`
  animation: ${fadeIn} ${props => 0.3 + props.index * 0.1}s ease-out;
`;

const GrammarDisplay: React.FC<GrammarDisplayProps> = ({ content }) => {
  const [selectedKeys, setSelectedKeys] = useState<Set<string>>(new Set());

  const cleanKey = (text: string): string => {
    const cleanedText = text
      .replace(/:/g, '')
      .replace(/^[*+\s]+/, '')
      .replace(/\([^)]*\)/g, '')
      .replace(/^Ví dụ/i, '')
      .trim();

    return cleanedText.charAt(0).toUpperCase() + cleanedText.slice(1).toLowerCase();
  };

  const cleanValue = (text: string): string => {
    const cleanedText = text
      .replace(/^[*:+\s]+/, '')
      .replace(/^Ví dụ:?\s*/i, '')
      .trim();
  
    return cleanedText.charAt(0).match(/[A-Za-z]/) 
      ? cleanedText.charAt(0).toUpperCase() + cleanedText.slice(1).toLowerCase()
      : cleanedText;
  };
  
  const toggleKey = (key: string) => {
    const newSelectedKeys = new Set(selectedKeys);
    selectedKeys.has(key) ? newSelectedKeys.delete(key) : newSelectedKeys.add(key);
    setSelectedKeys(newSelectedKeys);
  };

  const renderWords = (text: string): JSX.Element[] => {
    return text.split(" ").map((word, index) => (
      <span key={index} className="selectable-text mx-1 hover:bg-blue-100">
        {word}
      </span>
    ));
  };

  const renderKeyButtons = () => (
    <KeyContainer>
      <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-2 md:gap-4">
        {Object.keys(content.detail).map((key) => (
          <KeyButton
            key={key}
            isSelected={selectedKeys.has(key)}
            onClick={() => toggleKey(key)}
          >
            {cleanKey(key)}
          </KeyButton>
        ))}
      </div>
    </KeyContainer>
  );

  const renderContent = () => {
    if (selectedKeys.size === 0) return null;

    return Array.from(selectedKeys).map((key) => {
      const values = content.detail[key];
      return (
        <ContentCard key={key}>
          <h3 className="font-bold text-blue-900 text-base md:text-lg mb-2 md:mb-4">
            {cleanKey(key)}
          </h3>
          <div className="space-y-2 md:space-y-4">
            {values.map((item, idx) => (
              <ContentItem key={idx} index={idx} className="ml-2 md:ml-4">
                <div className="bg-blue-50 p-2 md:p-3 rounded">
                  <div className="text-blue-700 text-sm md:text-base">
                    {renderWords(cleanValue(item))}
                  </div>
                </div>
              </ContentItem>
            ))}
          </div>
        </ContentCard>
      );
    });
  };

  return (
    <div>
      {renderKeyButtons()}
      {renderContent()}
    </div>
  );
};

export default GrammarDisplay;
