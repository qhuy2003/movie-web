<template>
  <Teleport to="body">
    <div v-if="movie" class="modal-backdrop" @click.self="$emit('close')">
      <div class="modal-container">
        <!-- Close Button -->
        <button class="modal-close-btn" @click="$emit('close')">
          <X :size="24" />
        </button>

        <!-- Header Backdrop Image -->
        <div
          class="modal-header-backdrop"
          :class="{ 'has-trailer': activeMovie.trailer_url }"
          @click="activeMovie.trailer_url && playTrailer()"
        >
          <img :src="activeMovie.thumb_url" :alt="activeMovie.title" class="header-img" />
          <div class="header-overlay"></div>

          <button
            v-if="activeMovie.trailer_url"
            class="trailer-preview-btn"
            type="button"
            aria-label="Xem trailer"
            @click.stop="playTrailer"
          >
            <span class="preview-play"><Play :size="30" fill="currentColor" /></span>
            <span class="preview-copy">
              <strong>Xem Trailer</strong>
              <small>Phát video giới thiệu</small>
            </span>
          </button>

          <!-- Quick Title Info Over Image -->
          <div class="header-content">
            <div class="badge-group">
              <span class="badge badge-rating">⭐ {{ activeMovie.rating }} TMDB</span>
              <span class="badge badge-quality">{{ activeMovie.quality || '4K UHD' }}</span>
              <span class="badge badge-genre">{{ activeMovie.year }}</span>
            </div>
            <h2 class="modal-title">{{ activeMovie.title }}</h2>
            <p class="modal-origin">{{ activeMovie.origin_name }}</p>
          </div>
        </div>

        <!-- Body Details -->
        <div class="modal-body">
          <div class="body-grid">
            <!-- Left Poster -->
            <div class="poster-col">
              <img :src="activeMovie.poster_url" :alt="activeMovie.title" class="poster-card-img" />
              <button class="btn btn-primary full-width" @click="playTrailer">
                <Play :size="18" fill="currentColor" /> Xem Trailer
              </button>
              <button class="btn btn-secondary full-width" @click="$emit('book', activeMovie)">
                Đặt Vé Tại Rạp
              </button>
              <button 
                class="btn btn-secondary full-width"
                :class="{ 'active': isSaved }"
                @click="toggleWatchlist"
              >
                <Heart :size="18" :fill="isSaved ? '#e50914' : 'none'" :color="isSaved ? '#e50914' : '#fff'" />
                {{ isSaved ? 'Đã Lưu Yêu Thích' : 'Lưu Yêu Thích' }}
              </button>
            </div>

            <!-- Right Info -->
            <div class="info-col">
              <!-- Meta Rows -->
              <div class="meta-list">
                <div class="meta-item">
                  <span class="meta-label"><Clock :size="16" /> Thời lượng:</span>
                  <span class="meta-val">{{ activeMovie.duration || '120 phút' }}</span>
                </div>
                <div class="meta-item">
                  <span class="meta-label"><User :size="16" /> Đạo diễn:</span>
                  <span class="meta-val">{{ activeMovie.director || 'N/A' }}</span>
                </div>
                <div class="meta-item">
                  <span class="meta-label"><Tag :size="16" /> Thể loại:</span>
                  <div class="genres-chips">
                    <span v-for="g in activeMovie.genres" :key="g" class="chip">{{ g }}</span>
                  </div>
                </div>
              </div>

              <!-- Cast Members -->
              <div v-if="activeMovie.cast && activeMovie.cast.length" class="cast-section">
                <h4 class="sub-heading">Diễn Viên</h4>
                <div class="cast-tags">
                  <span v-for="c in activeMovie.cast" :key="c" class="cast-pill">{{ c }}</span>
                </div>
              </div>

              <!-- Synopsis -->
              <div class="story-section">
                <h4 class="sub-heading">Nội Dung Phim</h4>
                <p class="story-text" v-html="activeMovie.description || 'Chưa có bài viết mô tả.'"></p>
              </div>

              <!-- Episodes List -->
              <div v-if="activeMovie.episodes && activeMovie.episodes.length" class="episodes-section">
                <h4 class="sub-heading">Danh Sách Tập / Bản Chiếu</h4>
                <div class="episodes-grid">
                  <button 
                    v-for="ep in activeMovie.episodes" 
                    :key="ep.id" 
                    class="ep-btn"
                    @click="$emit('play', activeMovie, ep)"
                  >
                    <Play :size="14" /> {{ ep.name }}
                  </button>
                </div>
              </div>
            </div>
          </div>

          <section class="reviews-section" aria-labelledby="reviews-title">
            <div class="reviews-heading">
              <div>
                <h3 id="reviews-title">Đánh giá & bình luận</h3>
                <p>Chia sẻ cảm nhận của bạn sau khi xem phim.</p>
              </div>
              <div class="reviews-summary">
                <strong>★ {{ averageRating.toFixed(1) }}</strong>
                <span>{{ totalReviews }} đánh giá</span>
              </div>
            </div>

            <p v-if="reviewsError" class="reviews-message error">{{ reviewsError }}</p>
            <p v-else-if="reviewsLoading" class="reviews-message">Đang tải bình luận…</p>

            <template v-else>
              <form v-if="authState.isAuthenticated" class="review-form" @submit.prevent="submitReview">
                <h4>Đánh giá phim này</h4>
                <p class="review-form-hint">Mỗi tài khoản có một đánh giá cho mỗi phim. Gửi lại sẽ cập nhật đánh giá trước đó.</p>
                <div class="rating-picker" role="radiogroup" aria-label="Chọn số sao">
                  <button
                    v-for="star in 5"
                    :key="star"
                    type="button"
                    class="rating-choice"
                    :class="{ selected: star <= reviewRating }"
                    :aria-label="`${star} sao`"
                    :aria-pressed="reviewRating === star"
                    @click="reviewRating = star"
                  >★</button>
                </div>
                <textarea
                  v-model="reviewComment"
                  maxlength="1000"
                  rows="4"
                  placeholder="Viết cảm nhận của bạn về bộ phim…"
                  required
                ></textarea>
                <div class="review-form-footer">
                  <span>{{ reviewComment.length }}/1000 ký tự</span>
                  <button class="btn btn-primary" type="submit" :disabled="reviewBusy || !reviewComment.trim()">
                    {{ reviewBusy ? 'Đang gửi…' : 'Gửi đánh giá' }}
                  </button>
                </div>
                <p v-if="reviewNotice" class="reviews-message success">{{ reviewNotice }}</p>
              </form>

              <div v-else class="review-login-prompt">
                <span>Đăng nhập để chấm điểm và bình luận về phim.</span>
                <button class="btn btn-secondary" type="button" @click="$emit('open-auth')">Đăng nhập</button>
              </div>

              <div v-if="reviews.length" class="review-list">
                <article v-for="review in reviews" :key="review.id" class="review-card">
                  <div class="review-card-heading">
                    <div class="review-author">
                      <span class="review-avatar">{{ review.username?.charAt(0)?.toUpperCase() || '?' }}</span>
                      <div><strong>{{ review.username }}</strong><time>{{ formatReviewDate(review.updated_at || review.created_at) }}</time></div>
                    </div>
                    <span class="review-stars" :aria-label="`${review.rating} trên 5 sao`">{{ '★'.repeat(review.rating) }}<span>{{ '★'.repeat(5 - review.rating) }}</span></span>
                  </div>
                  <p class="review-comment">{{ review.comment }}</p>
                </article>
              </div>
              <p v-else class="reviews-message empty-reviews">Chưa có bình luận nào. Hãy là người đầu tiên đánh giá phim này!</p>
            </template>
          </section>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup>
