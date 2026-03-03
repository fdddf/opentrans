<template>
  <div class="space-y-6">
    <header class="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
      <div>
        <p class="text-xs uppercase tracking-[0.3em] text-slate-500">{{ t('nav.apps') }}</p>
        <h1 class="text-2xl font-semibold">{{ t('whatsnew.title') }}</h1>
        <p class="text-sm text-slate-400">{{ t('whatsnew.subtitle') }}</p>
      </div>
      <div class="flex items-center gap-2">
        <button
          class="rounded-lg border border-white/20 px-3 py-2 text-sm hover:border-mint/60 hover:text-mint"
          @click="showPushModal = true"
          :disabled="pushing"
        >
          {{ pushing ? t('whatsnew.pushing') : t('whatsnew.pushToApple') }}
        </button>
        <button
          class="rounded-lg border border-white/20 px-3 py-2 text-sm hover:border-mint/60 hover:text-mint"
          @click="showBatchTranslateModal = true"
          :disabled="batchTranslating"
        >
          {{ batchTranslating ? t('whatsnew.translateProgress', { progress: batchTranslateProgress }) : t('whatsnew.batchTranslate') }}
        </button>
      </div>
    </header>

    <!-- Main Content with Two Columns -->
    <div class="flex gap-6">
      <!-- Left: Language Sidebar -->
      <aside class="w-80 flex-shrink-0">
        <div class="sticky top-6 rounded-2xl border border-white/10 bg-white/5 p-4">
          <h2 class="mb-4 text-sm font-semibold text-slate-300">{{ t('whatsnew.languages') }}</h2>
          
          <!-- Language Dropdown for Adding -->
          <div class="mb-4">
            <select
              v-model="selectedLanguageToAdd"
              class="w-full rounded-lg bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint"
              @change="handleLanguageSelect"
            >
              <option value="">{{ t('whatsnew.selectLanguage') }}</option>
              <option v-for="lang in availableLanguages" :key="lang.code" :value="lang.code">
                {{ getLanguageFlag(lang.code) }} {{ lang.name }} ({{ lang.code }})
              </option>
            </select>
          </div>
          
          <!-- Language List -->
          <div class="space-y-1 max-h-[calc(100vh-300px)] overflow-y-auto">
            <button
              v-for="loc in localizations"
              :key="loc.id"
              @click="selectedLocalizationId = loc.id"
              :class="[
                'w-full flex items-center gap-2 rounded-lg px-3 py-2 text-left text-sm transition-all',
                selectedLocalizationId === loc.id
                  ? 'bg-mint text-midnight font-medium'
                  : 'hover:bg-white/10 text-slate-300'
              ]"
            >
              <span class="text-lg">{{ getLanguageFlag(loc.languageCode) }}</span>
              <span class="flex-1 truncate">{{ getLanguageName(loc.languageCode) }}</span>
              <span v-if="loc.id < 0" class="text-xs text-orange-400">*</span>
              <span v-if="loc.syncStatus" :class="[
                'w-2 h-2 rounded-full flex-shrink-0',
                loc.syncStatus === 'synced' ? 'bg-green-500' :
                loc.syncStatus === 'pending' ? 'bg-yellow-500' :
                loc.syncStatus === 'failed' ? 'bg-red-500' : 'bg-gray-500'
              ]"></span>
            </button>
            
            <!-- Empty State -->
            <div v-if="localizations.length === 0" class="py-8 text-center text-sm text-slate-500">
              {{ t('whatsnew.noLanguages') }}
            </div>
          </div>
        </div>
      </aside>

      <!-- Right: Edit Panel -->
      <main class="flex-1">
        <div v-if="selectedLocalization" class="rounded-2xl border border-white/10 bg-white/5 p-6">
          <!-- Header -->
          <div class="flex items-center justify-between mb-6">
            <div class="flex items-center gap-3">
              <span class="text-4xl">{{ getLanguageFlag(selectedLocalization.languageCode) }}</span>
              <div>
                <h2 class="text-xl font-semibold">
                  {{ getLanguageName(selectedLocalization.languageCode) }} ({{ selectedLocalization.languageCode }})
                  <span v-if="selectedLocalization.id < 0" class="text-sm font-normal text-orange-400 ml-2">{{ t('whatsnew.new') }}</span>
                </h2>
              </div>
            </div>
            <button class="text-slate-400 hover:text-rose-500" @click="deleteLocalization(selectedLocalization.id)">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
              </svg>
            </button>
          </div>

          <!-- Sync Info -->
          <div v-if="selectedLocalization.id > 0" class="flex flex-wrap gap-4 text-xs text-slate-400 mb-6 pb-4 border-b border-white/10">
            <div>{{ t('whatsnew.syncStatus') }}: {{ getSyncStatusText(selectedLocalization.syncStatus) }}</div>
            <div>{{ t('whatsnew.version') }}: {{ selectedLocalization.version || '-' }}</div>
            <div>{{ t('whatsnew.lastSynced') }}: {{ formatDate(selectedLocalization.syncedAt) }}</div>
          </div>

          <!-- Edit Form -->
          <form class="space-y-6" @submit.prevent="updateLocalization">
            <!-- What's New Section -->
            <div class="p-4 rounded-xl bg-white/5 border border-white/10">
              <div class="flex items-center justify-between mb-3">
                <label class="text-sm font-medium text-slate-300">{{ t('whatsnew.whatsNew') }}</label>
                <span class="text-xs text-slate-500">{{ (editLocalizationData.whatsNew || '').length }} / 4000</span>
              </div>
              <textarea 
                v-model="editLocalizationData.whatsNew" 
                rows="6" 
                class="mt-1 w-full rounded-lg bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint whitespace-pre-wrap" 
                :placeholder="t('whatsnew.whatsNewPlaceholder')"
                maxlength="4000"
              ></textarea>
              <p class="text-xs text-slate-500 mt-2">{{ t('whatsnew.whatsNewDesc') }}</p>
            </div>

            <!-- Promotional Text Section -->
            <div class="p-4 rounded-xl bg-white/5 border border-white/10">
              <div class="flex items-center justify-between mb-3">
                <label class="text-sm font-medium text-slate-300">{{ t('whatsnew.promotionalText') }}</label>
                <span class="text-xs text-slate-500">{{ (editLocalizationData.promotionalText || '').length }} / 170</span>
              </div>
              <textarea 
                v-model="editLocalizationData.promotionalText" 
                rows="3" 
                class="mt-1 w-full rounded-lg bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint whitespace-pre-wrap" 
                :placeholder="t('whatsnew.promotionalTextPlaceholder')"
                maxlength="170"
              ></textarea>
              <p class="text-xs text-slate-500 mt-2">{{ t('whatsnew.promotionalTextDesc') }}</p>
            </div>

            <div class="flex justify-end gap-2 pt-4">
              <button
                type="button"
                class="rounded-lg border border-white/20 px-4 py-2 text-sm hover:border-mint/60 hover:text-mint"
                @click="translateLocalization"
                :disabled="translating"
              >
                {{ translating ? t('common.translating') : t('common.translate') }}
              </button>
              <button type="submit" class="rounded-lg bg-mint px-4 py-2 text-sm font-semibold text-midnight shadow hover:bg-mint/90">
                {{ t('common.save') }}
              </button>
            </div>
          </form>
        </div>

        <!-- Empty State when no language selected -->
        <div v-else-if="localizations.length === 0" class="rounded-2xl border border-white/10 bg-white/5 p-12 text-center">
          <div class="text-6xl mb-4">📝</div>
          <h3 class="text-lg font-semibold mb-2">{{ t('whatsnew.noLocalizations') }}</h3>
          <p class="text-sm text-slate-400">{{ t('whatsnew.noLocalizationsDesc') }}</p>
        </div>

        <!-- Select a language prompt -->
        <div v-else class="rounded-2xl border border-white/10 bg-white/5 p-12 text-center">
          <div class="text-6xl mb-4">👈</div>
          <h3 class="text-lg font-semibold mb-2">{{ t('whatsnew.selectLanguagePrompt') }}</h3>
          <p class="text-sm text-slate-400">{{ t('whatsnew.selectLanguageDesc') }}</p>
        </div>
      </main>
    </div>

    <!-- Batch Translate Modal -->
    <div v-if="showBatchTranslateModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 px-4">
      <div class="w-full max-w-lg rounded-2xl border border-white/10 bg-midnight p-6 shadow-xl">
        <div class="flex items-center justify-between">
          <h2 class="text-lg font-semibold">{{ t('whatsnew.batchTranslate') }}</h2>
          <button class="text-slate-400 hover:text-white" @click="closeBatchTranslateModal">×</button>
        </div>

        <div class="mt-4 space-y-4">
          <div class="flex items-center gap-2 p-3 rounded-lg bg-mint/10 border border-mint/20">
            <span class="text-lg">🦙</span>
            <div>
              <p class="text-sm font-medium text-mint">腾讯混元 Hunyuan 模型</p>
              <p class="text-xs text-slate-400">本地模型翻译，数据安全</p>
            </div>
          </div>

          <!-- Overwrite option -->
          <div class="flex items-center gap-2 p-3 rounded-lg bg-white/5 border border-white/10">
            <input
              type="checkbox"
              id="overwriteExistingWhatsNew"
              v-model="overwriteExisting"
              class="rounded bg-white/10 border-white/20 text-mint focus:ring-mint"
            />
            <label for="overwriteExistingWhatsNew" class="text-sm cursor-pointer">
              <span class="text-slate-300">{{ t('applocalizations.overwriteExisting') }}</span>
              <p class="text-xs text-slate-500 mt-0.5">{{ t('applocalizations.overwriteExistingDesc') }}</p>
            </label>
          </div>

          <div>
            <div class="flex items-center justify-between mb-3">
              <p class="text-sm text-slate-400">{{ t('whatsnew.selectTargetLanguages') }}</p>
              <label class="flex items-center gap-2 text-sm cursor-pointer hover:text-mint text-slate-400">
                <input
                  type="checkbox"
                  :checked="selectAllLanguages"
                  @change="toggleSelectAllLanguages"
                  class="rounded bg-white/10 border-white/20 text-mint focus:ring-mint"
                />
                <span>{{ t('whatsnew.selectAll') }}</span>
              </label>
            </div>
            <div class="max-h-60 overflow-y-auto rounded-lg bg-white/5 border border-white/10 p-2 space-y-1">
              <label v-for="lang in availableLanguages" :key="lang.code" class="flex items-center gap-2 text-sm cursor-pointer hover:bg-white/5 p-1 rounded">
                <input
                  type="checkbox"
                  :value="lang.code"
                  v-model="selectedBatchTranslateLanguages"
                  :disabled="app?.primaryLocale === lang.code"
                  class="rounded bg-white/10 border-white/20 text-mint focus:ring-mint"
                />
                <span>{{ lang.name }} ({{ lang.code }})</span>
              </label>
            </div>
            <p class="text-xs text-slate-500 mt-1">{{ t('whatsnew.selectTargetLanguagesDesc') }}</p>
          </div>

          <div v-if="batchTranslateResult" class="p-3 rounded-lg bg-white/5 border border-white/10">
            <p class="text-sm">{{ batchTranslateResult.message }}</p>
            <div v-if="batchTranslateResult.progress !== undefined" class="mt-2">
              <div class="h-2 bg-white/10 rounded-full overflow-hidden">
                <div class="h-full bg-mint transition-all duration-300" :style="{ width: batchTranslateResult.progress + '%' }"></div>
              </div>
              <div class="flex justify-between text-xs text-slate-400 mt-1">
                <span>{{ t('whatsnew.translateProgress', { progress: batchTranslateResult.progress }) }}</span>
                <span v-if="batchTranslateResult.done !== undefined && batchTranslateResult.total !== undefined">
                  {{ batchTranslateResult.done }} / {{ batchTranslateResult.total }}
                </span>
              </div>
            </div>
          </div>

          <div class="flex justify-end gap-2 pt-4">
            <button
              type="button"
              class="rounded-lg border border-white/20 px-3 py-2 text-sm hover:border-mint/60 hover:text-mint"
              @click="closeBatchTranslateModal"
              :disabled="batchTranslating"
            >
              {{ t('common.cancel') }}
            </button>
            <button
              type="button"
              class="rounded-lg bg-mint px-3 py-2 text-sm font-semibold text-midnight shadow"
              @click="submitBatchTranslate"
              :disabled="selectedBatchTranslateLanguages.length === 0 || batchTranslating"
            >
              <span v-if="!batchTranslating">{{ t('whatsnew.translateAll') }}</span>
              <span v-else class="flex items-center gap-1">
                <span class="h-3 w-3 rounded-full bg-midnight animate-pulse"></span>
                {{ t('whatsnew.translateProgress', { progress: batchTranslateProgress }) }}
              </span>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Push to Apple Modal -->
    <div v-if="showPushModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 px-4">
      <div class="w-full max-w-lg rounded-2xl border border-white/10 bg-midnight p-6 shadow-xl">
        <div class="flex items-center justify-between">
          <h2 class="text-lg font-semibold">{{ t('whatsnew.pushToApple') }}</h2>
          <button class="text-slate-400 hover:text-white" @click="closePushModal">×</button>
        </div>

        <div class="mt-4 space-y-4">
          <div>
            <p class="text-sm text-slate-400 mb-3">{{ t('applocalizations.selectProviderConfig') }}</p>
            <select
              v-model="selectedPushConfigId"
              class="w-full rounded-lg bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint"
              required
            >
              <option value="">{{ t('applocalizations.chooseConfig') }}</option>
              <option v-for="config in appleConnectConfigs" :key="config.id" :value="config.id">
                {{ config.providerType }} (ID: {{ config.id }})
              </option>
            </select>
          </div>

          <div>
            <div class="flex items-center justify-between mb-3">
              <p class="text-sm text-slate-400">{{ t('whatsnew.selectTargetLanguages') }}</p>
              <label class="flex items-center gap-2 text-sm cursor-pointer hover:text-mint text-slate-400">
                <input
                  type="checkbox"
                  :checked="selectAllPushLanguages"
                  @change="toggleSelectAllPushLanguages"
                  class="rounded bg-white/10 border-white/20 text-mint focus:ring-mint"
                />
                <span>{{ t('whatsnew.selectAll') }}</span>
              </label>
            </div>
            <div class="max-h-60 overflow-y-auto rounded-lg bg-white/5 border border-white/10 p-2 space-y-1">
              <label v-for="lang in availableLanguages" :key="lang.code" class="flex items-center gap-2 text-sm cursor-pointer hover:bg-white/5 p-1 rounded">
                <input
                  type="checkbox"
                  :value="lang.code"
                  v-model="selectedPushLanguages"
                  class="rounded bg-white/10 border-white/20 text-mint focus:ring-mint"
                />
                <span>{{ lang.name }} ({{ lang.code }})</span>
              </label>
            </div>
            <p class="text-xs text-slate-500 mt-1">{{ t('whatsnew.pushLanguagesDesc') }}</p>
          </div>

          <div v-if="pushResult" class="p-3 rounded-lg bg-white/5 border border-white/10">
            <p class="text-sm">{{ pushResult.message }}</p>
          </div>

          <div class="flex justify-end gap-2 pt-4">
            <button
              type="button"
              class="rounded-lg border border-white/20 px-3 py-2 text-sm hover:border-mint/60 hover:text-mint"
              @click="closePushModal"
              :disabled="pushing"
            >
              {{ t('common.cancel') }}
            </button>
            <button
              type="button"
              class="rounded-lg bg-mint px-3 py-2 text-sm font-semibold text-midnight shadow"
              @click="submitPushToApple"
              :disabled="!selectedPushConfigId || selectedPushLanguages.length === 0 || pushing"
            >
              <span v-if="!pushing">{{ t('whatsnew.pushNow') }}</span>
              <span v-else class="flex items-center gap-1">
                <span class="h-3 w-3 rounded-full bg-midnight animate-pulse"></span>
                {{ t('whatsnew.pushing') }}
              </span>
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'
import { useApi } from '../composables/useApi'
import { useToast } from '../composables/useToast'
import type { AppLocalization } from '../composables/useApi'

