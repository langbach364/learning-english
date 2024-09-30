export function highlightSelection(state, type) {
  if (!state.isScanning || !state.selectionStart || !state.selectionEnd) return;

  const container = type === 'original' ? document.querySelector(".original-sentence") : document.querySelector(".edited-sentence");
  const words = Array.from(container.querySelectorAll("span"));
  const startIndex = words.indexOf(state.selectionStart);
  const endIndex = words.indexOf(state.selectionEnd);

  const start = Math.min(startIndex, endIndex);
  const end = Math.max(startIndex, endIndex);

  const highlightedWords = words
    .slice(start, end + 1)
    .map((span) => span.textContent.trim());

  if (type === 'original') {
    state.highlightedOriginalWords = highlightedWords;
    state.highlightedEditedWords = [];
  } else {
    state.highlightedEditedWords = highlightedWords;
    state.highlightedOriginalWords = [];
  }

  updateTooltipPosition(state, state.selectionEnd);
}

export function clearHighlight(state) {
  state.highlightedOriginalWords = [];
  state.highlightedEditedWords = [];
  state.selectionStart = null;
  state.selectionEnd = null;
}

function updateTooltipPosition(state, element) {
  const rect = element.getBoundingClientRect();
  state.tooltipPosition = {
    x: rect.right,
    y: rect.bottom,
  };
}
