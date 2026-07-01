import { createRouter, createWebHistory } from 'vue-router'
import { storage } from '../services/storage'
import { useAdminStore } from '../store/admin'

const routes = [
  { path: '/',            name: 'Home',     component: () => import('../views/HomeView.vue') },
  { path: '/trainer',     name: 'Trainer',  component: () => import('../views/TrainerView.vue'), meta: { requiresAuth: true } },
  { path: '/trainer/:id', name: 'Game',     component: () => import('../views/GameView.vue'),    meta: { requiresAuth: true }, props: true },
  { path: '/profile',     name: 'Profile',  component: () => import('../views/ProfileView.vue'),  meta: { requiresAuth: true } },
  { path: '/settings',    name: 'Settings', component: () => import('../views/SettingsView.vue'), meta: { requiresAuth: true } },

  {
    path: '/admin',
    component: () => import('../views/admin/AdminLayout.vue'),
    children: [
      { path: '',          redirect: '/admin/login' },
      { path: 'login',     name: 'AdminLogin',     component: () => import('../views/admin/AdminLoginView.vue') },
      { path: 'dashboard', name: 'AdminDashboard', component: () => import('../views/admin/AdminDashboardView.vue'), meta: { requiresAdmin: true } },
      { path: 'users',     name: 'AdminUsers',     component: () => import('../views/admin/AdminUsersView.vue'),     meta: { requiresAdmin: true } },
    ],
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior: () => ({ top: 0 }),
})

router.beforeEach((to) => {
  if (to.meta.requiresAuth && !storage.get('token')) {
    return { name: 'Home' }
  }
  if (to.meta.requiresAdmin) {
    const adminStore = useAdminStore()
    if (!adminStore.isAuthed) return { name: 'AdminLogin' }
  }
})

export default router