const { t } = useI18n()
const route = useRoute()
const { api } = useApi()
const toast = useToast()

const appId = computed(() => Number(route.params.id))
const hasValidAppId = computed(() => Number.isFinite(appId.value) && appId.value > 0)

const localizations = ref<AppLocalization[]>([])
const app = ref<{ id: number; primaryLocale: string } | null>(null)
const selectedLocalizationId = ref<number | null>(null)
const selectedLanguageToAdd = ref<string>('')
const isNewLocalization = ref(false)
const availableLanguages = ref<{ code: string; name: string; native_name: string; region?: string; direction: string; emoji?: string }[]>([])
const loadingLanguages = ref(false)
const translating = ref(false)

// Batch translate state
const showBatchTranslateModal = ref(false)
const selectedBatchTranslateLanguages = ref<string[]>([])
const batchTranslating = ref(false)
const batchTranslateProgress = ref(0)
const batchTranslateResult = ref<{ message: string; progress?: number; done?: number; total?: number } | null>(null)
const overwriteExisting = ref(false) // Whether to overwrite existing translations

// Push to Apple state
const showPushModal = ref(false)
const selectedPushConfigId = ref<number | null>(null)
const selectedPushLanguages = ref<string[]>([])
const pushing = ref(false)
const pushResult = ref<{ message: string } | null>(null)
const appleConnectConfigs = ref<any[]>([])

