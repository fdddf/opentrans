<template>
  <div class="space-y-6">
    <header class="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
      <div>
        <p class="text-xs uppercase tracking-[0.3em] text-slate-500">{{ t('nav.applocalizations') }}</p>
        <h1 class="text-2xl font-semibold">{{ t('applocalizations.title') }}</h1>
      </div>
      <div class="flex items-center gap-2">
        <button 
          class="rounded-lg border border-white/20 px-3 py-2 text-sm hover:border-mint/60 hover:text-mint" 
          @click="showSyncModal = true"
        >
          {{ t('applocalizations.sync') }}
        </button>
      </div>
    </header>

    <!-- Main Content with Two Columns -->
    <div class="flex gap-6">
      <!-- Left: Language Sidebar -->
      <aside class="w-80 flex-shrink-0">
        <div class="sticky top-6 rounded-2xl border border-white/10 bg-white/5 p-4">
          <h2 class="mb-4 text-sm font-semibold text-slate-300">{{ t('applocalizations.languages') }}</h2>
          
          <!-- Language Dropdown for Adding -->
          <div class="mb-4">
            <select
              v-model="selectedLanguageToAdd"
              class="w-full rounded-lg bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint"
              @change="handleLanguageSelect"
            >
              <option value="">{{ t('applocalizations.selectLanguage') }}</option>
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
              {{ t('applocalizations.noLanguages') }}
            </div>
          </div>
        </div>
      </aside>

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

          <div>
            <p class="text-sm text-slate-400 mb-3">{{ t('applocalizations.syncDirection') }}</p>
            <select
              v-model="syncDirection"
              class="w-full rounded-lg bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint"
            >
              <option value="pull">{{ t('applocalizations.pullFromApple') }}</option>
              <option value="push">{{ t('applocalizations.pushToApple') }}</option>
              <option value="both">{{ t('applocalizations.bidirectional') }}</option>
            </select>
          </div>

          <div>
            <p class="text-sm text-slate-400 mb-3">{{ t('applocalizations.conflictStrategy') }}</p>
            <select
              v-model="syncStrategy"
              class="w-full rounded-lg bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint"
            >
              <option value="apple_first">{{ t('applocalizations.appleFirst') }}</option>
              <option value="local_first">{{ t('applocalizations.localFirst') }}</option>
              <option value="manual">{{ t('applocalizations.manual') }}</option>
            </select>
          </div>

          <div>
            <p class="text-sm text-slate-400 mb-3">{{ t('applocalizations.selectLanguages') }} ({{ t('applocalizations.optional') }})</p>
            <div class="max-h-40 overflow-y-auto rounded-lg bg-white/5 border border-white/10 p-2 space-y-1">
              <label v-for="lang in availableLanguages" :key="lang.code" class="flex items-center gap-2 text-sm cursor-pointer hover:bg-white/5 p-1 rounded">
                <input
                  type="checkbox"
                  :value="lang.code"
                  v-model="selectedSyncLanguages"
                  class="rounded bg-white/10 border-white/20 text-mint focus:ring-mint"
                />
                <span>{{ lang.name }} ({{ lang.code }})</span>
              </label>
            </div>
            <p class="text-xs text-slate-500 mt-1">{{ t('applocalizations.selectLanguagesDesc') }}</p>
          </div>

          <div v-if="syncResult" class="p-3 rounded-lg bg-white/5 border border-white/10">
            <p class="text-sm">{{ syncResult.message }}</p>
            <p v-if="syncResult.count" class="text-sm text-mint">{{ t('applocalizations.syncedCount', { count: syncResult.count }) }}</p>
            <div v-if="syncResult.conflicts && syncResult.conflicts.length > 0" class="mt-2 p-2 rounded bg-yellow-500/10 border border-yellow-500/20">
              <p class="text-xs text-yellow-500">{{ t('applocalizations.conflictsDetected', { count: syncResult.conflicts.length }) }}</p>
            </div>
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
                  <span v-if="selectedLocalization.id < 0" class="text-sm font-normal text-orange-400 ml-2">{{ t('applocalizations.add') }}</span>
                </h2>
              </div>
            </div>
            <button class="text-slate-400 hover:text-rose-500" @click="deleteLocalization(selectedLocalization.id)">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
              </svg>
            </button>
          </div>

          <!-- Sync Actions -->
          <div v-if="selectedLocalization.id > 0" class="flex gap-2 mb-6">
            <button class="rounded-lg border border-white/20 px-3 py-2 text-sm hover:border-mint/60 hover:text-mint" @click="pullFromApple(selectedLocalization)">
              {{ t('applocalizations.pullFromApple') }}
            </button>
            <button class="rounded-lg border border-white/20 px-3 py-2 text-sm hover:border-mint/60 hover:text-mint" @click="pushToApple(selectedLocalization)">
              {{ t('applocalizations.pushToApple') }}
            </button>
          </div>

          <!-- Sync Info -->
          <div v-if="selectedLocalization.id > 0" class="flex flex-wrap gap-4 text-xs text-slate-400 mb-6 pb-4 border-b border-white/10">
            <div>{{ t('applocalizations.syncStatus') }}: {{ getSyncStatusText(selectedLocalization.syncStatus) }}</div>
            <div>{{ t('applocalizations.version') }}: {{ selectedLocalization.version || '-' }}</div>
            <div>{{ t('applocalizations.lastSynced') }}: {{ formatDate(selectedLocalization.syncedAt) }}</div>
          </div>

          <!-- Edit Form -->
          <form class="space-y-4" @submit.prevent="updateLocalization">
            <div>
              <label class="text-sm text-slate-400">{{ t('applocalizations.appName') }}</label>
              <input v-model="editLocalizationData.name" class="mt-1 w-full rounded-lg bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint" :placeholder="t('applocalizations.appNamePlaceholder')" />
            </div>

            <div>
              <label class="text-sm text-slate-400">{{ t('applocalizations.subtitle') }}</label>
              <input v-model="editLocalizationData.subtitle" class="mt-1 w-full rounded-lg bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint" :placeholder="t('applocalizations.subtitlePlaceholder')" />
            </div>

            <div>
              <label class="text-sm text-slate-400">{{ t('applocalizations.description') }}</label>
              <textarea v-model="editLocalizationData.longDescription" rows="6" class="mt-1 w-full rounded-lg bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint" :placeholder="t('applocalizations.descriptionPlaceholder')"></textarea>
            </div>

            <div>
              <label class="text-sm text-slate-400">{{ t('applocalizations.keywords') }}</label>
              <input v-model="editLocalizationData.keywords" class="mt-1 w-full rounded-lg bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint" :placeholder="t('applocalizations.keywordsPlaceholder')" />
            </div>

            <div>
              <label class="text-sm text-slate-400">{{ t('applocalizations.promotionalText') }}</label>
              <textarea v-model="editLocalizationData.promotionalText" rows="3" class="mt-1 w-full rounded-lg bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint" :placeholder="t('applocalizations.promotionalTextPlaceholder')"></textarea>
            </div>

            <div>
              <label class="text-sm text-slate-400">{{ t('applocalizations.releaseNotes') }}</label>
              <textarea v-model="editLocalizationData.releaseNotes" rows="4" class="mt-1 w-full rounded-lg bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint" :placeholder="t('applocalizations.releaseNotesPlaceholder')"></textarea>
            </div>

            <div class="grid grid-cols-1 gap-4">
              <div>
                <label class="text-sm text-slate-400">{{ t('applocalizations.privacyUrl') }}</label>
                <input v-model="editLocalizationData.privacyUrl" class="mt-1 w-full rounded-lg bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint" :placeholder="t('applocalizations.privacyUrlPlaceholder')" />
              </div>
              <div>
                <label class="text-sm text-slate-400">{{ t('applocalizations.marketingUrl') }}</label>
                <input v-model="editLocalizationData.marketingUrl" class="mt-1 w-full rounded-lg bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint" :placeholder="t('applocalizations.marketingUrlPlaceholder')" />
              </div>
              <div>
                <label class="text-sm text-slate-400">{{ t('applocalizations.supportUrl') }}</label>
                <input v-model="editLocalizationData.supportUrl" class="mt-1 w-full rounded-lg bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint" :placeholder="t('applocalizations.supportUrlPlaceholder')" />
              </div>
            </div>

            <div class="flex justify-end gap-2 pt-4">
              <button type="submit" class="rounded-lg bg-mint px-4 py-2 text-sm font-semibold text-midnight shadow hover:bg-mint/90">
                {{ t('common.save') }}
              </button>
            </div>
          </form>
        </div>

        <!-- Empty State when no language selected -->
        <div v-else-if="localizations.length === 0" class="rounded-2xl border border-white/10 bg-white/5 p-12 text-center">
          <div class="text-6xl mb-4">🌐</div>
          <h3 class="text-lg font-semibold mb-2">{{ t('applocalizations.noLocalizations') }}</h3>
          <p class="text-sm text-slate-400">{{ t('applocalizations.noLocalizationsDesc') }}</p>
        </div>

        <!-- Select a language prompt -->
        <div v-else class="rounded-2xl border border-white/10 bg-white/5 p-12 text-center">
          <div class="text-6xl mb-4">👈</div>
          <h3 class="text-lg font-semibold mb-2">{{ t('applocalizations.selectLanguage') }}</h3>
          <p class="text-sm text-slate-400">{{ t('applocalizations.selectLanguageDesc') }}</p>
        </div>
      </main>
    </div>
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
const app = ref<{ id: number; primaryLocale: string } | null>(null)
const selectedLocalizationId = ref<number | null>(null)
const selectedLanguageToAdd = ref<string>('')
const isNewLocalization = ref(false)
const showSyncModal = ref(false)
const syncing = ref(false)
const syncResult = ref<{ message: string; count?: number; conflicts?: any[] } | null>(null)
const selectedConfigId = ref<number | null>(null)
const syncDirection = ref('pull')
const syncStrategy = ref('apple_first')
const selectedSyncLanguages = ref<string[]>([])
const appleConnectConfigs = ref<ProviderConfig[]>([])
const validationErrors = ref<Record<string, string>>({})
const availableLanguages = ref<{ code: string; name: string; native_name: string; region?: string; direction: string }[]>([])
const loadingLanguages = ref(false)

