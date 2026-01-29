<template>
  <div class="space-y-4">
    <div class="flex flex-wrap items-center gap-2 rounded-2xl border border-white/10 bg-white/5 px-4 py-3">
      <button class="flex items-center gap-2 rounded-lg border border-white/20 px-3 py-2 text-xs hover:border-mint/60 hover:text-mint" @click="showAddLanguageModal = true">
        <span>➕</span>
        <span>{{ t('workspace.addLang') }}</span>
      </button>
      <button class="flex items-center gap-2 rounded-lg border border-white/20 px-3 py-2 text-xs hover:border-mint/60 hover:text-mint" @click="translateAll">
        <span>⚡</span>
        <span>{{ t('workspace.translateAll') }}</span>
      </button>
      <button class="flex items-center gap-2 rounded-lg border border-white/20 px-3 py-2 text-xs hover:border-mint/60 hover:text-mint" @click="showSyncToAppleModal = true">
        <span>🔄</span>
        <span>{{ t('workspace.syncToApple') }}</span>
      </button>
    </div>

    <!-- Add Language Modal -->
    <div v-if="showAddLanguageModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 px-4">
      <div class="w-full max-w-lg rounded-2xl border border-white/10 bg-midnight p-6 shadow-xl">
        <div class="flex items-center justify-between">
          <h2 class="text-lg font-semibold">{{ t('workspace.addLang') }}</h2>
          <button class="text-slate-400 hover:text-white" @click="closeAddLanguageModal">×</button>
        </div>

        <div class="mt-4 space-y-4">
          <div>
            <p class="text-sm text-slate-400 mb-3">{{ t('workspace.selectLanguage') }}</p>
            <select 
              v-model="selectedAddLanguage" 
              class="w-full rounded-lg bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint" 
              required
            >
              <option value="">{{ t('workspace.chooseLanguage') }}</option>
              <option v-for="lang in availableLanguages" :key="lang.code" :value="lang.code">
                {{ lang.name }} ({{ lang.code }})
              </option>
            </select>
          </div>

          <div class="flex justify-end gap-2 pt-4">
            <button 
              type="button" 
              class="rounded-lg border border-white/20 px-3 py-2 text-sm hover:border-mint/60 hover:text-mint" 
              @click="closeAddLanguageModal"
            >
              {{ t('common.cancel') || 'Cancel' }}
            </button>
            <button 
              type="button" 
              class="rounded-lg bg-mint px-3 py-2 text-sm font-semibold text-midnight shadow" 
              @click="addLanguage"
              :disabled="!selectedAddLanguage"
            >
              {{ t('workspace.add') }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Sync to Apple Connect Modal -->
    <div v-if="showSyncToAppleModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 px-4">
      <div class="w-full max-w-lg rounded-2xl border border-white/10 bg-midnight p-6 shadow-xl">
        <div class="flex items-center justify-between">
          <h2 class="text-lg font-semibold">{{ t('workspace.syncToApple') }}</h2>
          <button class="text-slate-400 hover:text-white" @click="closeSyncToAppleModal">×</button>
        </div>

        <div class="mt-4 space-y-4" v-if="!hasAppleConnectConfig">
          <div class="p-4 rounded-lg bg-amber-900/20 border border-amber-500/30">
            <div class="flex items-start">
              <span class="text-amber-400">⚠</span>
              <div class="ml-3">
                <p class="text-sm font-medium text-amber-100">{{ t('workspace.noConfig.title') }}</p>
                <p class="text-sm mt-1 text-amber-200">{{ t('workspace.noConfig.message') }}</p>
              </div>
            </div>
            <div class="mt-4 flex justify-end">
              <RouterLink to="/apple-connect-config" class="rounded-lg bg-amber-600 px-3 py-2 text-sm font-semibold text-white">
                {{ t('workspace.setupConfig') }}
              </RouterLink>
            </div>
          </div>
        </div>
        <div class="mt-4 space-y-4" v-else>
          <div>
            <p class="text-sm text-slate-400 mb-3">{{ t('workspace.selectAppleConnect') }}</p>
            <select 
              v-model="selectedConfigId" 
              class="w-full rounded-lg bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint" 
              required
            >
              <option value="">{{ t('workspace.chooseConfig') }}</option>
              <option v-for="config in appleConnectConfigs" :key="config.id" :value="config.id">
                {{ t('common.appleConnectConfig') }} ({{ config.id }})
              </option>
            </select>
          </div>
          
          <div v-if="syncToAppleResult" class="p-3 rounded-lg bg-white/5 border border-white/10">
            <p class="text-sm">{{ syncToAppleResult.message }}</p>
          </div>

          <div class="flex justify-end gap-2 pt-4">
            <button 
              type="button" 
              class="rounded-lg border border-white/20 px-3 py-2 text-sm hover:border-mint/60 hover:text-mint" 
              @click="closeSyncToAppleModal"
            >
              {{ t('common.cancel') || 'Cancel' }}
            </button>
            <button 
              type="button" 
              class="rounded-lg bg-mint px-3 py-2 text-sm font-semibold text-midnight shadow" 
              @click="syncToAppleConnect"
              :disabled="!selectedConfigId || syncingToApple"
            >
              <span v-if="!syncingToApple">{{ t('workspace.syncNow') }}</span>
              <span v-else class="flex items-center gap-1">
                <span class="h-3 w-3 rounded-full bg-midnight animate-pulse"></span>
                {{ t('workspace.syncing') }}
              </span>
            </button>
          </div>
        </div>
      </div>
    </div>

    <div class="grid gap-4 lg:grid-cols-4">
      <section class="rounded-2xl border border-white/10 bg-white/5 p-4 lg:col-span-1 space-y-4">
        <div class="space-y-1">
          <p class="text-xs uppercase tracking-[0.2em] text-slate-500">{{ t('workspace.languages') }}</p>
          <h1 class="text-lg font-semibold">{{ t('workspace.title') }}</h1>
        </div>
        <div class="space-y-2">
          <div class="text-[11px] uppercase tracking-[0.2em] text-slate-500">{{ t('workspace.source') }}</div>
          <button
            class="w-full text-left rounded-lg px-3 py-2 border flex items-center justify-between transition text-xs"
            :class="[
              selectedLang?.code === sourceLanguage
                ? 'border-mint/60 bg-mint/10 text-mint'
                : 'border-white/10 bg-white/5 hover:border-mint/60 hover:text-mint'
            ]"
            @click="selectLang(sourceLanguage)"
          >
            <div>
              <div class="font-semibold text-sm">{{ sourceLanguageName }}</div>
              <div class="text-[11px] text-slate-500">{{ t('workspace.source') }}</div>
            </div>
          </button>
        </div>
        <div class="space-y-2">
          <div class="text-[11px] uppercase tracking-[0.2em] text-slate-500">{{ t('workspace.targets') }}</div>
          <button
            v-for="lang in translations"
            :key="lang.code"
            class="w-full text-left rounded-lg px-3 py-2 border flex items-center justify-between transition text-xs"
            :class="[
              selectedLang?.code === lang.code
                ? 'border-mint/60 bg-mint/10 text-mint'
                : 'border-white/10 bg-white/5 hover:border-mint/60 hover:text-mint'
            ]"
            @click="selectLang(lang.code)"
          >
            <div>
              <div class="font-semibold text-sm">{{ lang.name }} ({{ lang.code }})</div>
              <div class="text-[11px] text-slate-500">{{ lang.total }} {{ t('workspace.items') }}</div>
            </div>
            <span class="text-[11px] rounded-full bg-white/10 px-2 py-1">{{ Math.round((lang.done / lang.total) * 100) }}%</span>
          </button>
        </div>
      </section>

      <section class="rounded-2xl border border-white/10 bg-white/5 p-4 lg:col-span-3 space-y-4">
        <div class="flex items-center justify-between">
          <button class="flex items-center gap-2 rounded-lg border border-white/20 px-3 py-2 text-xs hover:border-mint/60 hover:text-mint" @click="translateCurrentLanguage">
            <span>💬</span>
            <span>{{ t('workspace.translateCurrent') }}</span>
          </button>
          <div class="flex items-center gap-2 text-xs text-slate-400">
            <span>{{ t('workspace.search') }}</span>
            <input class="rounded-lg bg-midnight/40 px-3 py-2 text-xs ring-1 ring-white/10" :placeholder="t('workspace.searchPlaceholder')" />
          </div>
        </div>

        <div class="overflow-x-auto">
          <table class="min-w-full text-sm">
            <thead class="text-left text-slate-500">
              <tr>
                <th class="py-2 px-2 w-1/6">{{ t('workspace.field') }}</th>
                <th class="py-2 px-2 w-1/6">{{ t('workspace.description') }}</th>
                <th class="py-2 px-2">{{ selectedLang ? selectedLang.name : t('workspace.selectedLanguage') }} ({{ selectedLang ? selectedLang.code : 'xx' }})</th>
                <th class="py-2 px-2 w-24">{{ t('common.status') }}</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-white/5">
              <tr v-for="item in metadataItems" :key="item.key" class="align-top">
                <td class="py-2 px-2 font-semibold">{{ item.key }}</td>
                <td class="py-2 px-2 text-xs text-slate-400">{{ item.description }}</td>
                <td class="py-2 px-2">
                  <textarea
                    class="w-full rounded-lg border border-white/10 bg-white/5 px-2 py-1 text-xs"
                    :placeholder="selectedLang && selectedLang.code === sourceLanguage ? item.source : t('workspace.translatePlaceholder')"
                    rows="2"
                    :value="item.translation"
                    @input="updateTranslation(item.key, ($event.target as HTMLTextAreaElement).value)"
                  />
                </td>
                <td class="py-2 px-2">
                  <span class="rounded-full bg-white/10 px-2 py-1 text-[11px]" 
                        :class="item.translation ? 'bg-emerald-900/40 text-emerald-200' : 'bg-amber-900/40 text-amber-200'">
                    {{ item.translation ? t('workspace.completed') : t('workspace.pending') }}
                  </span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'
