export function initScanningTooltip(state) {
  const tooltipContainer = document.querySelector('#scanningTooltip');

  function updateTooltip() {
    if (state.isScanning && state.currentScannedWord) {
      const { word } = state.currentScannedWord;
      let html = `<h3 class="text-xl font-bold mb-2">Đang quét: ${word}</h3>`;

      if (state.wordsInfo[word]) {
        html += `<div class="mt-4">
          <h4 class="font-semibold text-blue-600">Thông tin từ</h4>
          <ul class="list-disc list-inside">`;
        for (const definition of state.wordsInfo[word]) {
          html += `<li>${definition.definition}</li>`;
        }
        html += `</ul></div>`;
      }

      tooltipContainer.innerHTML = html;
      tooltipContainer.style.display = 'block';
      tooltipContainer.style.top = `${state.tooltipPosition.y}px`;
      tooltipContainer.style.left = `${state.tooltipPosition.x}px`;
    } else {
      tooltipContainer.style.display = 'none';
    }
  }

  setInterval(updateTooltip, 100);
}
