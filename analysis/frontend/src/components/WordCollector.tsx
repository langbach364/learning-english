import React, { useEffect, useState } from 'react';
import { DictionaryWordType, GrammarDetail } from '../types/dictionary';

declare global {
  interface Window {
    responsiveVoice: {
      speak: (text: string, voice: string, options?: any) => void;
      cancel: () => void;
      voiceSupport: () => boolean;
      isPlaying: () => boolean;
      init: (apiKey: string) => void;
    };
  }
}

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
  const [isVoiceReady, setIsVoiceReady] = useState(false);

  useEffect(() => {
    if (window.responsiveVoice) {
      window.responsiveVoice.init(process.env.REACT_APP_RESPONSIVE_VOICE_KEY || '');
      
      const checkVoiceReady = setInterval(() => {
        if (window.responsiveVoice.voiceSupport()) {
          setIsVoiceReady(true);
          clearInterval(checkVoiceReady);
        }
      }, 100);
    }
  }, []);

  const speakText = (text: string) => {
    if (isVoiceReady) {
      if (window.responsiveVoice.isPlaying()) {
        window.responsiveVoice.cancel();
      }
      window.responsiveVoice.speak(text, "US English Female", {
        pitch: 1,
        rate: 1,
        volume: 1
      });
    }
  };

  const clearSelection = () => {
    setSelectedWords([]);
    setLastSelectedPosition(null);
    setShowSelectedWords(false);
    if (isVoiceReady && window.responsiveVoice.isPlaying()) {
      window.responsiveVoice.cancel();
    }
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

  const toggleCollecting = () => {
    const newState = !isCollectingEnabled;
    setIsCollectingEnabled(newState);
    document.body.classList.toggle('collecting', newState);
    if (!newState) {
      clearSelection();
    }
  };

  useEffect(() => {
    const handleInteractionStart = (event: MouseEvent | TouchEvent) => {
      if (!isCollectingEnabled) return;
      const target = ('touches' in event ? event.touches[0].target : event.target) as HTMLElement;
      if (target.classList.contains('selectable-text')) {
        setIsDragging(true);
        setStartElement(target);
        target.classList.add('active-word');
        setActiveElements(new Set([target]));
      }
    };

    const handleInteractionMove = (event: MouseEvent | TouchEvent) => {
      if (!isCollectingEnabled || !isDragging || !startElement) return;
      
      const target = ('touches' in event 
        ? document.elementFromPoint(
            event.touches[0].clientX,
            event.touches[0].clientY
          )
        : event.target) as HTMLElement;

      if (target?.classList.contains('selectable-text')) {
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

    const handleInteractionEnd = (event: MouseEvent | TouchEvent) => {
      if (!isCollectingEnabled) return;
      
      const words = Array.from(activeElements)
        .map(element => element.textContent?.trim() || '')
        .filter(word => word.length > 0);
      
      if (words.length > 0) {
        setSelectedWords(words);
        setShowSelectedWords(true);
        speakText(words.join(' '));
        
        const position = 'changedTouches' in event
          ? event.changedTouches[0]
          : event;
          
        setLastSelectedPosition({
          x: position.clientX,
          y: position.clientY + window.scrollY
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
      document.addEventListener('mousedown', handleInteractionStart);
      document.addEventListener('mousemove', handleInteractionMove);
      document.addEventListener('mouseup', handleInteractionEnd);
      document.addEventListener('touchstart', handleInteractionStart);
      document.addEventListener('touchmove', handleInteractionMove);
      document.addEventListener('touchend', handleInteractionEnd);
    }

    return () => {
      document.removeEventListener('mousedown', handleInteractionStart);
      document.removeEventListener('mousemove', handleInteractionMove);
      document.removeEventListener('mouseup', handleInteractionEnd);
      document.removeEventListener('touchstart', handleInteractionStart);
      document.removeEventListener('touchmove', handleInteractionMove);
      document.removeEventListener('touchend', handleInteractionEnd);
    };
  }, [isDragging, activeElements, startElement, isCollectingEnabled, isVoiceReady]);

  return (
    <>
      <button
        onClick={toggleCollecting}
        className={`px-3 md:px-6 py-1 md:py-2 text-sm md:text-base rounded-lg transition-colors ${
          isCollectingEnabled 
            ? 'bg-blue-600 text-white hover:bg-blue-700' 
            : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
        }`}
      >
        {isCollectingEnabled ? 'Tắt quét từ' : 'Bật quét từ'}
      </button>

      {showSelectedWords && selectedWords.length > 0 && lastSelectedPosition && (
        <div 
          className="fixed bg-white rounded-xl shadow-lg border border-gray-100 p-3 md:p-4 max-w-[90vw] md:max-w-md"
          style={{
            top: `${lastSelectedPosition.y + 20}px`,
            left: `${lastSelectedPosition.x}px`,
            transform: 'translateX(-50%)',
            zIndex: 1000
          }}
        >
          <div className="flex justify-between items-start">
            <p className="text-gray-800 pr-4 text-sm md:text-base">
              {selectedWords.join(' ')}
            </p>
            <button 
              onClick={clearSelection}
              className="text-gray-500 hover:text-gray-700"
            >
              <svg xmlns="http://www.w3.org/2000/svg" className="h-4 w-4 md:h-5 md:w-5" viewBox="0 0 20 20" fill="currentColor">
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
