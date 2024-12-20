export const API_CONFIG = {
  BASE_URL: process.env.REACT_APP_BASE_URL,
  WS_URL: process.env.REACT_APP_WS_URL,
  AUTH: {
    LOGIN_URL: '/login',
    CREDENTIALS: {
      USERNAME: process.env.REACT_APP_USERNAME,
      PASSWORD: process.env.REACT_APP_PASSWORD,
    }
  }
};
