<template>
  <div class="admin-view">
    <div class="admin-container">
      <!-- Title & Header -->
      <div class="admin-header">
        <div>
          <h1 class="page-title">⚙️ Trang Quản Trị Hệ Thống Admin</h1>
          <p class="page-sub">Quản lý cơ sở dữ liệu phim, danh mục thể loại và phân quyền tài khoản người dùng.</p>
        </div>
      </div>

      <!-- Stats Cards Row -->
      <div class="stats-grid">
        <div class="stat-card glass-panel">
          <div class="stat-icon red"><Film :size="24" /></div>
          <div class="stat-info">
            <span class="stat-num">{{ movies.length }}</span>
            <span class="stat-label">Tổng Số Phim</span>
          </div>
        </div>

        <div class="stat-card glass-panel">
          <div class="stat-icon cyan"><Tag :size="24" /></div>
          <div class="stat-info">
            <span class="stat-num">{{ genres.length }}</span>
            <span class="stat-label">Thể Loại Phim</span>
          </div>
        </div>

        <div class="stat-card glass-panel">
          <div class="stat-icon gold"><Users :size="24" /></div>
          <div class="stat-info">
            <span class="stat-num">{{ users.length }}</span>
            <span class="stat-label">Tài Khoản User</span>
          </div>
        </div>
      </div>

      <!-- Admin Tabs Navigation -->
      <div class="admin-tabs">
        <button 
          class="tab-btn" 
          :class="{ active: activeTab === 'movies' }" 
          @click="activeTab = 'movies'"
        >
          <Film :size="18" /> Quản Lý Phim
        </button>

        <button 
          class="tab-btn" 
          :class="{ active: activeTab === 'genres' }" 
          @click="activeTab = 'genres'"
        >
          <Tag :size="18" /> Quản Lý Thể Loại
        </button>

        <button 
          class="tab-btn" 
          :class="{ active: activeTab === 'users' }" 
          @click="activeTab = 'users'"
        >
          <Users :size="18" /> Quản Lý Người Dùng
        </button>
      </div>

      <!-- Alert Notification -->
      <div v-if="alertMessage" class="alert-box" :class="alertType">
        <span>{{ alertMessage }}</span>
        <button class="alert-close" @click="alertMessage = ''"><X :size="16" /></button>
      </div>

      <!-- TAB 1: MOVIES MANAGEMENT -->
      <div v-if="activeTab === 'movies'" class="tab-content glass-panel">
        <div class="tab-toolbar">
          <h2 class="section-title">Danh Sách Phim TRONG POSTGRESQL</h2>
          <button class="btn btn-primary" @click="openMovieModal()">
            <Plus :size="18" /> Thêm Phim Mới
          </button>
        </div>

        <div class="table-responsive">
          <table class="admin-table">
            <thead>
              <tr>
                <th>ID</th>
                <th>Hình Ảnh</th>
                <th>Tên Phim</th>
                <th>Năm</th>
                <th>Đánh Giá</th>
                <th>Chất Lượng</th>
                <th>Thể Loại</th>
                <th>Thao Tác</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="m in movies" :key="m.id">
                <td>#{{ m.id }}</td>
                <td>
                  <img :src="m.thumb_url || m.poster_url" class="thumb-img" alt="poster" />
                </td>
                <td>
                  <div class="font-bold">{{ m.title }}</div>
                  <div class="text-muted text-sm">{{ m.origin_name }}</div>
                </td>
                <td>{{ m.year || '2024' }}</td>
                <td>⭐ {{ m.rating || 8.0 }}</td>
                <td><span class="badge badge-quality">{{ m.quality || 'HD' }}</span></td>
                <td>{{ m.genres }}</td>
                <td>
                  <div class="action-btns">
                    <button class="btn-icon-sm edit" @click="openMovieModal(m)" title="Chỉnh Sửa">
                      <Edit3 :size="16" />
                    </button>
                    <button class="btn-icon-sm delete" @click="handleDeleteMovie(m.id)" title="Xóa Phim">
                      <Trash2 :size="16" />
                    </button>
                  </div>
                </td>
              </tr>
              <tr v-if="!movies.length">
                <td colspan="8" class="empty-cell">Chưa có phim nào được tạo trong cơ sở dữ liệu. Nút Thêm Phim Mới để khởi tạo.</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- TAB 2: GENRES MANAGEMENT -->
      <div v-if="activeTab === 'genres'" class="tab-content glass-panel">
        <div class="tab-toolbar">
          <h2 class="section-title">Danh Sách Thể Loại Phim</h2>
          <button class="btn btn-primary" @click="openGenreModal()">
            <Plus :size="18" /> Thêm Thể Loại
          </button>
        </div>

        <div class="table-responsive">
          <table class="admin-table">
            <thead>
              <tr>
                <th>ID</th>
                <th>Tên Thể Loại</th>
                <th>Slug Đường Dẫn</th>
                <th>Thao Tác</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="g in genres" :key="g.id">
                <td>#{{ g.id }}</td>
                <td class="font-bold">{{ g.name }}</td>
                <td><code>{{ g.slug }}</code></td>
                <td>
                  <div class="action-btns">
                    <button class="btn-icon-sm edit" @click="openGenreModal(g)" title="Chỉnh Sửa">
                      <Edit3 :size="16" />
                    </button>
                    <button class="btn-icon-sm delete" @click="handleDeleteGenre(g.id)" title="Xóa">
                      <Trash2 :size="16" />
                    </button>
                  </div>
                </td>
              </tr>
              <tr v-if="!genres.length">
                <td colspan="4" class="empty-cell">Chưa có thể loại tùy chỉnh nào trong PostgreSQL.</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- TAB 3: USERS MANAGEMENT -->
      <div v-if="activeTab === 'users'" class="tab-content glass-panel">
        <div class="tab-toolbar">
          <h2 class="section-title">Danh Sách Tài Khoản Người Dùng</h2>
        </div>

        <div class="table-responsive">
          <table class="admin-table">
            <thead>
              <tr>
                <th>ID</th>
                <th>Email</th>
                <th>Tên Hiển Thị</th>
                <th>Quyền Hạn (Role)</th>
                <th>Ngày Đăng Ký</th>
                <th>Thao Tác</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="u in users" :key="u.id">
                <td>#{{ u.id }}</td>
                <td class="font-bold">{{ u.email }}</td>
                <td>{{ u.username }}</td>
                <td>
                  <span class="role-badge" :class="u.role === 'admin' ? 'role-admin' : 'role-user'">
                    {{ u.role === 'admin' ? '👑 ADMIN' : '👤 USER' }}
                  </span>
                </td>
                <td>{{ formatDate(u.created_at) }}</td>
                <td>
                  <div class="action-btns">
                    <button 
                      class="btn-sm" 
                      :class="u.role === 'admin' ? 'btn-demote' : 'btn-promote'"
                      @click="handleToggleRole(u)"
                    >
                      {{ u.role === 'admin' ? 'Hạ Quyền User' : 'Nâng Quyền Admin' }}
                    </button>
                    <button class="btn-icon-sm delete" @click="handleDeleteUser(u.id)" title="Xóa Tài Khoản">
                      <Trash2 :size="16" />
                    </button>
                  </div>
                </td>
              </tr>
              <tr v-if="!users.length">
                <td colspan="6" class="empty-cell">Chưa có dữ liệu người dùng.</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- MOVIE FORM MODAL -->
      <div v-if="showMovieModal" class="modal-backdrop" @click.self="showMovieModal = false">
        <div class="modal-card glass-panel">
          <div class="modal-header">
            <h3>{{ editingMovieId ? 'Chỉnh Sửa Phim' : 'Thêm Phim Mới Vào PostgreSQL' }}</h3>
            <button class="alert-close" @click="showMovieModal = false"><X :size="20" /></button>
          </div>

          <form @submit.prevent="saveMovie" class="modal-form">
            <div class="form-grid">
              <div class="form-group">
                <label>Tên Phim (Tiếng Việt) *</label>
                <input type="text" v-model="movieForm.title" required placeholder="Ví dụ: Dune: Phần 2" />
              </div>

              <div class="form-group">
                <label>Tên Gốc (Origin Name)</label>
                <input type="text" v-model="movieForm.origin_name" placeholder="Ví dụ: Dune: Part Two" />
              </div>

              <div class="form-group">
                <label>Năm Phát Hành</label>
                <input type="number" v-model.number="movieForm.year" placeholder="2024" />
              </div>

              <div class="form-group">
                <label>Điểm TMDB (0 - 10)</label>
                <input type="number" step="0.1" v-model.number="movieForm.rating" placeholder="8.5" />
              </div>

              <div class="form-group">
                <label>Chất Lượng</label>
                <input type="text" v-model="movieForm.quality" placeholder="4K UHD, Full HD..." />
              </div>

              <div class="form-group">
                <label>Thời Lượng</label>
                <input type="text" v-model="movieForm.duration" placeholder="120 phút" />
              </div>

              <div class="form-group">
                <label>Đạo Diễn</label>
                <input type="text" v-model="movieForm.director" placeholder="Christopher Nolan" />
              </div>

              <div class="form-group">
                <label>Thể Loại (Phân cách bằng dấu phẩy)</label>
                <input type="text" v-model="movieForm.genres" placeholder="Hành Động, Viễn Tưởng, Phiêu Lưu" />
              </div>
            </div>

            <div class="form-group margin-top">
              <label>Diễn Viên Chính</label>
              <input type="text" v-model="movieForm.cast" placeholder="Timothée Chalamet, Zendaya..." />
            </div>

            <div class="form-group">
              <label>Link Poster Ảnh Ngang (Poster URL)</label>
              <input type="url" v-model="movieForm.poster_url" placeholder="https://..." />
            </div>

            <div class="form-group">
              <label>Link Ảnh Dọc (Thumb URL)</label>
              <input type="url" v-model="movieForm.thumb_url" placeholder="https://..." />
            </div>

            <div class="form-group">
              <label>Link Stream Video (.mp4 / .m3u8)</label>
              <input type="url" v-model="movieForm.video_url" placeholder="https://commondatastorage.googleapis.com/..." />
            </div>

            <div class="form-group">
              <label>Mô Tả Nội Dung Phim</label>
              <textarea v-model="movieForm.description" rows="3" placeholder="Tóm tắt kịch bản phim..."></textarea>
            </div>

            <div class="modal-actions">
              <button type="button" class="btn btn-secondary" @click="showMovieModal = false">Hủy</button>
              <button type="submit" class="btn btn-primary">Lưu Thông Tin Phim</button>
            </div>
          </form>
        </div>
      </div>

      <!-- GENRE FORM MODAL -->
      <div v-if="showGenreModal" class="modal-backdrop" @click.self="showGenreModal = false">
        <div class="modal-card glass-panel small">
          <div class="modal-header">
            <h3>{{ editingGenreId ? 'Chỉnh Sửa Thể Loại' : 'Thêm Thể Loại Mới' }}</h3>
            <button class="alert-close" @click="showGenreModal = false"><X :size="20" /></button>
          </div>

          <form @submit.prevent="saveGenre" class="modal-form">
            <div class="form-group">
              <label>Tên Thể Loại *</label>
              <input type="text" v-model="genreForm.name" required placeholder="Tên thể loại (Ví dụ: Tâm Lý Học)" />
            </div>

            <div class="form-group">
              <label>Slug Đường Dẫn (Để trống để tự tạo)</label>
              <input type="text" v-model="genreForm.slug" placeholder="tam-ly-hoc" />
            </div>

            <div class="modal-actions">
              <button type="button" class="btn btn-secondary" @click="showGenreModal = false">Hủy</button>
              <button type="submit" class="btn btn-primary">Lưu Thể Loại</button>
            </div>
          </form>
        </div>
      </div>

    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { Film, Tag, Users, Plus, Edit3, Trash2, X } from '@lucide/vue';
