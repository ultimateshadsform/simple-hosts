import { createApp } from 'vue';
import App from '@/App.vue';
import '@/sass/main.scss';

// If dev then import dev css
if (process.env.NODE_ENV === 'development') {
    import('@/sass/dev/dev.scss');
}

createApp(App).mount('#app');