// Computed property for selected localization
const selectedLocalization = computed(() => {
  return localizations.value.find(loc => loc.id === selectedLocalizationId.value) || null
})

async function fetchLanguages() {
  loadingLanguages.value = true
  try {
    const response = await api.getSupportedLanguages()
    if (response.success) {
      availableLanguages.value = response.languages.map(lang => ({
        code: lang.code,
        name: lang.name,
        native_name: lang.native_name,
        direction: lang.direction,
        emoji: lang.emoji || '🌐'
      }))
    }
  } catch (error) {
    console.error('Failed to load languages:', error)
  } finally {
    loadingLanguages.value = false
  }
}

function getLanguageFlag(code: string): string {
  const lang = availableLanguages.value.find(l => l.code === code)
  if (lang && lang.emoji) {
    return lang.emoji
  }
  
  const baseCode = code.split('-')[0]
  const baseLang = availableLanguages.value.find(l => l.code === baseCode)
  if (baseLang && baseLang.emoji) {
    return baseLang.emoji
  }
  
  return '🌐'
}

function getLanguageName(code: string): string {
  const lang = availableLanguages.value.find(l => l.code === code)
  return lang ? lang.name : code
}

const editLocalizationData = reactive({
  id: 0,
  languageCode: '',
  whatsNew: '',
  promotionalText: ''
})

