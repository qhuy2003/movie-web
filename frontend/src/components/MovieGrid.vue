<template>
  <section v-reveal class="movie-grid-section">
    <!-- Header Row -->
    <div v-if="title" class="section-header">
      <h2 class="section-title">
        {{ title }}
        <span v-if="movies.length > 0" class="count-chip">{{ movies.length }}</span>
      </h2>
      <slot name="action"></slot>
    </div>

    <!-- Skeleton Loaders -->
    <div v-if="loading" class="grid-container">
      <div v-for="n in 8" :key="n" class="skeleton-card">
        <div class="skeleton-poster"></div>
        <div class="skeleton-line title"></div>
        <div class="skeleton-line sub"></div>
      </div>
    </div>

    <!-- Empty State -->
    <div v-else-if="movies.length === 0" class="empty-state">
      <Film :size="48" class="empty-icon" />
      <h3>Không tìm thấy phim phù hợp</h3>
      <p>Thử tìm kiếm với từ khóa khác hoặc thay đổi bộ lọc thể loại.</p>
    </div>

    <!-- Grid List -->
    <div v-else class="grid-container">
      <MovieCard
        v-for="(movie, index) in movies"
        :key="movie.id" 
        :movie="movie" 
        v-reveal="{ delay: Math.min(index, 7) * 55, distance: 22 }"
        @play="$emit('play', $event)" 
        @open-detail="$emit('open-detail', $event)"
      />
    </div>
  </section>
</template>

<script setup>
import { Film } from '@lucide/vue';
import MovieCard from './MovieCard.vue';

defineProps({
  title: {
    type: String,
    default: ''
  },
  movies: {
    type: Array,
    default: () => []
  },
  loading: {
    type: Boolean,
    default: false
  }
});

defineEmits(['play', 'open-detail']);
</script>

<style scoped>
.movie-grid-section {
  margin-bottom: 48px;
}

.count-chip {
  font-size: 0.85rem;
  font-weight: 600;
  color: var(--text-muted);
  background: rgba(255, 255, 255, 0.08);
  padding: 2px 10px;
  border-radius: 12px;
  margin-left: 6px;
}

.grid-container {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 22px;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  text-align: center;
  background: var(--bg-card);
  border-radius: var(--radius-md);
  border: 1px dashed var(--border-subtle);
}

.empty-icon {
  color: var(--text-muted);
  margin-bottom: 16px;
}

.empty-state h3 {
  font-size: 1.2rem;
  margin-bottom: 6px;
  color: var(--text-primary);
}

.empty-state p {
  color: var(--text-secondary);
  font-size: 0.9rem;
}

/* Skeleton Pulse Loading */
.skeleton-card {
  border-radius: var(--radius-md);
  overflow: hidden;
  background: var(--bg-card);
  padding-bottom: 14px;
}

.skeleton-poster {
  width: 100%;
  padding-top: 145%;
  background: linear-gradient(90deg, #111726 25%, #1d263b 50%, #111726 75%);
  background-size: 200% 100%;
  animation: loadingPulse 1.5s infinite;
}

.skeleton-line {
  height: 12px;
  margin: 10px 14px 0 14px;
  border-radius: 6px;
  background: linear-gradient(90deg, #111726 25%, #1d263b 50%, #111726 75%);
  background-size: 200% 100%;
  animation: loadingPulse 1.5s infinite;
}

.skeleton-line.title {
  width: 80%;
}

.skeleton-line.sub {
  width: 50%;
}

@keyframes loadingPulse {
  0% { background-position: 200% 0; }
  100% { background-position: -200% 0; }
}

@media (max-width: 640px) {
  .grid-container {
    grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
    gap: 14px;
  }
}
</style>
