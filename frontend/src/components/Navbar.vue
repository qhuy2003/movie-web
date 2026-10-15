<template>
  <header class="navbar" :class="{ 'scrolled': isScrolled }">
    <div class="nav-container">
      <!-- Brand Logo -->
      <router-link to="/" class="brand-logo">
        <Film class="logo-icon" :size="28" />
        <span class="logo-text">Cine<span class="highlight">Stream</span></span>
      </router-link>

      <!-- Navigation Links -->
      <nav class="nav-menu" :class="{ 'open': isMobileMenuOpen }">
        <router-link to="/" class="nav-link" exact-active-class="active" @click="isMobileMenuOpen = false">
          <Home :size="18" /> Trang Chủ
        </router-link>
        <router-link to="/movies" class="nav-link" exact-active-class="active" @click="isMobileMenuOpen = false">
          <Compass :size="18" /> Khám Phá
        </router-link>
        <router-link to="/booking" class="nav-link" exact-active-class="active" @click="isMobileMenuOpen = false">
          <TicketCheck :size="18" /> Đặt Vé
        </router-link>
        <router-link to="/watchlist" class="nav-link" exact-active-class="active" @click="isMobileMenuOpen = false">
          <Heart :size="18" /> Phim Yêu Thích
          <span v-if="watchlistCount > 0" class="badge-count">{{ watchlistCount }}</span>
        </router-link>
      </nav>

      <!-- Quick Search & Auth Buttons -->
      <div class="nav-actions">
        <div class="search-box">
          <Search class="search-icon" :size="18" />
          <input 
            type="text" 
            v-model="searchQuery" 
            placeholder="Tìm kiếm phim..." 
            @keyup.enter="handleSearch"
          />
          <button v-if="searchQuery" class="clear-btn" @click="searchQuery = ''">
            <X :size="14" />
          </button>
        </div>

        <!-- User Authentication Section -->
        <div v-if="authState.isAuthenticated" class="user-menu-wrapper">
          <button class="user-avatar-btn" @click="isUserDropdownOpen = !isUserDropdownOpen">
            <div class="avatar-circle">
              {{ userInitial }}
            </div>
            <span class="user-name">{{ authState.user?.username }}</span>
            <ChevronDown :size="16" />
          </button>

          <!-- User Dropdown Menu -->
          <div v-if="isUserDropdownOpen" class="user-dropdown glass-panel">
            <div class="dropdown-header">
              <p class="dropdown-username">{{ authState.user?.username }}</p>
              <p class="dropdown-email">{{ authState.user?.email }}</p>
            </div>
            <div class="dropdown-divider"></div>
            
            <router-link 
              v-if="authState.user?.role === 'admin'" 
              to="/admin" 
              class="dropdown-item admin-link" 
              @click="isUserDropdownOpen = false"
            >
              <Shield :size="16" /> Trang Quản Trị Admin
            </router-link>

            <router-link to="/watchlist" class="dropdown-item" @click="isUserDropdownOpen = false">
              <Heart :size="16" /> Tủ Phim Yêu Thích
            </router-link>
            <router-link to="/my-tickets" class="dropdown-item" @click="isUserDropdownOpen = false">
              <TicketCheck :size="16" /> Vé của tôi
            </router-link>
            <router-link to="/my-reviews" class="dropdown-item" @click="isUserDropdownOpen = false">
              <Star :size="16" /> Đánh giá của tôi
            </router-link>
            <button class="dropdown-item logout" @click="handleLogout">
              <LogOut :size="16" /> Đăng Xuất
            </button>
          </div>
        </div>

        <button v-else class="btn btn-primary login-btn" @click="showAuthModal = true">
          <User :size="18" /> Đăng Nhập
        </button>

        <!-- Mobile Menu Toggle Button -->
        <button class="mobile-toggle" @click="isMobileMenuOpen = !isMobileMenuOpen">
          <span></span>
          <span></span>
          <span></span>
        </button>
      </div>
    </div>

    <!-- Auth Modal Popup -->
    <AuthModal v-if="showAuthModal" @close="showAuthModal = false" />
  </header>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue';
import { useRouter } from 'vue-router';
import { Film, Search, Heart, Compass, Home, X, User, LogOut, ChevronDown, Shield, TicketCheck, Star } from '@lucide/vue';
import { watchlistStore } from '../services/watchlist.js';
import { authState, fetchMe } from '../services/auth.js';
import AuthModal from './AuthModal.vue';

const router = useRouter();
const searchQuery = ref('');
const isScrolled = ref(false);
const isMobileMenuOpen = ref(false);
const isUserDropdownOpen = ref(false);
const showAuthModal = ref(false);

const watchlistCount = computed(() => watchlistStore.items.length);

const userInitial = computed(() => {
  if (authState.user && authState.user.username) {
    return authState.user.username.charAt(0).toUpperCase();
  }
  return 'U';
});

const handleScroll = () => {
  isScrolled.value = window.scrollY > 40;
};

const handleSearch = () => {
  if (searchQuery.value.trim()) {
    router.push({ path: '/movies', query: { q: searchQuery.value.trim() } });
    isMobileMenuOpen.value = false;
  }
};

const handleLogout = () => {
  authState.logout();
  isUserDropdownOpen.value = false;
};

onMounted(() => {
  window.addEventListener('scroll', handleScroll);
  fetchMe();
});

onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll);
});
</script>

<style scoped>
.navbar {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: 80px;
  z-index: 1000;
  transition: all 0.35s ease;
  background: linear-gradient(180deg, rgba(7, 9, 14, 0.9) 0%, rgba(7, 9, 14, 0) 100%);
  display: flex;
  align-items: center;
}