function getSyncStatusText(status?: string): string {
  if (!status) return 'Unknown';
  const statusMap: Record<string, string> = {
    'synced': 'Synced',
    'pending': 'Pending',
    'failed': 'Failed'
  };
  return statusMap[status] || status || 'Unknown';
}

function formatDate(dateString?: string): string {
  if (!dateString) return 'N/A';
  const date = new Date(dateString);
  return date.toLocaleDateString() + ' ' + date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
}

// Watch for selection changes to populate edit form
watch(selectedLocalizationId, (newId) => {
  if (newId) {
    const loc = localizations.value.find(l => l.id === newId)
    if (loc) {
      editLocalization(loc)
    }
  }
})

async function fetchApp() {
  if (!hasValidAppId.value) return
  try {
    const response = await api.getApp(appId.value)
    if (response.success) {
      app.value = response.app
    }
  } catch (error) {
    console.error('Failed to fetch app:', error)
  }
}

async function fetchLocalizations() {
  if (!hasValidAppId.value) return
  try {
    const response = await api.getAppLocalizations(appId.value)
    if (response.success) {
      // Filter out duplicates by languageCode
      const seen = new Map<string, AppLocalization>()
      for (const loc of response.localizations) {
        const existing = seen.get(loc.languageCode)
        if (!existing || loc.id > existing.id) {
          seen.set(loc.languageCode, loc)
        }
      }
      localizations.value = Array.from(seen.values())
    }
  } catch (error) {
    console.error('Failed to fetch localizations:', error)
  }
}

