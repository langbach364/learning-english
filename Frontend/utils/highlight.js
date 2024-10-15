export function highlightWord(state, word, index, sentenceType, itemIndex = null) {
  state.currentScannedWord = { word: word.trim(), index, sentenceType, itemIndex };
  state.highlightedWords = [word.trim()];
  updateTooltipPosition(state, document.querySelector(`[data-word-index="${index}"][data-sentence-type="${sentenceType}"]${itemIndex !== null ? `[data-item-index="${itemIndex}"]` : ''}`));
}

export function clearHighlight(state) {
  state.currentScannedWord = null;
  state.highlightedWords = [];
}

export function highlightSentence(state, sentence) {
  state.sentence = sentence;
  state.currentScannedWord = null;
  state.highlightedWords = sentence.split(' ');
}

export function updateTooltipPosition(state, element) {
  if (element && state.isScanning) {
    const rect = element.getBoundingClientRect();
    state.tooltipPosition = {
      x: rect.right,
      y: rect.bottom,
    };
  }
}

export function getWordClass(state, word, index, sentenceType, itemIndex = null) {
  const isCurrentWord = state.currentScannedWord &&
    state.currentScannedWord.word === word &&
    state.currentScannedWord.sentenceType === sentenceType &&
    (itemIndex === null || state.currentScannedWord.itemIndex === itemIndex);

  return {
    'cursor-pointer': state.isScanning,
    'bg-green-200': state.isScanning && isCurrentWord,
  };
}