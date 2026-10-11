import { authState } from './auth.js';
import { API_BASE_URL } from './config.js';

async function request(path, options = {}) {
  const headers = { 'Content-Type': 'application/json', ...options.headers };
  if (authState.token) headers.Authorization = `Bearer ${authState.token}`;
  const response = await fetch(`${API_BASE_URL}${path}`, { ...options, headers });
  const text = await response.text();
  let data;
  try {
    data = text ? JSON.parse(text) : {};
  } catch {
    throw new Error(`Backend trả về dữ liệu không hợp lệ (HTTP ${response.status}).`);
  }
  if (!response.ok) throw new Error(data.error || 'Không thể thực hiện thao tác đặt vé.');
  return data;
}

export function fetchCinemas() {
  return request('/cinemas');
}

export function fetchShowtimes(movie, date, cinemaId = '') {
  const params = new URLSearchParams({
    movie_id: movie.id,
    movie_title: movie.title,
    poster_url: movie.poster_url || '',
    date
  });
  if (cinemaId) params.set('cinema_id', cinemaId);
  return request(`/showtimes?${params}`);
}

export function fetchShowtimeSeats(showtimeId) {
  return request(`/showtimes/${showtimeId}/seats`);
}

export function createBooking(showtimeId, seatCodes) {
  return request('/api/bookings', {
    method: 'POST',
    body: JSON.stringify({ showtime_id: showtimeId, seat_codes: seatCodes })
  });
}

export function startZaloPayPayment(code) {
  return request(`/api/bookings/${encodeURIComponent(code)}/pay`, { method: 'POST' });
}

export function verifyBookingPayment(code) {
  return request(`/api/bookings/${encodeURIComponent(code)}/verify-payment`, { method: 'POST' });
}

export function cancelBooking(code) {
  return request(`/api/bookings/${encodeURIComponent(code)}/cancel`, { method: 'POST' });
}

export function fetchMyBookings() {
  return request('/api/bookings');
}
