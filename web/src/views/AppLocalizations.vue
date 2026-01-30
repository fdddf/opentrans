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

            <div>
              <label class="text-sm text-slate-400">
                {{ t('applocalizations.promotionalText') }}
                <span class="text-xs text-slate-500">({{ newLocalization.promotionalText.length }}/170)</span>
              </label>
              <textarea v-model="newLocalization.promotionalText" rows="3" class="mt-1 w-full rounded-lg bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint" :placeholder="t('applocalizations.promotionalTextPlaceholder')"></textarea>
            </div>

            <div>
              <label class="text-sm text-slate-400">
                {{ t('applocalizations.releaseNotes') }}
                <span class="text-xs text-slate-500">({{ newLocalization.releaseNotes.length }}/4000)</span>
              </label>
              <textarea v-model="newLocalization.releaseNotes" rows="4" class="mt-1 w-full rounded-lg bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint" :placeholder="t('applocalizations.releaseNotesPlaceholder')"></textarea>
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

    <!-- Edit Localization Modal -->
    <div v-if="showEditLocalizationModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 px-4">
      <div class="w-full max-w-2xl rounded-2xl border border-white/10 bg-midnight p-6 shadow-xl max-h-[90vh] overflow-y-auto">
        <div class="flex items-center justify-between">
          <h2 class="text-lg font-semibold">{{ t('applocalizations.edit') }}</h2>
          <button class="text-slate-400 hover:text-white" @click="closeEditLocalizationModal">×</button>
        </div>

        <form class="mt-4 space-y-4" @submit.prevent="updateLocalization">
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label class="text-sm text-slate-400">{{ t('applocalizations.language') }}</label>
              <select v-model="editLocalizationData.languageCode" class="mt-1 w-full rounded-lg bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint" disabled>
                <option v-for="lang in availableLanguages" :key="lang.code" :value="lang.code">{{ lang.name }} ({{ lang.code }})</option>
              </select>
            </div>
          </div>

          <div class="space-y-4">
            <div>
              <label class="text-sm text-slate-400">
                {{ t('applocalizations.appName') }}
                <span class="text-xs text-slate-500">({{ editLocalizationData.name.length }}/30)</span>
              </label>
              <input v-model="editLocalizationData.name" class="mt-1 w-full rounded-lg bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint" :placeholder="t('applocalizations.appNamePlaceholder')" :class="{ 'ring-rose-500': getFieldError('name') }" />
              <p v-if="getFieldError('name')" class="mt-1 text-xs text-rose-500">{{ getFieldError('name') }}</p>
            </div>

            <div>
              <label class="text-sm text-slate-400">
                {{ t('applocalizations.subtitle') }}
                <span class="text-xs text-slate-500">({{ editLocalizationData.subtitle.length }}/30)</span>
              </label>
              <input v-model="editLocalizationData.subtitle" class="mt-1 w-full rounded-lg bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint" :placeholder="t('applocalizations.subtitlePlaceholder')" :class="{ 'ring-rose-500': getFieldError('subtitle') }" />
              <p v-if="getFieldError('subtitle')" class="mt-1 text-xs text-rose-500">{{ getFieldError('subtitle') }}</p>
            </div>

            <div>
              <label class="text-sm text-slate-400">
                {{ t('applocalizations.shortDescription') }}
                <span class="text-xs text-slate-500">({{ editLocalizationData.shortDescription.length }}/80)</span>
              </label>
              <textarea v-model="editLocalizationData.shortDescription" rows="2" class="mt-1 w-full rounded-lg bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint" :placeholder="t('applocalizations.shortDescriptionPlaceholder')" :class="{ 'ring-rose-500': getFieldError('shortDescription') }"></textarea>
              <p v-if="getFieldError('shortDescription')" class="mt-1 text-xs text-rose-500">{{ getFieldError('shortDescription') }}</p>
            </div>

            <div>
              <label class="text-sm text-slate-400">{{ t('applocalizations.longDescription') }}</label>
              <textarea v-model="editLocalizationData.longDescription" rows="4" class="mt-1 w-full rounded-lg bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint" :placeholder="t('applocalizations.longDescriptionPlaceholder')"></textarea>
            </div>

            <div>
              <label class="text-sm text-slate-400">
                {{ t('applocalizations.keywords') }}
                <span class="text-xs text-slate-500">({{ editLocalizationData.keywords.length }}/100)</span>
              </label>
              <input v-model="editLocalizationData.keywords" class="mt-1 w-full rounded-lg bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint" :placeholder="t('applocalizations.keywordsPlaceholder')" :class="{ 'ring-rose-500': getFieldError('keywords') }" />
              <p v-if="getFieldError('keywords')" class="mt-1 text-xs text-rose-500">{{ getFieldError('keywords') }}</p>
            </div>

            <div>
              <label class="text-sm text-slate-400">
                {{ t('applocalizations.promotionalText') }}
                <span class="text-xs text-slate-500">({{ editLocalizationData.promotionalText.length }}/170)</span>
              </label>
              <textarea v-model="editLocalizationData.promotionalText" rows="3" class="mt-1 w-full rounded-lg bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint" :placeholder="t('applocalizations.promotionalTextPlaceholder')" :class="{ 'ring-rose-500': getFieldError('promotionalText') }"></textarea>
              <p v-if="getFieldError('promotionalText')" class="mt-1 text-xs text-rose-500">{{ getFieldError('promotionalText') }}</p>
            </div>

            <div>
              <label class="text-sm text-slate-400">
                {{ t('applocalizations.releaseNotes') }}
                <span class="text-xs text-slate-500">({{ editLocalizationData.releaseNotes.length }}/4000)</span>
              </label>
              <textarea v-model="editLocalizationData.releaseNotes" rows="4" class="mt-1 w-full rounded-lg bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint" :placeholder="t('applocalizations.releaseNotesPlaceholder')" :class="{ 'ring-rose-500': getFieldError('releaseNotes') }"></textarea>
              <p v-if="getFieldError('releaseNotes')" class="mt-1 text-xs text-rose-500">{{ getFieldError('releaseNotes') }}</p>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div>
                <label class="text-sm text-slate-400">
                  {{ t('applocalizations.privacyUrl') }}
                  <span class="text-xs text-slate-500">({{ editLocalizationData.privacyUrl.length }}/255)</span>
                </label>
                <input v-model="editLocalizationData.privacyUrl" class="mt-1 w-full rounded-lg bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint" :placeholder="t('applocalizations.privacyUrlPlaceholder')" :class="{ 'ring-rose-500': getFieldError('privacyUrl') }" />
                <p v-if="getFieldError('privacyUrl')" class="mt-1 text-xs text-rose-500">{{ getFieldError('privacyUrl') }}</p>
              </div>
              <div>
                <label class="text-sm text-slate-400">
                  {{ t('applocalizations.marketingUrl') }}
                  <span class="text-xs text-slate-500">({{ editLocalizationData.marketingUrl.length }}/255)</span>
                </label>
                <input v-model="editLocalizationData.marketingUrl" class="mt-1 w-full rounded-lg bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint" :placeholder="t('applocalizations.marketingUrlPlaceholder')" :class="{ 'ring-rose-500': getFieldError('marketingUrl') }" />
                <p v-if="getFieldError('marketingUrl')" class="mt-1 text-xs text-rose-500">{{ getFieldError('marketingUrl') }}</p>
              </div>
              <div>
                <label class="text-sm text-slate-400">
                  {{ t('applocalizations.supportUrl') }}
                  <span class="text-xs text-slate-500">({{ editLocalizationData.supportUrl.length }}/255)</span>
                </label>
                <input v-model="editLocalizationData.supportUrl" class="mt-1 w-full rounded-lg bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint" :placeholder="t('applocalizations.supportUrlPlaceholder')" :class="{ 'ring-rose-500': getFieldError('supportUrl') }" />
                <p v-if="getFieldError('supportUrl')" class="mt-1 text-xs text-rose-500">{{ getFieldError('supportUrl') }}</p>
              </div>
            </div>
          </div>

          <div class="flex justify-end gap-2 pt-4">
            <button type="button" class="rounded-lg border border-white/20 px-3 py-2 text-sm hover:border-mint/60 hover:text-mint" @click="closeEditLocalizationModal">{{ t('common.cancel') || 'Cancel' }}</button>
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
          
          <!-- Sync Status and Timestamps -->
          <div class="flex flex-wrap gap-3 pt-2 border-t border-white/10 mt-2">
            <div class="flex items-center gap-1 text-xs">
              <span :class="{
                'w-2 h-2 rounded-full': true,
                'bg-green-500': loc.syncStatus === 'synced',
                'bg-yellow-500': loc.syncStatus === 'pending',
                'bg-red-500': loc.syncStatus === 'failed',
                'bg-gray-500': !loc.syncStatus
              }"></span>
              <span class="text-slate-400">{{ t('applocalizations.syncStatus') }}: {{ getSyncStatusText(loc.syncStatus) }}</span>
            </div>
            
            <div v-if="loc.syncedAt" class="text-xs text-slate-400">
              {{ t('applocalizations.lastSynced') }}: {{ formatDate(loc.syncedAt) }}
            </div>
            
            <div class="text-xs text-slate-400">
              {{ t('applocalizations.source') }}: {{ getSourceText(loc.source) }}
            </div>
          </div>
        </div>
        
        <!-- Sync Actions -->
        <div class="flex justify-end gap-2 pt-3 mt-3 border-t border-white/10">
          <button class="rounded border border-white/20 px-2 py-1 text-xs hover:border-mint/60 hover:text-mint" @click="pullFromApple(loc)">{{ t('applocalizations.pullFromApple') }}</button>
          <button class="rounded border border-white/20 px-2 py-1 text-xs hover:border-mint/60 hover:text-mint" @click="pushToApple(loc)">{{ t('applocalizations.pushToApple') }}</button>
          <button class="rounded border border-white/20 px-2 py-1 text-xs hover:border-mint/60 hover:text-mint" @click="editLocalization(loc)">{{ t('common.edit') }}</button>
          <button class="rounded border border-white/20 px-2 py-1 text-xs hover:border-rose-600/60 hover:text-rose-500" @click="deleteLocalization(loc.id)">{{ t('common.delete') }}</button>
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
const showEditLocalizationModal = ref(false)
const showSyncModal = ref(false)
const syncing = ref(false)
const syncResult = ref<{ message: string; count?: number } | null>(null)
const selectedConfigId = ref<number | null>(null)
const appleConnectConfigs = ref<ProviderConfig[]>([])
const validationErrors = ref<Record<string, string>>({})

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
  releaseNotes: '',
  promotionalText: ''
})

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
  shortDescription: { maxLength: 80, required: false },
  longDescription: { maxLength: 4000, required: false },
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
  newLocalization.promotionalText = ''
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
  if (!validateLocalization(newLocalization)) {
    return
  }

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
      downloadDescription: '',
      shortDescription: newLocalization.shortDescription,
      longDescription: newLocalization.longDescription,
      keywords: newLocalization.keywords,
      releaseNotes: newLocalization.releaseNotes,
      promotionalText: newLocalization.promotionalText
    })

    if (response.success) {
      localizations.value.push(response.localization)
      closeAddLocalizationModal()
    }
  } catch (error) {
    console.error('Failed to add localization:', error)
    alert('Failed to add localization. Please try again.')
  }
}

