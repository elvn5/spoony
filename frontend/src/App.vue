<template>
  <div class="min-h-screen bg-background font-sans">
    <RouterView v-slot="{ Component }">
      <Transition name="page" mode="out-in">
        <component :is="Component" />
      </Transition>
    </RouterView>

    <BottomNav v-if="userStore.isAuthenticated && !route.path.startsWith('/admin')" />
    <NotificationStack />
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useUserStore } from './store/user'
import { initTelegram } from './services/telegram'
import BottomNav from './components/BottomNav.vue'
import NotificationStack from './components/NotificationStack.vue'

const route = useRoute()
const userStore = useUserStore()

onMounted(() => {
  initTelegram()
  if (userStore.isAuthenticated) {
    userStore.fetchMe()
  }
})
</script>

<style>
.page-enter-active,
.page-leave-active { transition: opacity 0.15s ease; }
.page-enter-from,
.page-leave-to     { opacity: 0; }
</style>