.navbar.scrolled {
  background: var(--bg-glass);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  height: 70px;
  border-bottom: 1px solid var(--border-subtle);
  box-shadow: 0 4px 30px rgba(0, 0, 0, 0.5);
}

.nav-container {
  max-width: 1320px;
  width: 100%;
  margin: 0 auto;
  padding: 0 24px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 20px;
}

.brand-logo {
  display: flex;
  align-items: center;
  gap: 10px;
  text-decoration: none;
  font-family: var(--font-heading);
  font-size: 1.6rem;
  font-weight: 800;
  color: #fff;
  letter-spacing: -0.5px;
}

.logo-icon {
  color: var(--accent-red);
  filter: drop-shadow(0 0 10px rgba(229, 9, 20, 0.8));
}

.highlight {
  color: var(--accent-red);
}

.nav-menu {
  display: flex;
  align-items: center;
  gap: 28px;
}

.nav-link {
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--text-secondary);
  text-decoration: none;
  font-weight: 600;
  font-size: 0.95rem;
  transition: all 0.25s ease;
  position: relative;
}

.nav-link:hover,
.nav-link.active {
  color: var(--text-primary);
}

.nav-link.active::after {
  content: '';
  position: absolute;
  bottom: -8px;
  left: 0;
  right: 0;
  height: 2px;
  background: var(--accent-gradient);
  border-radius: 2px;
}

.badge-count {
  background: var(--accent-red);
  color: #fff;
  font-size: 0.7rem;
  font-weight: 800;
  padding: 2px 6px;
  border-radius: 10px;
  line-height: 1;
}

.nav-actions {
  display: flex;
  align-items: center;
  gap: 16px;
}

.search-box {
  position: relative;
  display: flex;
  align-items: center;
  background: rgba(255, 255, 255, 0.07);
  border: 1px solid var(--border-subtle);
  border-radius: 24px;
  padding: 6px 14px;
  transition: all 0.3s ease;
  width: 220px;
}

.search-box:focus-within {
  background: rgba(255, 255, 255, 0.12);
  border-color: rgba(255, 255, 255, 0.3);
  width: 260px;
}

.search-icon {
  color: var(--text-secondary);
  margin-right: 8px;
}

.search-box input {
  background: transparent;
  border: none;
  outline: none;
  color: #fff;
  font-size: 0.88rem;
  font-family: var(--font-main);
  width: 100%;
}

.clear-btn {
  background: none;
  border: none;
  color: var(--text-secondary);
  cursor: pointer;
  display: flex;
  align-items: center;
}

.login-btn {
  padding: 8px 18px;
  font-size: 0.88rem;
}

/* User Menu & Avatar */
.user-menu-wrapper {
  position: relative;
}

.user-avatar-btn {
  display: flex;
  align-items: center;
  gap: 10px;
  background: rgba(255, 255, 255, 0.08);
  border: 1px solid var(--border-subtle);
  padding: 4px 12px 4px 4px;
  border-radius: 30px;
  color: #fff;
  cursor: pointer;
  transition: all 0.25s ease;
}

.user-avatar-btn:hover {
  background: rgba(255, 255, 255, 0.15);
  border-color: rgba(255, 255, 255, 0.3);
}

.avatar-circle {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: var(--accent-gradient);
  color: #fff;
  font-weight: 800;
  font-size: 0.9rem;
  display: flex;
  align-items: center;
  justify-content: center;
}

.user-name {
  font-weight: 600;
  font-size: 0.88rem;
  max-width: 100px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.user-dropdown {
  position: absolute;
  top: 48px;
  right: 0;
  width: 220px;
  padding: 12px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.7);
  z-index: 100;
  animation: fadeIn 0.2s ease;
}

.dropdown-header {
  padding: 6px 8px 10px 8px;
}

.dropdown-username {
  font-weight: 700;
  color: #fff;
  font-size: 0.95rem;
}

.dropdown-email {
  font-size: 0.78rem;
  color: var(--text-muted);
}

.dropdown-divider {
  height: 1px;
  background: var(--border-subtle);
  margin: 6px 0;
}

.dropdown-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 12px;
  color: var(--text-secondary);
  text-decoration: none;
  font-size: 0.88rem;
  font-weight: 600;
  border-radius: var(--radius-sm);
  transition: all 0.2s ease;
  background: none;
  border: none;
  width: 100%;
  cursor: pointer;
}

.dropdown-item:hover {
  background: rgba(255, 255, 255, 0.08);
  color: #fff;
}

.dropdown-item.logout:hover {
  background: rgba(229, 9, 20, 0.15);
  color: #ff6b6b;
}

.mobile-toggle {
  display: none;
  flex-direction: column;
  gap: 5px;
  background: none;
  border: none;
  cursor: pointer;
  padding: 4px;
}

.mobile-toggle span {
  width: 22px;
  height: 2px;
  background: #fff;
  border-radius: 2px;
}

@media (max-width: 768px) {
  .mobile-toggle {
    display: flex;
  }
  .nav-menu {
    position: fixed;
    top: 80px;
    left: 0;
    right: 0;
    background: var(--bg-modal);
    flex-direction: column;
    padding: 24px;
    gap: 20px;
    transform: translateY(-150%);
    transition: transform 0.3s ease;
    border-bottom: 1px solid var(--border-subtle);
  }
  .nav-menu.open {
    transform: translateY(0);
  }
  .search-box {
    width: 140px;
  }
}
</style>
