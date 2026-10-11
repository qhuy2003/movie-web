import { reactive } from 'vue';
import { API_BASE_URL } from './config.js';

const savedToken = localStorage.getItem('cinestream_token') || '';
const savedUser = JSON.parse(localStorage.getItem('cinestream_user') || 'null');

export const authState = reactive({
  token: savedToken,
  user: savedUser,
  isAuthenticated: !!savedToken && !!savedUser,

  setAuth(user, token) {
    this.user = user;
    this.token = token;
    this.isAuthenticated = true;
    localStorage.setItem('cinestream_token', token);
    localStorage.setItem('cinestream_user', JSON.stringify(user));
  },

  logout() {
    this.user = null;
    this.token = '';
    this.isAuthenticated = false;
    localStorage.removeItem('cinestream_token');
    localStorage.removeItem('cinestream_user');
  }
});

async function parseResponseBody(res) {
  const text = await res.text();
  try {
    return JSON.parse(text);
  } catch (e) {
    if (!res.ok) {
      throw new Error(`Lỗi máy chủ (${res.status}): ${text || 'Không thể kết nối đến Go Backend Server'}`);
    }
    throw new Error('Dữ liệu trả về từ server không phải định dạng JSON hợp lệ.');
  }
}

export async function registerUser(email, username, password) {
  try {
    const res = await fetch(`${API_BASE_URL}/auth/register`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email, username, password })
    });

    const data = await parseResponseBody(res);
    if (!res.ok) {
      throw new Error(data.error || 'Đăng ký thất bại');
    }

    authState.setAuth(data.user, data.token);
    return data;
  } catch (err) {
    if (err.name === 'TypeError' && err.message.includes('fetch')) {
      throw new Error('Không thể kết nối tới Go Backend Server. Hãy kiểm tra URL backend và trạng thái dịch vụ.');
    }
    throw err;
  }
}

export async function loginUser(email, password) {
  try {
    const res = await fetch(`${API_BASE_URL}/auth/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email, password })
    });

    const data = await parseResponseBody(res);
    if (!res.ok) {
      throw new Error(data.error || 'Đăng nhập thất bại');
    }

    authState.setAuth(data.user, data.token);
    return data;
  } catch (err) {
    if (err.name === 'TypeError' && err.message.includes('fetch')) {
      throw new Error('Không thể kết nối tới Go Backend Server. Hãy kiểm tra URL backend và trạng thái dịch vụ.');
    }
    throw err;
  }
}

export async function fetchMe() {
  if (!authState.token) return null;
  try {
    const res = await fetch(`${API_BASE_URL}/auth/me`, {
      headers: { 'Authorization': `Bearer ${authState.token}` }
    });
    if (res.ok) {
      const data = await parseResponseBody(res);
      authState.user = data;
      localStorage.setItem('cinestream_user', JSON.stringify(data));
      return data;
    } else {
      authState.logout();
    }
  } catch (e) {
    console.warn('Không thể xác thực token từ server:', e);
  }
  return null;
}