function editLocalization(localization: AppLocalization) {
  editLocalizationData.id = localization.id
  editLocalizationData.languageCode = localization.languageCode
  editLocalizationData.whatsNew = localization.whatsNew || ''
  editLocalizationData.promotionalText = localization.promotionalText || ''
}

async function handleLanguageSelect() {
  if (!selectedLanguageToAdd.value) return

  // Check if language already exists
  const existing = localizations.value.find(loc => loc.languageCode === selectedLanguageToAdd.value)
  if (existing) {
    selectedLocalizationId.value = existing.id
    selectedLanguageToAdd.value = ''
    return
  }

  // Get en-US localization data for default values
  const enUsLocalization = localizations.value.find(loc => loc.languageCode === 'en-US' || loc.languageCode === 'en')

  editLocalizationData.id = 0
  editLocalizationData.languageCode = selectedLanguageToAdd.value
  editLocalizationData.whatsNew = enUsLocalization?.whatsNew || ''
  editLocalizationData.promotionalText = enUsLocalization?.promotionalText || ''

  isNewLocalization.value = true

  // Add temporary entry to list
  const tempId = -Date.now()
  const tempLocalization: AppLocalization = {
    id: tempId,
    appId: appId.value,
    languageCode: selectedLanguageToAdd.value,
    whatsNew: editLocalizationData.whatsNew,
    promotionalText: editLocalizationData.promotionalText,
    syncStatus: 'pending',
    version: '',
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString(),
    versionState: '',
    syncedAt: '',
    source: 'local'
  }
  localizations.value.push(tempLocalization)
  selectedLocalizationId.value = tempId
  selectedLanguageToAdd.value = ''
}

