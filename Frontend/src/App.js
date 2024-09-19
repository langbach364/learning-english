const app = Vue.createApp({
    data() {
      return {
        sentenceInfo: {}
      }
    },
    mounted() {
      fetch('../Handler/sourcegraph-cody/answer.txt')
        .then(response => response.text())
        .then(data => {
          const parsedData = window.parseDefinitions(data);
          this.sentenceInfo = parsedData.sentence;
        })
        .catch(error => console.error('Error:', error));
    }
  });
  
  app.mount('#app');
  