import {
  fetchAdminUsers, updateUserRole, deleteUser,
  fetchAdminGenres, createGenre, updateGenre, deleteGenre, importGenresAdmin,
  fetchAdminMovies, createMovieAdmin, updateMovieAdmin, deleteMovieAdmin, importMoviesAdmin
} from '../services/admin.js';

const activeTab = ref('movies');

const movies = ref([]);
const genres = ref([]);
const users = ref([]);

const alertMessage = ref('');
const alertType = ref('success');
const importingMovies = ref(false);

// Modals State
const showMovieModal = ref(false);
const editingMovieId = ref(null);
const movieForm = ref({
  title: '', origin_name: '', year: 2024, rating: 8.5,
  quality: '4K UHD', duration: '120 phút', director: '', cast: '',
  genres: '', poster_url: '', thumb_url: '', description: '', video_url: ''
});

const showGenreModal = ref(false);
const editingGenreId = ref(null);
const genreForm = ref({ name: '', slug: '' });

const showAlert = (msg, type = 'success') => {
  alertMessage.value = msg;
  alertType.value = type;
  setTimeout(() => {
    if (alertMessage.value === msg) alertMessage.value = '';
  }, 4000);
};

const formatDate = (dateStr) => {
  if (!dateStr) return 'N/A';
  return new Date(dateStr).toLocaleDateString('vi-VN');
};