async function updateLocalization() {
  try {
    if (!hasValidAppId.value) {
      throw new Error('Invalid app ID')
    }

    let response
    if (isNewLocalization.value) {
      // Create new localization with only whatsNew and promotionalText
      response = await api.createAppLocalization(appId.value, {
        languageCode: editLocalizationData.languageCode,
        whatsNew: editLocalizationData.whatsNew,
        promotionalText: editLocalizationData.promotionalText
      })
    } else {
      // Update existing localization
      response = await api.updateAppLocalization(appId.value, editLocalizationData.languageCode, {
        whatsNew: editLocalizationData.whatsNew,
        promotionalText: editLocalizationData.promotionalText
      })
    }

    if (response.success) {
      isNewLocalization.value = false
      await fetchLocalizations()
      if (response.localization) {
        selectedLocalizationId.value = response.localization.id
      }
      toast.success(isNewLocalization.value ? t('whatsnew.addSuccess') : t('whatsnew.updateSuccess'))
    } else {
      toast.error(response.message || (isNewLocalization.value ? t('whatsnew.addFailed') : t('whatsnew.updateFailed')))
    }
  } catch (error) {
    console.error('Failed to save:', error)
    toast.error(isNewLocalization.value ? t('whatsnew.addFailed') : t('whatsnew.updateFailed'))
  }
}

async function deleteLocalization(id: number) {
  const localization = localizations.value.find(loc => loc.id === id)
  
  if (localization && app.value?.primaryLocale === localization.languageCode) {
    toast.warning(t('whatsnew.cannotDeletePrimary'))
    return
  }
  
  if (!confirm(t('whatsnew.confirmDelete'))) {
    return
  }
  
  // If it's a new unsaved localization (negative ID), just remove from list
  if (id < 0) {
    localizations.value = localizations.value.filter(loc => loc.id !== id)
    selectedLocalizationId.value = null
    isNewLocalization.value = false
    return
  }
  
  try {
    if (!hasValidAppId.value) {
      throw new Error('Invalid app ID')
    }
    const languageCode = localizations.value.find(loc => loc.id === id)?.languageCode
    
    if (languageCode) {
      await api.deleteAppLocalization(appId.value, languageCode)
      localizations.value = localizations.value.filter(loc => loc.id !== id)
      selectedLocalizationId.value = null
      toast.success(t('whatsnew.deleteSuccess'))
    }
  } catch (error) {
    console.error('Failed to delete:', error)
    toast.error(t('whatsnew.deleteFailed'))
  }
}