import { ref, computed, watch } from 'vue';
import { X, Play, Heart, Clock, User, Tag } from '@lucide/vue';
import { watchlistStore } from '../services/watchlist.js';
import { fetchMovieDetailBySlug } from '../services/api.js';
import { authState } from '../services/auth.js';
import { fetchMovieReviews, saveMovieReview } from '../services/reviews.js';

const props = defineProps({
  movie: {
    type: Object,
    default: null
  }
});

const emit = defineEmits(['close', 'play', 'book', 'open-auth']);

const detailedMovieInfo = ref(null);
const reviews = ref([]);
const averageRating = ref(0);
const totalReviews = ref(0);
const reviewsLoading = ref(true);
const reviewsError = ref('');
const reviewRating = ref(5);
const reviewComment = ref('');
const reviewBusy = ref(false);
const reviewNotice = ref('');

const apiEpisodes = computed(() => {
  const groups = detailedMovieInfo.value?.episodes;
  if (!Array.isArray(groups)) return [];
  return groups.flatMap(group => (group.server_data || group.items || []).map((episode, index) => ({
    ...episode,
    id: episode.id || episode.slug || `${group.server_name || 'server'}-${index}`,
    name: episode.name || `Tập ${index + 1}`,
    server_name: group.server_name || ''
  })));
});

