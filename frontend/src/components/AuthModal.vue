<template>
  <Teleport to="body">
    <div class="auth-backdrop" @click.self="$emit('close')">
      <div class="auth-card glass-panel">
        <!-- Close Button -->
        <button class="close-btn" @click="$emit('close')">
          <X :size="20" />
        </button>

        <!-- Brand Icon & Header -->
        <div class="auth-header">
          <div class="logo-circle">
            <Film :size="28" class="logo-icon" />
          </div>
          <h2 class="auth-title">{{ isLoginTab ? 'Đăng Nhập CineStream' : 'Tạo Tài Khoản Mới' }}</h2>
          <p class="auth-sub">Đặt vé xem phim, khám phá thể loại, năm phát hành và điểm đánh giá của các bộ phim yêu thích.</p>
        </div>

        <!-- Mode Tabs -->
        <div class="tabs-bar">
          <button 
            class="tab-btn" 
            :class="{ 'active': isLoginTab }" 
            @click="switchTab(true)"
          >
            Đăng Nhập
          </button>
          <button 
            class="tab-btn" 
            :class="{ 'active': !isLoginTab }" 
            @click="switchTab(false)"
          >
            Đăng Ký
          </button>
        </div>

        <!-- Error / Success Alert -->
        <div v-if="errorMessage" class="alert alert-error">
          <AlertCircle :size="18" /> {{ errorMessage }}
        </div>
        <div v-if="successMessage" class="alert alert-success">
          <CheckCircle2 :size="18" /> {{ successMessage }}
        </div>

        <!-- Forms -->
        <form @submit.prevent="handleSubmit" class="auth-form">
          <!-- Username field (Only for Register) -->
          <div v-if="!isLoginTab" class="form-group">
            <label><User :size="16" /> Tên hiển thị</label>
            <input 
              type="text" 
              v-model="username" 
              placeholder="Nhập tên của bạn (vd: AnhTuan99)" 
              required 
            />
          </div>

          <!-- Email field -->
          <div class="form-group">
            <label><Mail :size="16" /> Địa chỉ Email</label>
            <input 
              type="email" 
              v-model="email" 
              placeholder="example@gmail.com" 
              required 
            />
          </div>

          <!-- Password field -->
          <div class="form-group">
            <label><Lock :size="16" /> Mật khẩu</label>
            <div class="input-password">
              <input 
                :type="showPassword ? 'text' : 'password'" 
                v-model="password" 
                placeholder="Nhập mật khẩu (tối thiểu 6 ký tự)" 
                required 
              />
              <button type="button" class="eye-btn" @click="showPassword = !showPassword">
                <Eye v-if="!showPassword" :size="16" />
                <EyeOff v-else :size="16" />
              </button>
            </div>
          </div>

          <!-- Confirm Password (Only for Register) -->
          <div v-if="!isLoginTab" class="form-group">
            <label><Lock :size="16" /> Xác nhận mật khẩu</label>
            <input 
              :type="showPassword ? 'text' : 'password'" 
              v-model="confirmPassword" 
              placeholder="Nhập lại mật khẩu" 
              required 
            />
          </div>

          <!-- Submit Button -->
          <button type="submit" class="btn btn-primary full-width submit-btn" :disabled="loading">
            <span v-if="loading" class="spinner"></span>
            <span v-else>{{ isLoginTab ? 'Đăng Nhập Ngay' : 'Hoàn Tất Đăng Ký' }}</span>
          </button>
        </form>
      </div>
    </div>
  </Teleport>
</template>

<script setup>
import { ref } from 'vue';
import { Film, X, User, Mail, Lock, Eye, EyeOff, AlertCircle, CheckCircle2 } from '@lucide/vue';
import { loginUser, registerUser } from '../services/auth.js';

const props = defineProps({
  initialIsLogin: {
    type: Boolean,
    default: true
  }
});

const emit = defineEmits(['close', 'success']);

const isLoginTab = ref(props.initialIsLogin);
const email = ref('');
const username = ref('');
const password = ref('');
const confirmPassword = ref('');
const showPassword = ref(false);

const loading = ref(false);
const errorMessage = ref('');
const successMessage = ref('');

