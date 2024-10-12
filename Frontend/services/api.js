export function sendWordToServer(word) {
  return fetch("http://localhost:7089/word", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ data: word }),
  });
}

export function sendScannedWordsToServer(words) {
  return fetch("http://localhost:7089/listen_word", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ data: words.join(' ') }),
  });
}

export function fetchDefinitionsFromServer() {
  return fetch("http://localhost:7089/read_word", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
  })
  .then(response => response.json())
  .then(data => JSON.parse(data.data));
}

