<template>
  <div class="space-y-6">
    <header class="flex flex-col gap-4 rounded-3xl border border-white/10 bg-white/5 p-5 lg:flex-row lg:items-start lg:justify-between">
      <div class="space-y-2">
        <p class="text-xs uppercase tracking-[0.3em] text-slate-500">{{ t('nav.apps') }}</p>
        <div>
          <h1 class="text-2xl font-semibold">{{ app?.name || t('workspace.title') }}</h1>
          <p class="mt-1 text-sm text-slate-400">
            {{ workspaceSubtitle }}
          </p>
        </div>
        <div class="flex flex-wrap gap-2 text-xs text-slate-400">
          <span class="rounded-full border border-white/10 bg-white/5 px-3 py-1">
            {{ t('workspace.source') }}: {{ sourceLanguageLabel }}
          </span>
          <span v-if="workspaceProject" class="rounded-full border border-white/10 bg-white/5 px-3 py-1">
            {{ workspaceProject.fileName || 'Localizable.xcstrings' }}
          </span>
          <span v-if="stringEntries.length" class="rounded-full border border-white/10 bg-white/5 px-3 py-1">
            {{ stringEntries.length }} {{ t('workspace.items') }}
          </span>
        </div>
      </div>

      <div class="grid gap-3 sm:grid-cols-2">
        <div class="rounded-2xl border border-white/10 bg-midnight/40 p-3">
          <label class="block text-xs uppercase tracking-[0.2em] text-slate-500">{{ t('workspace.selectProviderConfig') }}</label>
          <select
            v-model="selectedProviderConfigId"
            class="mt-2 w-full rounded-xl bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint"
          >
            <option :value="null">{{ t('workspace.chooseConfig') }}</option>
            <option v-for="config in translationProviderConfigs" :key="config.id" :value="config.id">
              {{ providerLabel(config) }}
            </option>
          </select>
        </div>

        <div class="rounded-2xl border border-white/10 bg-midnight/40 p-3">
          <label class="block text-xs uppercase tracking-[0.2em] text-slate-500">{{ t('workspace.search') }}</label>
          <input
            v-model.trim="searchQuery"
            class="mt-2 w-full rounded-xl bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint"
            :placeholder="t('workspace.searchPlaceholder')"
          />
        </div>
      </div>
    </header>

    <section class="rounded-3xl border border-dashed border-white/10 bg-white/5 p-5">
      <div class="flex flex-col gap-4 xl:flex-row xl:items-end xl:justify-between">
        <div class="space-y-2">
          <h2 class="text-lg font-semibold">{{ t('workspace.title') }}</h2>
          <p class="text-sm text-slate-400">
            {{ workspaceProject ? t('workspace.selectTargetLanguagesDesc') : t('workspace.noLocalizationsDesc') }}
          </p>
        </div>

        <div class="flex flex-wrap items-center gap-2">
          <input
            ref="fileInputRef"
            type="file"
            accept=".xcstrings,application/json"
            class="hidden"
            @change="handleFileSelected"
          />
          <button class="rounded-xl border border-white/20 px-3 py-2 text-sm hover:border-mint/60 hover:text-mint" @click="fileInputRef?.click()">
            {{ workspaceProject ? '替换 xcstrings 文件' : '上传 xcstrings 文件' }}
          </button>
          <button
            class="rounded-xl border border-white/20 px-3 py-2 text-sm hover:border-mint/60 hover:text-mint disabled:cursor-not-allowed disabled:opacity-50"
            :disabled="!workspaceProject || savingDrafts"
            @click="saveSelectedLanguage"
          >
            {{ savingDrafts ? '保存中...' : (t('common.save') || 'Save') }}
          </button>
          <button
            class="rounded-xl bg-mint px-3 py-2 text-sm font-semibold text-midnight shadow disabled:cursor-not-allowed disabled:opacity-50"
            :disabled="!canTranslate || translationRunning"
            @click="translateSelectedLanguages"
          >
            {{ translationRunning ? t('workspace.syncing') : t('workspace.translateAll') }}
          </button>
          <button
            class="rounded-xl border border-white/20 px-3 py-2 text-sm hover:border-mint/60 hover:text-mint disabled:cursor-not-allowed disabled:opacity-50"
            :disabled="!workspaceProject"
            @click="downloadXcstrings"
          >
            导出 xcstrings
          </button>
        </div>
      </div>

      <div class="mt-4 grid gap-4 lg:grid-cols-[minmax(0,1fr)_320px]">
        <div class="rounded-2xl border border-white/10 bg-midnight/30 p-4">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm font-medium">{{ t('workspace.selectTargetLanguages') }}</p>
              <p class="text-xs text-slate-500">{{ t('workspace.selectTargetLanguagesDesc') }}</p>
            </div>
            <button class="text-xs text-slate-400 hover:text-mint" @click="showAddLanguageModal = true">
              {{ t('workspace.addLanguage') }}
            </button>
          </div>

          <div class="mt-3 flex flex-wrap gap-2">
            <button
              v-for="lang in targetLanguageOptions"
              :key="lang.code"
              class="rounded-full border px-3 py-1 text-xs transition"
              :class="selectedTargetLanguages.includes(lang.code)
                ? 'border-mint/60 bg-mint/10 text-mint'
                : 'border-white/10 bg-white/5 text-slate-300 hover:border-mint/60 hover:text-mint'"
              @click="toggleTargetLanguage(lang.code)"
            >
              {{ languageLabel(lang.code) }}
            </button>
            <span v-if="!targetLanguageOptions.length" class="text-sm text-slate-500">
              先上传 `Localizable.xcstrings`，再选择目标语言。
            </span>
          </div>
        </div>

        <div class="rounded-2xl border border-white/10 bg-midnight/30 p-4">
          <p class="text-sm font-medium">翻译状态</p>
          <div v-if="activeJob" class="mt-3 space-y-2">
            <div class="flex items-center justify-between text-xs text-slate-400">
              <span>{{ activeJob.status }}</span>
              <span>{{ activeJob.done }}/{{ activeJob.total || stringEntries.length }}</span>
            </div>
            <div class="h-2 rounded-full bg-white/10">
              <div class="h-2 rounded-full bg-mint transition-all" :style="{ width: `${activeJob.progress}%` }"></div>
            </div>
            <p class="text-xs text-slate-400">{{ activeJob.progress }}%</p>
          </div>
          <p v-else class="mt-3 text-sm text-slate-500">暂无进行中的翻译任务。</p>
        </div>
      </div>
    </section>

    <div v-if="showAddLanguageModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 px-4">
      <div class="w-full max-w-lg rounded-2xl border border-white/10 bg-midnight p-6 shadow-xl">
        <div class="flex items-center justify-between">
          <h2 class="text-lg font-semibold">{{ t('workspace.addLanguage') }}</h2>
          <button class="text-slate-400 hover:text-white" @click="closeAddLanguageModal">×</button>
        </div>

        <div class="mt-4 space-y-4">
          <div>
            <p class="mb-3 text-sm text-slate-400">{{ t('workspace.selectLanguage') }}</p>
            <select
              v-model="selectedAddLanguage"
              class="w-full rounded-lg bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint"
            >
              <option value="">{{ t('workspace.chooseLanguage') }}</option>
              <option v-for="lang in addableLanguages" :key="lang.code" :value="lang.code">
                {{ languageLabel(lang.code) }}
              </option>
            </select>
          </div>

          <div class="flex justify-end gap-2 pt-4">
            <button type="button" class="rounded-lg border border-white/20 px-3 py-2 text-sm hover:border-mint/60 hover:text-mint" @click="closeAddLanguageModal">
              {{ t('common.cancel') || 'Cancel' }}
            </button>
            <button
              type="button"
              class="rounded-lg bg-mint px-3 py-2 text-sm font-semibold text-midnight shadow disabled:cursor-not-allowed disabled:opacity-50"
              :disabled="!selectedAddLanguage"
              @click="confirmAddLanguage"
            >
              {{ t('workspace.add') }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <div class="grid gap-4 lg:grid-cols-[280px_minmax(0,1fr)]">
      <aside class="rounded-3xl border border-white/10 bg-white/5 p-4">
        <div class="space-y-1">
          <p class="text-xs uppercase tracking-[0.2em] text-slate-500">{{ t('workspace.languages') }}</p>
          <h2 class="text-lg font-semibold">{{ t('workspace.selectLanguage') }}</h2>
        </div>

        <div class="mt-4 space-y-2">
          <button
            v-for="lang in visibleLanguages"
            :key="lang.code"
            class="w-full rounded-2xl border px-3 py-3 text-left transition"
            :class="selectedLanguage === lang.code
              ? 'border-mint/60 bg-mint/10 text-mint'
              : 'border-white/10 bg-white/5 hover:border-mint/60 hover:text-mint'"
            @click="selectedLanguage = lang.code"
          >
            <div class="flex items-center justify-between gap-3">
              <div>
                <div class="font-medium">{{ languageLabel(lang.code) }}</div>
                <div class="mt-1 text-[11px] text-slate-500">
                  {{ lang.code === sourceLanguage ? t('workspace.source') : t('workspace.targets') }}
                </div>
              </div>
              <span class="rounded-full bg-white/10 px-2 py-1 text-[11px]">
                {{ lang.done }}/{{ lang.total }}
              </span>
            </div>
          </button>
        </div>
      </aside>

      <section class="rounded-3xl border border-white/10 bg-white/5 p-4">
        <div v-if="!workspaceProject" class="flex min-h-[320px] items-center justify-center rounded-2xl border border-dashed border-white/10 bg-midnight/20 p-8 text-center">
          <div>
            <p class="text-lg font-medium">上传 `Localizable.xcstrings` 开始编辑</p>
            <p class="mt-2 text-sm text-slate-400">
              页面会为当前 app 建立一个专属工作区，支持上传、翻译、手动修改和导出。
            </p>
          </div>
        </div>

        <template v-else>
          <div class="flex flex-col gap-2 border-b border-white/10 pb-4 sm:flex-row sm:items-center sm:justify-between">
            <div>
              <h2 class="text-lg font-semibold">{{ languageLabel(selectedLanguage) }}</h2>
              <p class="text-sm text-slate-400">
                {{ selectedLanguage === sourceLanguage ? '源文案只读预览。' : '右侧内容可直接编辑并保存到项目。' }}
              </p>
            </div>
            <div class="text-xs text-slate-500">
              {{ filteredEntries.length }} / {{ stringEntries.length }} 项
            </div>
          </div>

          <div class="mt-4 overflow-x-auto">
            <table class="min-w-full text-sm">
              <thead class="text-left text-slate-500">
                <tr>
                  <th class="w-[28%] px-3 py-2">Key</th>
                  <th class="w-[32%] px-3 py-2">{{ sourceLanguageLabel }}</th>
                  <th class="px-3 py-2">{{ languageLabel(selectedLanguage) }}</th>
                  <th class="w-24 px-3 py-2">{{ t('common.status') }}</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-white/5">
                <tr v-for="entry in filteredEntries" :key="entry.key" class="align-top">
                  <td class="px-3 py-3">
                    <div class="font-medium">{{ entry.key }}</div>
                  </td>
                  <td class="px-3 py-3 text-slate-300">
                    <div class="max-h-36 overflow-auto whitespace-pre-wrap rounded-xl bg-midnight/40 px-3 py-2 text-xs">
                      {{ entry.sourceText || '-' }}
                    </div>
                  </td>
                  <td class="px-3 py-3">
                    <textarea
                      v-if="selectedLanguage !== sourceLanguage"
                      class="min-h-[88px] w-full rounded-xl border border-white/10 bg-white/5 px-3 py-2 text-xs"
                      :value="draftValue(entry.key)"
                      :placeholder="t('workspace.translatePlaceholder')"
                      @input="updateDraft(entry.key, ($event.target as HTMLTextAreaElement).value)"
                    />
                    <div v-else class="min-h-[88px] whitespace-pre-wrap rounded-xl border border-white/10 bg-white/5 px-3 py-2 text-xs text-slate-300">
                      {{ entry.sourceText || '-' }}
                    </div>
                  </td>
                  <td class="px-3 py-3">
                    <span
                      class="rounded-full px-2 py-1 text-[11px]"
                      :class="entryStatus(entry.key).filled
                        ? 'bg-emerald-900/40 text-emerald-200'
                        : 'bg-amber-900/40 text-amber-200'"
                    >
                      {{ entryStatus(entry.key).filled ? t('workspace.completed') : t('workspace.pending') }}
                    </span>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </template>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'
import { useApi } from '../composables/useApi'
import { useToast } from '../composables/useToast'
import type { LanguageMetadata, Project, ProviderConfig, Translation, TranslationQueue, App } from '../composables/useApi'

type XCStrings = {
  sourceLanguage?: string
  strings?: Record<string, XCStringEntry>
  version?: string
}

type XCStringEntry = {
  localizations?: Record<string, { stringUnit?: { state?: string; value?: string } }>
}

type WorkspaceEntry = {
  key: string
  sourceText: string
}

type LanguageSummary = {
  code: string
  total: number
  done: number
}

const WORKSPACE_MARKER = '__app_workspace__:'

const route = useRoute()
const { t } = useI18n()
const { api } = useApi()
const toast = useToast()

const appId = Number(route.params.id)

const app = ref<App | null>(null)
const workspaceProject = ref<Project | null>(null)
const availableLanguages = ref<LanguageMetadata[]>([])
const translationProviderConfigs = ref<ProviderConfig[]>([])
const translations = ref<Translation[]>([])
const stringEntries = ref<WorkspaceEntry[]>([])
const draftTranslations = ref<Record<string, Record<string, string>>>({})
const selectedLanguage = ref('en')
const selectedTargetLanguages = ref<string[]>([])
const selectedProviderConfigId = ref<number | null>(null)
const searchQuery = ref('')
const showAddLanguageModal = ref(false)
const selectedAddLanguage = ref('')
const savingDrafts = ref(false)
const translationRunning = ref(false)
const fileInputRef = ref<HTMLInputElement | null>(null)
const activeJob = ref<TranslationQueue | null>(null)

let jobPollTimer: number | null = null

const sourceLanguage = computed(() => workspaceProject.value?.sourceLanguage || app.value?.primaryLocale || 'en')

const sourceLanguageLabel = computed(() => languageLabel(sourceLanguage.value))

const targetLanguageOptions = computed(() =>
  availableLanguages.value.filter((lang) => lang.code !== sourceLanguage.value),
)

const addableLanguages = computed(() =>
  targetLanguageOptions.value.filter((lang) => !selectedTargetLanguages.value.includes(lang.code)),
)

const languageStats = computed<Record<string, LanguageSummary>>(() => {
  const stats: Record<string, LanguageSummary> = {
    [sourceLanguage.value]: {
      code: sourceLanguage.value,
      total: stringEntries.value.length,
      done: stringEntries.value.filter((entry) => entry.sourceText.trim() !== '').length,
    },
  }

  for (const code of selectedTargetLanguages.value) {
    stats[code] = {
      code,
      total: stringEntries.value.length,
      done: stringEntries.value.filter((entry) => getTranslationValue(entry.key, code).trim() !== '').length,
    }
  }

  return stats
})

const visibleLanguages = computed(() => {
  const order = [sourceLanguage.value, ...selectedTargetLanguages.value]
  return order.map((code) => languageStats.value[code]).filter(Boolean)
})

const filteredEntries = computed(() => {
  const query = searchQuery.value.trim().toLowerCase()
  if (!query) {
    return stringEntries.value
  }

  return stringEntries.value.filter((entry) => {
    const target = getTranslationValue(entry.key, selectedLanguage.value)
    return [entry.key, entry.sourceText, target].some((value) => value.toLowerCase().includes(query))
  })
})

const workspaceSubtitle = computed(() => {
  if (!workspaceProject.value) {
    return '上传 `Localizable.xcstrings` 后，可在这里选择目标语言、执行机器翻译并手动修订。'
  }
  return '类似 Xcode 的 `Localizable.xcstrings` 工作区，支持上传、翻译、编辑和导出。'
})

const canTranslate = computed(() => {
  return Boolean(
    workspaceProject.value &&
      selectedTargetLanguages.value.length &&
      selectedProviderConfig.value &&
      stringEntries.value.length,
  )
})

const selectedProviderConfig = computed(() =>
  translationProviderConfigs.value.find((config) => config.id === selectedProviderConfigId.value) || null,
)

watch(sourceLanguage, (nextSource) => {
  if (!selectedLanguage.value || selectedLanguage.value === 'en') {
    selectedLanguage.value = nextSource
  }
})

watch(
  translationProviderConfigs,
  (configs) => {
    if (!configs.length) {
      selectedProviderConfigId.value = null
      return
    }

    if (!configs.some((config) => config.id === selectedProviderConfigId.value)) {
      selectedProviderConfigId.value = configs.find((config) => config.isDefault)?.id || configs[0].id
    }
  },
  { immediate: true },
)

onMounted(async () => {
  if (!Number.isFinite(appId) || appId <= 0) {
    toast.error('无效的应用 ID')
    return
  }

  await Promise.all([fetchApp(), fetchLanguages(), fetchProviderConfigs()])
  await loadWorkspaceProject()
})

onBeforeUnmount(() => {
  stopPollingJob()
})

function workspaceDescription(appIdValue: number) {
  return `${WORKSPACE_MARKER}${appIdValue}`
}

function languageLabel(code: string) {
  const found = availableLanguages.value.find((lang) => lang.code === code)
  if (!found) {
    return code
  }
  return `${found.name} (${found.code})`
}

function providerLabel(config: ProviderConfig) {
  return `${config.providerType}${config.isDefault ? ' · default' : ''}`
}

function draftValue(key: string) {
  return getTranslationValue(key, selectedLanguage.value)
}

function getTranslationValue(key: string, language: string) {
  return draftTranslations.value[language]?.[key] || ''
}

function entryStatus(key: string) {
  const value = getTranslationValue(key, selectedLanguage.value)
  return { filled: value.trim().length > 0 }
}

function toggleTargetLanguage(code: string) {
  if (selectedTargetLanguages.value.includes(code)) {
    selectedTargetLanguages.value = selectedTargetLanguages.value.filter((lang) => lang !== code)
    if (selectedLanguage.value === code) {
      selectedLanguage.value = sourceLanguage.value
    }
    return
  }

  selectedTargetLanguages.value = [...selectedTargetLanguages.value, code]
}

function closeAddLanguageModal() {
  showAddLanguageModal.value = false
  selectedAddLanguage.value = ''
}

function confirmAddLanguage() {
  if (!selectedAddLanguage.value) {
    return
  }
  if (!selectedTargetLanguages.value.includes(selectedAddLanguage.value)) {
    selectedTargetLanguages.value = [...selectedTargetLanguages.value, selectedAddLanguage.value]
  }
  selectedLanguage.value = selectedAddLanguage.value
  closeAddLanguageModal()
}

function updateDraft(key: string, value: string) {
  if (selectedLanguage.value === sourceLanguage.value) {
    return
  }

  draftTranslations.value = {
    ...draftTranslations.value,
    [selectedLanguage.value]: {
      ...(draftTranslations.value[selectedLanguage.value] || {}),
      [key]: value,
    },
  }
}

function parseXCStrings(content: string) {
  const parsed = JSON.parse(content) as XCStrings
  const source = parsed.sourceLanguage || sourceLanguage.value
  const entries = Object.entries(parsed.strings || {}).map(([key, entry]) => ({
    key,
    sourceText: entry.localizations?.[source]?.stringUnit?.value || '',
  }))
  return { parsed, entries, source }
}

function syncDraftsFromProject(project: Project, projectTranslations: Translation[]) {
  const fileContent = project.fileContent || ''
  if (!fileContent) {
    stringEntries.value = []
    draftTranslations.value = {}
    selectedTargetLanguages.value = []
    selectedLanguage.value = sourceLanguage.value
    return
  }

  const { parsed, entries, source } = parseXCStrings(fileContent)
  stringEntries.value = entries

  const nextDrafts: Record<string, Record<string, string>> = {
    [source]: Object.fromEntries(entries.map((entry) => [entry.key, entry.sourceText])),
  }

  const targetLangs = new Set<string>()

  for (const [key, item] of Object.entries(parsed.strings || {})) {
    for (const [lang, localization] of Object.entries(item.localizations || {})) {
      if (lang === source) {
        continue
      }
      targetLangs.add(lang)
      if (!nextDrafts[lang]) {
        nextDrafts[lang] = {}
      }
      nextDrafts[lang][key] = localization.stringUnit?.value || ''
    }
  }

  for (const translation of projectTranslations) {
    if (!nextDrafts[translation.targetLanguage]) {
      nextDrafts[translation.targetLanguage] = {}
    }
    nextDrafts[translation.targetLanguage][translation.key] = translation.targetText || ''
    targetLangs.add(translation.targetLanguage)
  }

  draftTranslations.value = nextDrafts
  selectedTargetLanguages.value = Array.from(targetLangs)
  selectedLanguage.value = selectedTargetLanguages.value[0] || source
}

async function fetchApp() {
  const response = await api.getApp(appId)
  if (response.success) {
    app.value = response.app
  }
}

async function fetchLanguages() {
  try {
    const response = await api.getSupportedLanguages()
    if (response.success) {
      availableLanguages.value = response.languages
    }
  } catch (error) {
    console.error('Failed to fetch languages:', error)
  }
}

async function fetchProviderConfigs() {
  try {
    const response = await api.getProviderConfigs()
    if (response.success) {
      translationProviderConfigs.value = response.configs.filter(
        (config) => config.providerType !== 'appleconnect',
      )
    }
  } catch (error) {
    console.error('Failed to fetch provider configs:', error)
  }
}

async function loadWorkspaceProject() {
  try {
    const response = await api.getProjects('xcstrings')
    if (!response.success) {
      return
    }

    const project = response.projects.find((item) => item.description === workspaceDescription(appId)) || null
    workspaceProject.value = project
    if (!project) {
      stringEntries.value = []
      draftTranslations.value = {}
      selectedTargetLanguages.value = []
      selectedLanguage.value = sourceLanguage.value
      return
    }

    await refreshProjectData(project.id)
  } catch (error) {
    console.error('Failed to load workspace project:', error)
    toast.error('加载工作区失败')
  }
}

async function refreshProjectData(projectId: number) {
  const [projectResponse, translationsResponse] = await Promise.all([
    api.getProject(projectId),
    api.getTranslations(projectId),
  ])

  if (!projectResponse.success) {
    return
  }

  workspaceProject.value = projectResponse.project
  translations.value = translationsResponse.success ? translationsResponse.translations : []
  syncDraftsFromProject(projectResponse.project, translations.value)
}

async function ensureWorkspaceProject(file: File) {
  if (workspaceProject.value) {
    return workspaceProject.value
  }

  const fileContent = await file.text()
  const parsed = JSON.parse(fileContent) as XCStrings
  const response = await api.createProject({
    name: `${app.value?.name || 'App'} Localizable`,
    description: workspaceDescription(appId),
    fileName: file.name,
    fileContent,
    sourceLanguage: parsed.sourceLanguage || app.value?.primaryLocale || 'en',
    projectType: 'xcstrings',
  })

  if (!response.success) {
    throw new Error('Failed to create workspace project')
  }

  workspaceProject.value = response.project
  return response.project
}

async function handleFileSelected(event: Event) {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) {
    return
  }

  try {
    const project = await ensureWorkspaceProject(file)
    if (workspaceProject.value) {
      await api.uploadToProject(project.id, file)
    }
    await refreshProjectData(project.id)
    toast.success('xcstrings 文件已上传')
  } catch (error) {
    console.error('Failed to upload xcstrings file:', error)
    toast.error(error instanceof Error ? error.message : '上传 xcstrings 文件失败')
  } finally {
    input.value = ''
  }
}

