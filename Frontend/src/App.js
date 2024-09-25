const app = Vue.createApp({
  data() {
    return {
      sentenceInfo: {},
      wordsInfo: {},
      newWord: "",
      highlightedOriginalWords: [],
      highlightedEditedWords: [],
      tooltipPosition: { x: 0, y: 0 },
      isScanning: false,
      selectionStart: null,
      selectionEnd: null,
    };
  },
  methods: {
    sendWordToServer() {
      fetch("http://localhost:7089/word", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ data: this.newWord }),
      })
        .then((response) => response.json())
        .then((data) => {
          console.log("Success:", data);
          this.newWord = "";
          this.fetchDefinitions();
        })
        .catch((error) => console.error("Error:", error));
    },
    fetchDefinitions() {
      fetch("../Handler/sourcegraph-cody/answer.txt")
        .then((response) => response.text())
        .then((data) => {
          const parsedData = window.parseDefinitions(data);
          this.sentenceInfo = parsedData.sentence;
          this.wordsInfo = parsedData.words;
        })
        .catch((error) => console.error("Error:", error));
    },
    toggleScanning() {
      this.isScanning = !this.isScanning;
      if (!this.isScanning) {
        this.clearHighlight();
      }
    },
    startSelection(event, type) {
      if (!this.isScanning) return;
      event.preventDefault();
      this.selectionStart = event.target;
      if (type === 'original') {
        this.highlightedOriginalWords = [event.target.textContent.trim()];
        this.highlightedEditedWords = [];
      } else {
        this.highlightedEditedWords = [event.target.textContent.trim()];
        this.highlightedOriginalWords = [];
      }
      this.updateTooltipPosition(event.target);
    },
    updateSelection(event, type) {
      if (!this.isScanning || !this.selectionStart) return;
      event.preventDefault();
      this.selectionEnd = event.target;
      this.highlightSelection(type);
    },
    endSelection(event, type) {
      if (!this.isScanning) return;
      event.preventDefault();
      this.selectionStart = null;
      this.selectionEnd = null;
    },
    highlightSelection(type) {
      if (!this.isScanning || !this.selectionStart || !this.selectionEnd) return;

      const words = Array.from(
        this.$el.querySelectorAll(type === 'original' ? ".original-sentence span" : ".edited-sentence span")
      );
      const startIndex = words.indexOf(this.selectionStart);
      const endIndex = words.indexOf(this.selectionEnd);

      const start = Math.min(startIndex, endIndex);
      const end = Math.max(startIndex, endIndex);

      const highlightedWords = words
        .slice(start, end + 1)
        .map((span) => span.textContent.trim());

      if (type === 'original') {
        this.highlightedOriginalWords = highlightedWords;
        this.highlightedEditedWords = [];
      } else {
        this.highlightedEditedWords = highlightedWords;
        this.highlightedOriginalWords = [];
      }

      this.updateTooltipPosition(this.selectionEnd);
    },
    clearHighlight() {
      this.highlightedOriginalWords = [];
      this.highlightedEditedWords = [];
      this.selectionStart = null;
      this.selectionEnd = null;
    },
    handleClick(event, type) {
      if (!this.isScanning) return;
      if (event.target.classList.contains("word-span")) {
        const clickedWord = event.target.textContent.trim();
        if (type === 'original') {
          this.highlightedOriginalWords = [clickedWord];
          this.highlightedEditedWords = [];
        } else {
          this.highlightedEditedWords = [clickedWord];
          this.highlightedOriginalWords = [];
        }
        this.updateTooltipPosition(event.target);
      } else if (event.target.classList.contains("original-sentence") || event.target.classList.contains("edited-sentence")) {
        if (type === 'original') {
          this.highlightedOriginalWords = this.sentenceInfo["Câu gốc"].split(" ");
          this.highlightedEditedWords = [];
        } else {
          this.highlightedEditedWords = this.sentenceInfo["Ghi lại câu đã sửa"].split(" ");
          this.highlightedOriginalWords = [];
        }
        this.updateTooltipPosition(event.target);
      } else {
        this.highlightedOriginalWords = [];
        this.highlightedEditedWords = [];
      }
    },
    updateTooltipPosition(element) {
      const rect = element.getBoundingClientRect();
      this.tooltipPosition = {
        x: rect.right,
        y: rect.bottom,
      };
    },
  },
  mounted() {
    this.fetchDefinitions();
  },
});

app.mount("#app");
