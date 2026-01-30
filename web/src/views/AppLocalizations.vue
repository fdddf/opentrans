<template>
  <div class="space-y-6">
    <header class="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
      <div>
        <p class="text-xs uppercase tracking-[0.3em] text-slate-500">{{ t('nav.applocalizations') }}</p>
        <h1 class="text-2xl font-semibold">{{ t('applocalizations.title') }}</h1>
        <p class="text-sm text-slate-400">{{ t('applocalizations.subtitle') }}</p>
      </div>
      <div class="flex items-center gap-2">
        <button 
          class="rounded-lg border border-white/20 px-3 py-2 text-sm hover:border-mint/60 hover:text-mint" 
          @click="showSyncModal = true"
        >
          {{ t('applocalizations.sync') }}
        </button>
        <button class="rounded-lg bg-mint px-3 py-2 text-sm font-semibold text-midnight shadow" @click="showAddLocalizationModal = true">{{ t('applocalizations.add') }}</button>
      </div>
    </header>

    <!-- Sync Localizations Modal -->
    <div v-if="showSyncModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 px-4">
      <div class="w-full max-w-lg rounded-2xl border border-white/10 bg-midnight p-6 shadow-xl">
        <div class="flex items-center justify-between">
          <h2 class="text-lg font-semibold">{{ t('applocalizations.sync') }}</h2>
          <button class="text-slate-400 hover:text-white" @click="closeSyncModal">×</button>
        </div>

        <div class="mt-4 space-y-4">
          <div>
            <p class="text-sm text-slate-400 mb-3">{{ t('applocalizations.selectProviderConfig') }}</p>
            <select 
              v-model="selectedConfigId" 
              class="w-full rounded-lg bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint" 
              required
            >
              <option value="">{{ t('applocalizations.chooseConfig') }}</option>
              <option v-for="config in appleConnectConfigs" :key="config.id" :value="config.id">
                {{ config.providerType }} ({{ getProviderDisplayName(config.providerType) }})
              </option>
            </select>
          </div>
          
          <div v-if="syncResult" class="p-3 rounded-lg bg-white/5 border border-white/10">
            <p class="text-sm">{{ syncResult.message }}</p>
            <p v-if="syncResult.count" class="text-sm text-mint">{{ t('applocalizations.syncedCount', { count: syncResult.count }) }}</p>
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
              @click="syncLocalizations"
              :disabled="!selectedConfigId || syncing"
            >
              <span v-if="!syncing">{{ t('applocalizations.syncNow') }}</span>
              <span v-else class="flex items-center gap-1">
                <span class="h-3 w-3 rounded-full bg-midnight animate-pulse"></span>
                {{ t('applocalizations.syncing') }}
              </span>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Add Localization Modal -->
    <div v-if="showAddLocalizationModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 px-4">
      <div class="w-full max-w-2xl rounded-2xl border border-white/10 bg-midnight p-6 shadow-xl">
        <div class="flex items-center justify-between">
          <h2 class="text-lg font-semibold">{{ t('applocalizations.add') }}</h2>
          <button class="text-slate-400 hover:text-white" @click="closeAddLocalizationModal">×</button>
        </div>

        <form class="mt-4 space-y-4" @submit.prevent="addLocalization">
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label class="text-sm text-slate-400">{{ t('applocalizations.language') }}</label>
              <select v-model="newLocalization.languageCode" class="mt-1 w-full rounded-lg bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint" required>
                <option v-for="lang in availableLanguages" :key="lang.code" :value="lang.code">{{ lang.name }} ({{ lang.code }})</option>
              </select>
            </div>
          </div>
          
          <div class="space-y-4">
            <div>
              <label class="text-sm text-slate-400">{{ t('applocalizations.appName') }}</label>
              <input v-model="newLocalization.name" class="mt-1 w-full rounded-lg bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint" :placeholder="t('applocalizations.appNamePlaceholder')" />
            </div>
            
            <div>
              <label class="text-sm text-slate-400">{{ t('applocalizations.subtitle') }}</label>
              <input v-model="newLocalization.subtitle" class="mt-1 w-full rounded-lg bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint" :placeholder="t('applocalizations.subtitlePlaceholder')" />
            </div>
            
            <div>
              <label class="text-sm text-slate-400">{{ t('applocalizations.shortDescription') }}</label>
              <textarea v-model="newLocalization.shortDescription" rows="2" class="mt-1 w-full rounded-lg bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint" :placeholder="t('applocalizations.shortDescriptionPlaceholder')"></textarea>
            </div>
            
            <div>
              <label class="text-sm text-slate-400">{{ t('applocalizations.longDescription') }}</label>
              <textarea v-model="newLocalization.longDescription" rows="4" class="mt-1 w-full rounded-lg bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint" :placeholder="t('applocalizations.longDescriptionPlaceholder')"></textarea>
            </div>
            
            <div>
              <label class="text-sm text-slate-400">{{ t('applocalizations.keywords') }}</label>
              <input v-model="newLocalization.keywords" class="mt-1 w-full rounded-lg bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint" :placeholder="t('applocalizations.keywordsPlaceholder')" />
            </div>
            
            <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div>
                <label class="text-sm text-slate-400">{{ t('applocalizations.privacyUrl') }}</label>
                <input v-model="newLocalization.privacyUrl" class="mt-1 w-full rounded-lg bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint" :placeholder="t('applocalizations.privacyUrlPlaceholder')" />
              </div>
              <div>
                <label class="text-sm text-slate-400">{{ t('applocalizations.marketingUrl') }}</label>
                <input v-model="newLocalization.marketingUrl" class="mt-1 w-full rounded-lg bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint" :placeholder="t('applocalizations.marketingUrlPlaceholder')" />
              </div>
              <div>
                <label class="text-sm text-slate-400">{{ t('applocalizations.supportUrl') }}</label>
                <input v-model="newLocalization.supportUrl" class="mt-1 w-full rounded-lg bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint" :placeholder="t('applocalizations.supportUrlPlaceholder')" />
              </div>
            </div>
          </div>

          <div class="flex justify-end gap-2 pt-4">
            <button type="button" class="rounded-lg border border-white/20 px-3 py-2 text-sm hover:border-mint/60 hover:text-mint" @click="closeAddLocalizationModal">{{ t('common.cancel') || 'Cancel' }}</button>
            <button type="submit" class="rounded-lg bg-mint px-3 py-2 text-sm font-semibold text-midnight shadow">{{ t('common.save') || 'Save' }}</button>
          </div>
        </form>
      </div>
    </div>

    <!-- Localizations List -->
    <section class="space-y-4">
      <div class="rounded-2xl border border-white/10 bg-white/5 p-4" v-for="loc in localizations" :key="loc.id">
        <div class="flex items-center justify-between">
          <div>
            <h3 class="font-semibold">{{ loc.name || loc.languageCode }}</h3>
            <p class="text-xs text-slate-500">{{ getLanguageName(loc.languageCode) }} ({{ loc.languageCode }})</p>
          </div>
          <div class="flex gap-2">
            <button class="rounded border border-white/20 px-2 py-1 text-xs hover:border-mint/60 hover:text-mint" @click="editLocalization(loc)">{{ t('common.edit') }}</button>
            <button class="rounded border border-white/20 px-2 py-1 text-xs hover:border-rose-600/60 hover:text-rose-500" @click="deleteLocalization(loc.id)">{{ t('common.delete') }}</button>
          </div>
        </div>
        
        <div class="mt-3 space-y-2 text-sm">
          <p v-if="loc.subtitle"><span class="text-slate-500">{{ t('applocalizations.subtitle') }}:</span> {{ loc.subtitle }}</p>
          <p v-if="loc.shortDescription"><span class="text-slate-500">{{ t('applocalizations.shortDescription') }}:</span> {{ loc.shortDescription }}</p>
          <p v-if="loc.longDescription"><span class="text-slate-500">{{ t('applocalizations.longDescription') }}:</span> {{ loc.longDescription }}</p>
          <p v-if="loc.keywords"><span class="text-slate-500">{{ t('applocalizations.keywords') }}:</span> {{ loc.keywords }}</p>
          <div v-if="loc.privacyUrl || loc.marketingUrl || loc.supportUrl" class="flex flex-wrap gap-4">
            <p v-if="loc.privacyUrl"><span class="text-slate-500">{{ t('applocalizations.privacyUrl') }}:</span> <a :href="loc.privacyUrl" target="_blank" class="text-mint hover:underline">{{ loc.privacyUrl }}</a></p>
            <p v-if="loc.marketingUrl"><span class="text-slate-500">{{ t('applocalizations.marketingUrl') }}:</span> <a :href="loc.marketingUrl" target="_blank" class="text-mint hover:underline">{{ loc.marketingUrl }}</a></p>
            <p v-if="loc.supportUrl"><span class="text-slate-500">{{ t('applocalizations.supportUrl') }}:</span> <a :href="loc.supportUrl" target="_blank" class="text-mint hover:underline">{{ loc.supportUrl }}</a></p>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'
