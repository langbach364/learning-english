interface Config {
  BASE_URL: string;
  WS_URL: string;
  API_TOKEN: string;
}

export const API_CONFIG: Config = {
  BASE_URL: process.env.base_url || '',
  WS_URL: process.env.ws_url || '',
  API_TOKEN: process.env.API_TOKEN || ''
};
