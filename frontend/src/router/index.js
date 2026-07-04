import { createRouter, createWebHistory } from 'vue-router'
import { storage } from '../services/storage'
import { useAdminStore } from '../features/admin/store'
import { isAlphabetCompleted } from '../features/alphabet/progress'
import { i18n } from '../i18n'
import { showAlert } from '../services/telegram'

const routes = [
  { path: '/',            name: 'Home',     component: () => import('../features/news/HomeView.vue') },
  { path: '/trainer',     name: 'Trainer',  component: () => import('../features/trainer/TrainerView.vue'), meta: { requiresAuth: true, requiresAlphabet: true } },
  { path: '/trainer/:id', name: 'Game',     component: () => import('../features/trainer/GameView.vue'),    meta: { requiresAuth: true, requiresAlphabet: true }, props: true },
  { path: '/alphabet',     name: 'Alphabet',     component: () => import('../features/alphabet/AlphabetView.vue'),     meta: { requiresAuth: true } },
  { path: '/alphabet/:level', name: 'AlphabetGame', component: () => import('../features/alphabet/AlphabetGameView.vue'), meta: { requiresAuth: true }, props: true },
  { path: '/alphabet/match/play', name: 'AlphabetMatch', component: () => import('../features/alphabet/AlphabetMatchView.vue'), meta: { requiresAuth: true } },
  { path: '/alphabet/words/play', name: 'AlphabetWords', component: () => import('../features/alphabet/AlphabetWordsView.vue'), meta: { requiresAuth: true } },
  { path: '/alphabet/combos/:group', name: 'AlphabetCombos', component: () => import('../features/alphabet/AlphabetCombosView.vue'), meta: { requiresAuth: true }, props: true },
  { path: '/profile',     name: 'Profile',  component: () => import('../features/profile/ProfileView.vue'),  meta: { requiresAuth: true } },
  { path: '/settings',    name: 'Settings', component: () => import('../features/settings/SettingsView.vue'), meta: { requiresAuth: true } },

  {
    path: '/admin',
    component: () => import('../features/admin/AdminLayout.vue'),
    children: [
      { path: '',          redirect: '/admin/login' },
      { path: 'login',     name: 'AdminLogin',     component: () => import('../features/admin/AdminLoginView.vue') },
      { path: 'dashboard', name: 'AdminDashboard', component: () => import('../features/admin/AdminDashboardView.vue'), meta: { requiresAdmin: true } },
      { path: 'users',     name: 'AdminUsers',     component: () => import('../features/admin/AdminUsersView.vue'),     meta: { requiresAdmin: true } },
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
  if (to.meta.requiresAlphabet && !isAlphabetCompleted()) {
    try { showAlert(i18n.global.t('trainer.lockedByAlphabet')) } catch {}
    return { name: 'Alphabet' }
  }
  if (to.meta.requiresAdmin) {
    const adminStore = useAdminStore()
    if (!adminStore.isAuthed) return { name: 'AdminLogin' }
  }
})

export default router