const activeMovie = computed(() => {
  if (!detailedMovieInfo.value) return props.movie;
  const m = detailedMovieInfo.value.movie;
  return {
    ...props.movie,
    description: m.content || props.movie?.description,
    cast: m.actor && m.actor.length ? m.actor : props.movie?.cast,
    director: m.director && m.director.length ? m.director.join(', ') : props.movie?.director,
    genres: m.category && m.category.length ? m.category.map(c => c.name) : props.movie?.genres,
    quality: m.quality || props.movie?.quality,
    duration: m.time ? `${m.time} phút` : props.movie?.duration,
    trailer_url: m.trailer_url || props.movie?.trailer_url,
    episodes: apiEpisodes.value.length ? apiEpisodes.value : props.movie?.episodes
  };
});

const isSaved = computed(() => {
  return props.movie ? watchlistStore.isSaved(props.movie.id) : false;
});

const toggleWatchlist = () => {
  if (props.movie) {
    watchlistStore.toggle(props.movie);
  }
};

const playTrailer = () => {
  if (!activeMovie.value?.trailer_url) {
    const title = activeMovie.value?.origin_name || activeMovie.value?.title || '';
    if (!title) return;
    const query = encodeURIComponent(`${title} official trailer`);
    window.open(`https://www.youtube.com/results?search_query=${query}`, '_blank', 'noopener,noreferrer');
    return;
  }
  emit('play', activeMovie.value, {
    id: `trailer-${activeMovie.value.id}`,
    name: 'Trailer chính thức',
    link_embed: activeMovie.value.trailer_url,
    is_trailer: true
  });
};

const loadDetailData = async () => {
  detailedMovieInfo.value = null;
  reviews.value = [];
  averageRating.value = 0;
  totalReviews.value = 0;
  reviewRating.value = 5;
  reviewComment.value = '';
  reviewNotice.value = '';
  reviewsError.value = '';
  reviewsLoading.value = true;

  if (props.movie?.id) {
    try {
      const data = await fetchMovieReviews(props.movie.id);
      reviews.value = data.reviews || [];
      averageRating.value = Number(data.average_rating) || 0;
      totalReviews.value = Number(data.total_reviews) || 0;
    } catch (error) {
      reviewsError.value = error.message || 'Không thể tải đánh giá phim.';
    } finally {
      reviewsLoading.value = false;
    }
  } else {
    reviewsLoading.value = false;
  }

  if (props.movie && props.movie.slug) {
    const res = await fetchMovieDetailBySlug(props.movie.slug);
    if (res && res.movie) {
      detailedMovieInfo.value = res;
    }
  }
};

async function submitReview() {
  if (!props.movie || !reviewComment.value.trim()) return;
  reviewBusy.value = true;
  reviewNotice.value = '';
  reviewsError.value = '';
  try {
    await saveMovieReview(props.movie, reviewRating.value, reviewComment.value.trim());
    reviewNotice.value = 'Đánh giá của bạn đã được lưu.';
    const data = await fetchMovieReviews(props.movie.id);
    reviews.value = data.reviews || [];
    averageRating.value = Number(data.average_rating) || 0;
    totalReviews.value = Number(data.total_reviews) || 0;
  } catch (error) {
    reviewsError.value = error.message || 'Không thể gửi đánh giá.';
  } finally {
    reviewBusy.value = false;
  }
}

function formatReviewDate(value) {
  return new Intl.DateTimeFormat('vi-VN', {
    day: '2-digit', month: '2-digit', year: 'numeric',
    timeZone: 'Asia/Ho_Chi_Minh',
  }).format(new Date(value));
}

watch(() => props.movie, loadDetailData, { immediate: true });
</script>

<style scoped>
.modal-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(4, 6, 11, 0.85);
  backdrop-filter: blur(12px);
  z-index: 2000;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px;
  animation: fadeIn 0.3s ease;
}

