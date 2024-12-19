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
    console.log('Getting API token from config:', API_CONFIG.API_TOKEN);
    const token = API_CONFIG.API_TOKEN;
    if (!token) {
      console.error('Token không tồn tại trong config');
      throw new Error('Token không tồn tại');
    }
    const cleanToken = token.replace(/['"]+/g, '').trim();
    console.log('Clean token:', cleanToken);
    return cleanToken;
  }

  private static getHeaders(): HeadersInit {
    const headers = {
      'Content-Type': 'application/json',
      'Authorization': this.getCleanToken()
    };
    console.log('Request headers:', headers);
    return headers;
  }

  private static getWebSocketUrl(): string {
    console.log('WS_URL from config:', API_CONFIG.WS_URL);
    const wsUrl = new URL(API_CONFIG.WS_URL);
    wsUrl.searchParams.append('token', this.getCleanToken());
    console.log('Final WebSocket URL:', wsUrl.toString());
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
      const errorMessage = errors[response.status] || 'Lỗi không xác định';
      console.error(`API Error ${response.status}:`, errorMessage);
      throw new Error(errorMessage);
    }
  }

  static async searchWord(word: string): Promise<void> {
    try {
      console.log('Searching word:', word);
      const url = `${API_CONFIG.BASE_URL}/word`;
      console.log('Search API URL:', url);

      const requestBody = JSON.stringify({ data: word });
      console.log('Request body:', requestBody);

      const response = await fetch(url, {
        method: 'POST',
        headers: this.getHeaders(),
        body: requestBody,
      });

      console.log('Response status:', response.status);
      this.handleResponseError(response);
      
      const data = await response.json();
      console.log('Response data:', data);
      return data;
    } catch (error) {
      console.error('Search word error:', error);
      throw error;
    }
  }

  private static startReconnectInterval(
    onMessage: (data: AnswerData, type: string) => void,
    onConnectionChange: (status: boolean) => void
  ): void {
    if (this.reconnectInterval) {
      console.log('Reconnect interval already running');
      return;
    }

    console.log('Starting reconnect interval');
    this.reconnectInterval = setInterval(() => {
      if (this.retryCount >= this.MAX_RETRIES) {
        console.log('Max retries reached, cleaning up');
        this.cleanup();
        onConnectionChange(false);
        return;
      }

      if (!this.ws || this.ws.readyState === WebSocket.CLOSED) {
        console.log(`Retry attempt ${this.retryCount + 1}/${this.MAX_RETRIES}`);
        this.retryCount++;
        this.connectWebSocket(onMessage, onConnectionChange);
      }
    }, this.RECONNECT_DELAY);
  }

  private static cleanup(): void {
    console.log('Cleaning up WebSocket connection');
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
    console.log('Initializing WebSocket connection');
    console.log('Current connection state:', {
      isConnecting: this.isConnecting,
      wsState: this.ws?.readyState
    });

    if (this.ws?.readyState === WebSocket.OPEN || this.isConnecting) {
      console.log('WebSocket already connected or connecting');
      return () => this.cleanup();
    }

    try {
      this.isConnecting = true;
      const wsUrl = this.getWebSocketUrl();
      console.log('Connecting to WebSocket:', wsUrl);
      this.ws = new WebSocket(wsUrl);

      this.ws.onopen = () => {
        console.log('WebSocket connected successfully');
        this.isConnecting = false;
        this.retryCount = 0;
        onConnectionChange(true);
      };

      this.ws.onmessage = (event) => {
        try {
          console.log('WebSocket message received:', event.data);
          const message = JSON.parse(event.data);
          if (message) {
            const processedData: AnswerData = {
              detail: message.detail || {},
              structure: message.structure || "Sentence"
            };
            console.log('Processed WebSocket data:', processedData);
            onMessage(processedData, message.type || 'dictionary');
          }
        } catch (error) {
          console.error('WebSocket message processing error:', error);
        }
      };

      this.ws.onclose = () => {
        console.log('WebSocket connection closed');
        this.isConnecting = false;
        onConnectionChange(false);
        this.startReconnectInterval(onMessage, onConnectionChange);
      };

      this.ws.onerror = (error) => {
        console.error('WebSocket error:', error);
        this.isConnecting = false;
        onConnectionChange(false);
        this.ws?.close();
      };

    } catch (error) {
      console.error('WebSocket initialization error:', error);
      this.isConnecting = false;
      onConnectionChange(false);
    }

    return () => this.cleanup();
  }
}
