<template>
  <aside class="hidden md:flex md:flex-col md:w-60 md:shrink-0 md:sticky md:top-0 md:h-screen border-r border-border bg-popover/50 backdrop-blur-md px-4 py-6">
    <RouterLink to="/" class="flex items-center gap-2.5 mb-8 px-2">
      <div class="h-11 w-11 rounded-2xl bg-primary/15 flex items-center justify-center text-2xl shadow-sm">🥄</div>
      <div class="min-w-0">
        <p class="text-xl font-extrabold leading-none text-gradient">Spoony</p>
        <p class="text-[11px] text-muted-foreground mt-1 truncate">{{ t('app.tagline') }}</p>
      </div>
    </RouterLink>

    <nav class="flex flex-col gap-1.5">
      <SideItem to="/"        :icon="HomeIcon" :label="t('nav.home')" />
      <SideItem to="/trainer" :icon="MapIcon"  :label="t('nav.trainer')" />
      <SideItem to="/profile" :icon="UserIcon" :label="t('nav.profile')" />
    </nav>

    <div v-if="userStore.user" class="mt-auto flex items-center gap-3 px-2 pt-4 border-t border-border">
      <Avatar :src="userStore.user.avatar_url" :alt="userStore.user.first_name" size="md" />
      <div class="min-w-0">
        <p class="text-sm font-semibold truncate">{{ userStore.user.first_name }}</p>
        <p class="text-xs text-muted-foreground truncate">@{{ userStore.user.username || 'user' }}</p>
      </div>
    </div>
  </aside>
</template>

<script setup>
import { h } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useUserStore } from '../store/user'
import Avatar from './ui/avatar.vue'
import { Home as HomeIcon, Map as MapIcon, User as UserIcon } from 'lucide-vue-next'

const { t } = useI18n()
const userStore = useUserStore()

function SideItem({ to, icon, label }) {
  const route = useRoute()
  const isActive = to === '/' ? route.path === '/' : route.path.startsWith(to)

  return h(RouterLink, {
    to,
    class: [
      'flex items-center gap-3 px-3 py-2.5 rounded-xl font-semibold text-sm transition-colors',
      isActive ? 'bg-primary/12 text-primary' : 'text-muted-foreground hover:bg-secondary hover:text-foreground',
    ].join(' '),
  }, () => [
    h(icon, { class: 'h-5 w-5 shrink-0' }),
    h('span', label),
  ])
}
SideItem.props = ['to', 'icon', 'label']
</script>