import { useApi } from '../composables/useApi'
import type { AppLocalization, ProviderConfig } from '../composables/useApi'

const { t } = useI18n()
const route = useRoute()
const { api } = useApi()

const appId = computed(() => Number(route.params.id))
const hasValidAppId = computed(() => Number.isFinite(appId.value) && appId.value > 0)

const localizations = ref<AppLocalization[]>([])
const showAddLocalizationModal = ref(false)
const showSyncModal = ref(false)
const syncing = ref(false)
const syncResult = ref<{ message: string; count?: number } | null>(null)
const selectedConfigId = ref<number | null>(null)
const appleConnectConfigs = ref<ProviderConfig[]>([])

const newLocalization = reactive({
  languageCode: '',
  name: '',
  subtitle: '',
  privacyUrl: '',
  marketingUrl: '',
  supportUrl: '',
  shortDescription: '',
  longDescription: '',
  keywords: '',
  releaseNotes: ''
})

const availableLanguages = [
  { code: 'en-US', name: 'English (US)' },
  { code: 'en-GB', name: 'English (UK)' },
  { code: 'zh-Hans', name: 'Chinese (Simplified)' },
  { code: 'zh-Hant', name: 'Chinese (Traditional)' },
  { code: 'ja', name: 'Japanese' },
  { code: 'ko', name: 'Korean' },
  { code: 'fr', name: 'French' },
  { code: 'de', name: 'German' },
  { code: 'es', name: 'Spanish' },
  { code: 'ru', name: 'Russian' },
  { code: 'ar', name: 'Arabic' },
  { code: 'pt-BR', name: 'Portuguese (Brazil)' },
  { code: 'pt-PT', name: 'Portuguese (Portugal)' },
  { code: 'it', name: 'Italian' },
  { code: 'nl', name: 'Dutch' },
  { code: 'sv', name: 'Swedish' },
  { code: 'da', name: 'Danish' },
  { code: 'fi', name: 'Finnish' },
  { code: 'no', name: 'Norwegian' },
  { code: 'pl', name: 'Polish' },
  { code: 'tr', name: 'Turkish' },
  { code: 'th', name: 'Thai' },
  { code: 'vi', name: 'Vietnamese' },
  { code: 'hi', name: 'Hindi' },
  { code: 'id', name: 'Indonesian' },
  { code: 'cs', name: 'Czech' },
  { code: 'el', name: 'Greek' },
  { code: 'he', name: 'Hebrew' },
  { code: 'hu', name: 'Hungarian' },
  { code: 'ro', name: 'Romanian' },
  { code: 'sk', name: 'Slovak' },
  { code: 'uk', name: 'Ukrainian' }
]