import { useApi } from '../composables/useApi'
import type { AppLocalization, ProviderConfig } from '../composables/useApi'

const route = useRoute();
const { t } = useI18n()
const { api } = useApi()

const appId = ref(Number(route.params.id));

// State for app metadata
const app = ref<any>(null);

// State for language management
const sourceLanguage = ref('en');
const sourceLanguageName = computed(() => {
  const lang = availableLanguages.find(l => l.code === sourceLanguage.value);
  return lang ? lang.name : sourceLanguage.value;
});
const translations = ref([
  { code: 'zh-CN', name: '简体中文', total: 10, done: 4 },
  { code: 'ja', name: '日本語', total: 10, done: 1 }
]);
const selectedLang = ref<any>(translations.value[0]);

// State for UI modals
const showAddLanguageModal = ref(false);
const showSyncToAppleModal = ref(false);
const syncingToApple = ref(false);
const syncToAppleResult = ref<{ message: string } | null>(null);
const selectedConfigId = ref<number | null>(null);
const appleConnectConfigs = ref<ProviderConfig[]>([]);
const selectedAddLanguage = ref('');

// App metadata items
const metadataItems = ref([
  { key: 'name', source: 'My App', translation: '', description: t('workspace.appNameDesc') },
  { key: 'subtitle', source: 'Subtitle for the app', translation: '', description: t('workspace.appSubtitleDesc') },
  { key: 'shortDescription', source: 'A brief description of the app', translation: '', description: t('workspace.shortDescDesc') },
  { key: 'longDescription', source: 'A detailed description of the app', translation: '', description: t('workspace.longDescDesc') },
  { key: 'keywords', source: 'keyword1, keyword2, keyword3', translation: '', description: t('workspace.keywordsDesc') },
  { key: 'privacyUrl', source: 'https://example.com/privacy', translation: '', description: t('workspace.privacyUrlDesc') },
  { key: 'marketingUrl', source: 'https://example.com', translation: '', description: t('workspace.marketingUrlDesc') },
  { key: 'supportUrl', source: 'https://example.com/support', translation: '', description: t('workspace.supportUrlDesc') },
  { key: 'releaseNotes', source: 'What\'s new in this version', translation: '', description: t('workspace.releaseNotesDesc') }
]);

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
];

