import { fetchDefinitionsFromServer } from "../services/api.js";

const app = Vue.createApp({
  data() {
    return {
      definitions: {},
      newWord: "",
      highlightedWords: [],
      tooltipPosition: { x: 0, y: 0 },
      isScanning: false,
      isSending: false,
    };
  },
  methods: {
    sendWordToServer() {
      if (this.newWord.trim() === "") return;

      this.isSending = true;
      fetch("http://localhost:7089/word", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ data: this.newWord }),
      })
        .then(() => {
          console.log("Word sent successfully");
          this.newWord = "";
          this.fetchDefinitions();
        })
        .catch((error) => console.error("Error sending word:", error))
        .finally(() => {
          this.isSending = false;
        });
    },
    fetchDefinitions() {
      console.log("Fetching definitions...");
      fetchDefinitionsFromServer()
        .then((data) => {
          this.definitions = this.processDefinitions(data);
          console.log("Definitions fetched and processed");
        })
        .catch((error) => console.error("Error fetching definitions:", error));
    },
    processDefinitions(data) {
      const processed = {};
      for (const [type, definitions] of Object.entries(data)) {
        processed[type.replace('*', '').trim()] = {};
        const sortedDefinitions = Object.entries(definitions).sort((a, b) => {
          const numA = parseInt(a[0].match(/\((\d+)\)/)[1]);
          const numB = parseInt(b[0].match(/\((\d+)\)/)[1]);
          return numA - numB;
        });
        for (const [def, examples] of sortedDefinitions) {
          const [definition, number] = def.replace('+', '').split('(');
          const key = definition.trim();
          processed[type.replace('*', '').trim()][key] = {
            number: number ? number.replace(')', '').trim() : '',
            examples: {
              VI: examples.VI ? examples.VI.map(e => e.replace(/^Ví dụ: /, '').trim()) : [],
              EN: examples.EN ? examples.EN.map(e => e.replace(/^Ví dụ: /, '').trim()) : []
            }
          };
        }
      }
      return processed;
    }
    ,
    toggleScanning() {
      this.isScanning = !this.isScanning;
      if (!this.isScanning) {
        this.clearHighlight();
      }
      console.log("Scanning mode:", this.isScanning ? "ON" : "OFF");
    },
    handleWordClick(word) {
      if (!this.isScanning) return;
      this.highlightedWords = [word];
      this.updateTooltipPosition(event.target);
      this.sendScannedWordsToServer([word]);
    },
    sendScannedWordsToServer(words) {
      if (this.isSending) return;
      this.isSending = true;
      fetch("http://localhost:7089/listen_word", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ data: words.join(" ") }),
      })
        .then((response) => {
          if (!response.ok) {
            throw new Error("Network response was not ok");
          }
          return response.text();
        })
        .then((data) => {
          console.log("Server response:", data);
        })
        .catch((error) => console.error("Error sending scanned words:", error))
        .finally(() => {
          this.isSending = false;
        });
    },
    clearHighlight() {
      this.highlightedWords = [];
    },
    updateTooltipPosition(element) {
      const rect = element.getBoundingClientRect();
      this.tooltipPosition = {
        x: rect.right,
        y: rect.bottom,
      };
    },
  },
  created() {
    this.fetchDefinitions();
    this.isScanning = JSON.parse(localStorage.getItem("isScanning")) || false;
  },
  watch: {
    isScanning(newValue) {
      localStorage.setItem("isScanning", JSON.stringify(newValue));
    },
  },
});

app.mount("#app");