function getLanguageName(code: string): string {
  const lang = availableLanguages.find(l => l.code === code)
  return lang ? lang.name : code
}

function resetLocalizationForm() {
  newLocalization.languageCode = ''
  newLocalization.name = ''
  newLocalization.subtitle = ''
  newLocalization.privacyUrl = ''
  newLocalization.marketingUrl = ''
  newLocalization.supportUrl = ''
  newLocalization.shortDescription = ''
  newLocalization.longDescription = ''
  newLocalization.keywords = ''
  newLocalization.releaseNotes = ''
}

function closeAddLocalizationModal() {
  showAddLocalizationModal.value = false
  resetLocalizationForm()
}

function closeSyncModal() {
  showSyncModal.value = false
  selectedConfigId.value = null
  syncResult.value = null
}

function getProviderDisplayName(providerType: string): string {
  const names: Record<string, string> = {
    'openai': 'OpenAI',
    'google': 'Google',
    'deepl': 'DeepL',
    'baidu': 'Baidu',
    'appleconnect': 'Apple Connect'
  }
  return names[providerType] || providerType
}

async function fetchProviderConfigs() {
  try {
    const response = await api.getProviderConfigs()
    if (response.success) {
      appleConnectConfigs.value = response.configs.filter(config => config.providerType === 'appleconnect')
    }
  } catch (error) {
    console.error('Failed to fetch provider configs:', error)
  }
}