// Computed property for selected localization
const selectedLocalization = computed(() => {
  return localizations.value.find(loc => loc.id === selectedLocalizationId.value) || null
})

// Language code to flag emoji mapping
const languageFlagMap: Record<string, string> = {
  'en': '🇺🇸',
  'en-US': '🇺🇸',
  'en-GB': '🇬🇧',
  'zh': '🇨🇳',
  'zh-Hans': '🇨🇳',
  'zh-Hant': '🇹🇼',
  'zh-TW': '🇹🇼',
  'zh-HK': '🇭🇰',
  'ja': '🇯🇵',
  'ko': '🇰🇷',
  'es': '🇪🇸',
  'fr': '🇫🇷',
  'de': '🇩🇪',
  'it': '🇮🇹',
  'pt': '🇵🇹',
  'pt-BR': '🇧🇷',
  'ru': '🇷🇺',
  'ar': '🇸🇦',
  'hi': '🇮🇳',
  'th': '🇹🇭',
  'vi': '🇻🇳',
  'id': '🇮🇩',
  'ms': '🇲🇾',
  'nl': '🇳🇱',
  'pl': '🇵🇱',
  'tr': '🇹🇷',
  'uk': '🇺🇦',
  'cs': '🇨🇿',
  'da': '🇩🇰',
  'fi': '🇫🇮',
  'no': '🇳🇴',
  'sv': '🇸🇪',
  'el': '🇬🇷',
  'he': '🇮🇱',
  'ro': '🇷🇴',
  'hu': '🇭🇺',
  'sk': '🇸🇰',
  'bg': '🇧🇬',
  'hr': '🇭🇷',
  'ca': '🇪🇸',
  'fil': '🇵🇭',
  'ms-MY': '🇲🇾',
  'pt-PT': '🇵🇹',
  'es-MX': '🇲🇽',
  'es-ES': '🇪🇸',
  'fr-CA': '🇨🇦',
  'fr-FR': '🇫🇷',
  'de-DE': '🇩🇪',
  'it-IT': '🇮🇹',
  'nl-NL': '🇳🇱',
  'en-CA': '🇨🇦',
  'en-AU': '🇦🇺',
  'en-IN': '🇮🇳',
  'en-SG': '🇸🇬',
  'ja-JP': '🇯🇵',
  'ko-KR': '🇰🇷',
  'zh-CN': '🇨🇳',
  'zh-SG': '🇸🇬',
  'ar-SA': '🇸🇦',
  'ar-AE': '🇦🇪',
  'ar-EG': '🇪🇬',
  'tr-TR': '🇹🇷',
  'pl-PL': '🇵🇱',
  'ru-RU': '🇷🇺',
  'uk-UA': '🇺🇦',
  'th-TH': '🇹🇭',
  'vi-VN': '🇻🇳',
  'id-ID': '🇮🇩',
  'hi-IN': '🇮🇳',
  'bn-IN': '🇮🇳',
  'ta-IN': '🇮🇳',
  'te-IN': '🇮🇳',
  'mr-IN': '🇮🇳',
  'gu-IN': '🇮🇳',
  'kn-IN': '🇮🇳',
  'ml-IN': '🇮🇳',
  'pa-IN': '🇮🇳',
  'or-IN': '🇮🇳',
  'as-IN': '🇮🇳',
  'fa': '🇮🇷',
  'fa-IR': '🇮🇷',
  'ur': '🇵🇰',
  'ur-PK': '🇵🇰',
  'my': '🇲🇲',
  'my-MM': '🇲🇲',
  'km': '🇰🇭',
  'lo': '🇱🇦',
  'ne': '🇳🇵',
  'ne-NP': '🇳🇵',
  'si': '🇱🇰',
  'si-LK': '🇱🇰',
  'sw': '🇰🇪',
  'sw-KE': '🇰🇪',
  'af': '🇿🇦',
  'af-ZA': '🇿🇦',
  'zu': '🇿🇦',
  'zu-ZA': '🇿🇦',
  'xh': '🇿🇦',
  'xh-ZA': '🇿🇦',
  'is': '🇮🇸',
  'is-IS': '🇮🇸',
  'ga': '🇮🇪',
  'ga-IE': '🇮🇪',
  'cy': '🇬🇧',
  'cy-GB': '🇬🇧',
  'eu': '🇪🇸',
  'eu-ES': '🇪🇸',
  'gl': '🇪🇸',
  'gl-ES': '🇪🇸',
  'sq': '🇦🇱',
  'sq-AL': '🇦🇱',
  'mk': '🇲🇰',
  'mk-MK': '🇲🇰',
  'sr': '🇷🇸',
  'sr-RS': '🇷🇸',
  'bs': '🇧🇦',
  'bs-BA': '🇧🇦',
  'me': '🇲🇪',
  'me-ME': '🇲🇪',
  'sl': '🇸🇮',
  'sl-SI': '🇸🇮',
  'hr-HR': '🇭🇷',
  'mt': '🇲🇹',
  'mt-MT': '🇲🇹',
  'lb': '🇱🇺',
  'lb-LU': '🇱🇺',
  'lt': '🇱🇹',
  'lt-LT': '🇱🇹',
  'lv': '🇱🇻',
  'lv-LV': '🇱🇻',
  'et': '🇪🇪',
  'et-EE': '🇪🇪',
  'kk': '🇰🇿',
  'kk-KZ': '🇰🇿',
  'ky': '🇰🇬',
  'ky-KG': '🇰🇬',
  'uz': '🇺🇿',
  'uz-UZ': '🇺🇿',
  'tg': '🇹🇯',
  'tg-TJ': '🇹🇯',
  'hy': '🇦🇲',
  'hy-AM': '🇦🇲',
  'az': '🇦🇿',
  'az-AZ': '🇦🇿',
  'ka': '🇬🇪',
  'ka-GE': '🇬🇪',
  'mn': '🇲🇳',
  'mn-MN': '🇲🇳',
  'bo': '🇨🇳',
  'bo-CN': '🇨🇳',
  'dz': '🇧🇹',
  'dz-BT': '🇧🇹'
}

