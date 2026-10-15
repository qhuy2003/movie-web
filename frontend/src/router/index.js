import { createRouter, createWebHistory } from 'vue-router';
import HomeView from '../views/HomeView.vue';
import CatalogView from '../views/CatalogView.vue';
import WatchlistView from '../views/WatchlistView.vue';
import LoginView from '../views/LoginView.vue';
import RegisterView from '../views/RegisterView.vue';
import AdminView from '../views/AdminView.vue';
import BookingView from '../views/BookingView.vue';
import MyTicketsView from '../views/MyTicketsView.vue';
import MyReviewsView from '../views/MyReviewsView.vue';
import { authState, fetchMe } from '../services/auth.js';

const routes = [
  {
    path: '/',
    name: 'home',
    component: HomeView
  },
  {
    path: '/movies',
    name: 'catalog',
    component: CatalogView
  },
  {
    path: '/watchlist',
    name: 'watchlist',
    component: WatchlistView
  },
  {
    path: '/booking',
    name: 'booking',
    component: BookingView
  },
  {
    path: '/my-tickets',
    name: 'my-tickets',
    component: MyTicketsView,
    meta: { requiresAuth: true }
  },
  {
    path: '/my-reviews',
    name: 'my-reviews',
    component: MyReviewsView,
    meta: { requiresAuth: true }
  },
  {
    path: '/login',
    name: 'login',
    component: LoginView
  },
  {
    path: '/register',
    name: 'register',
    component: RegisterView
  },
  {
    path: '/admin',
    name: 'admin',
    component: AdminView,
    meta: { requiresAdmin: true }
  }
];

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior() {
    return { top: 0 };
  }
});

router.beforeEach(async (to) => {
  if (to.meta.requiresAuth && !authState.isAuthenticated) {
    return { name: 'login', query: { redirect: to.fullPath } };
  }
  if (to.meta.requiresAdmin) {
    if (!authState.isAuthenticated) {
      return { name: 'home' };
    }
    const user = await fetchMe();
    if (!user || user.role !== 'admin') return { name: 'home' };
  }
  return true;
});

export default router;
