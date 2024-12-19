import { API_CONFIG } from '../constants/config';
import { AnswerData } from '../types/dictionary';

export class APIService {
  private static ws: WebSocket | null = null;
  private static isConnecting = false;
  private static reconnectInterval: NodeJS.Timeout | null = null;
  private static readonly RECONNECT_DELAY = 1000;

  private static getHeaders(): HeadersInit {
    return {
      'Content-Type': 'application/json',
      'Authorization': API_CONFIG.API_TOKEN || ''
    };
  }

  private static handleResponseError(response: Response): void {
    if (response.status === 401) {
      throw new Error('Token không hợp lệ hoặc đã hết hạn');
    }
    if (response.status === 403) {
      throw new Error('Không có quyền truy cập');
    }
    if (!response.ok) {
      throw new Error('Lỗi kết nối mạng');
    }
  }

  private static getWebSocketUrl(): string {
    if (!API_CONFIG.WS_URL) {
      throw new Error('WebSocket URL chưa được cấu hình');
    }
    return `${API_CONFIG.WS_URL}?token=${API_CONFIG.API_TOKEN}`;
  }

  private static startReconnectInterval(
    onMessage: (data: AnswerData, type: string) => void,
    onConnectionChange: (status: boolean) => void
  ): void {
    if (this.reconnectInterval) return;

    this.reconnectInterval = setInterval(() => {
      if (!this.ws || this.ws.readyState === WebSocket.CLOSED) {
        this.connectWebSocket(onMessage, onConnectionChange);
      }
    }, this.RECONNECT_DELAY);
  }

  private static stopReconnectInterval(): void {
    if (this.reconnectInterval) {
      clearInterval(this.reconnectInterval);
      this.reconnectInterval = null;
    }
  }

  static async searchWord(word: string): Promise<void> {
    try {
      const response = await fetch(`${API_CONFIG.BASE_URL}/word`, {
        method: 'POST',
        headers: this.getHeaders(),
        body: JSON.stringify({ data: word }),
      });

      this.handleResponseError(response);
    } catch (error) {
      console.error('Lỗi khi tìm kiếm từ:', error);
      throw error;
    }
  }

  static connectWebSocket(
    onMessage: (data: AnswerData, type: string) => void,
    onConnectionChange: (status: boolean) => void
  ): () => void {
    if (this.ws?.readyState === WebSocket.OPEN || this.isConnecting) {
      return () => this.cleanup();
    }

    try {
      this.isConnecting = true;
      this.ws = new WebSocket(this.getWebSocketUrl());

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

    } catch (error) {
      console.error('Lỗi kết nối WebSocket:', error);
      this.isConnecting = false;
      onConnectionChange(false);
    }

    return () => this.cleanup();
  }

  private static cleanup(): void {
    this.stopReconnectInterval();
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
  }
}