async function syncLocalizations() {
  if (!selectedConfigId.value || !hasValidAppId.value) return

  syncing.value = true
  syncResult.value = null

  try {
    // Get the selected config to extract credentials
    const selectedConfig = appleConnectConfigs.value.find(config => config.id === selectedConfigId.value);
    if (!selectedConfig) {
      throw new Error('Selected configuration not found');
    }

    const response = await api.syncAppleAppLocalizations(appId.value, {
      configId: selectedConfigId.value
    })
    
    if (response.success) {
      syncResult.value = {
        message: response.message || t('applocalizations.syncSuccess'),
        count: response.count
      }
      // Refresh localizations list
      await fetchLocalizations()
    } else {
      syncResult.value = {
        message: response.message || t('applocalizations.syncFailed')
      }
    }
  } catch (error) {
    console.error('Failed to sync localizations:', error)
    syncResult.value = {
      message: t('applocalizations.syncError')
    }
  } finally {
    syncing.value = false
  }
}

async function addLocalization() {
  try {
    if (!hasValidAppId.value) {
      throw new Error('Invalid app ID')
    }
    const response = await api.createAppLocalization(appId.value, {
      languageCode: newLocalization.languageCode,
      name: newLocalization.name,
      subtitle: newLocalization.subtitle,
      privacyUrl: newLocalization.privacyUrl,
      marketingUrl: newLocalization.marketingUrl,
      supportUrl: newLocalization.supportUrl,
      downloadDescription: '', // Not using this field in the form
      shortDescription: newLocalization.shortDescription,
      longDescription: newLocalization.longDescription,
      keywords: newLocalization.keywords,
      releaseNotes: newLocalization.releaseNotes
    })
    
    if (response.success) {
      localizations.value.push(response.localization)
      closeAddLocalizationModal()
    }
  } catch (error) {
    console.error('Failed to add localization:', error)
  }
}

async function fetchLocalizations() {
  if (!hasValidAppId.value) return
  try {
    const response = await api.getAppLocalizations(appId.value)
    if (response.success) {
      localizations.value = response.localizations
    }
  } catch (error) {
    console.error('Failed to fetch localizations:', error)
  }
}

async function deleteLocalization(id: number) {
  if (!confirm(t('applocalizations.confirmDelete'))) {
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
    }
  } catch (error) {
    console.error('Failed to delete localization:', error)
  }
}

function editLocalization(localization: AppLocalization) {
  // This would open an edit modal in a real implementation
  console.log('Edit localization:', localization)
}

watch(appId, () => {
  fetchLocalizations()
})

onMounted(() => {
  fetchLocalizations()
  fetchProviderConfigs()
})
</script>
