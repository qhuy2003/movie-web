<template>
  <div class="catalog-view">
    <div class="catalog-container">
      <!-- Title Header -->
      <div class="catalog-header">
        <h1 class="page-title">Khám Phá Thư Viện Phim</h1>
        <p class="page-sub">Tìm kiếm phim yêu thích theo từ khóa, thể loại, quốc gia hoặc trạng thái phát sóng.</p>
      </div>

      <!-- Controls & Filter Toolbar -->
      <div class="filter-toolbar glass-panel">
        <!-- Search Input -->
        <div class="filter-search">
          <Search :size="18" class="search-icon" />
          <input 
            type="text" 
            v-model="searchQuery" 
            placeholder="Nhập tên phim cần tìm..." 
          />
          <button v-if="searchQuery" class="clear-btn" @click="searchQuery = ''">
            <X :size="16" />
          </button>
        </div>

        <!-- Genre Dropdown Selector -->
        <div class="filter-select-group">
          <label><SlidersHorizontal :size="16" /> Thể Loại:</label>
          <select v-model="selectedGenre" class="custom-select">
            <option v-for="g in genresList" :key="g" :value="g">{{ g }}</option>
          </select>
        </div>

        <!-- Country Dropdown Selector -->
        <div class="filter-select-group">
          <label><Globe :size="16" /> Quốc Gia:</label>
          <select v-model="selectedCountry" class="custom-select" @change="handleCountryChange">
            <option value="all">Tất Cả Quốc Gia</option>
            <option v-for="c in countriesList" :key="c._id" :value="c.slug">{{ c.name }}</option>
          </select>
        </div>

        <!-- Type Selector (Phim bộ / Phim lẻ) -->
        <div class="filter-select-group">
          <label><CalendarDays :size="16" /> Năm Phát Hành:</label>
          <select v-model="selectedYear" class="custom-select year-select" @change="handleYearChange">
            <option value="all">Tất Cả Năm</option>
            <option v-for="year in yearsList" :key="year.slug || year.name || year" :value="year.slug || year.name || year">
              {{ year.name || year.slug || year }}
            </option>
          </select>
        </div>

        <div class="filter-select-group">
          <label><Film :size="16" /> Loại Phim:</label>
          <select v-model="selectedType" class="custom-select" @change="selectedYear !== 'all' ? handleYearChange() : handleCountryChange()">
            <option value="all">Tất Cả Loại</option>
            <option value="series">Phim Bộ (Series)</option>
            <option value="single">Phim Lẻ (Single)</option>
          </select>
        </div>

        <!-- Sort Selector -->
        <div class="filter-select-group">
          <label><ArrowUpDown :size="16" /> Sắp Xếp:</label>
          <select v-model="sortBy" class="custom-select">
            <option value="rating">Điểm TMDB Cao Nhất</option>
            <option value="year">Năm Phát Hành Mới Nhất</option>
            <option value="title">Tên Phim (A-Z)</option>
          </select>
        </div>
      </div>

      <!-- Active Filters Display Chips -->
      <div v-if="searchQuery || selectedGenre !== 'Tất Cả' || selectedCountry !== 'all' || selectedYear !== 'all' || selectedType !== 'all'" class="active-chips">
        <span class="active-label">Bộ lọc đang mở:</span>
        <span v-if="searchQuery" class="filter-chip">
          Từ khóa: "{{ searchQuery }}" <X :size="14" @click="searchQuery = ''" class="chip-x" />
        </span>
        <span v-if="selectedGenre !== 'Tất Cả'" class="filter-chip">
          Thể loại: {{ selectedGenre }} <X :size="14" @click="selectedGenre = 'Tất Cả'" class="chip-x" />
        </span>
        <span v-if="selectedCountry !== 'all'" class="filter-chip">
          Quốc gia: {{ getCountryName(selectedCountry) }} <X :size="14" @click="selectedCountry = 'all'; reloadDefaultMovies();" class="chip-x" />
        </span>
        <span v-if="selectedYear !== 'all'" class="filter-chip">
          Năm: {{ selectedYear }} <X :size="14" @click="selectedYear = 'all'; reloadDefaultMovies();" class="chip-x" />
        </span>
        <span v-if="selectedType !== 'all'" class="filter-chip">
          Loại: {{ selectedType === 'series' ? 'Phim Bộ' : 'Phim Lẻ' }} <X :size="14" @click="selectedType = 'all'" class="chip-x" />
        </span>
        <button class="reset-all" @click="resetFilters">Xóa tất cả bộ lọc</button>
      </div>

      <!-- Movie Grid Results -->
      <MovieGrid 
        :movies="filteredMovies" 
        :loading="loading" 
        @play="$emit('play', $event)" 
        @open-detail="$emit('open-detail', $event)" 
      />
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue';
import { useRoute } from 'vue-router';
import { Search, X, SlidersHorizontal, ArrowUpDown, Globe, Film, CalendarDays } from '@lucide/vue';
import MovieGrid from '../components/MovieGrid.vue';
import { fetchMovies, fetchGenres, fetchCountries, fetchMoviesByCountry, fetchYears, fetchMoviesByYear, GENRES_LIST } from '../services/api.js';

