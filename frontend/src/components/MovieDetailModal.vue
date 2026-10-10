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
              <button class="btn btn-primary full-width" @click="$emit('play', activeMovie)">
                <Play :size="18" fill="currentColor" /> Xem Phim
              </button>
              <button
                v-if="activeMovie.trailer_url"
                class="btn btn-trailer full-width"
                @click="playTrailer"
              >
                <Play :size="18" /> Xem Trailer
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
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue';
import { X, Play, Heart, Clock, User, Tag } from '@lucide/vue';
import { watchlistStore } from '../services/watchlist.js';
import { fetchMovieDetailBySlug } from '../services/api.js';

const props = defineProps({
  movie: {
    type: Object,
    default: null
  }
});

const emit = defineEmits(['close', 'play']);

const detailedMovieInfo = ref(null);

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
  if (!activeMovie.value?.trailer_url) return;
  emit('play', activeMovie.value, {
    id: `trailer-${activeMovie.value.id}`,
    name: 'Trailer chính thức',
    link_embed: activeMovie.value.trailer_url,
    is_trailer: true
  });
};

const loadDetailData = async () => {
  detailedMovieInfo.value = null;
  if (props.movie && props.movie.slug) {
    const res = await fetchMovieDetailBySlug(props.movie.slug);
    if (res && res.movie) {
      detailedMovieInfo.value = res;
    }
  }
};

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

@media (max-width: 768px) {
  .body-grid {
    grid-template-columns: 1fr;
  }
  .poster-col {
    max-width: 240px;
    margin: 0 auto;
  }
}
</style>