// Load All Admin Data
const loadData = async () => {
  try {
    const [uList, gList, mList] = await Promise.all([
      fetchAdminUsers(),
      fetchAdminGenres(),
      fetchAdminMovies()
    ]);
    users.value = uList || [];
    genres.value = gList || [];
    movies.value = mList || [];
  } catch (err) {
    showAlert(err.message || 'Không thể tải dữ liệu admin', 'error');
  }
};

onMounted(async () => {
  importingMovies.value = true;
  try {
    const [movieResult, genreResult] = await Promise.all([
      importMoviesAdmin(),
      importGenresAdmin()
    ]);
    await loadData();
    const messages = [movieResult, genreResult]
      .filter((result) => result.count > 0)
      .map((result) => result.message);
    if (messages.length) showAlert(messages.join(' '));
  } catch (err) {
    await loadData();
    showAlert(err.message || 'Không thể nhập phim từ API', 'error');
  } finally {
    importingMovies.value = false;
  }
});

// ================= MOVIES ACTIONS =================
const openMovieModal = (movie = null) => {
  if (movie) {
    editingMovieId.value = movie.id;
    movieForm.value = { ...movie };
  } else {
    editingMovieId.value = null;
    movieForm.value = {
      title: '', origin_name: '', year: 2024, rating: 8.5,
      quality: '4K UHD', duration: '120 phút', director: '', cast: '',
      genres: '', poster_url: '', thumb_url: '', description: '', video_url: ''
    };
  }
  showMovieModal.value = true;
};

