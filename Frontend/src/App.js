const app = Vue.createApp({
  data() {
    return {
      sentenceInfo: {},
      wordsInfo: {}
    }
  },
  mounted() {
    fetch('../Handler/sourcegraph-cody/answer.txt')
      .then(response => response.text())
      .then(data => {
        const parsedData = window.parseDefinitions(data);
        this.sentenceInfo = parsedData.sentence;
        this.wordsInfo = parsedData.words;
      })
      .catch(error => console.error('Error:', error));
  }
});

app.mount('#app');