const switchTab = (toLogin) => {
  isLoginTab.value = toLogin;
  errorMessage.value = '';
  successMessage.value = '';
};

const handleSubmit = async () => {
  errorMessage.value = '';
  successMessage.value = '';

  if (!isLoginTab.value && password.value !== confirmPassword.value) {
    errorMessage.value = 'Mật khẩu xác nhận không khớp!';
    return;
  }

  loading.value = true;

  try {
    if (isLoginTab.value) {
      await loginUser(email.value, password.value);
      successMessage.value = 'Đăng nhập thành công!';
    } else {
      await registerUser(email.value, username.value, password.value);
      successMessage.value = 'Đăng ký tài khoản thành công!';
    }

    setTimeout(() => {
      emit('success');
      emit('close');
    }, 600);
  } catch (err) {
    errorMessage.value = err.message || 'Đã có lỗi xảy ra. Vui lòng thử lại.';
  } finally {
    loading.value = false;
  }
};
</script>

<style scoped>
.auth-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(4, 6, 11, 0.85);
  backdrop-filter: blur(12px);
  z-index: 2500;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px;
  animation: fadeIn 0.25s ease;
}

.auth-card {
  width: 100%;
  max-width: 440px;
  padding: 36px 30px;
  position: relative;
  box-shadow: 0 25px 60px rgba(0, 0, 0, 0.8), 0 0 30px rgba(229, 9, 20, 0.15);
}

.close-btn {
  position: absolute;
  top: 16px;
  right: 16px;
  background: none;
  border: none;
  color: var(--text-secondary);
  cursor: pointer;
  transition: color 0.2s ease;
}
.close-btn:hover {
  color: var(--accent-red);
}

.auth-header {
  text-align: center;
  margin-bottom: 24px;
}

.logo-circle {
  width: 56px;
  height: 56px;
  border-radius: 50%;
  background: rgba(229, 9, 20, 0.12);
  border: 1px solid var(--accent-red);
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 12px auto;
}

.logo-icon {
  color: var(--accent-red);
}

.auth-title {
  font-family: var(--font-heading);
  font-size: 1.5rem;
  font-weight: 800;
  color: #fff;
  margin-bottom: 4px;
}

.auth-sub {
  font-size: 0.85rem;
  color: var(--text-secondary);
}

.tabs-bar {
  display: flex;
  background: rgba(255, 255, 255, 0.06);
  padding: 4px;
  border-radius: var(--radius-sm);
  margin-bottom: 20px;
}

.tab-btn {
  flex: 1;
  padding: 8px;
  border: none;
  background: none;
  color: var(--text-secondary);
  font-weight: 600;
  font-size: 0.9rem;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.25s ease;
}

.tab-btn.active {
  background: var(--bg-modal);
  color: #fff;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.4);
}

.alert {
  padding: 10px 14px;
  border-radius: var(--radius-sm);
  font-size: 0.85rem;
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 16px;
}

.alert-error {
  background: rgba(229, 9, 20, 0.15);
  border: 1px solid var(--accent-red);
  color: #ff6b6b;
}

.alert-success {
  background: rgba(16, 185, 129, 0.15);
  border: 1px solid #10b981;
  color: #34d399;
}

.auth-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.form-group label {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 0.85rem;
  font-weight: 600;
  color: var(--text-secondary);
}

.form-group input {
  background: rgba(13, 18, 31, 0.9);
  border: 1px solid var(--border-subtle);
  padding: 11px 14px;
  border-radius: var(--radius-sm);
  color: #fff;
  font-family: var(--font-main);
  font-size: 0.9rem;
  outline: none;
  transition: border-color 0.2s ease;
}

.form-group input:focus {
  border-color: var(--accent-red);
}

.input-password {
  position: relative;
  display: flex;
  align-items: center;
}

.input-password input {
  width: 100%;
  padding-right: 40px;
}

.eye-btn {
  position: absolute;
  right: 12px;
  background: none;
  border: none;
  color: var(--text-secondary);
  cursor: pointer;
}

.full-width {
  width: 100%;
}

.submit-btn {
  margin-top: 8px;
}

.spinner {
  width: 18px;
  height: 18px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: #fff;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>
