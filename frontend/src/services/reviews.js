import { authState } from './auth.js';
import { API_BASE_URL } from './config.js';

async function readResponse(response, fallbackMessage) {
  const body = await response.text();
  let data;
  try {
    data = body ? JSON.parse(body) : {};
  } catch {
    if (!response.ok) {
      throw new Error(`${fallbackMessage} (HTTP ${response.status}). Backend có thể chưa được deploy phiên bản mới.`);
    }
    throw new Error('Backend trả về dữ liệu không đúng định dạng.');
  }
  if (!response.ok) throw new Error(data.error || `${fallbackMessage} (HTTP ${response.status})`);
  return data;
}

export async function fetchMovieReviews(movieId) {
  const response = await fetch(`${API_BASE_URL}/movies/${encodeURIComponent(movieId)}/reviews`);
  return readResponse(response, 'Không thể tải đánh giá phim');
}

export async function fetchMyReviews() {
  const response = await fetch(`${API_BASE_URL}/api/reviews`, {
    headers: { Authorization: `Bearer ${authState.token}` },
  });
  return readResponse(response, 'Không thể tải đánh giá của bạn');
}

export async function saveMovieReview(movie, rating, comment) {
  const response = await fetch(`${API_BASE_URL}/api/movies/${encodeURIComponent(movie.id)}/reviews`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${authState.token}`,
    },
    body: JSON.stringify({ movie_title: movie.title, rating, comment }),
  });
  return readResponse(response, 'Không thể gửi đánh giá');
}