async function saveSelectedLanguage() {
  if (!workspaceProject.value || selectedLanguage.value === sourceLanguage.value) {
    return
  }

  savingDrafts.value = true
  try {
    const updates = stringEntries.value.map((entry) => ({
      key: entry.key,
      targetLanguage: selectedLanguage.value,
      targetText: getTranslationValue(entry.key, selectedLanguage.value),
      state: getTranslationValue(entry.key, selectedLanguage.value).trim() ? 'translated' : 'new',
    }))

    await api.bulkUpdateTranslations(workspaceProject.value.id, updates)
    await refreshProjectData(workspaceProject.value.id)
    toast.success('当前语言已保存')
  } catch (error) {
    console.error('Failed to save language drafts:', error)
    toast.error('保存翻译失败')
  } finally {
    savingDrafts.value = false
  }
}

async function translateSelectedLanguages() {
  if (!workspaceProject.value || !selectedProviderConfig.value) {
    toast.warning('请先上传 xcstrings 文件并选择翻译配置')
    return
  }

  if (!selectedTargetLanguages.value.length) {
    toast.warning('请至少选择一个目标语言')
    return
  }

  translationRunning.value = true

  try {
    await savePendingBeforeTranslate()
    const response = await api.createQueueJob({
      jobType: 'xcstrings',
      projectId: workspaceProject.value.id,
      providerType: selectedProviderConfig.value.providerType,
      sourceLanguage: sourceLanguage.value,
      targetLanguages: selectedTargetLanguages.value,
      configData: selectedProviderConfig.value.configData || {},
    })

    if (!response.success) {
      throw new Error('Failed to submit translation job')
    }

    activeJob.value = response.job
    toast.success('翻译任务已提交')
    startPollingJob(response.job.id)
  } catch (error) {
    console.error('Failed to translate selected languages:', error)
    translationRunning.value = false
    toast.error(error instanceof Error ? error.message : '提交翻译任务失败')
  }
}

