<template>
  <div class="space-y-4">
    <div class="flex flex-wrap items-center gap-2 rounded-2xl border border-white/10 bg-white/5 px-4 py-3">
      <button class="flex items-center gap-2 rounded-lg border border-white/20 px-3 py-2 text-xs hover:border-mint/60 hover:text-mint">
        <span>📤</span>
        <span>{{ t('workspace.upload') }}</span>
      </button>
      <button class="flex items-center gap-2 rounded-lg border border-white/20 px-3 py-2 text-xs hover:border-mint/60 hover:text-mint">
        <span>➕</span>
        <span>{{ t('workspace.addLang') }}</span>
      </button>
      <button class="flex items-center gap-2 rounded-lg border border-white/20 px-3 py-2 text-xs hover:border-mint/60 hover:text-mint">
        <span>⚡</span>
        <span>{{ t('workspace.translateAll') }}</span>
      </button>
      <button class="flex items-center gap-2 rounded-lg border border-white/20 px-3 py-2 text-xs hover:border-mint/60 hover:text-mint">
        <span>⬇️</span>
        <span>{{ t('workspace.export') }}</span>
      </button>
    </div>

    <div class="grid gap-4 lg:grid-cols-4">
      <section class="rounded-2xl border border-white/10 bg-white/5 p-4 lg:col-span-1 space-y-4">
        <div class="space-y-1">
          <p class="text-xs uppercase tracking-[0.2em] text-slate-500">Languages</p>
          <h1 class="text-lg font-semibold">{{ t('workspace.title') }}</h1>
        </div>
        <div class="space-y-2">
          <button
            class="flex w-full items-center justify-between rounded-lg border border-white/10 bg-white/5 px-3 py-2 text-left text-xs text-slate-200 hover:border-mint/60 hover:text-mint"
            @click="showSourceModal = true"
          >
            <span>Source</span>
            <span>{{ sourceLanguage }}</span>
          </button>
          <div class="flex items-center justify-between text-xs text-slate-400">
            <span>Total Strings</span>
            <span>{{ stats.total }}</span>
          </div>
        </div>
        <div class="space-y-2">
          <div class="text-[11px] uppercase tracking-[0.2em] text-slate-500">Targets</div>
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
              <div class="text-[11px] text-slate-500">{{ lang.done }}/{{ stats.total }}</div>
            </div>
            <span class="text-[11px] rounded-full bg-white/10 px-2 py-1">{{ Math.round((lang.done / stats.total) * 100) }}%</span>
          </button>
        </div>
      </section>

      <section class="rounded-2xl border border-white/10 bg-white/5 p-4 lg:col-span-3 space-y-4">
        <div class="flex items-center justify-between">
          <button class="flex items-center gap-2 rounded-lg border border-white/20 px-3 py-2 text-xs hover:border-mint/60 hover:text-mint">
            <span>💬</span>
            <span>{{ t('workspace.translateLang') }}</span>
          </button>
          <div class="flex items-center gap-2 text-xs text-slate-400">
            <span>{{ t('workspace.search') }}</span>
            <input class="rounded-lg bg-midnight/40 px-3 py-2 text-xs ring-1 ring-white/10" :placeholder="t('workspace.search')" />
          </div>
        </div>

        <div class="overflow-x-auto">
          <table class="min-w-full text-sm">
            <thead class="text-left text-slate-500">
              <tr>
                <th class="py-2 px-2 w-1/6">Key</th>
                <th class="py-2 px-2 w-1/6">Comment</th>
                <th class="py-2 px-2" v-if="selectedLang">{{ selectedLang.name }} ({{ selectedLang.code }})</th>
                <th class="py-2 px-2 w-24">Status</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-white/5">
              <tr v-for="item in stringItems" :key="item.key" class="align-top">
                <td class="py-2 px-2 font-semibold">{{ item.key }}</td>
                <td class="py-2 px-2 text-xs text-slate-400">{{ item.comment }}</td>
                <td class="py-2 px-2" v-if="selectedLang">
                  <textarea
                    class="w-full rounded-lg border border-white/10 bg-white/5 px-2 py-1 text-xs"
                    :placeholder="selectedLang.code === sourceLanguage ? item.source : '…'"
                    rows="2"
                  />
                </td>
                <td class="py-2 px-2">
                  <span class="rounded-full bg-white/10 px-2 py-1 text-[11px]">Pending</span>
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
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const sourceLanguage = ref('en')
const stats = { total: 8 }
const translations = [
  { code: 'en', name: 'English', done: 8 },
  { code: 'zh-CN', name: '简体中文', done: 4 },
  { code: 'ja', name: '日本語', done: 1 }
]

const selectedLang = computed(() => translations.find((l) => l.code === selectedCode.value) || translations[0])
const selectedCode = ref(translations[0].code)
const showSourceModal = ref(false)

function selectLang(code: string) {
  selectedCode.value = code
}

const stringItems = [
  { key: 'welcome_title', comment: 'Welcome headline', source: 'Welcome to the app' },
  { key: 'cta_start', comment: 'Primary CTA', source: 'Get started' },
  { key: 'menu_home', comment: 'Menu item', source: 'Home' },
  { key: 'menu_settings', comment: 'Menu item', source: 'Settings' }
]

const sourceOptions = ['en', 'zh-CN', 'ja', 'fr']

function setSource(lang: string) {
  sourceLanguage.value = lang
  showSourceModal.value = false
}
</script>
