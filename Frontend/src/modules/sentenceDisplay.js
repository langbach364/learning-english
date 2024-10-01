import { highlightSelection, clearHighlight } from '../utils/highlight.js';
import { sendScannedWordsToServer } from '../services/api.js';

function debounce(func, wait) {
  let timeout;
  return function executedFunction(...args) {
    const later = () => {
      clearTimeout(timeout);
      func(...args);
    };
    clearTimeout(timeout);
    timeout = setTimeout(later, wait);
  };
}

const debouncedSendScannedWordsToServer = debounce(sendScannedWordsToServer, 300);

export function initSentenceDisplay(state) {
  const originalSentence = document.querySelector('.original-sentence');
  const editedSentence = document.querySelector('.edited-sentence');
  const scanButton = document.querySelector('#scanButton');

  scanButton.addEventListener('click', () => {
    state.isScanning = !state.isScanning;
    scanButton.textContent = state.isScanning ? 'Dừng quét' : 'Bắt đầu quét';
    if (!state.isScanning) {
      clearHighlight(state);
    }
  });

  [originalSentence, editedSentence].forEach(sentence => {
    sentence.addEventListener('mousedown', (event) => startSelection(event, state));
    sentence.addEventListener('mousemove', (event) => updateSelection(event, state));
    sentence.addEventListener('mouseup', (event) => endSelection(event, state));
    sentence.addEventListener('click', (event) => handleClick(event, state));
  });

  Vue.watch(() => [state.highlightedOriginalWords, state.highlightedEditedWords, state.tooltipPosition], updateTooltip);
}

function startSelection(event, state) {
  if (!state.isScanning) return;
  event.preventDefault();
  event.stopPropagation();
  state.selectionStart = event.target;
  state.selectionEnd = event.target;

  const type = event.target.closest('.original-sentence') ? 'original' : 'edited';
  if (type === 'original') {
    state.highlightedOriginalWords = [event.target.textContent.trim()];
    state.highlightedEditedWords = [];
  } else {
    state.highlightedEditedWords = [event.target.textContent.trim()];
    state.highlightedOriginalWords = [];
  }
  updateTooltipPosition(state, event.target);
}

function updateSelection(event, state) {
  if (!state.isScanning || !state.selectionStart) return;
  event.preventDefault();
  event.stopPropagation();
  state.selectionEnd = event.target;
  const type = event.target.closest('.original-sentence') ? 'original' : 'edited';
  highlightSelection(state, type);
}

function endSelection(event, state) {
  if (!state.isScanning) return;
  event.preventDefault();
  event.stopPropagation();

  const selectedWords = state.highlightedOriginalWords.length ? state.highlightedOriginalWords : state.highlightedEditedWords;
  if (selectedWords.length > 0) {
    debouncedSendScannedWordsToServer(selectedWords);
    state.suppressClickHandler = true;
  }

  state.selectionStart = null;
  state.selectionEnd = null;
}

function handleClick(event, state) {
  if (state.suppressClickHandler) {
    state.suppressClickHandler = false;
    return;
  }

  if (!state.isScanning) return;
  event.preventDefault();
  event.stopPropagation();

  if (event.target.classList.contains("word-span")) {
    const clickedWord = event.target.textContent.trim();
    const type = event.target.closest('.original-sentence') ? 'original' : 'edited';
    if (type === 'original') {
      state.highlightedOriginalWords = [clickedWord];
      state.highlightedEditedWords = [];
    } else {
      state.highlightedEditedWords = [clickedWord];
      state.highlightedOriginalWords = [];
    }
    updateTooltipPosition(state, event.target);

    debouncedSendScannedWordsToServer([clickedWord]);
  } else if (event.target.classList.contains("original-sentence") || event.target.classList.contains("edited-sentence")) {
    const type = event.target.classList.contains("original-sentence") ? 'original' : 'edited';
    if (type === 'original') {
      state.highlightedOriginalWords = state.sentenceInfo["Câu gốc"].split(" ");
      state.highlightedEditedWords = [];
    } else {
      state.highlightedEditedWords = state.sentenceInfo["Ghi lại câu đã sửa"].split(" ");
      state.highlightedOriginalWords = [];
    }
    updateTooltipPosition(state, event.target);

    const fullSentenceWords = state.highlightedOriginalWords.length > 0 ? state.highlightedOriginalWords : state.highlightedEditedWords;
    debouncedSendScannedWordsToServer(fullSentenceWords);
  } else {
    state.highlightedOriginalWords = [];
    state.highlightedEditedWords = [];
  }
}

function updateTooltipPosition(state, element) {
  const rect = element.getBoundingClientRect();
  state.tooltipPosition = {
    x: rect.right,
    y: rect.bottom,
  };
}

function updateTooltip(state) {
  // Implement the logic to update the tooltip here
  // This function will be called whenever the watched values change
}
