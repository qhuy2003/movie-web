<template>
  <div v-if="featuredMovie" class="hero-banner">
    <!-- Backdrop Image with Gradient Overlay -->
    <div class="backdrop-wrapper">
      <img :src="featuredMovie.thumb_url" :alt="featuredMovie.title" class="backdrop-img" />
      <div class="backdrop-overlay"></div>
    </div>

    <!-- Banner Content -->
    <div class="hero-container">
      <div class="hero-content">
        <!-- Badges Bar -->
        <div v-reveal="{ delay: 80 }" class="badge-group">
          <span class="badge badge-rating">⭐ {{ featuredMovie.rating }} TMDB</span>
          <span class="badge badge-quality">{{ featuredMovie.quality || '4K UHD' }}</span>
          <span class="badge badge-genre">{{ featuredMovie.year }}</span>
        </div>

        <!-- Movie Title -->
        <h1 v-reveal="{ delay: 150, distance: 38 }" class="hero-title">{{ featuredMovie.title }}</h1>
        <h3 v-reveal="{ delay: 210 }" class="hero-subtitle">{{ featuredMovie.origin_name }}</h3>

        <!-- Genres List -->
        <div v-reveal="{ delay: 260 }" class="genres-row">
          <span v-for="g in featuredMovie.genres" :key="g" class="genre-tag">
            {{ g }}
          </span>
        </div>

        <!-- Synopsis -->
        <p v-reveal="{ delay: 320 }" class="hero-desc">{{ featuredMovie.description }}</p>

        <!-- Actions -->
        <div v-reveal="{ delay: 380 }" class="hero-actions">
          <button class="btn btn-primary" @click="$emit('open-detail', featuredMovie)">
            <Play :size="20" fill="currentColor" /> Xem Trailer
          </button>
          <button class="btn btn-secondary" @click="$emit('open-detail', featuredMovie)">
            <Info :size="20" /> Chi Tiết Phim
          </button>
          <button 
            class="btn btn-icon" 
            :class="{ 'active': isSaved }"
            @click="toggleWatchlist"
            title="Lưu vào phim yêu thích"
          >
            <Heart :size="20" :fill="isSaved ? '#e50914' : 'none'" :color="isSaved ? '#e50914' : '#fff'" />
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue';
import { Play, Info, Heart } from '@lucide/vue';
import { watchlistStore } from '../services/watchlist.js';

const props = defineProps({
  featuredMovie: {
    type: Object,
    default: null
  }
});

defineEmits(['play', 'open-detail']);

const isSaved = computed(() => {
  return props.featuredMovie ? watchlistStore.isSaved(props.featuredMovie.id) : false;
});

const toggleWatchlist = () => {
  if (props.featuredMovie) {
    watchlistStore.toggle(props.featuredMovie);
  }
};
</script>

<style scoped>
.hero-banner {
  position: relative;
  width: 100%;
  height: 80vh;
  min-height: 550px;
  max-height: 750px;
  display: flex;
  align-items: center;
  overflow: hidden;
}

.backdrop-wrapper {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
}

.backdrop-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  object-position: center 20%;
  filter: brightness(0.75) contrast(1.1);
  transform: scale(1.02);
  transition: transform 10s ease-out;
}

.hero-banner:hover .backdrop-img {
  transform: scale(1.06);
}

.backdrop-overlay {
  position: absolute;
  inset: 0;
  background: linear-gradient(
    90deg, 
    #07090e 0%, 
    rgba(7, 9, 14, 0.85) 45%, 
    rgba(7, 9, 14, 0.2) 100%
  ),
  linear-gradient(
    0deg, 
    #07090e 0%, 
    rgba(7, 9, 14, 0) 50%,
    rgba(7, 9, 14, 0.5) 100%
  );
}

.hero-container {
  max-width: 1320px;
  width: 100%;
  margin: 0 auto;
  padding: 0 24px;
  position: relative;
  z-index: 10;
}

.hero-content { max-width: 620px; }

.badge-group {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 16px;
}

.hero-title {
  font-family: var(--font-heading);
  font-size: 3.2rem;
  font-weight: 800;
  line-height: 1.1;
  color: #fff;
  margin-bottom: 6px;
  text-shadow: 0 4px 20px rgba(0, 0, 0, 0.8);
}

.hero-subtitle {
  font-size: 1.1rem;
  font-weight: 500;
  color: var(--text-secondary);
  margin-bottom: 16px;
}

.genres-row {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 20px;
}

.genre-tag {
  background: rgba(255, 255, 255, 0.12);
  color: var(--text-primary);
  font-size: 0.82rem;
  font-weight: 600;
  padding: 4px 12px;
  border-radius: 20px;
  backdrop-filter: blur(8px);
}

.hero-desc {
  font-size: 0.98rem;
  color: var(--text-secondary);
  line-height: 1.6;
  margin-bottom: 28px;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
  text-shadow: 0 2px 10px rgba(0,0,0,0.8);
}

.hero-actions {
  display: flex;
  align-items: center;
  gap: 14px;
}

@media (max-width: 768px) {
  .hero-banner {
    height: 70vh;
    min-height: 480px;
  }
  .hero-title {
    font-size: 2.2rem;
  }
  .backdrop-overlay {
    background: linear-gradient(
      180deg, 
      rgba(7, 9, 14, 0.4) 0%, 
      rgba(7, 9, 14, 0.95) 80%,
      #07090e 100%
    );
  }
}
</style>