const hasAppleConnectConfig = computed(() => appleConnectConfigs.value.length > 0);

function selectLang(code: string) {
  const lang = [...translations.value, { code: sourceLanguage.value, name: sourceLanguageName.value, total: 0, done: 0 }].find(l => l.code === code);
  if (lang) {
    selectedLang.value = lang;
  }
}

function updateTranslation(key: string, value: string) {
  const item = metadataItems.value.find(i => i.key === key);
  if (item) {
    item.translation = value;
    
    // Update the done count for the selected language
    const selected = translations.value.find(t => t === selectedLang.value);
    if (selected) {
      const total = metadataItems.value.length;
      const done = metadataItems.value.filter(i => i.translation.trim() !== '').length;
      selected.done = done;
      selected.total = total;
    }
  }
}

function closeAddLanguageModal() {
  showAddLanguageModal.value = false;
  selectedAddLanguage.value = '';
}

function closeSyncToAppleModal() {
  showSyncToAppleModal.value = false
  selectedConfigId.value = null
  syncToAppleResult.value = null
}

async function addLanguage() {
  if (!selectedAddLanguage.value) return;
  
  // Check if language already exists
  if (!translations.value.some(lang => lang.code === selectedAddLanguage.value)) {
    const lang = availableLanguages.find(l => l.code === selectedAddLanguage.value);
    if (lang) {
      translations.value.push({
        code: selectedAddLanguage.value,
        name: lang.name,
        total: metadataItems.value.length,
        done: 0
      });
    }
  }
  
  closeAddLanguageModal();
}

