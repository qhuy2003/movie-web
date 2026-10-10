<template>
  <div class="watchlist-view">
    <div class="watchlist-container">
      <div class="header-bar">
        <div>
          <h1 class="page-title">Danh Sách Phim Yêu Thích</h1>
          <p class="page-sub">Các bộ phim bạn đã lưu lại để xem về sau.</p>
        </div>
        <button 
          v-if="watchlistStore.items.length > 0" 
          class="btn btn-secondary clear-btn"
          @click="watchlistStore.clear()"
        >
          <Trash2 :size="16" /> Xóa Tất Cả
        </button>
      </div>

      <MovieGrid 
        :movies="watchlistStore.items" 
        @play="$emit('play', $event)" 
        @open-detail="$emit('open-detail', $event)" 
      />
    </div>
  </div>
</template>

<script setup>
import { Trash2 } from '@lucide/vue';
import MovieGrid from '../components/MovieGrid.vue';
import { watchlistStore } from '../services/watchlist.js';

defineEmits(['play', 'open-detail']);
</script>

<style scoped>
.watchlist-view {
  padding-top: 30px;
  min-height: 80vh;
}

.watchlist-container {
  max-width: 1320px;
  margin: 0 auto;
  padding: 0 24px;
}

.header-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 32px;
}

.page-title {
  font-family: var(--font-heading);
  font-size: 2.2rem;
  font-weight: 800;
  color: #fff;
  margin-bottom: 6px;
}

.page-sub {
  color: var(--text-secondary);
  font-size: 0.95rem;
}

.clear-btn:hover {
  color: var(--accent-red);
  border-color: var(--accent-red);
}
</style>
