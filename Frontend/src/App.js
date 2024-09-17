const app = Vue.createApp({
  data() {
      return {
          definitions: {},
          currentWord: ''
      }
  },
  computed: {
      groupedDefinitions() {
          if (!this.currentWord) return {};
          return this.definitions[this.currentWord].reduce((acc, def) => {
              if (!acc[def['Từ loại']]) {
                  acc[def['Từ loại']] = [];
              }
              acc[def['Từ loại']].push(def);
              return acc;
          }, {});
      }
  },
  mounted() {
      fetch('../Handler/sourcegraph-cody/answer.txt')
          .then(response => response.text())
          .then(data => {
              this.definitions = window.parseDefinitions(data);
              this.currentWord = Object.keys(this.definitions)[0];
          })
          .catch(error => console.error('Error:', error));
  }
});

app.mount('#app');
