import { API_CONFIG } from '../constants/config';
import { AnswerData } from '../types/dictionary';

export class APIService {
  private static ws: WebSocket | null = null;
  private static isConnecting = false;
  private static reconnectInterval: NodeJS.Timeout | null = null;
  private static readonly RECONNECT_DELAY = 1000;
  private static readonly MAX_RETRIES = 3;
  private static retryCount = 0;


  private static getCleanToken(): string {
    const token = API_CONFIG.API_TOKEN;
    if (!token) {
      throw new Error('Token không tồn tại');
    }
    return token.replace(/['"]+/g, '').trim();
  }
  


  private static getHeaders(): HeadersInit {
    return {
      'Content-Type': 'application/json',
      'Authorization': this.getCleanToken()
    };
  }

  private static getWebSocketUrl(): string {
    const wsUrl = new URL(API_CONFIG.WS_URL);
    wsUrl.searchParams.append('token', this.getCleanToken());
    return wsUrl.toString();
  }

  private static handleResponseError(response: Response): void {
    const errors: Record<number, string> = {
      401: 'Token không hợp lệ hoặc đã hết hạn',
      403: 'Không có quyền truy cập',
      404: 'Không tìm thấy tài nguyên',
      500: 'Lỗi server',
      502: 'Bad Gateway',
      503: 'Dịch vụ không khả dụng'
    };

    if (!response.ok) {
      throw new Error(errors[response.status] || 'Lỗi không xác định');
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
      return await response.json();
    } catch (error) {
      console.error('Lỗi khi tìm kiếm từ:', error);
      throw error;
    }
  }

  private static startReconnectInterval(
    onMessage: (data: AnswerData, type: string) => void,
    onConnectionChange: (status: boolean) => void
  ): void {
    if (this.reconnectInterval) return;

    this.reconnectInterval = setInterval(() => {
      if (this.retryCount >= this.MAX_RETRIES) {
        this.cleanup();
        onConnectionChange(false);
        return;
      }

      if (!this.ws || this.ws.readyState === WebSocket.CLOSED) {
        this.retryCount++;
        this.connectWebSocket(onMessage, onConnectionChange);
      }
    }, this.RECONNECT_DELAY);
  }

  private static cleanup(): void {
    if (this.reconnectInterval) {
      clearInterval(this.reconnectInterval);
      this.reconnectInterval = null;
    }
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
    this.retryCount = 0;
    this.isConnecting = false;
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
        this.retryCount = 0;
        onConnectionChange(true);
        console.log('WebSocket đã kết nối thành công');
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
          }
        } catch (error) {
          console.error('Lỗi xử lý dữ liệu WebSocket:', error);
        }
      };

      this.ws.onclose = () => {
        this.isConnecting = false;
        onConnectionChange(false);
        this.startReconnectInterval(onMessage, onConnectionChange);
        console.log('WebSocket đã đóng kết nối');
      };

      this.ws.onerror = (error) => {
        this.isConnecting = false;
        onConnectionChange(false);
        console.error('Lỗi WebSocket:', error);
        this.ws?.close();
      };

    } catch (error) {
      console.error('Lỗi khởi tạo WebSocket:', error);
      this.isConnecting = false;
      onConnectionChange(false);
    }

    return () => this.cleanup();
  }
}