async function translateCurrentLanguage() {
  // Simulate translation API call
  if (selectedLang.value && selectedLang.value.code !== sourceLanguage.value) {
    // In a real implementation, this would call the translation API
    // For now, we'll just simulate by copying the source values
    metadataItems.value.forEach(item => {
      if (!item.translation) {
        item.translation = `Translated ${item.key} to ${selectedLang.value.name}`;
      }
    });
    
    // Update done count
    const selected = translations.value.find(t => t === selectedLang.value);
    if (selected) {
      selected.done = metadataItems.value.filter(i => i.translation).length;
    }
  }
}

async function translateAll() {
  // Simulate translation for all languages
  // In a real implementation, this would call the translation API for each language
  translations.value.forEach(translation => {
    if (translation.code !== sourceLanguage.value) {
      metadataItems.value.forEach(item => {
        if (!item.translation) {
          item.translation = `Translated ${item.key} to ${translation.name}`;
        }
      });
      translation.done = metadataItems.value.filter(i => i.translation).length;
    }
  });
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

async function syncToAppleConnect() {
  if (!selectedConfigId.value) return

  syncingToApple.value = true
  syncToAppleResult.value = null

  try {
    // Get the selected config to extract credentials
    const selectedConfig = appleConnectConfigs.value.find(config => config.id === selectedConfigId.value);
    if (!selectedConfig) {
      throw new Error('Selected configuration not found');
    }

    const response = await api.syncAppToApple(appId.value, {
      issuerId: selectedConfig.configData.issuerID,
      keyId: selectedConfig.configData.keyID,
      privateKey: selectedConfig.configData.privateKey
    })

    if (response.success) {
      syncToAppleResult.value = {
        message: response.message || t('workspace.syncSuccess')
      }
    } else {
      syncToAppleResult.value = {
        message: response.message || t('workspace.syncFailed')
      }
    }
  } catch (error) {
    console.error('Failed to sync to Apple Connect:', error)
    syncToAppleResult.value = {
      message: t('workspace.syncError')
    }
  } finally {
    syncingToApple.value = false
  }
}

// Fetch app data on component mount
onMounted(async () => {
  try {
    const response = await api.getApp(appId.value);
    if (response.success) {
      app.value = response.app;
      sourceLanguage.value = response.app.primaryLocale || 'en';
    }
    
    // Fetch localizations for this app
    const localizationResponse = await api.getAppLocalizations(appId.value);
    if (localizationResponse.success) {
      // Process localizations to populate translations
      localizationResponse.localizations.forEach((loc: AppLocalization) => {
        if (loc.languageCode !== sourceLanguage.value) {
          const lang = availableLanguages.find(l => l.code === loc.languageCode);
          if (lang) {
            translations.value.push({
              code: loc.languageCode,
              name: lang.name,
              total: metadataItems.value.length,
              done: 0 // This would be calculated based on actual translated fields
            });
          }
        }
      });
    }
  } catch (error) {
    console.error('Failed to fetch app data:', error);
  }
  
  // Fetch provider configs
  fetchProviderConfigs();
});
</script>
