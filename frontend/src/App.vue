<template>
  <div class="app-container">
    <!-- Navbar Component -->
    <Navbar @open-auth="showAuthModal = true" />

    <!-- Main View Content -->
    <main class="main-content">
      <router-view v-slot="{ Component }">
        <transition name="page-motion" mode="out-in">
          <component 
            :is="Component" 
            @play="handlePlayMovie" 
            @open-detail="handleOpenDetail" 
          />
        </transition>
      </router-view>
    </main>

    <!-- Footer Component -->
    <Footer />

    <!-- Movie Detail Modal Popup -->
    <MovieDetailModal 
      v-if="showDetailModal" 
      :movie="selectedMovie" 
      @close="showDetailModal = false" 
      @play="handlePlayMovie" 
    />

    <!-- Video Player Modal Screen -->
    <VideoPlayerModal 
      v-if="showPlayerModal" 
      :movie="selectedMovie" 
      :episode="selectedEpisode"
      @close="showPlayerModal = false" 
    />

    <!-- Auth Modal Popup -->
    <AuthModal 
      v-if="showAuthModal" 
      @close="showAuthModal = false" 
    />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import Navbar from './components/Navbar.vue';
import Footer from './components/Footer.vue';
import MovieDetailModal from './components/MovieDetailModal.vue';
import VideoPlayerModal from './components/VideoPlayerModal.vue';
import AuthModal from './components/AuthModal.vue';
import { watchlistStore } from './services/watchlist.js';
import { fetchMovieDetailBySlug } from './services/api.js';

const selectedMovie = ref(null);
const selectedEpisode = ref(null);
const showDetailModal = ref(false);
const showPlayerModal = ref(false);
const showAuthModal = ref(false);

onMounted(() => {
  watchlistStore.onAuthRequired = () => {
    showAuthModal.value = true;
  };
});

const handleOpenDetail = (movie) => {
  selectedMovie.value = movie;
  showDetailModal.value = true;
};

const handlePlayMovie = async (movie, ep = null) => {
  let playableMovie = movie;
  let playableEpisode = ep;

  // A card only contains catalogue metadata. Fetch its fresh playback source by slug.
  if (!playableEpisode && movie.slug) {
    const detail = await fetchMovieDetailBySlug(movie.slug);
    const episodes = Array.isArray(detail?.episodes)
      ? detail.episodes.flatMap(group => (group.server_data || []).map((item, index) => ({
          ...item,
          id: item.id || item.slug || `${group.server_name || 'server'}-${index}`,
          server_name: group.server_name || ''
        })))
      : [];
    if (episodes.length) {
      playableMovie = { ...movie, episodes };
      playableEpisode = episodes[0];
    }
  }

  selectedMovie.value = playableMovie;
  selectedEpisode.value = playableEpisode || playableMovie.episodes?.[0] || null;
  showDetailModal.value = false;
  showPlayerModal.value = true;
};
</script>

<style>
/* Global App Styles Imports */
</style>
