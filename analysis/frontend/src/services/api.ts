import { API_CONFIG } from '../constants/config';
import { AnswerData } from '../types/dictionary';
import { LoginRequest, LoginResponse } from '../constants/auth';

export class APIService {
  private static token: string | null = localStorage.getItem('token');
  private static ws: WebSocket | null = null;
  private static isConnecting = false;
  private static reconnectInterval: NodeJS.Timeout | null = null;
  private static reconnectAttempts = 0;
  private static maxReconnectAttempts = 5;
  private static reconnectDelay = 5000;

  static getWebSocketState(): string {
    if (!this.ws) return 'CLOSED';
    const states: { [key: number]: string } = {
      [WebSocket.CONNECTING]: 'CONNECTING',
      [WebSocket.OPEN]: 'OPEN',
      [WebSocket.CLOSING]: 'CLOSING',
      [WebSocket.CLOSED]: 'CLOSED'
    };
    return states[this.ws.readyState] || 'UNKNOWN';
  }

  static async login(credentials: LoginRequest): Promise<string> {
    try {
      const response = await fetch(`${API_CONFIG.BASE_URL}/login`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json',
        },
        body: JSON.stringify(credentials)
      });

      if (!response.ok) {
        const error = await response.json();
        throw new Error(error.message || 'Đăng nhập thất bại');
      }

      const data: LoginResponse = await response.json();
      this.token = data.token;
      localStorage.setItem('token', data.token);
      return data.token;
    } catch (error) {
      if (error instanceof Error) {
        throw new Error(`Lỗi đăng nhập: ${error.message}`);
      }
      throw new Error('Đã xảy ra lỗi không xác định trong quá trình đăng nhập');
    }
  }

  static async searchWord(word: string): Promise<void> {
    if (!this.token) throw new Error('Vui lòng đăng nhập để tiếp tục');

    try {
      const response = await fetch(`${API_CONFIG.BASE_URL}/word`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${this.token}`,
          'Content-Type': 'application/json',
          'Accept': 'application/json',
        },
        body: JSON.stringify({ data: word }),
      });

      if (!response.ok) {
        const error = await response.json();
        throw new Error(error.message || 'Lỗi khi tìm từ');
      }
    } catch (error) {
      if (error instanceof Error) {
        throw new Error(`Lỗi tìm kiếm từ: ${error.message}`);
      }
      throw new Error('Đã xảy ra lỗi không xác định trong quá trình tìm kiếm');
    }
  }

  static connectWebSocket(
    onMessage: (data: AnswerData, type: string) => void,
    onConnectionChange: (status: boolean) => void
  ): () => void {
    if (!this.token) {
      throw new Error('Vui lòng đăng nhập để kết nối WebSocket');
    }

    if (this.ws?.readyState === WebSocket.OPEN) {
      onConnectionChange(true);
      return () => this.cleanup();
    }

    if (this.isConnecting) {
      return () => this.cleanup();
    }

    this.initializeWebSocket(onMessage, onConnectionChange);
    return () => this.cleanup();
  }

  private static initializeWebSocket(
    onMessage: (data: AnswerData, type: string) => void,
    onConnectionChange: (status: boolean) => void
  ): void {
    this.isConnecting = true;
    this.reconnectAttempts = 0;

    try {
      if (this.token === null) {
        throw new Error('Token không tồn tại');
      }
      if (!API_CONFIG.WS_URL) {
        throw new Error('Không tìm thấy địa chỉ WebSocket');
      }
      this.ws = new WebSocket(API_CONFIG.WS_URL, [this.token]);

      this.ws.onopen = () => {
        this.isConnecting = false;
        this.reconnectAttempts = 0;
        onConnectionChange(true);
        this.stopReconnectInterval();
      };

      this.ws.onmessage = (event) => {
        try {
          const message = JSON.parse(event.data);
          
          if (!message) {
            throw new Error('Dữ liệu nhận được không hợp lệ');
          }

          const processedData: AnswerData = {
            detail: message.detail || {},
            structure: message.structure || "Sentence"
          };
          onMessage(processedData, message.type || 'dictionary');
        } catch (error) {
          if (error instanceof Error) {
            throw new Error(`Lỗi xử lý tin nhắn WebSocket: ${error.message}`);
          }
          throw new Error('Lỗi không xác định khi xử lý tin nhắn WebSocket');
        }
      };

      this.ws.onclose = (event) => {
        this.isConnecting = false;
        onConnectionChange(false);
        
        if (this.reconnectAttempts < this.maxReconnectAttempts) {
          this.startReconnectInterval(onMessage, onConnectionChange);
        } else {
          throw new Error(`WebSocket đóng với mã: ${event.code}, lý do: ${event.reason}`);
        }
      };

      this.ws.onerror = (error) => {
        this.isConnecting = false;
        onConnectionChange(false);
        throw new Error(`Lỗi kết nối WebSocket: ${error}`);
      };

    } catch (error) {
      this.isConnecting = false;
      onConnectionChange(false);
      if (error instanceof Error) {
        throw new Error(`Lỗi khởi tạo WebSocket: ${error.message}`);
      }
      throw new Error('Lỗi không xác định khi khởi tạo WebSocket');
    }
  }

  private static startReconnectInterval(
    onMessage: (data: AnswerData, type: string) => void,
    onConnectionChange: (status: boolean) => void
  ): void {
    if (this.reconnectInterval) return;

    this.reconnectInterval = setInterval(() => {
      if (!this.ws || this.ws.readyState === WebSocket.CLOSED) {
        this.reconnectAttempts++;
        
        if (this.reconnectAttempts < this.maxReconnectAttempts) {
          this.initializeWebSocket(onMessage, onConnectionChange);
        } else {
          this.stopReconnectInterval();
          throw new Error(`Không thể kết nối lại sau ${this.maxReconnectAttempts} lần thử`);
        }
      }
    }, this.reconnectDelay);
  }

  private static stopReconnectInterval(): void {
    if (this.reconnectInterval) {
      clearInterval(this.reconnectInterval);
      this.reconnectInterval = null;
    }
  }

  private static cleanup(): void {
    this.stopReconnectInterval();
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
  }

  static logout(): void {
    localStorage.removeItem('token');
    this.token = null;
    this.cleanup();
  }

  static isWebSocketConnected(): boolean {
    return this.ws?.readyState === WebSocket.OPEN;
  }
}
