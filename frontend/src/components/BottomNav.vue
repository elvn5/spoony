<template>
  <nav class="fixed bottom-0 left-0 right-0 z-40 safe-bottom">
    <div class="border-t border-border bg-popover/95 backdrop-blur-md">
      <div class="flex justify-around items-center h-16 max-w-md mx-auto px-2">
        <NavItem to="/"         :icon="HomeIcon"     label="Home"     />
        <NavItem to="/profile"  :icon="UserIcon"     label="Profile"  />
        <NavItem to="/settings" :icon="SettingsIcon" label="Settings" />
      </div>
    </div>
  </nav>
</template>

<script setup>
import { h } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import { Home as HomeIcon, User as UserIcon, Settings as SettingsIcon } from 'lucide-vue-next'

function NavItem({ to, icon, label }) {
  const route = useRoute()
  const isActive = to === '/' ? route.path === '/' : route.path.startsWith(to)

  return h(RouterLink, {
    to,
    class: [
      'flex flex-col items-center gap-1 px-4 py-1.5 rounded-xl transition-colors',
      isActive ? 'text-primary' : 'text-muted-foreground hover:text-foreground',
    ].join(' '),
  }, () => [
    h(icon, { class: 'h-5 w-5' }),
    h('span', { class: 'text-[10px]' }, label),
  ])
}

NavItem.props = ['to', 'icon', 'label']
</script>
