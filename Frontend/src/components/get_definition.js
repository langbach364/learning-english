fetch('../Handler/sourcegraph-cody/answer.txt')
  .then(response => response.text())
  .then(data => {
    const definitions = window.parseDefinitions(data);
    window.definitions = definitions;
  })
  .catch(error => console.error('Error:', error));
