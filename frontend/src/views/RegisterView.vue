<template>
  <div class="auth-page">
    <AuthModal :initial-is-login="false" @close="goHome" @success="goHome" />
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router';
import AuthModal from '../components/AuthModal.vue';

const router = useRouter();
const goHome = () => {
  router.push('/');
};
</script>

<style scoped>
.auth-page {
  min-height: 80vh;
}
</style>
