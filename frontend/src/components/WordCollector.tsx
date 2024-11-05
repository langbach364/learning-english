import React, { useEffect, useState } from 'react';
import { DictionaryWordType, GrammarDetail } from '../types/dictionary';
import { API_CONFIG } from '../constants/config';

interface WordCollectorProps {
  data: DictionaryWordType | GrammarDetail;
}

const WordCollector: React.FC<WordCollectorProps> = ({ data }) => {
  const [isDragging, setIsDragging] = useState(false);
  const [activeElements, setActiveElements] = useState<Set<HTMLElement>>(new Set());
  const [startElement, setStartElement] = useState<HTMLElement | null>(null);
  const [isCollectingEnabled, setIsCollectingEnabled] = useState(false);
  const [selectedWords, setSelectedWords] = useState<string[]>([]);
  const [lastSelectedPosition, setLastSelectedPosition] = useState<{ x: number; y: number } | null>(null);
  const [showSelectedWords, setShowSelectedWords] = useState(false);

  const clearSelection = () => {
    setSelectedWords([]);
    setLastSelectedPosition(null);
    setShowSelectedWords(false);
  };

  const getAllElementsBetween = (start: HTMLElement, end: HTMLElement) => {
    const elements: HTMLElement[] = [];
    const allSelectableElements = document.querySelectorAll('.selectable-text');
    const startIndex = Array.from(allSelectableElements).indexOf(start);
    const endIndex = Array.from(allSelectableElements).indexOf(end);
    
    const [fromIndex, toIndex] = startIndex < endIndex 
      ? [startIndex, endIndex] 
      : [endIndex, startIndex];
    
    for (let i = fromIndex; i <= toIndex; i++) {
      elements.push(allSelectableElements[i] as HTMLElement);
    }
    
    return elements;
  };

  useEffect(() => {
    const handleMouseDown = (event: MouseEvent) => {
      if (!isCollectingEnabled) return;
      const target = event.target as HTMLElement;
      if (target.classList.contains('selectable-text')) {
        setIsDragging(true);
        setStartElement(target);
        target.classList.add('active-word');
        setActiveElements(new Set([target]));
      }
    };

    const handleMouseEnter = (event: MouseEvent) => {
      if (!isCollectingEnabled || !isDragging || !startElement) return;
      
      const target = event.target as HTMLElement;
      if (target.classList.contains('selectable-text')) {
        activeElements.forEach(element => {
          element.classList.remove('active-word');
        });
        
        const elements = getAllElementsBetween(startElement, target);
        
        const newActiveElements = new Set<HTMLElement>();
        elements.forEach(element => {
          element.classList.add('active-word');
          newActiveElements.add(element);
        });
        
        setActiveElements(newActiveElements);
      }
    };

    const handleMouseUp = async (event: MouseEvent) => {
      if (!isCollectingEnabled) return;
      
      const words = Array.from(activeElements)
        .map(element => element.textContent?.trim() || '')
        .filter(word => word.length > 0);
      
      if (words.length > 0) {
        setSelectedWords(words);
        setShowSelectedWords(true);

        try {
          await fetch(`${API_CONFIG.BASE_URL}/listen_word`, {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify({ data: words.join(' ') }),
          });
        } catch (error) {
          console.error('Lỗi khi gửi từ đã quét:', error);
        }
        
        setLastSelectedPosition({
          x: event.clientX,
          y: event.clientY + window.scrollY
        });
      }

      setIsDragging(false);
      setStartElement(null);
      
      activeElements.forEach(element => {
        element.classList.remove('active-word');
      });
      setActiveElements(new Set());
    };

    if (isCollectingEnabled) {
      document.addEventListener('mousedown', handleMouseDown);
      document.addEventListener('mouseenter', handleMouseEnter, true);
      document.addEventListener('mouseup', handleMouseUp);
    }

    return () => {
      document.removeEventListener('mousedown', handleMouseDown);
      document.removeEventListener('mouseenter', handleMouseEnter, true);
      document.removeEventListener('mouseup', handleMouseUp);
    };
  }, [isDragging, activeElements, startElement, isCollectingEnabled]);

  useEffect(() => {
    if (!isCollectingEnabled) {
      clearSelection();
    }
  }, [isCollectingEnabled]);

  return (
    <>
      <button
        onClick={() => setIsCollectingEnabled(!isCollectingEnabled)}
        className={`px-6 py-2 rounded-lg transition-colors ${
          isCollectingEnabled 
            ? 'bg-blue-600 text-white hover:bg-blue-700' 
            : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
        }`}
      >
        {isCollectingEnabled ? 'Tắt quét từ' : 'Bật quét từ'}
      </button>

      {showSelectedWords && selectedWords.length > 0 && lastSelectedPosition && (
        <div 
          className="fixed bg-white rounded-xl shadow-lg border border-gray-100 p-4 max-w-md"
          style={{
            top: `${lastSelectedPosition.y + 20}px`,
            left: `${lastSelectedPosition.x}px`,
            transform: 'translateX(-50%)',
            zIndex: 1000
          }}
        >
          <div className="flex justify-between items-start">
            <p className="text-gray-800 pr-4">
              {selectedWords.join(' ')}
            </p>
            <button 
              onClick={clearSelection}
              className="text-gray-500 hover:text-gray-700"
            >
              <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                <path fillRule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clipRule="evenodd" />
              </svg>
            </button>
          </div>
        </div>
      )}
    </>
  );
};

export default WordCollector;