defineEmits(['play', 'open-detail']);

const route = useRoute();
const movies = ref([]);
const loading = ref(true);
const genresList = ref(GENRES_LIST);
const countriesList = ref([]);
const yearsList = ref([]);

const searchQuery = ref('');
const selectedGenre = ref('Tất Cả');
const selectedCountry = ref('all');
const selectedYear = ref('all');
const selectedType = ref('all');
const sortBy = ref('rating');

// Watch URL Query Params
watch(() => route.query, (newQuery) => {
  if (newQuery.q) searchQuery.value = newQuery.q;
  if (newQuery.genre) selectedGenre.value = newQuery.genre;
}, { immediate: true });

const getCountryName = (slug) => {
  const found = countriesList.value.find(c => c.slug === slug);
  return found ? found.name : slug;
};

const handleCountryChange = async () => {
  if (selectedCountry.value === 'all') {
    await reloadDefaultMovies();
    return;
  }
  loading.value = true;
  try {
    const res = await fetchMoviesByCountry(selectedCountry.value, {
      limit: 30,
      type: selectedType.value !== 'all' ? selectedType.value : undefined
    });
    if (res && res.items) {
      movies.value = res.items.map(item => ({
        id: item._id || item.id,
        title: item.name || item.title,
        origin_name: item.origin_name || '',
        year: item.year || 2024,
        rating: item.tmdb?.vote_average ? parseFloat(item.tmdb.vote_average) : 8.0,
        genres: [getCountryName(selectedCountry.value), 'Phim Quốc Gia'],
        poster_url: item.poster_url,
        thumb_url: item.thumb_url,
        quality: 'HD',
        slug: item.slug
      }));
    }
  } finally {
    loading.value = false;
  }
};

const handleYearChange = async () => {
  if (selectedYear.value === 'all') {
    await reloadDefaultMovies();
    return;
  }
  selectedCountry.value = 'all';
  loading.value = true;
  try {
    const payload = await fetchMoviesByYear(selectedYear.value, {
      limit: 30,
      page: 1,
      type: selectedType.value !== 'all' ? selectedType.value : undefined,
      status: 'completed'
    });
    const items = payload?.data?.items || payload?.items || [];
    movies.value = items.map(item => ({
      id: item._id || item.id,
      title: item.name || item.title,
      origin_name: item.origin_name || '',
      year: Number(item.year) || Number(selectedYear.value),
      rating: Number(item.tmdb?.vote_average || item.rating || 0),
      genres: Array.isArray(item.category || item.genres)
        ? (item.category || item.genres).map(g => g.name || g)
        : [],
      poster_url: item.poster_url,
      thumb_url: item.thumb_url,
      quality: item.quality || 'HD',
      slug: item.slug
    }));
  } finally {
    loading.value = false;
  }
};

const reloadDefaultMovies = async () => {
  loading.value = true;
  try {
    movies.value = await fetchMovies(1);
  } finally {
    loading.value = false;
  }
};

