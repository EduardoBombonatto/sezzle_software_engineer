import axios from 'axios';

const DEFAULT_API_BASE_URL = '/api/v1';
const REQUEST_TIMEOUT_MS = 5000;

export const httpClient = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || DEFAULT_API_BASE_URL,
  timeout: REQUEST_TIMEOUT_MS,
  headers: {
    'Content-Type': 'application/json',
    Accept: 'application/json',
  },
});
