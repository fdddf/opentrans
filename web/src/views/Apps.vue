<template>
  <div class="space-y-6">
    <header class="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
      <div>
        <p class="text-xs uppercase tracking-[0.3em] text-slate-500">{{ t('nav.apps') }}</p>
        <h1 class="text-2xl font-semibold">{{ t('apps.title') }}</h1>
        <p class="text-sm text-slate-400">{{ t('apps.subtitle') }}</p>
      </div>
      <div class="flex items-center gap-2">
        <button 
          class="rounded-lg border border-white/20 px-3 py-2 text-sm hover:border-mint/60 hover:text-mint" 
          @click="showSyncModal = true"
          :disabled="syncing"
        >
          <span v-if="!syncing">{{ t('apps.sync') }}</span>
          <span v-else class="flex items-center gap-1">
            <span class="h-3 w-3 rounded-full bg-mint animate-pulse"></span>
            {{ t('apps.syncing') }}
          </span>
        </button>
        <button class="rounded-lg bg-mint px-3 py-2 text-sm font-semibold text-midnight shadow" @click="showModal = true">{{ t('apps.add') }}</button>
      </div>
    </header>

    <!-- Sync Apps Modal -->
    <div v-if="showSyncModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 px-4">
      <div class="w-full max-w-lg rounded-2xl border border-white/10 bg-midnight p-6 shadow-xl">
        <div class="flex items-center justify-between">
          <h2 class="text-lg font-semibold">{{ t('apps.sync') }}</h2>
          <button class="text-slate-400 hover:text-white" @click="closeSyncModal">×</button>
        </div>

        <div class="mt-4 space-y-4" v-if="!hasAppleConnectConfig">
          <div class="p-4 rounded-lg bg-amber-900/20 border border-amber-500/30">
            <div class="flex items-start">
              <span class="text-amber-400">⚠</span>
              <div class="ml-3">
                <p class="text-sm font-medium text-amber-100">{{ t('apps.noConfig.title') }}</p>
                <p class="text-sm mt-1 text-amber-200">{{ t('apps.noConfig.message') }}</p>
              </div>
            </div>
            <div class="mt-4 flex justify-end">
              <RouterLink to="/apple-connect-config" class="rounded-lg bg-amber-600 px-3 py-2 text-sm font-semibold text-white">
                {{ t('apps.setupConfig') }}
              </RouterLink>
            </div>
          </div>
        </div>
        <div class="mt-4 space-y-4" v-else>
          <div>
            <p class="text-sm text-slate-400 mb-3">{{ t('apps.selectAppleConnect') }}</p>
            <select 
              v-model="selectedConfigId" 
              class="w-full rounded-lg bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint" 
              required
            >
              <option value="">{{ t('apps.chooseConfig') }}</option>
              <option v-for="config in appleConnectConfigs" :key="config.id" :value="config.id">
                {{ t('common.appleConnectConfig') }} ({{ config.id }})
              </option>
            </select>
          </div>
          
          <div v-if="syncResult" class="p-3 rounded-lg bg-white/5 border border-white/10">
            <p class="text-sm">{{ syncResult.message }}</p>
            <p v-if="syncResult.count" class="text-sm text-mint">{{ t('apps.syncedCount', { count: syncResult.count }) }}</p>
          </div>

          <div class="flex justify-end gap-2 pt-4">
            <button 
              type="button" 
              class="rounded-lg border border-white/20 px-3 py-2 text-sm hover:border-mint/60 hover:text-mint" 
              @click="closeSyncModal"
            >
              {{ t('common.cancel') || 'Cancel' }}
            </button>
            <button 
              type="button" 
              class="rounded-lg bg-mint px-3 py-2 text-sm font-semibold text-midnight shadow" 
              @click="syncApps"
              :disabled="!selectedConfigId || syncing"
            >
              <span v-if="!syncing">{{ t('apps.syncNow') }}</span>
              <span v-else class="flex items-center gap-1">
                <span class="h-3 w-3 rounded-full bg-midnight animate-pulse"></span>
                {{ t('apps.syncing') }}
              </span>
            </button>
          </div>
        </div>
      </div>
    </div>

    <div v-if="showModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 px-4">
      <div class="w-full max-w-lg rounded-2xl border border-white/10 bg-midnight p-6 shadow-xl">
        <div class="flex items-center justify-between">
          <h2 class="text-lg font-semibold">{{ t('apps.add') }}</h2>
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
            <label class="text-sm text-slate-400">{{ t('apps.primaryLocale') }}</label>
            <input v-model="form.sourceLanguage" class="mt-1 w-full rounded-lg bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint" placeholder="en" />
          </div>
          <div>
            <label class="text-sm text-slate-400">{{ t('apps.appleId') }}</label>
            <input v-model="form.appleId" class="mt-1 w-full rounded-lg bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint" :placeholder="t('apps.appleIdPlaceholder')" />
          </div>
          <div>
            <label class="text-sm text-slate-400">{{ t('apps.projectGroup') }}</label>
            <select
              v-model="form.projectId"
              class="mt-1 w-full rounded-lg bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint"
            >
              <option :value="null">{{ t('apps.noGroup') }}</option>
              <option v-for="project in projectGroups" :key="project.id" :value="project.id">
                {{ project.name }}
              </option>
            </select>
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
          <span class="rounded-full bg-emerald-900/40 px-3 py-1 text-xs text-emerald-200" v-if="app.origin === 'synced'">{{ t('apps.synced') }}</span>
          <span class="rounded-full bg-amber-900/40 px-3 py-1 text-xs text-amber-200" v-else>{{ t('apps.manual') }}</span>
        </div>
        <p class="mt-2 text-sm text-slate-400">{{ t('apps.bundleId') }}: {{ app.bundleId }}</p>
        <p class="text-sm text-slate-400">{{ t('apps.primaryLocale') }}: {{ app.primaryLocale }}</p>
        
        <!-- App metadata preview -->
        <div class="mt-3 text-xs text-slate-400 space-y-1">
          <p v-if="app.shortDescription"><span class="text-slate-500">{{ t('apps.shortDescription') }}:</span> {{ app.shortDescription }}</p>
          <p v-if="app.keywords"><span class="text-slate-500">{{ t('apps.keywords') }}:</span> {{ app.keywords }}</p>
          <p v-if="app.supportUrl"><span class="text-slate-500">{{ t('apps.supportUrl') }}:</span> {{ app.supportUrl }}</p>
        </div>
        
        <div class="mt-3 flex gap-2 text-xs">
      <RouterLink :to="`/apps/${app.id}`" class="rounded border border-white/20 px-2 py-1 hover:border-mint/60 hover:text-mint">{{ t('apps.manage') }}</RouterLink>
      <RouterLink :to="`/apps/${app.id}/localizations`" class="rounded border border-white/20 px-2 py-1 hover:border-mint/60 hover:text-mint">{{ t('apps.localizations') }}</RouterLink>
      <button v-if="app.origin === 'manual'" class="rounded border border-white/20 px-2 py-1 hover:border-rose-600/60 hover:text-rose-500" @click="deleteApp(app.id)">{{ t('common.delete') }}</button>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useApi } from '../composables/useApi'