const resetFilters = async () => {
  searchQuery.value = '';
  selectedGenre.value = 'Tất Cả';
  selectedCountry.value = 'all';
  selectedYear.value = 'all';
  selectedType.value = 'all';
  sortBy.value = 'rating';
  await reloadDefaultMovies();
};

const filteredMovies = computed(() => {
  let result = [...movies.value];

  // 1. Filter by Search Query
  if (searchQuery.value.trim()) {
    const q = searchQuery.value.toLowerCase().trim();
    result = result.filter(m => 
      m.title.toLowerCase().includes(q) || 
      (m.origin_name && m.origin_name.toLowerCase().includes(q))
    );
  }

  // 2. Filter by Genre
  if (selectedGenre.value !== 'Tất Cả') {
    result = result.filter(m => m.genres && m.genres.includes(selectedGenre.value));
  }

  // 3. Sort
  if (sortBy.value === 'rating') {
    result.sort((a, b) => b.rating - a.rating);
  } else if (sortBy.value === 'year') {
    result.sort((a, b) => b.year - a.year);
  } else if (sortBy.value === 'title') {
    result.sort((a, b) => a.title.localeCompare(b.title));
  }

  return result;
});

onMounted(async () => {
  try {
    const [fetchedMoviesList, fetchedGenresList, fetchedCountriesList, fetchedYearsList] = await Promise.all([
      fetchMovies(1),
      fetchGenres(),
      fetchCountries(),
      fetchYears()
    ]);
    movies.value = fetchedMoviesList;
    genresList.value = fetchedGenresList;
    countriesList.value = fetchedCountriesList;
    yearsList.value = fetchedYearsList;
  } finally {
    loading.value = false;
  }
});
</script>

<style scoped>
.catalog-view {
  padding-top: 30px;
  min-height: 80vh;
}

.catalog-container {
  max-width: 1320px;
  margin: 0 auto;
  padding: 0 24px;
}

.catalog-header {
  margin-bottom: 28px;
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

.filter-toolbar {
  padding: 16px 24px;
  display: flex;
  align-items: center;
  gap: 20px;
  flex-wrap: wrap;
  margin-bottom: 24px;
}

.filter-search {
  position: relative;
  display: flex;
  align-items: center;
  background: rgba(255, 255, 255, 0.08);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-sm);
  padding: 8px 14px;
  flex: 1;
  min-width: 250px;
}

.filter-search input {
  background: transparent;
  border: none;
  outline: none;
  color: #fff;
  font-size: 0.9rem;
  font-family: var(--font-main);
  width: 100%;
  margin-left: 8px;
}

.search-icon {
  color: var(--text-secondary);
}

.clear-btn {
  background: none;
  border: none;
  color: var(--text-secondary);
  cursor: pointer;
}

.filter-select-group {
  display: flex;
  align-items: center;
  gap: 10px;
}

.filter-select-group label {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 0.88rem;
  font-weight: 600;
  color: var(--text-secondary);
  white-space: nowrap;
}

.custom-select {
  background: rgba(13, 18, 31, 0.9);
  color: #fff;
  border: 1px solid var(--border-subtle);
  padding: 10px 14px;
  border-radius: var(--radius-sm);
  font-family: var(--font-main);
  font-size: 0.88rem;
  outline: none;
  cursor: pointer;
}

.custom-select option {
  background: #0d121f;
  color: #fff;
}

.active-chips {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
  margin-bottom: 24px;
}

.active-label {
  font-size: 0.85rem;
  color: var(--text-muted);
}

.filter-chip {
  background: rgba(229, 9, 20, 0.15);
  border: 1px solid var(--accent-red);
  color: #fff;
  font-size: 0.82rem;
  font-weight: 600;
  padding: 4px 12px;
  border-radius: 16px;
  display: flex;
  align-items: center;
  gap: 6px;
}

.chip-x {
  cursor: pointer;
}

.reset-all {
  background: none;
  border: none;
  color: var(--accent-cyan);
  font-size: 0.85rem;
  font-weight: 600;
  cursor: pointer;
  text-decoration: underline;
}
</style>