function getLanguageFlag(code: string): string {
  // First try exact match
  if (languageFlagMap[code]) {
    return languageFlagMap[code]
  }
  
  // Try to extract the base language code (e.g., 'en-US' -> 'en')
  const baseCode = code.split('-')[0]
  if (languageFlagMap[baseCode]) {
    return languageFlagMap[baseCode]
  }
  
  // Default flag for unknown languages
  return '🌐'
}

const editLocalizationData = reactive({
  id: 0,
  languageCode: '',
  name: '',
  subtitle: '',
  privacyUrl: '',
  marketingUrl: '',
  supportUrl: '',
  shortDescription: '',
  longDescription: '',
  keywords: '',
  releaseNotes: '',
  promotionalText: ''
})

// Validation rules based on Apple App Store Connect requirements
const validationRules = {
  name: { maxLength: 30, required: true },
  subtitle: { maxLength: 30, required: false },
  description: { maxLength: 4000, required: false },
  keywords: { maxLength: 100, required: false },
  privacyUrl: { maxLength: 255, required: false },
  marketingUrl: { maxLength: 255, required: false },
  supportUrl: { maxLength: 255, required: false },
  promotionalText: { maxLength: 170, required: false },
  releaseNotes: { maxLength: 4000, required: false }
}

function validateField(field: string, value: string): string | null {
  const rule = validationRules[field as keyof typeof validationRules]
  if (!rule) return null

  if (rule.required && !value.trim()) {
    return 'This field is required'
  }

  if (rule.maxLength && value.length > rule.maxLength) {
    return `Must not exceed ${rule.maxLength} characters`
  }

  return null
}

