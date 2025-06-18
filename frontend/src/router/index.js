import { createRouter, createWebHistory } from 'vue-router';
import Chat from '@/components/Chat.vue';
import MemoryTable from '@/components/MemoryTable.vue';
import Login from '@/views/Login.vue';
import Register from '@/views/Register.vue';

const routes = [
  { path: '/login', component: Login },
  { path: '/register', component: Register },
  { path: '/', component: Chat },
  { path: '/memory', component: MemoryTable }
];

const router = createRouter({
  history: createWebHistory(),
  routes
});

export default router;