async function translateLocalization() {
  if (!editLocalizationData.languageCode) {
    toast.warning(t('whatsnew.selectLanguageFirst'))
    return
  }

  translating.value = true

  try {
    const sourceLanguage = app.value?.primaryLocale || 'en-US'
    const targetLanguage = editLocalizationData.languageCode

    // Translate what's new
    if (editLocalizationData.whatsNew) {
      const result = await api.translateText(editLocalizationData.whatsNew, sourceLanguage, targetLanguage)
      if (result.success) {
        editLocalizationData.whatsNew = result.text
      }
    }

    // Translate promotional text
    if (editLocalizationData.promotionalText) {
      const result = await api.translateText(editLocalizationData.promotionalText, sourceLanguage, targetLanguage)
      if (result.success) {
        editLocalizationData.promotionalText = result.text
      }
    }

    toast.success(t('whatsnew.translateSuccess'))
  } catch (error) {
    console.error('Translation failed:', error)
    toast.error(t('whatsnew.translateFailed'))
  } finally {
    translating.value = false
  }
}

// Computed property for select all languages
const selectAllLanguages = computed(() => {
  if (!app.value || availableLanguages.value.length === 0) {
    return false
  }
  const selectableLanguages = availableLanguages.value.filter(lang => lang.code !== app.value?.primaryLocale)
  return selectableLanguages.length > 0 &&
    selectableLanguages.every(lang => selectedBatchTranslateLanguages.value.includes(lang.code))
})

// Toggle select all languages
function toggleSelectAllLanguages() {
  if (!app.value) {
    return
  }

  const selectableLanguages = availableLanguages.value.filter(lang => lang.code !== app.value?.primaryLocale)
  const allSelected = selectAllLanguages.value

  if (allSelected) {
    selectedBatchTranslateLanguages.value = []
  } else {
    selectedBatchTranslateLanguages.value = selectableLanguages.map(lang => lang.code)
  }
}

function closeBatchTranslateModal() {
  showBatchTranslateModal.value = false
  selectedBatchTranslateLanguages.value = []
  batchTranslateResult.value = null
  batchTranslateProgress.value = 0
  overwriteExisting.value = false
}

// Push to Apple functions
const selectAllPushLanguages = computed(() => {
  if (availableLanguages.value.length === 0) return false
  return availableLanguages.value.every(lang => selectedPushLanguages.value.includes(lang.code))
})

function toggleSelectAllPushLanguages() {
  if (selectAllPushLanguages.value) {
    selectedPushLanguages.value = []
  } else {
    selectedPushLanguages.value = availableLanguages.value.map(lang => lang.code)
  }
}

function closePushModal() {
  showPushModal.value = false
  selectedPushConfigId.value = null
  selectedPushLanguages.value = []
  pushResult.value = null
}

async function fetchAppleConnectConfigs() {
  try {
    const response = await api.getAppleConnectConfigs()
    if (response.success) {
      appleConnectConfigs.value = response.data || []
      // Auto-select the first config if none is selected
      if (appleConnectConfigs.value.length > 0 && !selectedPushConfigId.value) {
        const defaultConfig = appleConnectConfigs.value.find((c: any) => c.isDefault)
        selectedPushConfigId.value = defaultConfig ? defaultConfig.id : appleConnectConfigs.value[0].id
      }
    }
  } catch (error) {
    console.error('Failed to fetch Apple Connect configs:', error)
  }
}

