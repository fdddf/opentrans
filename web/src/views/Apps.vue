<template>
  <div class="space-y-6">
    <header class="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
      <div>
        <p class="text-xs uppercase tracking-[0.3em] text-slate-500">{{ t('nav.apps') }}</p>
        <h1 class="text-2xl font-semibold">{{ t('apps.title') }}</h1>
        <p class="text-sm text-slate-400">{{ t('apps.subtitle') }}</p>
      </div>
      <div class="flex items-center gap-2">
        <button class="rounded-lg border border-white/20 px-3 py-2 text-sm hover:border-mint/60 hover:text-mint">{{ t('apps.sync') }}</button>
        <button class="rounded-lg bg-mint px-3 py-2 text-sm font-semibold text-midnight shadow" @click="showModal = true">{{ t('apps.manual') }}</button>
      </div>
    </header>

    <div v-if="showModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 px-4">
      <div class="w-full max-w-lg rounded-2xl border border-white/10 bg-midnight p-6 shadow-xl">
        <div class="flex items-center justify-between">
          <h2 class="text-lg font-semibold">{{ t('apps.manual') }}</h2>
          <button class="text-slate-400 hover:text-white" @click="closeModal">×</button>
        </div>

        <form class="mt-4 space-y-4" @submit.prevent="createApp">
          <div>
            <label class="text-sm text-slate-400">{{ t('common.name') || 'Name' }}</label>
            <input v-model="form.name" class="mt-1 w-full rounded-lg bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint" required />
          </div>
          <div>
            <label class="text-sm text-slate-400">Platform</label>
            <select v-model="form.platform" class="mt-1 w-full rounded-lg bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint">
              <option value="iOS">iOS</option>
              <option value="Android">Android</option>
              <option value="Web">Web</option>
            </select>
          </div>
          <div>
            <label class="text-sm text-slate-400">Bundle ID</label>
            <input v-model="form.bundleId" class="mt-1 w-full rounded-lg bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint" placeholder="com.example.app" required />
          </div>
          <div>
            <label class="text-sm text-slate-400">{{ t('workspace.addLang') }}</label>
            <input v-model="form.sourceLanguage" class="mt-1 w-full rounded-lg bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint" placeholder="en" />
          </div>
          <div>
            <label class="text-sm text-slate-400">Apple ID (optional)</label>
            <input v-model="form.appleId" class="mt-1 w-full rounded-lg bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint" placeholder="123456789" />
          </div>

          <div class="flex justify-end gap-2 pt-2">
            <button type="button" class="rounded-lg border border-white/20 px-3 py-2 text-sm hover:border-mint/60 hover:text-mint" @click="closeModal">{{ t('common.cancel') || 'Cancel' }}</button>
            <button type="submit" class="rounded-lg bg-mint px-3 py-2 text-sm font-semibold text-midnight shadow">{{ t('common.save') || 'Save' }}</button>
          </div>
        </form>
      </div>
    </div>

    <section class="grid grid-cols-1 lg:grid-cols-3 gap-4">
      <div class="rounded-2xl border border-white/10 bg-white/5 p-4" v-for="app in apps" :key="app.id">
        <div class="flex items-center justify-between">
          <div>
            <h3 class="font-semibold">{{ app.name }}</h3>
            <p class="text-xs text-slate-500">{{ app.platform || 'iOS' }}</p>
          </div>
          <span class="rounded-full bg-emerald-900/40 px-3 py-1 text-xs text-emerald-200" v-if="app.synced">同步</span>
          <span class="rounded-full bg-amber-900/40 px-3 py-1 text-xs text-amber-200" v-else>手动</span>
        </div>
        <p class="mt-2 text-sm text-slate-400">Bundle ID：{{ app.bundleId }}</p>
        <p class="text-sm text-slate-400">来源语言：{{ app.sourceLanguage }}</p>
        <div class="mt-3 flex gap-2 text-xs">
          <RouterLink :to="`/apps/${app.id}`" class="rounded border border-white/20 px-2 py-1 hover:border-mint/60 hover:text-mint">进入工作区</RouterLink>
          <button class="rounded border border-white/20 px-2 py-1 hover:border-mint/60 hover:text-mint" :disabled="app.synced">删除</button>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useApi } from '../composables/useApi'
import type { App } from '../composables/useApi'

const { t } = useI18n()
const { api } = useApi()

const apps = ref<App[]>([])
const showModal = ref(false)
const form = reactive({
  name: '',
  platform: 'iOS',
  bundleId: '',
  sourceLanguage: 'en',
  appleId: ''
})

const nextId = computed(() => (apps.value.length ? Math.max(...apps.value.map((a) => a.id)) + 1 : 1))

function resetForm() {
  form.name = ''
  form.platform = 'iOS'
  form.bundleId = ''
  form.sourceLanguage = 'en'
  form.appleId = ''
}

function closeModal() {
  showModal.value = false
  resetForm()
}

async function createApp() {
  if (!form.name.trim() || !form.bundleId.trim()) return
  
  try {
    const response = await api.createApp({
      name: form.name.trim(),
      bundleId: form.bundleId.trim(),
      appleId: form.appleId,
      primaryLocale: form.sourceLanguage || 'en'
    })
    
    if (response.success) {
      apps.value.unshift({
        id: response.app.id,
        name: form.name.trim(),
        platform: form.platform,
        synced: false,
        bundleId: form.bundleId.trim(),
        sourceLanguage: form.sourceLanguage || 'en',
        appleId: form.appleId,
        userId: response.app.userId,
        createdAt: response.app.createdAt,
        updatedAt: response.app.updatedAt,
        description: response.app.description || '',
        isReadyForReview: response.app.isReadyForReview || false,
        primaryLocale: form.sourceLanguage || 'en',
        shortDescription: response.app.shortDescription || '',
        longDescription: response.app.longDescription || '',
        keywords: response.app.keywords || '',
        supportUrl: response.app.supportUrl || '',
        marketingUrl: response.app.marketingUrl || '',
        privacyUrl: response.app.privacyUrl || '',
        version: response.app.version || '',
        appCategory: response.app.appCategory || ''
      })
      closeModal()
    }
  } catch (error) {
    console.error('Failed to create app:', error)
  }
}

async function fetchApps() {
  try {
    const response = await api.getApps()
    if (response.success) {
      apps.value = response.apps
    }
  } catch (error) {
    console.error('Failed to fetch apps:', error)
  }
}

onMounted(() => {
  fetchApps()
})
</script>
