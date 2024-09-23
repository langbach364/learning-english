const app = Vue.createApp({
  data() {
    return {
      sentenceInfo: {},
      wordsInfo: {},
      newWord: ''
    }
  },
  methods: {
    sendWordToServer() {
      fetch('http://localhost:7089/word', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ data: this.newWord }),
      })
      .then(response => response.json())
      .then(data => {
        console.log('Success:', data);
        this.newWord = '';
        this.fetchDefinitions();
      })
      .catch((error) => console.error('Error:', error));
    },
    fetchDefinitions() {
      fetch('../Handler/sourcegraph-cody/answer.txt')
        .then(response => response.text())
        .then(data => {
          const parsedData = window.parseDefinitions(data);
          this.sentenceInfo = parsedData.sentence;
          this.wordsInfo = parsedData.words;
        })
        .catch(error => console.error('Error:', error));
    }
  },
  mounted() {
    this.fetchDefinitions();
  }
});

app.mount('#app');