const saveMovie = async () => {
  try {
    if (editingMovieId.value) {
      await updateMovieAdmin(editingMovieId.value, movieForm.value);
      showAlert('Đã cập nhật phim thành công!');
    } else {
      await createMovieAdmin(movieForm.value);
      showAlert('Đã thêm phim mới vào PostgreSQL!');
    }
    showMovieModal.value = false;
    await loadData();
  } catch (err) {
    showAlert(err.message, 'error');
  }
};

const handleDeleteMovie = async (id) => {
  if (!confirm('Bạn có chắc chắn muốn xóa phim này khỏi PostgreSQL?')) return;
  try {
    await deleteMovieAdmin(id);
    showAlert('Đã xóa phim thành công');
    await loadData();
  } catch (err) {
    showAlert(err.message, 'error');
  }
};

// ================= GENRES ACTIONS =================
const openGenreModal = (genre = null) => {
  if (genre) {
    editingGenreId.value = genre.id;
    genreForm.value = { name: genre.name, slug: genre.slug };
  } else {
    editingGenreId.value = null;
    genreForm.value = { name: '', slug: '' };
  }
  showGenreModal.value = true;
};

const saveGenre = async () => {
  try {
    if (editingGenreId.value) {
      await updateGenre(editingGenreId.value, genreForm.value.name, genreForm.value.slug);
      showAlert('Đã cập nhật thể loại!');
    } else {
      await createGenre(genreForm.value.name, genreForm.value.slug);
      showAlert('Đã thêm thể loại mới!');
    }
    showGenreModal.value = false;
    await loadData();
  } catch (err) {
    showAlert(err.message, 'error');
  }
};

const handleDeleteGenre = async (id) => {
  if (!confirm('Bạn có chắc chắn muốn xóa thể loại này?')) return;
  try {
    await deleteGenre(id);
    showAlert('Đã xóa thể loại');
    await loadData();
  } catch (err) {
    showAlert(err.message, 'error');
  }
};

// ================= USERS ACTIONS =================
const handleToggleRole = async (user) => {
  const newRole = user.role === 'admin' ? 'user' : 'admin';
  try {
    await updateUserRole(user.id, newRole);
    showAlert(`Đã đổi quyền của ${user.username} thành ${newRole.toUpperCase()}`);
    await loadData();
  } catch (err) {
    showAlert(err.message, 'error');
  }
};

const handleDeleteUser = async (id) => {
  if (!confirm('Bạn có chắc chắn muốn xóa tài khoản người dùng này?')) return;
  try {
    await deleteUser(id);
    showAlert('Đã xóa người dùng thành công');
    await loadData();
  } catch (err) {
    showAlert(err.message, 'error');
  }
};
</script>

<style scoped>
.admin-view {
  padding: 40px 24px 80px;
  max-width: 1280px;
  margin: 0 auto;
}

.admin-header {
  margin-bottom: 32px;
}

.page-title {
  font-family: var(--font-heading);
  font-size: 2.2rem;
  font-weight: 800;
  color: #fff;
}

.page-sub {
  color: var(--text-secondary);
  font-size: 1rem;
}

/* Stats Cards */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
  gap: 20px;
  margin-bottom: 36px;
}

.stat-card {
  padding: 24px;
  display: flex;
  align-items: center;
  gap: 20px;
}

.stat-icon {
  width: 54px;
  height: 54px;
  border-radius: var(--radius-sm);
  display: flex;
  align-items: center;
  justify-content: center;
}