function validateLocalization(data: typeof newLocalization): boolean {
  const errors: Record<string, string> = {}
  let isValid = true

  for (const [field, value] of Object.entries(data)) {
    const error = validateField(field, value)
    if (error) {
      errors[field] = error
      isValid = false
    }
  }

  validationErrors.value = errors
  return isValid
}

function getFieldError(field: string): string | null {
  return validationErrors.value[field] || null
}

function clearValidationErrors() {
  validationErrors.value = {}
}

function getLanguageName(code: string): string {
  const lang = availableLanguages.value.find(l => l.code === code)
  return lang ? lang.name : code
}

function closeSyncModal() {
  showSyncModal.value = false
  selectedConfigId.value = null
  syncDirection.value = 'pull'
  syncStrategy.value = 'apple_first'
  selectedSyncLanguages.value = []
  syncResult.value = null
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

function getSyncStatusText(status?: string): string {
  if (!status) return 'Unknown';
  const statusMap: Record<string, string> = {
    'synced': 'Synced',
    'pending': 'Pending',
    'failed': 'Failed'
  };
  return statusMap[status] || status || 'Unknown';
}

function getSourceText(source?: string): string {
  if (!source) return 'Unknown';
  const sourceMap: Record<string, string> = {
    'apple': 'Apple',
    'local': 'Local'
  };
  return sourceMap[source] || source || 'Unknown';
}

function getVersionStateText(state?: string): string {
  if (!state) return 'Unknown';
  const stateMap: Record<string, string> = {
    'PREPARE_FOR_SUBMISSION': 'Prepare for Submission',
    'WAITING_FOR_REVIEW': 'Waiting for Review',
    'IN_REVIEW': 'In Review',
    'PENDING_DEVELOPER_RELEASE': 'Pending Developer Release',
    'PENDING_APPLE_RELEASE': 'Pending Apple Release',
    'READY_FOR_SALE': 'Ready for Sale',
    'REJECTED': 'Rejected',
    'REMOVED_FROM_SALE': 'Removed from Sale',
    'DEVELOPER_REJECTED': 'Developer Rejected',
    'METADATA_REJECTED': 'Metadata Rejected'
  };
  return stateMap[state] || state || 'Unknown';
}

function isLocalizationEditable(loc: AppLocalization): boolean {
  // Only allow editing if:
  // 1. Source is local (manually created)
  // 2. Or source is apple but version state is not READY_FOR_SALE
  if (loc.source === 'local') return true;
  if (loc.source === 'apple') {
    return loc.versionState !== 'READY_FOR_SALE';
  }
  return true;
}

function formatDate(dateString?: string): string {
  if (!dateString) return 'N/A';
  const date = new Date(dateString);
  return date.toLocaleDateString() + ' ' + date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
}

async function pullFromApple(localization: AppLocalization) {
  if (!confirm(`Pull localization for ${localization.languageCode} from Apple?`)) {
    return;
  }
  
  if (!selectedConfigId.value || !hasValidAppId.value) {
    alert('Please select an Apple Connect configuration first.');
    return;
  }
  
  try {
    const response = await api.syncAppleAppLocalizations(appId.value, {
      configId: selectedConfigId.value
    });
    
    if (response.success) {
      await fetchLocalizations();
      alert(`Successfully pulled localization for ${localization.languageCode} from Apple.`);
    } else {
      alert(`Failed to pull localization: ${response.message || 'Unknown error'}`);
    }
  } catch (error) {
    console.error('Failed to pull from Apple:', error);
    alert('Failed to pull localization from Apple. See console for details.');
  }
}

async function pushToApple(localization: AppLocalization) {
  if (!confirm(`Push localization for ${localization.languageCode} to Apple?`)) {
    return;
  }
  
  if (!selectedConfigId.value || !hasValidAppId.value) {
    alert('Please select an Apple Connect configuration first.');
    return;
  }
  
  try {
    const response = await api.syncAppToApple(appId.value, {
      configId: selectedConfigId.value
    });
    
    if (response.success) {
      await fetchLocalizations();
      alert(`Successfully pushed localization for ${localization.languageCode} to Apple.`);
    } else {
      alert(`Failed to push localization: ${response.message || 'Unknown error'}`);
    }
  } catch (error) {
    console.error('Failed to push to Apple:', error);
    alert('Failed to push localization to Apple. See console for details.');
  }
}

async function fetchProviderConfigs() {
  try {
    const response = await api.getAppleConnectConfigs()
    if (response.success) {
      appleConnectConfigs.value = response.data || []
    }
  } catch (error) {
    console.error('Failed to fetch Apple Connect configs:', error)
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
      configId: selectedConfigId.value,
      languageCodes: selectedSyncLanguages.value.length > 0 ? selectedSyncLanguages.value : undefined,
      direction: syncDirection.value,
      strategy: syncStrategy.value
    })

    if (response.success) {
      syncResult.value = {
        message: response.message || t('applocalizations.syncSuccess'),
        count: response.count,
        conflicts: response.conflicts
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

async function updateLocalization() {
  if (!validateLocalization(editLocalizationData)) {
    return
  }

  try {
    if (!hasValidAppId.value) {
      throw new Error('Invalid app ID')
    }

    let response
    if (isNewLocalization.value) {
      // Create new localization
      response = await api.createAppLocalization(appId.value, {
        languageCode: editLocalizationData.languageCode,
        name: editLocalizationData.name,
        subtitle: editLocalizationData.subtitle,
        privacyUrl: editLocalizationData.privacyUrl,
        marketingUrl: editLocalizationData.marketingUrl,
        supportUrl: editLocalizationData.supportUrl,
        downloadDescription: '',
        shortDescription: editLocalizationData.shortDescription,
        longDescription: editLocalizationData.longDescription,
        keywords: editLocalizationData.keywords,
        releaseNotes: editLocalizationData.releaseNotes,
        promotionalText: editLocalizationData.promotionalText
      })
    } else {
      // Update existing localization
      response = await api.updateAppLocalization(appId.value, editLocalizationData.languageCode, {
        name: editLocalizationData.name,
        subtitle: editLocalizationData.subtitle,
        privacyUrl: editLocalizationData.privacyUrl,
        marketingUrl: editLocalizationData.marketingUrl,
        supportUrl: editLocalizationData.supportUrl,
        shortDescription: editLocalizationData.shortDescription,
        longDescription: editLocalizationData.longDescription,
        keywords: editLocalizationData.keywords,
        releaseNotes: editLocalizationData.releaseNotes,
        promotionalText: editLocalizationData.promotionalText
      })
    }

    if (response.success) {
      isNewLocalization.value = false
      await fetchLocalizations()
      // Select the newly created/updated localization
      if (response.localization) {
        selectedLocalizationId.value = response.localization.id
      }
    } else {
      alert(response.message || (isNewLocalization.value ? 'Failed to add localization' : 'Failed to update localization'))
    }
  } catch (error) {
    console.error('Failed to save localization:', error)
    alert(isNewLocalization.value ? 'Failed to add localization. Please try again.' : 'Failed to update localization. Please try again.')
  }
}

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
      localizations.value = response.localizations
    }
  } catch (error) {
    console.error('Failed to fetch localizations:', error)
  }
}

