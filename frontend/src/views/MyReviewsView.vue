<template>
  <section class="my-reviews-page">
    <div class="my-reviews-container">
      <header class="reviews-page-header">
        <span class="eyebrow">CINESTREAM</span>
        <h1>Đánh giá của tôi</h1>
        <p>Xem lại các bộ phim bạn đã đánh giá và chia sẻ.</p>
      </header>

      <div v-if="loading" class="reviews-state">Đang tải đánh giá…</div>
      <div v-else-if="error" class="reviews-state error-state">
        <span>{{ error }}</span>
        <button class="retry-button" @click="loadReviews">Thử lại</button>
      </div>
      <div v-else-if="!reviews.length" class="reviews-state empty-state">
        <Star :size="36" />
        <strong>Bạn chưa đánh giá phim nào</strong>
        <span>Mở chi tiết một bộ phim để chấm sao và viết bình luận.</span>
        <router-link to="/movies" class="browse-button">Khám phá phim</router-link>
      </div>

      <div v-else class="my-review-list">
        <article v-for="review in reviews" :key="review.id" class="my-review-card">
          <div class="my-review-heading">
            <div>
              <span class="review-kicker">PHIM ĐÃ ĐÁNH GIÁ</span>
              <h2>{{ review.movie_title }}</h2>
            </div>
            <div class="my-review-rating" :aria-label="`${review.rating} trên 5 sao`">
              <span>{{ '★'.repeat(review.rating) }}</span><i>{{ '★'.repeat(5 - review.rating) }}</i>
              <strong>{{ review.rating }}/5</strong>
            </div>
          </div>
          <blockquote>{{ review.comment }}</blockquote>
          <time :datetime="review.updated_at || review.created_at">{{ formatDate(review.updated_at || review.created_at) }}</time>
        </article>
      </div>
    </div>
  </section>
</template>

<script setup>
import { onMounted, ref } from 'vue';
import { Star } from '@lucide/vue';
import { fetchMyReviews } from '../services/reviews.js';

const reviews = ref([]);
const loading = ref(true);
const error = ref('');

async function loadReviews() {
  loading.value = true;
  error.value = '';
  try {
    const data = await fetchMyReviews();
    reviews.value = data.reviews || [];
  } catch (err) {
    error.value = err.message || 'Không thể tải đánh giá của bạn.';
  } finally {
    loading.value = false;
  }
}

function formatDate(value) {
  return new Intl.DateTimeFormat('vi-VN', {
    day: '2-digit', month: '2-digit', year: 'numeric',
    timeZone: 'Asia/Ho_Chi_Minh',
  }).format(new Date(value));
}

onMounted(loadReviews);
</script>

<style scoped>
.my-reviews-page { min-height: 80vh; padding: 116px 24px 64px; color: var(--text-primary); }
.my-reviews-container { width: min(100%, 900px); margin: 0 auto; }
.reviews-page-header { margin-bottom: 28px; }
.eyebrow, .review-kicker { color: var(--accent-cyan); font-size: .75rem; font-weight: 800; letter-spacing: .15em; }
.reviews-page-header h1 { margin: 8px 0; font-family: var(--font-heading); font-size: clamp(2rem, 5vw, 3rem); }
.reviews-page-header p { color: var(--text-secondary); }
.reviews-state { display: flex; min-height: 180px; align-items: center; justify-content: center; gap: 14px; padding: 24px; border: 1px solid var(--border-subtle); border-radius: 16px; background: var(--bg-card); color: var(--text-secondary); }
.empty-state { flex-direction: column; text-align: center; }
.empty-state strong { color: var(--text-primary); font-size: 1.1rem; }
.empty-state > svg { color: #fbbf24; }
.my-review-list { display: grid; gap: 16px; }
.my-review-card { padding: 22px 24px; border: 1px solid var(--border-subtle); border-radius: 16px; background: var(--bg-card); }
.my-review-heading { display: flex; align-items: flex-start; justify-content: space-between; gap: 18px; }
.my-review-heading h2 { margin: 7px 0 0; font-size: 1.2rem; }
.my-review-rating { display: flex; align-items: center; gap: 2px; white-space: nowrap; color: #fbbf24; font-size: 1.1rem; }
.my-review-rating i { color: #475569; font-style: normal; }
.my-review-rating strong { margin-left: 7px; font-size: .9rem; }
.my-review-card blockquote { margin: 20px 0 12px; padding-left: 14px; border-left: 3px solid var(--accent-cyan); color: var(--text-secondary); line-height: 1.7; white-space: pre-wrap; overflow-wrap: anywhere; }
.my-review-card time { color: var(--text-secondary); font-size: .8rem; }
.browse-button, .retry-button { padding: 10px 16px; border: 0; border-radius: 10px; background: var(--accent-primary, #6d5dfc); color: white; font-weight: 700; text-decoration: none; cursor: pointer; }
.error-state { flex-wrap: wrap; color: #fca5a5; }
@media (max-width: 600px) { .my-reviews-page { padding: 100px 14px 44px; } .my-review-card { padding: 18px 16px; } .my-review-heading { flex-direction: column; } }
</style>
