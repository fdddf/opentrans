<template>
  <div class="min-h-screen bg-midnight text-slate-100">
    <div class="flex min-h-screen">
      <aside class="w-64 border-r border-white/10 bg-midnight/80 backdrop-blur-lg hidden lg:flex flex-col">
        <div class="p-6 border-b border-white/5">
          <div class="text-sm uppercase tracking-[0.25em] text-slate-500">XCStrings</div>
          <div class="mt-1 text-xl font-semibold">{{ t('common.brand') }}</div>
        </div>
        <nav class="flex-1 p-4 space-y-2">
          <RouterLink
            v-for="item in navItems"
            :key="item.to"
            :to="item.to"
            class="flex items-center gap-3 rounded-xl px-3 py-2 text-sm transition"
            :class="isActive(item.to) ? 'bg-mint/10 text-mint border border-mint/40' : 'text-slate-300 hover:text-white hover:bg-white/5'"
          >
            <span class="inline-flex h-6 w-6 items-center justify-center rounded-lg bg-white/5 text-xs">{{ item.icon }}</span>
            <span>{{ typeof item.label === 'function' ? item.label() : item.label }}</span>
          </RouterLink>
        </nav>
        <div class="p-4 border-t border-white/5 text-xs text-slate-500">
          Multi-tenant & multilingual workspace
        </div>
      </aside>

      <div class="flex-1 flex flex-col">
        <header class="sticky top-0 z-10 bg-midnight/90 backdrop-blur-xl border-b border-white/10">
          <div class="mx-auto max-w-7xl px-4 py-4 flex items-center justify-between">
            <div>
              <div class="text-xs uppercase tracking-[0.2em] text-slate-500">Admin Console</div>
              <div class="text-lg font-semibold">Localization Ops</div>
            </div>
            <div class="flex items-center gap-3">
              <div class="hidden sm:flex items-center gap-2 rounded-full bg-white/5 px-3 py-1 text-xs">
                <span class="h-2 w-2 rounded-full bg-emerald-400"></span>
                <span>Data isolated per tenant</span>
              </div>
              <select v-model="locale" class="rounded-lg bg-white/5 px-2 py-1 text-xs border border-white/10">
                <option value="en">English</option>
                <option value="zh">中文</option>
              </select>
              <button class="rounded-full border border-white/20 px-3 py-1 text-xs hover:border-mint/60 hover:text-mint">Profile</button>
              <button class="rounded-full bg-mint px-3 py-1 text-xs font-semibold text-midnight">{{ t('common.logout') }}</button>
            </div>
          </div>
        </header>

        <main class="flex-1 mx-auto w-full max-w-7xl px-4 py-8">
          <router-view />
        </main>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { RouterLink, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'

const route = useRoute()
const { t, locale } = useI18n()

const navItems = [
  { to: '/dashboard', label: () => t('nav.dashboard'), icon: '🏠' },
  { to: '/apps', label: () => t('nav.apps'), icon: '📱' },
  { to: '/languages', label: () => t('nav.languages'), icon: '🌐' },
  { to: '/users', label: () => t('nav.users'), icon: '👥' },
  { to: '/subscriptions', label: () => t('nav.subscriptions'), icon: '💳' }
]

const isActive = (path: string) => route.path.startsWith(path)
</script>