async function deleteLocalization(id: number) {
  const localization = localizations.value.find(loc => loc.id === id)
  
  // Check if it's the primary language
  if (localization && app.value?.primaryLocale === localization.languageCode) {
    alert('Cannot delete the primary language.')
    return
  }
  
  if (!confirm(t('applocalizations.confirmDelete'))) {
    return
  }
  
  // If it's a new unsaved localization (negative ID), just remove from list
  if (id < 0) {
    localizations.value = localizations.value.filter(loc => loc.id !== id)
    selectedLocalizationId.value = null
    isNewLocalization.value = false
    return
  }
  
  // For existing localizations, delete via API
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
  // Populate edit form with existing data
  editLocalizationData.id = localization.id
  editLocalizationData.languageCode = localization.languageCode
  editLocalizationData.name = localization.name || ''
  editLocalizationData.subtitle = localization.subtitle || ''
  editLocalizationData.privacyUrl = localization.privacyUrl || ''
  editLocalizationData.marketingUrl = localization.marketingUrl || ''
  editLocalizationData.supportUrl = localization.supportUrl || ''
  editLocalizationData.shortDescription = localization.shortDescription || ''
  editLocalizationData.longDescription = localization.longDescription || localization.description || ''
  editLocalizationData.keywords = localization.keywords || ''
  editLocalizationData.releaseNotes = localization.releaseNotes || ''
  editLocalizationData.promotionalText = localization.promotionalText || ''

  clearValidationErrors()
}