function closeEditLocalizationModal() {
  showEditLocalizationModal.value = false
  clearValidationErrors()
}

async function updateLocalization() {
  if (!validateLocalization(editLocalizationData)) {
    return
  }

  try {
    if (!hasValidAppId.value) {
      throw new Error('Invalid app ID')
    }

    const response = await api.updateAppLocalization(appId.value, editLocalizationData.languageCode, {
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

    if (response.success) {
      await fetchLocalizations()
      closeEditLocalizationModal()
    } else {
      alert(response.message || 'Failed to update localization')
    }
  } catch (error) {
    console.error('Failed to update localization:', error)
    alert('Failed to update localization. Please try again.')
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
  // Populate edit form with existing data
  editLocalizationData.id = localization.id
  editLocalizationData.languageCode = localization.languageCode
  editLocalizationData.name = localization.name || ''
  editLocalizationData.subtitle = localization.subtitle || ''
  editLocalizationData.privacyUrl = localization.privacyUrl || ''
  editLocalizationData.marketingUrl = localization.marketingUrl || ''
  editLocalizationData.supportUrl = localization.supportUrl || ''
  editLocalizationData.shortDescription = localization.shortDescription || ''
  editLocalizationData.longDescription = localization.longDescription || ''
  editLocalizationData.keywords = localization.keywords || ''
  editLocalizationData.releaseNotes = localization.releaseNotes || ''
  editLocalizationData.promotionalText = localization.promotionalText || ''

  clearValidationErrors()
  showEditLocalizationModal.value = true
}

watch(appId, () => {
  fetchLocalizations()
})

onMounted(() => {
  fetchLocalizations()
  fetchProviderConfigs()
})
</script>
