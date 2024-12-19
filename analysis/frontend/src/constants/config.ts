interface Config {
  BASE_URL: string;
  WS_URL: string;
  API_TOKEN: string;
}

export const API_CONFIG = {
  API_TOKEN: process.env.REACT_APP_API_TOKEN || "",
  BASE_URL: process.env.REACT_APP_BASE_URL || "",
  WS_URL: process.env.REACT_APP_WS_URL || "",
}
