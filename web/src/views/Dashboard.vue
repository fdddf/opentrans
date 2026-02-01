<template>
  <div class="space-y-6">
    <header class="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
      <div>
        <p class="text-xs uppercase tracking-[0.3em] text-slate-500">{{ t('common.dashboard') }}</p>
        <h1 class="text-2xl font-semibold">{{ t('dashboard.title') }}</h1>
        <p class="text-sm text-slate-400">{{ t('dashboard.subtitle') }}</p>
      </div>
      <div class="flex flex-wrap items-center gap-2">
        <router-link to="/apps" class="rounded-full bg-mint px-4 py-2 text-sm font-semibold text-midnight shadow-glow">{{ t('nav.apps') }}</router-link>
        <router-link to="/languages" class="rounded-full border border-white/20 px-4 py-2 text-sm hover:border-mint/60 hover:text-mint">{{ t('nav.languages') }}</router-link>
      </div>
    </header>

    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
      <div class="glass rounded-2xl p-4">
        <p class="text-xs text-slate-500">{{ t('dashboard.activeUsers') }}</p>
        <p class="mt-1 text-2xl font-semibold">{{ metrics.activeUsers }}</p>
      </div>
      <div class="glass rounded-2xl p-4">
        <p class="text-xs text-slate-500">{{ t('dashboard.activeSubscriptions') }}</p>
        <p class="mt-1 text-2xl font-semibold">{{ metrics.activeSubscriptions }}</p>
      </div>
      <div class="glass rounded-2xl p-4">
        <p class="text-xs text-slate-500">{{ t('dashboard.apps') }}</p>
        <p class="mt-1 text-2xl font-semibold">{{ metrics.apps }}</p>
      </div>
      <div class="glass rounded-2xl p-4">
        <p class="text-xs text-slate-500">{{ t('dashboard.pendingStrings') }}</p>
        <p class="mt-1 text-2xl font-semibold">{{ metrics.pendingStrings }}</p>
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
      <section class="glass rounded-2xl p-4">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-xs uppercase tracking-[0.2em] text-slate-500">{{ t('dashboard.recentActivity') }}</p>
            <h2 class="text-lg font-semibold">{{ t('dashboard.opsLog') }}</h2>
          </div>
        </div>
        <div class="mt-4 space-y-3">
          <div v-for="activity in recentActivities" :key="activity.id" class="flex items-start gap-3 rounded-xl bg-white/5 p-3">
            <div class="rounded-lg bg-mint/10 px-2 py-1 text-xs text-mint">{{ activity.scope }}</div>
            <div>
              <div class="text-sm font-semibold">{{ activity.action }}</div>
              <div class="text-xs text-slate-500">{{ activity.timestamp }}</div>
            </div>
          </div>
        </div>
      </section>

      <section class="glass rounded-2xl p-4">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-xs uppercase tracking-[0.2em] text-slate-500">{{ t('dashboard.translationProgress') }}</p>
            <h2 class="text-lg font-semibold">{{ t('dashboard.translationProgress') }}</h2>
          </div>
          <router-link to="/apps" class="text-xs text-mint hover:underline">{{ t('dashboard.viewDetails') }}</router-link>
        </div>
        <div class="mt-4 space-y-3">
          <div v-for="lang in languageProgress" :key="lang.code" class="rounded-xl border border-white/10 bg-white/5 p-3">
            <div class="flex items-center justify-between">
              <div class="font-semibold">{{ lang.name }} ({{ lang.code }})</div>
              <div class="text-xs text-slate-500">{{ lang.done }}/{{ lang.total }}</div>
            </div>
            <div class="mt-2 h-2 rounded-full bg-white/10">
              <div class="h-2 rounded-full bg-mint" :style="{ width: `${Math.round((lang.done / lang.total) * 100)}%` }"></div>
            </div>
          </div>
        </div>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { ref, onMounted } from 'vue'
import { useApi } from '@/composables/useApi'

const { t } = useI18n()
const { api } = useApi()

const metrics = ref({
  activeUsers: 0,
  activeSubscriptions: 0,
  apps: 0,
  pendingStrings: 0
})

const recentActivities = ref<any[]>([])
const languageProgress = ref([
  { code: 'zh-CN', name: '简体中文', done: 0, total: 0 },
  { code: 'ja', name: '日本語', done: 0, total: 0 },
  { code: 'fr', name: 'Français', done: 0, total: 0 }
])

const formatTimeAgo = (date: string) => {
  const now = new Date()
  const past = new Date(date)
  const diffMs = now.getTime() - past.getTime()
  const diffMins = Math.floor(diffMs / 60000)

  if (diffMins < 1) return '刚刚'
  if (diffMins < 60) return `${diffMins} 分钟前`
  const diffHours = Math.floor(diffMins / 60)
  if (diffHours < 24) return `${diffHours} 小时前`
  const diffDays = Math.floor(diffHours / 24)
  return `${diffDays} 天前`
}

const loadActivities = async () => {
  try {
    const response = await api.getUserActivities()
    if (response.success) {
      recentActivities.value = response.activities.map((activity: any) => ({
        id: activity.id,
        scope: activity.details || '操作',
        action: activity.action,
        timestamp: formatTimeAgo(activity.createdAt)
      }))
    }
  } catch (error) {
    console.error('Failed to load activities:', error)
  }
}

onMounted(() => {
  loadActivities()
})
</script>