.modal-container {
  width: 100%;
  max-width: 900px;
  max-height: 90vh;
  background: var(--bg-modal);
  border-radius: var(--radius-lg);
  border: 1px solid var(--border-subtle);
  overflow-y: auto;
  position: relative;
  box-shadow: 0 25px 60px rgba(0, 0, 0, 0.8), 0 0 30px rgba(229, 9, 20, 0.15);
  animation: springScale 0.45s var(--ease-spring);
}

.modal-close-btn {
  position: absolute;
  top: 16px;
  right: 16px;
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: rgba(13, 18, 31, 0.8);
  backdrop-filter: blur(8px);
  border: 1px solid var(--border-subtle);
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  z-index: 20;
  transition: all 0.25s ease;
}

.modal-close-btn:hover {
  background: var(--accent-red);
  transform: rotate(90deg);
}

.modal-header-backdrop {
  position: relative;
  height: 280px;
  width: 100%;
  overflow: hidden;
}

.modal-header-backdrop.has-trailer {
  cursor: pointer;
}

.modal-header-backdrop.has-trailer .header-img {
  transition: transform 700ms var(--ease-spring), filter 350ms ease;
}

.modal-header-backdrop.has-trailer:hover .header-img {
  transform: scale(1.035);
  filter: brightness(0.72);
}

.trailer-preview-btn {
  position: absolute;
  inset: 50% auto auto 50%;
  z-index: 12;
  transform: translate(-50%, -50%);
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 16px 10px 10px;
  border: 1px solid rgba(255, 255, 255, 0.25);
  border-radius: 999px;
  background: rgba(7, 9, 14, 0.72);
  backdrop-filter: blur(14px);
  color: #fff;
  cursor: pointer;
  box-shadow: 0 14px 40px rgba(0, 0, 0, 0.45);
  transition: transform 300ms var(--ease-out-back), background 250ms ease;
}

.trailer-preview-btn:hover {
  transform: translate(-50%, -50%) scale(1.06);
  background: rgba(229, 9, 20, 0.9);
}

.preview-play {
  width: 48px;
  height: 48px;
  display: grid;
  place-items: center;
  border-radius: 50%;
  background: #fff;
  color: var(--accent-red);
}

.preview-copy {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  line-height: 1.2;
}

.preview-copy strong { font-size: 0.92rem; }
.preview-copy small { margin-top: 3px; color: rgba(255, 255, 255, 0.7); }

.header-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.header-overlay {
  position: absolute;
  inset: 0;
  background: linear-gradient(0deg, var(--bg-modal) 0%, rgba(13, 18, 31, 0.4) 60%, rgba(0,0,0,0.6) 100%);
}

.header-content {
  position: absolute;
  bottom: 20px;
  left: 32px;
  right: 32px;
  z-index: 10;
}

.badge-group {
  display: flex;
  gap: 8px;
  margin-bottom: 8px;
}

.modal-title {
  font-family: var(--font-heading);
  font-size: 2.2rem;
  font-weight: 800;
  color: #fff;
}

.modal-origin {
  font-size: 1rem;
  color: var(--text-secondary);
}

.modal-body {
  padding: 32px;
}

.body-grid {
  display: grid;
  grid-template-columns: 220px 1fr;
  gap: 32px;
}

.poster-col {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.poster-card-img {
  width: 100%;
  border-radius: var(--radius-md);
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.5);
  border: 1px solid var(--border-subtle);
}

.full-width {
  width: 100%;
}

.btn-trailer {
  background: rgba(255, 255, 255, 0.08);
  color: #fff;
  border: 1px solid rgba(255, 255, 255, 0.16);
}

.btn-trailer:hover {
  background: #ff0033;
  border-color: #ff0033;
  transform: translateY(-2px);
  box-shadow: 0 10px 24px rgba(255, 0, 51, 0.28);
}

.info-col {
  display: flex;
  flex-direction: column;
  gap: 22px;
}

.meta-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
  padding-bottom: 16px;
  border-bottom: 1px solid var(--border-subtle);
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 0.95rem;
}

.meta-label {
  color: var(--text-muted);
  display: flex;
  align-items: center;
  gap: 6px;
  width: 110px;
}

.meta-val {
  color: var(--text-primary);
  font-weight: 600;
}

