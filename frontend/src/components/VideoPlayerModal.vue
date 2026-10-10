<template>
  <Teleport to="body">
    <div v-if="movie" class="player-backdrop">
      <!-- Top Bar -->
      <div class="player-topbar">
        <div class="movie-meta">
          <button class="back-btn" @click="$emit('close')">
            <ArrowLeft :size="20" /> Quay Lại
          </button>
          <div class="title-info">
            <h3 class="playing-title">{{ movie.title }}</h3>
            <span class="playing-ep" v-if="activeEpisode">{{ activeEpisode.name }}</span>
          </div>
        </div>

        <button class="close-icon-btn" @click="$emit('close')">
          <X :size="24" />
        </button>
      </div>

      <!-- Player Content Wrapper -->
      <div class="video-container">
        <iframe
          v-if="isEmbedSource"
          class="video-element embed-player"
          :src="currentVideoUrl"
          :title="`Trình phát ${movie.title}`"
          allow="autoplay; fullscreen; picture-in-picture; encrypted-media"
          allowfullscreen
          referrerpolicy="origin-when-cross-origin"
        ></iframe>
        <video 
          v-else
          ref="videoRef"
          class="video-element"
          :src="currentVideoUrl" 
          controls 
          autoplay
        >
          Trình duyệt của bạn không hỗ trợ phát video HTML5.
        </video>
      </div>

      <!-- Episode Switcher Bar -->
      <div v-if="movie.episodes && movie.episodes.length > 1" class="episodes-bar">
        <span class="ep-label"><Film :size="16" /> Chọn Tập Phim:</span>
        <div class="ep-buttons-row">
          <button 
            v-for="ep in movie.episodes" 
            :key="ep.id"
            class="ep-tab"
            :class="{ 'active': activeEpisode && activeEpisode.id === ep.id }"
            @click="selectEpisode(ep)"
          >
            {{ ep.name }}
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue';
import { ArrowLeft, X, Film } from '@lucide/vue';

const props = defineProps({
  movie: {
    type: Object,
    default: null
  },
  episode: {
    type: Object,
    default: null
  }
});

defineEmits(['close']);

const videoRef = ref(null);
const activeEpisode = ref(props.episode || (props.movie?.episodes?.[0] || null));

const toEmbedUrl = (source = '') => {
  try {
    const url = new URL(source);
    let videoId = '';
    if (url.hostname.includes('youtu.be')) videoId = url.pathname.split('/').filter(Boolean)[0] || '';
    if (url.hostname.includes('youtube.com')) {
      videoId = url.searchParams.get('v') || '';
      if (!videoId && /\/(shorts|embed)\//.test(url.pathname)) {
        videoId = url.pathname.split('/').filter(Boolean)[1] || '';
      }
    }
    if (videoId) return `https://www.youtube.com/embed/${videoId}?autoplay=1&rel=0&modestbranding=1`;
  } catch (_) {
    // Keep non-YouTube sources unchanged.
  }
  return source;
};

const currentVideoUrl = computed(() => {
  if (activeEpisode.value) {
    return toEmbedUrl(activeEpisode.value.link_embed
      || activeEpisode.value.link_m3u8
      || activeEpisode.value.url
      || activeEpisode.value.filename);
  }
  return props.movie?.video_url || 'https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/TearsOfSteel.mp4';
});

const isEmbedSource = computed(() => {
  if (activeEpisode.value?.link_embed) return true;
  const url = currentVideoUrl.value || '';
  return /streamvsmov\.com\/video\//i.test(url) || !/\.(mp4|webm|ogg|m3u8)(?:[?#]|$)/i.test(url);
});

const selectEpisode = (ep) => {
  activeEpisode.value = ep;
  if (videoRef.value && !isEmbedSource.value) {
    videoRef.value.load();
    videoRef.value.play();
  }
};

watch(() => props.episode, (newEp) => {
  if (newEp) activeEpisode.value = newEp;
});
</script>

<style scoped>
.player-backdrop {
  position: fixed;
  inset: 0;
  background: #000;
  z-index: 3000;
  display: flex;
  flex-direction: column;
  animation: fadeIn 0.25s ease;
}

.player-topbar {
  height: 70px;
  background: rgba(10, 14, 23, 0.95);
  backdrop-filter: blur(16px);
  padding: 0 24px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid var(--border-subtle);
}

.movie-meta {
  display: flex;
  align-items: center;
  gap: 20px;
}

.back-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  background: rgba(255, 255, 255, 0.08);
  border: 1px solid var(--border-subtle);
  color: #fff;
  padding: 8px 16px;
  border-radius: var(--radius-sm);
  cursor: pointer;
  font-weight: 600;
  font-size: 0.9rem;
  transition: all 0.2s ease;
}

.back-btn:hover {
  background: rgba(255, 255, 255, 0.18);
  border-color: rgba(255, 255, 255, 0.3);
}

.title-info {
  display: flex;
  flex-direction: column;
}

.playing-title {
  font-family: var(--font-heading);
  font-size: 1.2rem;
  font-weight: 700;
  color: #fff;
}

.playing-ep {
  font-size: 0.82rem;
  color: var(--accent-cyan);
}

.close-icon-btn {
  background: none;
  border: none;
  color: var(--text-secondary);
  cursor: pointer;
  transition: color 0.2s ease;
}

.close-icon-btn:hover {
  color: var(--accent-red);
}

.video-container {
  flex: 1;
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #000;
}

.video-element {
  width: 100%;
  max-width: 1200px;
  max-height: 80vh;
  outline: none;
  box-shadow: 0 0 50px rgba(0, 0, 0, 0.9);
}

.embed-player {
  aspect-ratio: 16 / 9;
  height: min(80vh, 675px);
  border: 0;
  background: #000;
}

.episodes-bar {
  background: rgba(10, 14, 23, 0.95);
  border-top: 1px solid var(--border-subtle);
  padding: 16px 24px;
  display: flex;
  align-items: center;
  gap: 16px;
}

.ep-label {
  display: flex;
  align-items: center;
  gap: 6px;
  color: var(--text-secondary);
  font-size: 0.9rem;
  font-weight: 600;
}

.ep-buttons-row {
  display: flex;
  gap: 10px;
  overflow-x: auto;
  padding-bottom: 4px;
}

.ep-tab {
  background: rgba(255, 255, 255, 0.08);
  border: 1px solid var(--border-subtle);
  color: var(--text-primary);
  padding: 8px 16px;
  border-radius: var(--radius-sm);
  cursor: pointer;
  font-size: 0.85rem;
  font-weight: 600;
  white-space: nowrap;
  transition: all 0.25s ease;
}

.ep-tab.active {
  background: var(--accent-gradient);
  border-color: transparent;
  box-shadow: 0 0 15px rgba(229, 9, 20, 0.5);
}
</style>
