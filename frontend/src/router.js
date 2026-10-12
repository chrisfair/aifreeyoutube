import { createRouter, createWebHistory } from 'vue-router';
import HomeView from './views/HomeView.vue';
import SearchView from './views/SearchView.vue';
import WatchView from './views/WatchView.vue';

const routes = [
  {
    path: '/',
    name: 'home',
    component: HomeView,
  },
  {
    path: '/search',
    name: 'search',
    component: SearchView,
  },
  {
    path: '/watch/:id',
    name: 'watch',
    component: WatchView,
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: '/',
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior() {
    return { top: 0 };
  },
});

export default router;
