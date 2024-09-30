export function initScanningTooltip(state) {
  const tooltipContainer = document.querySelector('#scanningTooltip');

  function updateTooltip() {
    if (state.highlightedOriginalWords.length || state.highlightedEditedWords.length) {
      const words = state.highlightedOriginalWords.length ? state.highlightedOriginalWords : state.highlightedEditedWords;
      let html = `<h3 class="text-xl font-bold mb-2 tooltip-words">${words.join(' ')}</h3>`;
      for (const word of words) {
        for (const [type, typeInfo] of Object.entries(state.wordsInfo)) {
          if (typeInfo[word]) {
            html += `<div>
              <h4 class="font-semibold text-blue-600">${type}</h4>
              <ul class="list-disc list-inside">`;
            for (const definition of typeInfo[word]) {
              html += `<li>
                ${definition.definition}
                ${definition.examples.length ? `
                  <ul class="list-disc list-inside ml-4 text-gray-600">
                    ${definition.examples.map(example => `<li>${example}</li>`).join('')}
                  </ul>
                ` : ''}
              </li>`;
            }
            html += `</ul></div>`;
          }
        }
      }
      tooltipContainer.innerHTML = html;
      tooltipContainer.style.display = 'block';
      tooltipContainer.style.top = `${state.tooltipPosition.y}px`;
      tooltipContainer.style.left = `${state.tooltipPosition.x}px`;
    } else {
      tooltipContainer.style.display = 'none';
    }
  }

  // Observe state changes and update tooltip
  setInterval(() => updateTooltip(), 100);
}
