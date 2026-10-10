import { createApp } from 'vue';
import './style.css';
import App from './App.vue';
import router from './router';
import { reveal } from './directives/reveal.js';

const app = createApp(App);
app.use(router);
app.directive('reveal', reveal);
app.mount('#app');
