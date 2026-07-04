<template>
  <div class="min-h-screen bg-background font-sans">
    <div :class="showShell ? 'md:flex md:max-w-6xl md:mx-auto' : ''">
      <SideNav v-if="showShell" />

      <main class="flex-1 min-w-0">
        <RouterView v-slot="{ Component }">
          <Transition name="page" mode="out-in">
            <component :is="Component" />
          </Transition>
        </RouterView>
      </main>
    </div>

    <BottomNav v-if="showShell" />
    <NotificationStack />
    <WordLookupPopover />
  </div>
</template>

<script setup>
import { computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useUserStore } from './store/user'
import { initTelegram } from './services/telegram'
import SideNav from './components/SideNav.vue'
import BottomNav from './components/BottomNav.vue'
import NotificationStack from './components/NotificationStack.vue'
import WordLookupPopover from './components/WordLookupPopover.vue'

const route = useRoute()
const userStore = useUserStore()

const showShell = computed(() =>
  userStore.isAuthenticated && !route.path.startsWith('/admin')
)

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