import type { App, ProviderConfig, Project } from '../composables/useApi'

const { t } = useI18n()
const { api } = useApi()

const apps = ref<App[]>([])
const showModal = ref(false)
const showSyncModal = ref(false)
const syncing = ref(false)
const syncResult = ref<{ message: string; count?: number } | null>(null)
const selectedConfigId = ref<number | null>(null)
const appleConnectConfigs = ref<ProviderConfig[]>([])
const projectGroups = ref<Project[]>([])

const form = reactive({
  name: '',
  platform: 'iOS',
  bundleId: '',
  sourceLanguage: 'en',
  appleId: '',
  projectId: null as number | null
})

const hasAppleConnectConfig = computed(() => appleConnectConfigs.value.length > 0)

const nextId = computed(() => (apps.value.length ? Math.max(...apps.value.map((a) => a.id)) + 1 : 1))

function resetForm() {
  form.name = ''
  form.platform = 'iOS'
  form.bundleId = ''
  form.sourceLanguage = 'en'
  form.appleId = ''
  form.projectId = null
}

function closeModal() {
  showModal.value = false
  resetForm()
}

function closeSyncModal() {
  showSyncModal.value = false
  selectedConfigId.value = null
  syncResult.value = null
}

async function fetchProviderConfigs() {
  try {
    const response = await api.getAppleConnectConfigs()
    if (response.success) {
      appleConnectConfigs.value = response.configs || response.data || []
    }
  } catch (error) {
    console.error('Failed to fetch Apple Connect configs:', error)
  }
}

async function fetchProjectGroups() {
  try {
    const response = await api.getProjects('app_group')
    if (response.success) {
      projectGroups.value = response.projects
    }
  } catch (error) {
    console.error('Failed to fetch project groups:', error)
  }
}

async function syncApps() {
  if (!selectedConfigId.value) return

  syncing.value = true
  syncResult.value = null

  try {
    // Get the selected config to extract credentials
    const selectedConfig = appleConnectConfigs.value.find(config => config.id === selectedConfigId.value);
    if (!selectedConfig) {
      throw new Error('Selected configuration not found');
    }

    const response = await api.syncAppleApps({
      configId: selectedConfig.id
    })
    
    if (response.success) {
      syncResult.value = {
        message: response.message || t('apps.syncSuccess'),
        count: response.count
      }
      // Refresh apps list
      fetchApps()
    } else {
      syncResult.value = {
        message: response.message || t('apps.syncFailed')
      }
    }
  } catch (error) {
    console.error('Failed to sync apps:', error)
    syncResult.value = {
      message: t('apps.syncError')
    }
  } finally {
    syncing.value = false
  }
}

async function createApp() {
  if (!form.name.trim() || !form.bundleId.trim()) return
  
  try {
    const response = await api.createApp({
      name: form.name.trim(),
      bundleId: form.bundleId.trim(),
      appleId: form.appleId,
      primaryLocale: form.sourceLanguage || 'en',
      projectId: form.projectId || undefined
    })
    
    if (response.success) {
      await fetchApps()
      closeModal()
    }
  } catch (error) {
    console.error('Failed to create app:', error)
  }
}

async function deleteApp(appId: number) {
  if (!confirm(t('apps.confirmDelete'))) {
    return;
  }
  
  try {
    await api.deleteApp(appId);
    apps.value = apps.value.filter(app => app.id !== appId);
  } catch (error) {
    console.error('Failed to delete app:', error);
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
  fetchProviderConfigs()
  fetchProjectGroups()
})
</script>