async function handleLanguageSelect() {
  if (!selectedLanguageToAdd.value) return

  // Check if language already exists
  const existing = localizations.value.find(loc => loc.languageCode === selectedLanguageToAdd.value)
  if (existing) {
    // Select existing language
    selectedLocalizationId.value = existing.id
    selectedLanguageToAdd.value = ''
    return
  }

  // Prepare new localization form
  editLocalizationData.id = 0
  editLocalizationData.languageCode = selectedLanguageToAdd.value
  editLocalizationData.name = ''
  editLocalizationData.subtitle = ''
  editLocalizationData.privacyUrl = ''
  editLocalizationData.marketingUrl = ''
  editLocalizationData.supportUrl = ''
  editLocalizationData.shortDescription = ''
  editLocalizationData.longDescription = ''
  editLocalizationData.keywords = ''
  editLocalizationData.releaseNotes = ''
  editLocalizationData.promotionalText = ''

  clearValidationErrors()
  isNewLocalization.value = true
  
  // Add temporary entry to list
  const tempId = -Date.now()
  const tempLocalization: AppLocalization = {
    id: tempId,
    languageCode: selectedLanguageToAdd.value,
    name: '',
    subtitle: '',
    description: '',
    shortDescription: '',
    longDescription: '',
    keywords: '',
    privacyUrl: '',
    marketingUrl: '',
    supportUrl: '',
    promotionalText: '',
    releaseNotes: '',
    syncStatus: 'pending',
    version: '',
    versionState: '',
    syncedAt: '',
    source: 'local'
  }
  localizations.value.push(tempLocalization)
  selectedLocalizationId.value = tempId
  selectedLanguageToAdd.value = ''
}

watch(appId, () => {
  fetchLocalizations()
})

async function fetchLanguages() {
  loadingLanguages.value = true
  try {
    const response = await api.getSupportedLanguages()
    if (response.success) {
      availableLanguages.value = response.languages
    }
  } catch (error) {
    console.error('Failed to fetch languages:', error)
  } finally {
    loadingLanguages.value = false
  }
}

onMounted(() => {
  fetchApp()
  fetchLanguages()
  fetchLocalizations()
  fetchProviderConfigs()
})

// Watch for selection changes to populate edit form
watch(selectedLocalizationId, (newId) => {
  if (newId) {
    const loc = localizations.value.find(l => l.id === newId)
    if (loc) {
      editLocalization(loc)
      isNewLocalization.value = loc.id < 0
    }
  }
})
</script>
