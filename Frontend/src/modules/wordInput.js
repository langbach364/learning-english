import { sendWordToServer } from '../services/api.js';

export function initWordInput(state) {
  const input = document.querySelector('#newWord');
  const sendButton = document.querySelector('#sendWord');

  sendButton.addEventListener('click', (event) => {
    event.preventDefault();
    sendWordToServer(input.value)
      .then(() => {
        console.log("Success");
        input.value = "";
        fetchDefinitions(state);
      })
      .catch((error) => console.error("Error:", error));
  });
}