.stat-icon.red { background: rgba(229, 9, 20, 0.15); color: var(--accent-red); }
.stat-icon.cyan { background: rgba(0, 242, 254, 0.15); color: var(--accent-cyan); }
.stat-icon.gold { background: rgba(255, 184, 0, 0.15); color: var(--accent-gold); }

.stat-num {
  font-family: var(--font-heading);
  font-size: 2rem;
  font-weight: 800;
  color: #fff;
  display: block;
}

.stat-label {
  color: var(--text-muted);
  font-size: 0.88rem;
}

/* Tabs */
.admin-tabs {
  display: flex;
  gap: 12px;
  margin-bottom: 24px;
  border-bottom: 1px solid var(--border-subtle);
  padding-bottom: 16px;
}

.tab-btn {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 12px 24px;
  border-radius: var(--radius-sm);
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid var(--border-subtle);
  color: var(--text-secondary);
  font-weight: 600;
  cursor: pointer;
  transition: all 0.25s ease;
}

.tab-btn:hover,
.tab-btn.active {
  background: var(--accent-gradient);
  border-color: transparent;
  color: #fff;
}

/* Tab Content Box */
.tab-content {
  padding: 28px;
}

.tab-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 24px;
}

/* Table Styles */
.table-responsive {
  width: 100%;
  overflow-x: auto;
}

.admin-table {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
}

.admin-table th,
.admin-table td {
  padding: 14px 16px;
  border-bottom: 1px solid var(--border-subtle);
  font-size: 0.92rem;
}

.admin-table th {
  color: var(--text-muted);
  font-weight: 700;
  text-transform: uppercase;
  font-size: 0.78rem;
  letter-spacing: 0.5px;
}

.thumb-img {
  width: 44px;
  height: 60px;
  object-fit: cover;
  border-radius: 6px;
  border: 1px solid var(--border-subtle);
}

.action-btns {
  display: flex;
  align-items: center;
  gap: 8px;
}

.btn-icon-sm {
  width: 32px;
  height: 32px;
  border-radius: 6px;
  border: none;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: transform 0.2s ease;
}

.btn-icon-sm.edit { background: rgba(0, 242, 254, 0.15); color: var(--accent-cyan); }
.btn-icon-sm.delete { background: rgba(229, 9, 20, 0.15); color: var(--accent-red); }

.btn-icon-sm:hover { transform: scale(1.15); }

.role-badge {
  padding: 4px 10px;
  border-radius: 6px;
  font-size: 0.78rem;
  font-weight: 700;
}

.role-admin { background: rgba(255, 184, 0, 0.2); color: var(--accent-gold); border: 1px solid rgba(255, 184, 0, 0.4); }
.role-user { background: rgba(255, 255, 255, 0.08); color: var(--text-secondary); }

.btn-sm {
  padding: 6px 12px;
  border-radius: 6px;
  border: none;
  font-weight: 600;
  font-size: 0.8rem;
  cursor: pointer;
}
.btn-promote { background: rgba(0, 242, 254, 0.2); color: var(--accent-cyan); }
.btn-demote { background: rgba(255, 184, 0, 0.2); color: var(--accent-gold); }

.empty-cell {
  text-align: center;
  padding: 40px;
  color: var(--text-muted);
}

/* Modals */
.modal-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(4, 6, 11, 0.85);
  backdrop-filter: blur(10px);
  z-index: 2000;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px;
}

.modal-card {
  width: 100%;
  max-width: 760px;
  max-height: 90vh;
  overflow-y: auto;
  padding: 32px;
  background: var(--bg-modal);
  border-radius: var(--radius-lg);
  border: 1px solid var(--border-subtle);
  position: relative;
  box-shadow: 0 25px 60px rgba(0, 0, 0, 0.8);
}

.modal-card.small {
  max-width: 480px;
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 24px;
}

.modal-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.form-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.form-group label {
  font-size: 0.85rem;
  color: var(--text-secondary);
  font-weight: 600;
}

.form-group input,
.form-group textarea {
  padding: 10px 14px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid var(--border-subtle);
  border-radius: 8px;
  color: #fff;
  font-family: inherit;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 24px;
}

/* Alert Box */
.alert-box {
  padding: 14px 20px;
  border-radius: 8px;
  margin-bottom: 20px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-weight: 600;
}

.alert-box.success { background: rgba(0, 242, 254, 0.15); border: 1px solid var(--accent-cyan); color: #fff; }
.alert-box.error { background: rgba(229, 9, 20, 0.15); border: 1px solid var(--accent-red); color: #fff; }

.alert-close {
  background: none;
  border: none;
  color: #fff;
  cursor: pointer;
}
</style>
