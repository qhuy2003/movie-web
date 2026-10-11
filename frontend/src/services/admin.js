import { authState } from './auth.js';
import { API_BASE_URL } from './config.js';

const ADMIN_API_URL = `${API_BASE_URL}/admin`;

function getHeaders() {
  return {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${authState.token}`
  };
}

// ==================== USERS API ====================
export async function fetchAdminUsers() {
  const res = await fetch(`${ADMIN_API_URL}/users`, { headers: getHeaders() });
  if (!res.ok) throw new Error('Không thể tải danh sách người dùng');
  return await res.json();
}

export async function updateUserRole(id, role) {
  const res = await fetch(`${ADMIN_API_URL}/users/${id}/role`, {
    method: 'PUT',
    headers: getHeaders(),
    body: JSON.stringify({ role })
  });
  const data = await res.json();
  if (!res.ok) throw new Error(data.error || 'Cập nhật role thất bại');
  return data;
}

export async function deleteUser(id) {
  const res = await fetch(`${ADMIN_API_URL}/users/${id}`, {
    method: 'DELETE',
    headers: getHeaders()
  });
  const data = await res.json();
  if (!res.ok) throw new Error(data.error || 'Xóa tài khoản thất bại');
  return data;
}

// ==================== GENRES API ====================
export async function fetchAdminGenres() {
  const res = await fetch(`${ADMIN_API_URL}/genres`, { headers: getHeaders() });
  if (!res.ok) throw new Error('Không thể tải danh sách thể loại');
  return await res.json();
}

export async function createGenre(name, slug) {
  const res = await fetch(`${ADMIN_API_URL}/genres`, {
    method: 'POST',
    headers: getHeaders(),
    body: JSON.stringify({ name, slug })
  });
  const data = await res.json();
  if (!res.ok) throw new Error(data.error || 'Tạo thể loại thất bại');
  return data;
}

export async function updateGenre(id, name, slug) {
  const res = await fetch(`${ADMIN_API_URL}/genres/${id}`, {
    method: 'PUT',
    headers: getHeaders(),
    body: JSON.stringify({ name, slug })
  });
  const data = await res.json();
  if (!res.ok) throw new Error(data.error || 'Cập nhật thể loại thất bại');
  return data;
}

export async function deleteGenre(id) {
  const res = await fetch(`${ADMIN_API_URL}/genres/${id}`, {
    method: 'DELETE',
    headers: getHeaders()
  });
  const data = await res.json();
  if (!res.ok) throw new Error(data.error || 'Xóa thể loại thất bại');
  return data;
}

export async function importGenresAdmin() {
  const res = await fetch(`${ADMIN_API_URL}/genres/import`, {
    method: 'POST',
    headers: getHeaders()
  });
  const body = await res.text();
  let data;
  try {
    data = JSON.parse(body);
  } catch {
    const contentType = res.headers.get('content-type') || 'không xác định';
    const preview = body.trim().slice(0, 120);
    throw new Error(`API đồng bộ thể loại trả về dữ liệu không phải JSON (HTTP ${res.status}, ${contentType}): ${preview}`);
  }
  if (!res.ok) throw new Error(data.error || 'Không thể đồng bộ thể loại');
  return data;
}

// ==================== MOVIES API ====================
export async function fetchAdminMovies() {
  const res = await fetch(`${ADMIN_API_URL}/movies`, { headers: getHeaders() });
  if (!res.ok) throw new Error('Không thể tải danh sách phim');
  return await res.json();
}

export async function createMovieAdmin(movieData) {
  const res = await fetch(`${ADMIN_API_URL}/movies`, {
    method: 'POST',
    headers: getHeaders(),
    body: JSON.stringify(movieData)
  });
  const data = await res.json();
  if (!res.ok) throw new Error(data.error || 'Tạo phim mới thất bại');
  return data;
}

export async function updateMovieAdmin(id, movieData) {
  const res = await fetch(`${ADMIN_API_URL}/movies/${id}`, {
    method: 'PUT',
    headers: getHeaders(),
    body: JSON.stringify(movieData)
  });
  const data = await res.json();
  if (!res.ok) throw new Error(data.error || 'Cập nhật phim thất bại');
  return data;
}

export async function deleteMovieAdmin(id) {
  const res = await fetch(`${ADMIN_API_URL}/movies/${id}`, {
    method: 'DELETE',
    headers: getHeaders()
  });
  const data = await res.json();
  if (!res.ok) throw new Error(data.error || 'Xóa phim thất bại');
  return data;
}

export async function importMoviesAdmin() {
  const res = await fetch(`${ADMIN_API_URL}/movies/import`, {
    method: 'POST',
    headers: getHeaders()
  });
  const body = await res.text();
  let data;
  try {
    data = JSON.parse(body);
  } catch {
    const contentType = res.headers.get('content-type') || 'không xác định';
    throw new Error(`API nhập phim trả về dữ liệu không phải JSON (HTTP ${res.status}, ${contentType}). Hãy kiểm tra backend đã deploy.`);
  }
  if (!res.ok) throw new Error(data.error || 'Import movies failed');
  return data;
}
