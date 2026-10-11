<template>
  <footer class="footer">
    <div class="footer-container">
      <div class="footer-grid">
        <!-- Col 1: Brand Info -->
        <div class="footer-brand">
          <div class="brand-title">
            <Film class="logo-icon" :size="24" />
            <span>Cine<span class="highlight">Stream</span></span>
          </div>
          <p class="brand-desc">
            Khám phá thông tin, thể loại, năm phát hành và điểm đánh giá phim; chọn suất chiếu và đặt vé trực tuyến dễ dàng.
          </p>
        </div>

        <!-- Col 2: Categories -->
        <div class="footer-links">
          <h4>Thể Loại Phim</h4>
          <ul>
            <li><router-link to="/movies">Phim Hành Động</router-link></li>
            <li><router-link to="/movies">Phim Viễn Tưởng</router-link></li>
            <li><router-link to="/movies">Phim Hoạt Hình / Anime</router-link></li>
            <li><router-link to="/movies">Phim Chính Kịch</router-link></li>
          </ul>
        </div>

        <!-- Col 3: Support & Info -->
        <div class="footer-links">
          <h4>Thông Tin</h4>
          <ul>
            <li><a href="#">Giới Thiệu</a></li>
            <li><a href="#">Điều Khoản Sử Dụng</a></li>
            <li><a href="#">Chính Sách Báo Mật</a></li>
            <li><a href="#">Liên Hệ & Hỗ Trợ</a></li>
          </ul>
        </div>
      </div>

      <div class="footer-bottom">
        <p>© 2026 CineStream Vue. Kết Nối Go Backend API & Vue 3 Front-End.</p>
      </div>
    </div>
  </footer>
</template>

<script setup>
import { Film } from '@lucide/vue';
</script>

<style scoped>
.footer {
  background: #04060a;
  border-top: 1px solid var(--border-subtle);
  padding: 60px 0 30px 0;
  margin-top: 60px;
}

.footer-container {
  max-width: 1320px;
  margin: 0 auto;
  padding: 0 24px;
}

.footer-grid {
  display: grid;
  grid-template-columns: 2fr 1fr 1fr;
  gap: 40px;
  margin-bottom: 40px;
}

.brand-title {
  display: flex;
  align-items: center;
  gap: 10px;
  font-family: var(--font-heading);
  font-size: 1.5rem;
  font-weight: 800;
  color: #fff;
  margin-bottom: 12px;
}

.logo-icon {
  color: var(--accent-red);
}

.highlight {
  color: var(--accent-red);
}

.brand-desc {
  color: var(--text-secondary);
  font-size: 0.9rem;
  line-height: 1.6;
  max-width: 380px;
}

.footer-links h4 {
  font-family: var(--font-heading);
  font-size: 1.1rem;
  font-weight: 700;
  color: #fff;
  margin-bottom: 16px;
}

.footer-links ul {
  list-style: none;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.footer-links a {
  color: var(--text-secondary);
  text-decoration: none;
  font-size: 0.9rem;
  transition: color 0.2s ease;
}

.footer-links a:hover {
  color: var(--accent-red);
}

.footer-bottom {
  padding-top: 24px;
  border-top: 1px solid var(--border-subtle);
  text-align: center;
  color: var(--text-muted);
  font-size: 0.85rem;
}

@media (max-width: 768px) {
  .footer-grid {
    grid-template-columns: 1fr;
    gap: 30px;
  }
}
</style>
