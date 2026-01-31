<template>
  <div class="space-y-6">
    <header class="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
      <div>
        <p class="text-xs uppercase tracking-[0.3em] text-slate-500">{{ t('nav.languages') }}</p>
        <h1 class="text-2xl font-semibold">{{ t('languages.title') }}</h1>
        <p class="text-sm text-slate-400">{{ t('languages.subtitle') }}</p>
      </div>
      <div class="text-xs text-slate-500">{{ t('languages.readonly') }}</div>
    </header>

    <section class="rounded-2xl border border-white/10 bg-white/5 p-4">
      <div class="flex items-center justify-between mb-3">
        <div class="text-sm text-slate-400">{{ t('languages.languageList') }}</div>
        <input
          v-model="searchQuery"
          class="rounded-lg bg-midnight/40 px-3 py-2 text-xs ring-1 ring-white/10 w-64 focus:outline-none focus:ring-mint/60"
          :placeholder="t('languages.searchPlaceholder')"
        />
      </div>
      <div class="overflow-x-auto">
        <table class="min-w-full text-sm">
          <thead class="text-left text-slate-500">
            <tr>
              <th class="py-2 w-12"></th>
              <th class="py-2">{{ t('languages.name') }}</th>
              <th class="py-2">{{ t('languages.code') }}</th>
              <th class="py-2">{{ t('languages.nativeName') }}</th>
              <th class="py-2">{{ t('languages.direction') }}</th>
              <th class="py-2">{{ t('languages.available') }}</th>
              <th class="py-2">{{ t('common.actions') }}</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-white/5">
            <tr v-for="lang in filteredLanguages" :key="lang.code" class="hover:bg-white/5">
              <td class="py-2 text-center text-xl">{{ lang.emoji || '🌐' }}</td>
              <td class="py-2">{{ lang.name }}</td>
              <td class="py-2">{{ lang.code }}</td>
              <td class="py-2">{{ lang.native_name }}</td>
              <td class="py-2">{{ lang.direction === 'rtl' ? t('languages.rtl') : t('languages.ltr') }}</td>
              <td class="py-2">
                <span class="rounded-full px-2 py-1 text-xs" :class="lang.enabled ? 'bg-emerald-900/40 text-emerald-200' : 'bg-rose-900/40 text-rose-200'">
                  {{ lang.enabled ? t('languages.available') : t('languages.disabled') }}
                </span>
              </td>
              <td class="py-2 space-x-2 text-xs">
                <button class="rounded border border-white/20 px-2 py-1 hover:border-mint/60 hover:text-mint">{{ t('common.edit') }}</button>
                <button class="rounded border border-white/20 px-2 py-1 hover:border-mint/60 hover:text-mint">{{ t('languages.disable') }}</button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useApi } from '../composables/useApi'

const { t } = useI18n()
const { api } = useApi()

interface Language {
  code: string
  name: string
  native_name: string
  region?: string
  direction: string
  emoji?: string
  enabled?: boolean
}

const languages = ref<Language[]>([])
const searchQuery = ref('')
const loading = ref(false)

const filteredLanguages = computed(() => {
  if (!searchQuery.value) {
    return languages.value
  }
  const query = searchQuery.value.toLowerCase()
  return languages.value.filter(lang =>
    lang.name.toLowerCase().includes(query) ||
    lang.code.toLowerCase().includes(query) ||
    lang.native_name.toLowerCase().includes(query)
  )
})

async function fetchLanguages() {
  loading.value = true
  try {
    // Fetch Apple Connect supported languages from API
    const response = await api.getSupportedLanguages()
    if (response.success) {
      languages.value = response.languages.map(lang => ({
        ...lang,
        enabled: true
      }))
    }
  } catch (error) {
    console.error('Failed to fetch languages:', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchLanguages()
})
</script>