async function savePendingBeforeTranslate() {
  if (selectedLanguage.value !== sourceLanguage.value) {
    await saveSelectedLanguage()
  }
}

function stopPollingJob() {
  if (jobPollTimer !== null) {
    window.clearInterval(jobPollTimer)
    jobPollTimer = null
  }
}

function startPollingJob(jobId: number) {
  stopPollingJob()
  jobPollTimer = window.setInterval(async () => {
    try {
      const response = await api.getQueueJob(jobId)
      if (!response.success) {
        return
      }

      activeJob.value = response.job
      if (response.job.status === 'completed') {
        stopPollingJob()
        translationRunning.value = false
        if (workspaceProject.value) {
          await refreshProjectData(workspaceProject.value.id)
        }
        toast.success('翻译完成')
      }

      if (response.job.status === 'failed') {
        stopPollingJob()
        translationRunning.value = false
        toast.error(response.job.error || '翻译失败')
      }
    } catch (error) {
      console.error('Failed to poll translation job:', error)
      stopPollingJob()
      translationRunning.value = false
      toast.error('查询翻译进度失败')
    }
  }, 2000)
}

async function downloadXcstrings() {
  if (!workspaceProject.value) {
    return
  }

  try {
    const blob = await api.exportProject(workspaceProject.value.id)
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = workspaceProject.value.fileName || 'Localizable_translated.xcstrings'
    link.click()
    window.URL.revokeObjectURL(url)
  } catch (error) {
    console.error('Failed to export project:', error)
    toast.error('导出 xcstrings 失败')
  }
}
</script>
