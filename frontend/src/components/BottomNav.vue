<template>
  <nav class="fixed bottom-0 left-0 right-0 z-40 safe-bottom md:hidden">
    <div class="border-t border-border bg-popover/95 backdrop-blur-md">
      <div class="flex justify-around items-center h-16 max-w-md mx-auto px-2">
        <NavItem to="/"         :icon="HomeIcon"   :label="t('nav.home')" />
        <NavItem to="/alphabet" :icon="PuzzleIcon" :label="t('nav.alphabet')" />
        <NavItem to="/trainer"  :icon="MapIcon"    :label="t('nav.trainer')" :locked="trainerLocked" />
        <NavItem to="/profile"  :icon="UserIcon"   :label="t('nav.profile')" />
      </div>
    </div>
  </nav>
</template>

<script setup>
import { h, computed } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { isAlphabetCompleted } from '../features/alphabet/progress'
import { hapticFeedback, showAlert } from '../services/telegram'
import { Home as HomeIcon, Map as MapIcon, Puzzle as PuzzleIcon, User as UserIcon, Lock as LockIcon } from 'lucide-vue-next'

const { t } = useI18n()
const route = useRoute()
const trainerLocked = computed(() => {
  route.fullPath // recompute on every navigation
  return !isAlphabetCompleted()
})

function NavItem({ to, icon, label, locked }) {
  const route = useRoute()
  const isActive = to === '/' ? route.path === '/' : route.path.startsWith(to)

  if (locked) {
    return h('button', {
      type: 'button',
      class: 'flex flex-col items-center gap-1 px-4 py-1.5 rounded-xl text-muted-foreground opacity-50',
      onClick: () => {
        hapticFeedback('rigid')
        showAlert(t('trainer.lockedByAlphabet'))
      },
    }, [
      h('div', { class: 'relative' }, [
        h(icon, { class: 'h-5 w-5' }),
        h(LockIcon, { class: 'absolute -top-1.5 -right-2 h-3 w-3' }),
      ]),
      h('span', { class: 'text-[10px]' }, label),
    ])
  }

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

NavItem.props = ['to', 'icon', 'label', 'locked']
</script>
