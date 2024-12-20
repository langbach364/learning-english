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
      console.error('Lỗi đăng nhập:', error);
      throw error;
    }
  }

  static async searchWord(word: string): Promise<void> {
    if (!this.token) throw new Error('Vui lòng đăng nhập');

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
      console.error('Lỗi tìm kiếm:', error);
      throw error;
    }
  }

  static connectWebSocket(
    onMessage: (data: AnswerData, type: string) => void,
    onConnectionChange: (status: boolean) => void
  ): () => void {
    if (!this.token) {
      throw new Error('Vui lòng đăng nhập');
    }

    if (this.ws?.readyState === WebSocket.OPEN) {
      console.log('WebSocket đã được kết nối');
      onConnectionChange(true);
      return () => this.cleanup();
    }

    if (this.isConnecting) {
      console.log('Đang trong quá trình kết nối');
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
    console.log('Khởi tạo kết nối WebSocket');

    try {
      if (this.token === null) {
        throw new Error('Token is null');
      }
      if (!API_CONFIG.WS_URL) {
        throw new Error('WebSocket URL is undefined');
      }
      this.ws = new WebSocket(API_CONFIG.WS_URL, [this.token]);

      this.ws.onopen = () => {
        console.log('WebSocket kết nối thành công');
        this.isConnecting = false;
        this.reconnectAttempts = 0;
        onConnectionChange(true);
        this.stopReconnectInterval();
      };

      this.ws.onmessage = (event) => {
        try {
          const message = JSON.parse(event.data);
          console.log('Nhận được tin nhắn:', message);
          
          if (message) {
            const processedData: AnswerData = {
              detail: message.detail || {},
              structure: message.structure || "Sentence"
            };
            onMessage(processedData, message.type || 'dictionary');
          }
        } catch (error) {
          console.error('Lỗi xử lý tin nhắn:', error);
        }
      };

      this.ws.onclose = () => {
        console.log('WebSocket đã đóng kết nối');
        this.isConnecting = false;
        onConnectionChange(false);
        if (this.reconnectAttempts < this.maxReconnectAttempts) {
          this.startReconnectInterval(onMessage, onConnectionChange);
        }
      };

      this.ws.onerror = (error) => {
        console.error('Lỗi WebSocket:', error);
        this.isConnecting = false;
        onConnectionChange(false);
      };

    } catch (error) {
      console.error('Lỗi khởi tạo WebSocket:', error);
      this.isConnecting = false;
      onConnectionChange(false);
    }
  }  private static startReconnectInterval(
    onMessage: (data: AnswerData, type: string) => void,
    onConnectionChange: (status: boolean) => void
  ): void {
    if (this.reconnectInterval) return;

    console.log('Khởi động cơ chế kết nối lại');
    this.reconnectInterval = setInterval(() => {
      if (!this.ws || this.ws.readyState === WebSocket.CLOSED) {
        this.reconnectAttempts++;
        console.log(`Đang thử kết nối lại... (Lần ${this.reconnectAttempts}/${this.maxReconnectAttempts})`);
        
        if (this.reconnectAttempts < this.maxReconnectAttempts) {
          this.initializeWebSocket(onMessage, onConnectionChange);
        } else {
          this.stopReconnectInterval();
          console.log('Đã đạt giới hạn số lần thử kết nối lại');
        }
      }
    }, this.reconnectDelay);
  }

  private static stopReconnectInterval(): void {
    if (this.reconnectInterval) {
      console.log('Dừng cơ chế kết nối lại');
      clearInterval(this.reconnectInterval);
      this.reconnectInterval = null;
    }
  }

  private static cleanup(): void {
    console.log('Thực hiện cleanup');
    this.stopReconnectInterval();
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
  }

  static logout(): void {
    console.log('Thực hiện đăng xuất');
    localStorage.removeItem('token');
    this.token = null;
    this.cleanup();
  }

  static isWebSocketConnected(): boolean {
    return this.ws?.readyState === WebSocket.OPEN;
  }
}
