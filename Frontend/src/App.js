import { fetchDefinitionsFromServer, sendWordToServer, sendScannedWordsToServer } from "../services/api.js";

const app = Vue.createApp({
  data() {
    return {
      definitions: {},
      sentenceStructure: {},
      newWord: "",
      highlightedWords: [],
      tooltipPosition: { x: 0, y: 0 },
      isScanning: false,
      isSending: false,
      dataType: null,
    };
  },
  methods: {
    sendWordToServer() {
      if (this.newWord.trim() === "") return;

      this.isSending = true;
      sendWordToServer(this.newWord)
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
          if (data['* Câu:']) {
            this.sentenceStructure = this.processSentenceStructure(data);
            this.dataType = 'sentence';
          } else {
            this.definitions = this.processWordDefinitions(data);
            this.dataType = 'word';
          }
          console.log("Definitions fetched and processed");
        })
        .catch((error) => console.error("Error fetching definitions:", error));
    },
    processWordDefinitions(data) {
      const processed = {};
      for (const [type, definitions] of Object.entries(data)) {
        processed[type.replace(/^\*\s*/, '')] = {};
        for (const [def, examples] of Object.entries(definitions)) {
          const newDef = def.replace(/^[+*]\s*/, '');
          processed[type.replace(/^\*\s*/, '')][newDef] = {
            examples: {
              VI: examples && examples.VI ? examples.VI.map(e => e.replace(/^Ví dụ: /, '').trim()) : [],
              EN: examples && examples.EN ? examples.EN.map(e => e.replace(/^Ví dụ: /, '').trim()) : []
            }
          };
        }
      }
      return processed;
    
    },
    processSentenceStructure(data) {
      const processed = {};
      const order = ['* Câu:', '* Ghi lại câu chưa sửa:', '* Ghi lại câu đã sửa:'];
      
      // Xử lý các khóa ưu tiên trước
      for (const key of order) {
        if (data[key]) {
          processed[key.replace(/^\*\s*/, '')] = data[key];
        }
      }
      
      // Xử lý các khóa còn lại
      for (const [key, value] of Object.entries(data)) {
        if (!order.includes(key)) {
          const newKey = key.replace(/^\*\s*/, '');
          if (Array.isArray(value)) {
            processed[newKey] = value.map(item => item.replace(/^[+*]\s*/, '').replace(/^\{|\}$/g, ''));
          } else if (typeof value === 'string') {
            processed[newKey] = value.replace(/^\{|\}$/g, '');
          } else {
            processed[newKey] = value;
          }
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
      sendScannedWordsToServer(words)
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
