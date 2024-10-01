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
      suppressClickHandler: false,
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
      fetch("../Handler/sourcegraph-cody/answer.txt")
        .then((response) => response.text())
        .then((data) => {
          const parsedData = window.parseDefinitions(data);
          this.sentenceInfo = parsedData.sentence;
          this.wordsInfo = parsedData.words;
          console.log("Definitions fetched and parsed");
        })
        .catch((error) => console.error("Error fetching definitions:", error));
    },
    toggleScanning() {
      this.isScanning = !this.isScanning;
      if (!this.isScanning) {
        this.clearHighlight();
      }
      console.log("Scanning mode:", this.isScanning ? "ON" : "OFF");
    },
    startSelection(event, type) {
      if (!this.isScanning) return;
      event.preventDefault();
      event.stopPropagation();
      this.selectionStart = event.target;
      this.selectionEnd = event.target;

      const word = event.target.textContent.trim();
      if (type === 'original') {
        this.highlightedOriginalWords = [word];
        this.highlightedEditedWords = [];
      } else {
        this.highlightedEditedWords = [word];
        this.highlightedOriginalWords = [];
      }
      this.updateTooltipPosition(event.target);
    },
    updateSelection(event, type) {
      if (!this.isScanning || !this.selectionStart) return;
      event.preventDefault();
      event.stopPropagation();
      this.selectionEnd = event.target;
      this.highlightSelection(type);
    },
    endSelection(event, type) {
      if (!this.isScanning) return;
      event.preventDefault();
      event.stopPropagation();

      const selectedWords = type === 'original' ? this.highlightedOriginalWords : this.highlightedEditedWords;
      if (selectedWords.length > 0) {
        this.sendScannedWordsToServer(selectedWords);
        this.suppressClickHandler = true;
      }

      this.selectionStart = null;
      this.selectionEnd = null;
    },
    sendScannedWordsToServer(words) {
      if (this.isSending) return;
      this.isSending = true;
      fetch("http://localhost:7089/listen_word", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ data: words.join(' ') }),
      })
        .then(response => {
          if (!response.ok) {
            throw new Error('Network response was not ok');
          }
          return response.text();
        })
        .then(data => {
          console.log("Server response:", data);
        })
        .catch((error) => console.error("Error sending scanned words:", error))
        .finally(() => {
          this.isSending = false;
        });
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
      if (this.suppressClickHandler) {
        this.suppressClickHandler = false;
        return;
      }

      if (!this.isScanning) return;
      event.preventDefault();
      event.stopPropagation();

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

        this.sendScannedWordsToServer([clickedWord]);
      } else if (event.target.classList.contains("original-sentence") || event.target.classList.contains("edited-sentence")) {
        if (type === 'original') {
          this.highlightedOriginalWords = this.sentenceInfo["Câu gốc"].split(" ");
          this.highlightedEditedWords = [];
        } else {
          this.highlightedEditedWords = this.sentenceInfo["Ghi lại câu đã sửa"].split(" ");
          this.highlightedOriginalWords = [];
        }
        this.updateTooltipPosition(event.target);

        const fullSentenceWords = this.highlightedOriginalWords.length > 0 ? this.highlightedOriginalWords : this.highlightedEditedWords;
        this.sendScannedWordsToServer(fullSentenceWords);
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
  created() {
    this.fetchDefinitions();
    this.isScanning = JSON.parse(localStorage.getItem('isScanning')) || false;
  },
  watch: {
    isScanning(newValue) {
      localStorage.setItem('isScanning', JSON.stringify(newValue));
    }
  }
});

app.mount("#app");
