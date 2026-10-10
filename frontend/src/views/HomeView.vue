<template>
  <div class="home-view">
    <!-- Hero Banner for Featured Movie -->
    <HeroBanner 
      :featured-movie="featuredMovie"
      @play="openPlayer" 
      @open-detail="openDetail" 
    />

    <!-- Main Sections -->
    <div class="home-container">
      <!-- Quick Genre Chips Bar -->
      <div v-reveal="{ delay: 100 }" class="genre-bar">
        <span 
          v-for="genre in genresList" 
          :key="genre" 
          class="genre-chip"
          @click="navigateToGenre(genre)"
        >
          {{ genre }}
        </span>
      </div>

      <!-- Trending Movies Row -->
      <MovieGrid 
        title="🔥 Phim Thịnh Hành" 
        :movies="trendingMovies" 
        :loading="loading"
        @play="openPlayer" 
        @open-detail="openDetail" 
      />

      <!-- Top Rated Movies Row -->
      <MovieGrid 
        title="⭐ Top Đánh Giá Cao" 
        :movies="topRatedMovies" 
        :loading="loading"
        @play="openPlayer" 
        @open-detail="openDetail" 
      />

      <!-- All Latest Movies -->
      <MovieGrid 
        title="🎬 Mới Cập Nhật" 
        :movies="allMovies" 
        :loading="loading"
        @play="openPlayer" 
        @open-detail="openDetail" 
      />
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import HeroBanner from '../components/HeroBanner.vue';
import MovieGrid from '../components/MovieGrid.vue';
import { fetchMovies, fetchGenres, GENRES_LIST } from '../services/api.js';

const router = useRouter();
const emit = defineEmits(['play', 'open-detail']);

const allMovies = ref([]);
const loading = ref(true);
const genresList = ref(GENRES_LIST);

const featuredMovie = computed(() => {
  return allMovies.value.find(m => m.is_featured) || allMovies.value[0] || null;
});

const trendingMovies = computed(() => {
  return allMovies.value.filter(m => m.is_trending || m.rating >= 8.5);
});

const topRatedMovies = computed(() => {
  return [...allMovies.value].sort((a, b) => b.rating - a.rating);
});

const openPlayer = (movie, ep = null) => {
  emit('play', movie, ep);
};

const openDetail = (movie) => {
  emit('open-detail', movie);
};

const navigateToGenre = (genre) => {
  if (genre === 'Tất Cả') {
    router.push('/movies');
  } else {
    router.push({ path: '/movies', query: { genre } });
  }
};

onMounted(async () => {
  try {
    const [fetchedMoviesList, fetchedGenresList] = await Promise.all([
      fetchMovies(1),
      fetchGenres()
    ]);
    allMovies.value = fetchedMoviesList;
    genresList.value = fetchedGenresList;
  } finally {
    loading.value = false;
  }
});
</script>

<style scoped>
.home-container {
  max-width: 1320px;
  margin: 0 auto;
  padding: 30px 24px 0 24px;
}

.genre-bar {
  display: flex;
  gap: 10px;
  overflow-x: auto;
  padding-bottom: 14px;
  margin-bottom: 36px;
  scrollbar-width: none; /* Firefox */
}

.genre-bar::-webkit-scrollbar {
  display: none;
}

.genre-chip {
  background: var(--bg-card);
  border: 1px solid var(--border-subtle);
  color: var(--text-secondary);
  padding: 10px 20px;
  border-radius: 24px;
  font-size: 0.9rem;
  font-weight: 600;
  cursor: pointer;
  white-space: nowrap;
  transition: all 0.25s ease;
}

.genre-chip:hover {
  background: rgba(229, 9, 20, 0.15);
  border-color: var(--accent-red);
  color: #fff;
  transform: translateY(-2px);
}
</style>
