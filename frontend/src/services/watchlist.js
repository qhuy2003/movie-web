import { reactive, watch } from 'vue';
import { authState } from './auth.js';
import { API_BASE_URL } from './config.js';

const WATCHLIST_API_URL = `${API_BASE_URL}/api/watchlist`;
const LOCAL_STORAGE_KEY = 'cinestream_watchlist';

export const watchlistStore = reactive({
  items: JSON.parse(localStorage.getItem(LOCAL_STORAGE_KEY) || '[]'),
  onAuthRequired: null, // Callback function to trigger AuthModal

  async fetchWatchlistFromDB() {
    if (!authState.isAuthenticated || !authState.token) {
      return;
    }
    try {
      const res = await fetch(WATCHLIST_API_URL, {
        headers: { 'Authorization': `Bearer ${authState.token}` }
      });
      if (res.ok) {
        const dbItems = await res.json();
        if (Array.isArray(dbItems)) {
          this.items = dbItems.map(item => ({
            id: item.movie_id || item.id,
            movie_id: item.movie_id || item.id,
            title: item.title,
            poster_url: item.poster_url,
            rating: item.rating,
            quality: '4K UHD',
            genres: ['Phim Yêu Thích']
          }));
        }
      }
    } catch (e) {
      console.warn('Không thể đồng bộ watchlist từ PostgreSQL server:', e);
    }
  },

  async toggle(movie) {
    // ⚠️ REJECTION REQUIREMENT: Favorites require login!
    if (!authState.isAuthenticated) {
      if (typeof this.onAuthRequired === 'function') {
        this.onAuthRequired('Vui lòng đăng nhập để lưu phim vào tủ phim yêu thích!');
      }
      return { success: false, requireAuth: true, message: 'Vui lòng đăng nhập để lưu phim vào tủ phim yêu thích!' };
    }

    const movieId = movie.id || movie.movie_id;
    const index = this.items.findIndex(item => (item.id === movieId || item.movie_id === movieId));

    if (index >= 0) {
      // Remove item
      this.items.splice(index, 1);
      if (authState.token) {
        try {
          await fetch(`${WATCHLIST_API_URL}/${movieId}`, {
            method: 'DELETE',
            headers: { 'Authorization': `Bearer ${authState.token}` }
          });
        } catch (e) {
          console.warn('Lỗi xóa khỏi DB:', e);
        }
      }
    } else {
      // Add item
      const newItem = {
        id: movieId,
        movie_id: movieId,
        title: movie.title,
        poster_url: movie.poster_url,
        rating: movie.rating || 8.0,
        quality: movie.quality || '4K UHD',
        genres: movie.genres || ['Phim Yêu Thích']
      };
      this.items.push(newItem);

      if (authState.token) {
        try {
          await fetch(WATCHLIST_API_URL, {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
              'Authorization': `Bearer ${authState.token}`
            },
            body: JSON.stringify({
              movie_id: movieId,
              title: movie.title,
              poster_url: movie.poster_url,
              rating: movie.rating || 8.0
            })
          });
        } catch (e) {
          console.warn('Lỗi thêm vào DB:', e);
        }
      }
    }

    return { success: true };
  },

  isSaved(movieId) {
    return this.items.some(item => (item.id === movieId || item.movie_id === movieId));
  },

  clear() {
    this.items.length = 0;
  }
});

// Auto-sync with localStorage
watch(
  () => watchlistStore.items,
  (newItems) => {
    localStorage.setItem(LOCAL_STORAGE_KEY, JSON.stringify(newItems));
  },
  { deep: true }
);

// Auto-fetch watchlist from PostgreSQL when user logs in
watch(
  () => authState.isAuthenticated,
  (isAuth) => {
    if (isAuth) {
      watchlistStore.fetchWatchlistFromDB();
    } else {
      watchlistStore.items.length = 0;
    }
  },
  { immediate: true }
);