async function submitPushToApple() {
  if (!hasValidAppId.value) {
    toast.warning(t('whatsnew.invalidAppId'))
    return
  }

  if (!selectedPushConfigId.value) {
    toast.warning(t('applocalizations.selectProviderConfig'))
    return
  }

  if (selectedPushLanguages.value.length === 0) {
    toast.warning(t('whatsnew.selectAtLeastOne'))
    return
  }

  pushing.value = true
  pushResult.value = null

  try {
    const response = await api.pushUpdateContentToApple(appId.value, {
      configId: selectedPushConfigId.value,
      languageCodes: selectedPushLanguages.value
    })

    if (response.success) {
      pushResult.value = { message: response.message || t('whatsnew.pushSuccess') }
      toast.success(t('whatsnew.pushSuccess'))
      // Refresh localizations
      await fetchLocalizations()
    } else {
      pushResult.value = { message: response.message || t('whatsnew.pushFailed') }
      toast.error(t('whatsnew.pushFailed'))
    }
  } catch (error) {
    console.error('Failed to push to Apple:', error)
    pushResult.value = { message: t('whatsnew.pushError') }
    toast.error(t('whatsnew.pushError'))
  } finally {
    pushing.value = false
  }
}

async function submitBatchTranslate() {
  if (!hasValidAppId.value) {
    toast.warning(t('whatsnew.invalidAppId'))
    return
  }

  if (selectedBatchTranslateLanguages.value.length === 0) {
    toast.warning(t('whatsnew.selectAtLeastOne'))
    return
  }

  batchTranslating.value = true
  batchTranslateProgress.value = 0
  batchTranslateResult.value = { message: t('whatsnew.translateSubmitted'), progress: 0 }

  try {
    const sourceLanguage = app.value?.primaryLocale || 'en-US'

    const response = await api.translateAppLocalizations(appId.value, {
      providerType: 'llama',
      sourceLanguage: sourceLanguage,
      targetLanguages: selectedBatchTranslateLanguages.value,
      onlyTranslateWhatsNew: true,  // Always true for WhatsNew page
      configData: {
        threads: 4,
        temperature: 0.7,
        topP: 0.6,
        topK: 20,
        tokens: 4096,
        // Skip existing translations if overwrite is false
        skipExisting: !overwriteExisting.value
      }
    })

    if (response.success) {
      batchTranslateResult.value = {
        message: t('whatsnew.translateSubmitted'),
        progress: 0
      }

      pollTranslationProgress(response.job.id)
      toast.success(t('whatsnew.translateSubmitted'))
    } else {
      throw new Error(response.message || t('whatsnew.translateFailed'))
    }
  } catch (error) {
    console.error('Failed to submit batch translation:', error)
    batchTranslateResult.value = {
      message: t('whatsnew.translateError', { error: error instanceof Error ? error.message : 'Unknown error' })
    }
    toast.error(t('whatsnew.translateFailed'))
    batchTranslating.value = false
  }
}

async function pollTranslationProgress(jobId: number) {
  const pollInterval = setInterval(async () => {
    try {
      const response = await api.getQueueJob(jobId)
      if (response.success) {
        const job = response.job

        batchTranslateProgress.value = job.progress || 0
        batchTranslateResult.value = {
          message: t('whatsnew.translateSubmitted'),
          progress: batchTranslateProgress.value,
          done: job.done,
          total: job.total
        }

        if (job.status === 'completed') {
          clearInterval(pollInterval)
          batchTranslating.value = false
          batchTranslateResult.value = {
            message: t('whatsnew.translateCompleted'),
            progress: 100,
            done: job.done,
            total: job.total
          }

          await fetchLocalizations()
          toast.success(t('whatsnew.translateCompleted'))
        } else if (job.status === 'failed') {
          clearInterval(pollInterval)
          batchTranslating.value = false
          batchTranslateResult.value = {
            message: t('whatsnew.translateError', { error: job.error || 'Unknown error' }),
            done: job.done,
            total: job.total
          }
          toast.error(t('whatsnew.translateError', { error: job.error || 'Unknown error' }))
        }
      }
    } catch (error) {
      console.error('Failed to poll translation progress:', error)
      clearInterval(pollInterval)
      batchTranslating.value = false
    }
  }, 2000)
}

watch(appId, () => {
  fetchLocalizations()
})

onMounted(() => {
  fetchApp()
  fetchLanguages()
  fetchLocalizations()
  fetchAppleConnectConfigs()
})
</script>
