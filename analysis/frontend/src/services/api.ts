import { API_CONFIG } from '../constants/config';

type WebSocketMessageCallback = (data: any, type: string) => void;
type ConnectionStatusCallback = (status: boolean) => void;

class APIServiceClass {
  private ws: WebSocket | null = null;
  private messageCallback: WebSocketMessageCallback | null = null;
  private statusCallback: ConnectionStatusCallback | null = null;

  constructor() {
    this.ws = null;
    this.messageCallback = null;
    this.statusCallback = null;
  }

  public connectWebSocket = (
    onMessage: WebSocketMessageCallback,
    onConnectionStatus: ConnectionStatusCallback
  ) => {
    this.messageCallback = onMessage;
    this.statusCallback = onConnectionStatus;

    const token = localStorage.getItem('token');
    if (!token) return;

    const wsUrl = `${API_CONFIG.WS_URL}?token=${token}`;
    this.ws = new WebSocket(wsUrl);

    this.ws.onopen = () => {
      this.statusCallback?.(true);
    };

    this.ws.onclose = () => {
      this.statusCallback?.(false);
      setTimeout(() => this.reconnect(), 3000);
    };

    this.ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        this.messageCallback?.(data.data, data.type);
      } catch (error) {
        console.error('WebSocket message parse error:', error);
      }
    };

    return () => {
      if (this.ws && this.ws.readyState === WebSocket.OPEN) {
        this.ws.close();
      }
    };
  };

  private reconnect = () => {
    if (this.messageCallback && this.statusCallback) {
      this.connectWebSocket(this.messageCallback, this.statusCallback);
    }
  };

  public searchWord = async (word: string) => {
    if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
      throw new Error('WebSocket connection not established');
    }

    this.ws.send(JSON.stringify({
      type: 'search',
      data: { word }
    }));
  };

  public login = async (credentials: { username: string; password: string }) => {
    try {
      const response = await fetch(`${API_CONFIG.BASE_URL}/auth/login`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(credentials)
      });

      if (!response.ok) {
        throw new Error('Login failed');
      }

      const data = await response.json();
      localStorage.setItem('token', data.token);
      return data;
    } catch (error) {
      console.error('Login error:', error);
      throw error;
    }
  };

  public logout = () => {
    localStorage.removeItem('token');
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.ws.close();
    }
  };

  public verifyToken = async (): Promise<{ success: boolean }> => {
    try {
      const token = localStorage.getItem('token');
      if (!token) {
        return { success: false };
      }

      const response = await fetch(`${API_CONFIG.BASE_URL}/auth/verify`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        }
      });

      if (!response.ok) {
        throw new Error('Token verification failed');
      }

      return { success: true };
    } catch (error) {
      console.error('Token verification error:', error);
      return { success: false };
    }
  };
}

export const APIService = new APIServiceClass();
