const API_BASE_URL = "http://localhost:7089";

export function sendWordToServer(word) {
  return fetch(`${API_BASE_URL}/word`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ data: word }),
  });
}

export function sendScannedWordsToServer(words) {
  return fetch(`${API_BASE_URL}/listen_word`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ data: words.join(' ') }),
  });
}

export function fetchDefinitionsFromServer() {
  return fetch(`${API_BASE_URL}/read_word`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
  })
  .then(response => response.json())
  .then(data => JSON.parse(data.data));
}

export function fetchSentenceAnalysis(sentence) {
  return fetch(`${API_BASE_URL}/analyze_sentence`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ data: sentence }),
  })
  .then(response => response.json());
}

export function fetchWordInfo(word) {
  return fetch(`${API_BASE_URL}/word_info`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ data: word }),
  })
  .then(response => response.json());
}
