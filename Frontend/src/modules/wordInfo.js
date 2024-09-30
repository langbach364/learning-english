export function initWordInfo(state) {
  const wordInfoContainer = document.querySelector('#wordInfo');

  function updateWordInfo() {
    let html = '<h2 class="text-2xl font-bold mb-2">Thông tin từ:</h2>';
    for (const [type, typeInfo] of Object.entries(state.wordsInfo)) {
      html += `<div class="mb-6">
        <h3 class="text-xl font-semibold mb-2">${type}</h3>`;
      for (const [word, wordInfo] of Object.entries(typeInfo)) {
        html += `<div class="mb-4">
          <h4 class="text-lg font-medium mb-2">${word}</h4>`;
        for (const definition of wordInfo) {
          html += `<div class="mb-2">
            <p>${definition.definition}</p>`;
          if (definition.examples && definition.examples.length) {
            html += `<ul class="list-disc list-inside ml-4">
              ${definition.examples.map(example => `<li>${example}</li>`).join('')}
            </ul>`;
          }
          html += `</div>`;
        }
        html += `</div>`;
      }
      html += `</div>`;
    }
    wordInfoContainer.innerHTML = html;
  }

  // Observe state changes and update word info
  setInterval(() => updateWordInfo(), 1000);
}
