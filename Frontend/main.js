import { initWordInput } from './src/modules/wordInput.js';
import { initSentenceDisplay } from './src/modules/wordInput.js';
import { initWordInfo } from './src/modules/wordInput.js';
import { initScanningTooltip } from './src/modules/wordInput.js';
import { fetchDefinitions } from './src/services/api.js';

document.addEventListener('DOMContentLoaded', () => {
  const state = {
    sentenceInfo: {},
    wordsInfo: {},
    newWord: "",
    highlightedOriginalWords: [],
    highlightedEditedWords: [],
    tooltipPosition: { x: 0, y: 0 },
    isScanning: false,
    selectionStart: null,
    selectionEnd: null,
    suppressClickHandler: false,
    isSending: false,
  };

  initWordInput(state);
  initSentenceDisplay(state);
  initWordInfo(state);
  initScanningTooltip(state);

  fetchDefinitions(state);
});
