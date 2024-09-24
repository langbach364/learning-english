const app = Vue.createApp({
  data() {
    return {
      sentenceInfo: {},
      wordsInfo: {},
      newWord: "",
      highlightedWords: [],
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
    startSelection(event) {
      if (!this.isScanning) return;
      event.preventDefault();
      this.selectionStart = event.target;
      this.highlightedWords = [event.target.textContent.trim()];
      this.updateTooltipPosition(event.target);
    },
    updateSelection(event) {
      if (!this.isScanning || !this.selectionStart) return;
      event.preventDefault();
      this.selectionEnd = event.target;
      this.highlightSelection();
    },
    endSelection(event) {
      if (!this.isScanning) return;
      event.preventDefault();
      this.selectionStart = null;
      this.selectionEnd = null;
    },
    highlightSelection() {
      if (!this.isScanning || !this.selectionStart || !this.selectionEnd)
        return;

      const words = Array.from(
        this.$el.querySelectorAll(".original-sentence span")
      );
      const startIndex = words.indexOf(this.selectionStart);
      const endIndex = words.indexOf(this.selectionEnd);

      const start = Math.min(startIndex, endIndex);
      const end = Math.max(startIndex, endIndex);

      this.highlightedWords = words
        .slice(start, end + 1)
        .map((span) => span.textContent.trim());

      this.updateTooltipPosition(this.selectionEnd);
    },
    clearHighlight() {
      this.highlightedWords = [];
      this.selectionStart = null;
      this.selectionEnd = null;
    },
    handleClick(event) {
      if (!this.isScanning) return;
      if (event.target.classList.contains("word-span")) {
        const clickedWord = event.target.textContent.trim();
        this.highlightedWords = [clickedWord];
        this.updateTooltipPosition(event.target);
      } else if (event.target.classList.contains("original-sentence")) {
        this.highlightedWords = this.sentenceInfo["Câu gốc"].split(" ");
        this.updateTooltipPosition(event.target);
      } else {
        this.highlightedWords = [];
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
