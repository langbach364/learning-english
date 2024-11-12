import { API_CONFIG } from '../constants/config';
import { AnswerData } from '../types/dictionary';

export class APIService {
  private static ws: WebSocket | null = null;
  private static isConnecting = false;
  private static reconnectInterval: NodeJS.Timeout | null = null;

  private static startReconnectInterval(
    onMessage: (data: AnswerData, type: string) => void,
    onConnectionChange: (status: boolean) => void
  ) {
    if (this.reconnectInterval) return;

    this.reconnectInterval = setInterval(() => {
      if (!this.ws || this.ws.readyState === WebSocket.CLOSED) {
        this.connectWebSocket(onMessage, onConnectionChange);
      }
    }, 1000);
  }

  private static stopReconnectInterval() {
    if (this.reconnectInterval) {
      clearInterval(this.reconnectInterval);
      this.reconnectInterval = null;
    }
  }

  static async searchWord(word: string): Promise<void> {
    const response = await fetch(`${API_CONFIG.BASE_URL}/word`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ data: word }),
    });

    if (!response.ok) throw new Error('Network response was not ok');
  }

  static connectWebSocket(
    onMessage: (data: AnswerData, type: string) => void,
    onConnectionChange: (status: boolean) => void
  ): () => void {
    if (this.ws?.readyState === WebSocket.OPEN) {
      return () => this.cleanup();
    }

    if (this.isConnecting) {
      return () => this.cleanup();
    }

    this.isConnecting = true;
    this.ws = new WebSocket(API_CONFIG.WS_URL);

    this.ws.onopen = () => {
      this.isConnecting = false;
      onConnectionChange(true);
      this.stopReconnectInterval();
    };

    this.ws.onmessage = (event) => {
      try {
        const message = JSON.parse(event.data);
        if (message) {
          const processedData: AnswerData = {
            detail: message.detail || {},
            structure: message.structure || "Sentence"
          };
          onMessage(processedData, message.type || 'dictionary');
        } else {
          console.warn('Dữ liệu không hợp lệ:', message);
        }
      } catch (error) {
        console.error('Lỗi xử lý dữ liệu:', error);
      }
    };
    

    this.ws.onclose = () => {
      this.isConnecting = false;
      onConnectionChange(false);
      this.startReconnectInterval(onMessage, onConnectionChange);
    };

    this.ws.onerror = () => {
      this.isConnecting = false;
      onConnectionChange(false);
      this.ws?.close();
    };

    this.startReconnectInterval(onMessage, onConnectionChange);

    return () => this.cleanup();
  }

  private static cleanup() {
    this.stopReconnectInterval();
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
  }
}