.genres-chips {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.chip {
  background: rgba(255, 255, 255, 0.08);
  padding: 3px 10px;
  border-radius: 6px;
  font-size: 0.8rem;
  color: var(--text-primary);
}

.sub-heading {
  font-family: var(--font-heading);
  font-size: 1.1rem;
  font-weight: 700;
  color: #fff;
  margin-bottom: 10px;
}

.cast-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.cast-pill {
  background: rgba(0, 242, 254, 0.08);
  color: var(--accent-cyan);
  border: 1px solid rgba(0, 242, 254, 0.2);
  padding: 4px 12px;
  border-radius: 20px;
  font-size: 0.82rem;
}

.story-text {
  color: var(--text-secondary);
  line-height: 1.6;
  font-size: 0.95rem;
}

.episodes-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.ep-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  background: rgba(255, 255, 255, 0.07);
  border: 1px solid var(--border-subtle);
  color: var(--text-primary);
  padding: 8px 16px;
  border-radius: var(--radius-sm);
  cursor: pointer;
  font-size: 0.85rem;
  font-weight: 600;
  transition: all 0.25s ease;
}

.ep-btn:hover {
  background: var(--accent-gradient);
  border-color: transparent;
  transform: translateY(-2px);
}

.reviews-section {
  margin-top: 34px;
  padding-top: 26px;
  border-top: 1px solid var(--border-subtle);
}

.reviews-heading, .review-card-heading, .review-author, .review-form-footer, .review-login-prompt {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 14px;
}

.reviews-heading h3 { margin: 0 0 5px; font-size: 1.35rem; }
.reviews-heading p, .review-form-hint { margin: 0; color: var(--text-secondary); font-size: .88rem; }
.reviews-summary { display: grid; justify-items: end; gap: 3px; white-space: nowrap; }
.reviews-summary strong { color: #fbbf24; font-size: 1.2rem; }
.reviews-summary span, .review-author time, .review-form-footer > span { color: var(--text-secondary); font-size: .8rem; }
.review-form, .review-login-prompt { margin-top: 22px; padding: 18px; border: 1px solid var(--border-subtle); border-radius: 14px; background: rgba(255,255,255,.035); }
.review-form h4 { margin: 0 0 5px; font-size: 1rem; }
.rating-picker { display: flex; gap: 3px; margin: 12px 0; }
.rating-choice { padding: 0 3px; border: 0; color: #64748b; background: transparent; font-size: 1.7rem; cursor: pointer; }
.rating-choice.selected { color: #fbbf24; }
.rating-choice:focus-visible { outline: 2px solid var(--accent-cyan); outline-offset: 2px; border-radius: 4px; }
.review-form textarea { display: block; width: 100%; min-height: 100px; resize: vertical; padding: 12px; border: 1px solid var(--border-subtle); border-radius: 10px; color: var(--text-primary); background: rgba(0,0,0,.2); font: inherit; }
.review-form-footer { margin-top: 10px; }
.review-form-footer button { padding: 9px 14px; }
.review-form-footer button:disabled { opacity: .55; cursor: not-allowed; }
.review-login-prompt { color: var(--text-secondary); }
.review-list { display: grid; gap: 12px; margin-top: 18px; }
.review-card { padding: 16px 0; border-bottom: 1px solid var(--border-subtle); }
.review-author { justify-content: flex-start; }
.review-avatar { display: grid; width: 38px; height: 38px; flex: 0 0 38px; place-items: center; border-radius: 50%; color: #041014; background: var(--accent-cyan); font-weight: 800; }
.review-author > div { display: grid; gap: 3px; }
.review-author time { font-size: .76rem; }
.review-stars { color: #fbbf24; letter-spacing: 1px; white-space: nowrap; }
.review-stars > span { color: #475569; }
.review-comment { margin: 12px 0 0 52px; color: var(--text-secondary); line-height: 1.6; white-space: pre-wrap; overflow-wrap: anywhere; }
.reviews-message { margin: 16px 0; color: var(--text-secondary); font-size: .9rem; }
.reviews-message.error { color: #fca5a5; }
.reviews-message.success { margin-bottom: 0; color: #86efac; }
.empty-reviews { padding: 20px 0; }

@media (max-width: 768px) {
  .body-grid {
    grid-template-columns: 1fr;
  }
  .poster-col {
    max-width: 240px;
    margin: 0 auto;
  }
  .modal-body { padding: 22px 18px; }
  .reviews-heading, .review-login-prompt { align-items: flex-start; flex-direction: column; }
  .reviews-summary { justify-items: start; }
  .review-card-heading { align-items: flex-start; flex-direction: column; }
  .review-comment { margin-left: 0; }
}
</style>
