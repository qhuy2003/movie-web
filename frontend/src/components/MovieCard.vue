<template>
  <div class="movie-card" @click="$emit('open-detail', movie)">
    <!-- Poster Container -->
    <div class="poster-container">
      <img :src="movie.poster_url" :alt="movie.title" class="poster-img" loading="lazy" />
      <div class="shine-sweep"></div>

      <!-- Badges overlay -->
      <div class="badges-top">
        <span class="badge badge-quality">{{ movie.quality || 'HD' }}</span>
        <button 
          class="favorite-btn" 
          :class="{ 'saved': isSaved }" 
          @click.stop="toggleWatchlist"
          title="Lưu yêu thích"
        >
          <Heart :size="16" :fill="isSaved ? '#e50914' : 'none'" :color="isSaved ? '#e50914' : '#fff'" />
        </button>
      </div>

      <!-- Hover Play Overlay -->
      <div class="hover-overlay">
        <button class="play-circle-btn" @click.stop="$emit('play', movie)" title="Xem phim ngay">
          <Play :size="24" fill="currentColor" />
        </button>
      </div>

      <!-- Rating Badge at Bottom Left -->
      <div class="rating-badge">
        <Star :size="13" fill="#ffb800" color="#ffb800" />
        <span>{{ movie.rating }}</span>
      </div>
    </div>

    <!-- Movie Meta Info -->
    <div class="movie-info">
      <h3 class="movie-title" :title="movie.title">{{ movie.title }}</h3>
      <div class="movie-sub">
        <span class="movie-year">{{ movie.year }}</span>
        <span class="dot">•</span>
        <span class="movie-genre">{{ primaryGenre }}</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue';
import { Play, Heart, Star } from '@lucide/vue';
import { watchlistStore } from '../services/watchlist.js';

const props = defineProps({
  movie: {
    type: Object,
    required: true
  }
});

defineEmits(['play', 'open-detail']);

const isSaved = computed(() => watchlistStore.isSaved(props.movie.id));

const primaryGenre = computed(() => {
  if (props.movie.genres && props.movie.genres.length > 0) {
    return props.movie.genres[0];
  }
  return 'Phim Phố';
});

const toggleWatchlist = () => {
  watchlistStore.toggle(props.movie);
};
</script>

<style scoped>
.movie-card {
  position: relative;
  border-radius: var(--radius-md);
  overflow: hidden;
  background: var(--bg-card);
  border: 1px solid var(--border-subtle);
  cursor: pointer;
  transition: transform 0.4s var(--ease-spring), border-color 0.3s ease, box-shadow 0.4s var(--ease-spring), background 0.3s ease;
  display: flex;
  flex-direction: column;
}

.movie-card:hover {
  transform: translateY(-10px) scale(1.03);
  border-color: rgba(229, 9, 20, 0.6);
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.8), 0 0 30px rgba(229, 9, 20, 0.35);
  background: var(--bg-card-hover);
}

.poster-container {
  position: relative;
  width: 100%;
  padding-top: 145%; /* Aspect ratio ~ 2:3 */
  overflow: hidden;
  background: #111522;
}

.poster-img {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.6s var(--ease-spring);
}

.movie-card:hover .poster-img {
  transform: scale(1.09);
}

/* Framer Motion Style Glass Shine Sweep Beam */
.shine-sweep {
  position: absolute;
  top: 0;
  left: -150%;
  width: 100%;
  height: 100%;
  background: linear-gradient(
    90deg, 
    transparent 0%, 
    rgba(255, 255, 255, 0.25) 50%, 
    transparent 100%
  );
  transform: skewX(-25deg);
  transition: left 0.8s var(--ease-spring);
  z-index: 3;
  pointer-events: none;
}

.movie-card:hover .shine-sweep {
  left: 150%;
}

.badges-top {
  position: absolute;
  top: 10px;
  left: 10px;
  right: 10px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  z-index: 5;
}

.favorite-btn {
  width: 34px;
  height: 34px;
  border-radius: 50%;
  background: rgba(13, 18, 31, 0.8);
  backdrop-filter: blur(8px);
  border: 1px solid var(--border-subtle);
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: transform 0.3s var(--ease-out-back), background 0.25s ease;
}

.favorite-btn:hover {
  transform: scale(1.2);
  background: rgba(13, 18, 31, 0.95);
}

.favorite-btn.saved {
  border-color: var(--accent-red);
}

.hover-overlay {
  position: absolute;
  inset: 0;
  background: rgba(7, 9, 14, 0.7);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transition: opacity 0.3s var(--ease-spring);
  z-index: 4;
}

.movie-card:hover .hover-overlay {
  opacity: 1;
}

.play-circle-btn {
  width: 58px;
  height: 58px;
  border-radius: 50%;
  background: var(--accent-gradient);
  color: #fff;
  border: none;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 6px 25px rgba(229, 9, 20, 0.7);
  cursor: pointer;
  transform: scale(0.85);
  transition: transform 0.4s var(--ease-out-back);
}

.movie-card:hover .play-circle-btn {
  transform: scale(1);
}

.play-circle-btn:hover {
  transform: scale(1.18) !important;
}

.rating-badge {
  position: absolute;
  bottom: 10px;
  left: 10px;
  background: rgba(13, 18, 31, 0.85);
  backdrop-filter: blur(8px);
  padding: 3px 8px;
  border-radius: 6px;
  font-size: 0.75rem;
  font-weight: 700;
  color: var(--accent-gold);
  display: flex;
  align-items: center;
  gap: 4px;
  z-index: 5;
  border: 1px solid rgba(255, 184, 0, 0.25);
}

.movie-info {
  padding: 14px;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.movie-title {
  font-family: var(--font-heading);
  font-size: 0.98rem;
  font-weight: 700;
  color: var(--text-primary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  transition: color 0.2s ease;
}

.movie-card:hover .movie-title {
  color: #fff;
}

.movie-sub {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 0.8rem;
  color: var(--text-secondary);
}

.dot {
  color: var(--text-muted);
}
</style>
