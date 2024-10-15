import { fetchDefinitionsFromServer, sendWordToServer, sendScannedWordsToServer } from "../services/api.js";
import { highlightWord, clearHighlight, highlightSentence, getWordClass } from "../utils/highlight.js";
import { initScanningTooltip } from "./modules/scanningTooltip.js";
import { initWordInfo } from "./modules/wordInfo.js";

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
      currentScannedWord: null,
      sentence: "",
      wordsInfo: {},
    };
  },
  methods: {
    separateWords(text) {
      return text.split(/(?=[A-Z0-9])|(?<=[a-z])(?=[0-9])|(?<=[0-9])(?=[a-zA-Z])|[.,;:!?]/)
        .map(word => word.trim())
        .filter(word => word.length > 0)
        .join(' ')
        .split(/\s+/);
    },
    getWordClass(word, index, sentenceType, itemIndex = null) {
      return getWordClass(this, word, index, sentenceType, itemIndex);
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
      
      for (const key of order) {
        if (data[key]) {
          processed[key.replace(/^\*\s*/, '')] = data[key];
        }
      }
      
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
    },
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
            this.sentence = this.sentenceStructure['Câu'];
          } else {
            this.definitions = this.processWordDefinitions(data);
            this.dataType = 'word';
          }
          console.log("Definitions fetched and processed");
        })
        .catch((error) => console.error("Error fetching definitions:", error));
    },
    handleScannedWords(words) {
      if (this.isSending) return;
      this.isSending = true;
      sendScannedWordsToServer(words)
        .then((response) => {
          if (!response.ok) {
            throw new Error("Network response was not ok");
          }
          return response.json();
        })
        .then((data) => {
          console.log("Server response:", data);
          this.wordsInfo = data;
        })
        .catch((error) => console.error("Error sending scanned words:", error))
        .finally(() => {
          this.isSending = false;
        });
    },
    toggleScanning() {
      this.isScanning = !this.isScanning;
      if (!this.isScanning) {
        clearHighlight(this);
      }
      console.log("Scanning mode:", this.isScanning ? "ON" : "OFF");
      document.body.classList.toggle('scanning-active', this.isScanning);
    },
    handleWordClick(word, index, sentenceType, itemIndex = null) {
      if (!this.isScanning) return;
      if (this.currentScannedWord && 
          this.currentScannedWord.word === word && 
          this.currentScannedWord.sentenceType === sentenceType &&
          this.currentScannedWord.itemIndex === itemIndex) {
        return;
      }
      console.log(`Quét từ: ${word} tại vị trí ${index} trong ${sentenceType}${itemIndex !== null ? `, item ${itemIndex}` : ''}`);
      highlightWord(this, word, index, sentenceType, itemIndex);
      this.handleScannedWords([word]);
    },
    handleSentenceScan() {
      if (!this.isScanning) return;
      Object.entries(this.sentenceStructure).forEach(([key, value]) => {
        if (typeof value === 'string') {
          const words = this.separateWords(value);
          words.forEach((word, index) => {
            setTimeout(() => {
              this.handleWordClick(word, index, key);
            }, index * 500);
          });
        } else if (Array.isArray(value)) {
          value.forEach((item, itemIndex) => {
            const words = this.separateWords(item);
            words.forEach((word, wordIndex) => {
              setTimeout(() => {
                this.handleWordClick(word, wordIndex, key, itemIndex);
              }, (itemIndex * words.length + wordIndex) * 500);
            });
          });
        }
      });
    },
  },
  created() {
    this.fetchDefinitions();
    this.isScanning = JSON.parse(localStorage.getItem("isScanning")) || false;
    initScanningTooltip(this);
    initWordInfo(this);
  },
  watch: {
    isScanning(newValue) {
      localStorage.setItem("isScanning", JSON.stringify(newValue));
      this.$nextTick(() => {
        this.$forceUpdate();
      });
    },
  },
});

app.mount("#app");
