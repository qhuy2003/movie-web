const defaultApiUrl = import.meta.env.DEV
  ? 'http://localhost:8080'
  : 'https://movie-web-em6f.onrender.com';

export const API_BASE_URL = (import.meta.env.VITE_API_BASE_URL || defaultApiUrl).replace(/\/+$/, '');
