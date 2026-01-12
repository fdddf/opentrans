<template>
  <div class="min-h-screen text-slate-50 dark:text-slate-100">
    <div class="mx-auto max-w-6xl px-6 py-10 space-y-8">
      <header class="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
        <div>
          <p class="text-sm uppercase tracking-[0.3em] text-slate-400 dark:text-slate-500">PROJECT</p>
          <h1 class="font-display text-3xl font-semibold tracking-tight text-white dark:text-slate-200 sm:text-4xl">
            {{ project?.name || 'Untitled Project' }}
          </h1>
          <p class="mt-2 max-w-2xl text-slate-400 dark:text-slate-500">
            {{ project?.description || 'Manage your XCStrings localization project.' }}
          </p>
        </div>
        <div class="flex flex-wrap items-center gap-3">
          <button
            class="rounded-full bg-mint px-4 py-2 text-sm font-semibold text-midnight shadow-glow transition hover:shadow-neon"
            @click="startTranslation"
          >
            Start Translation
          </button>
          <button
            class="rounded-full border border-white/20 px-4 py-2 text-sm font-medium text-slate-100 transition hover:border-mint/60 hover:text-mint"
            @click="exportProject"
          >
            Export
          </button>
          <button
            class="rounded-full border border-white/20 px-4 py-2 text-sm font-medium text-slate-100 transition hover:border-mint/60 hover:text-mint"
            @click="loadProject"
          >
            Refresh
          </button>
          <button
            class="rounded-full border border-white/20 px-4 py-2 text-sm font-medium text-slate-100 transition hover:border-mint/60 hover:text-mint"
            @click="logout"
          >
            Logout
          </button>
        </div>
      </header>

      <div class="grid gap-4 lg:grid-cols-3">
        <section class="glass rounded-2xl p-6 lg:col-span-2">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm uppercase tracking-[0.2em] text-slate-400 dark:text-slate-500">Workspace</p>
              <h2 class="font-display text-xl font-semibold">File & Languages</h2>
            </div>
            <span v-if="project?.fileName" class="rounded-full bg-emerald-900/40 px-3 py-1 text-xs font-semibold text-mint">Loaded</span>
          </div>

          <div class="mt-5 grid gap-4 md:grid-cols-3">
            <div class="rounded-xl border border-white/10 bg-white/5 p-4">
              <p class="text-xs uppercase tracking-wide text-slate-400 dark:text-slate-500">File</p>
              <p class="mt-1 font-semibold" v-if="project?.fileName">{{ project.fileName }}</p>
              <p class="mt-1 text-slate-400 dark:text-slate-500" v-else>No file uploaded</p>
              <p class="mt-2 text-xs text-slate-400 dark:text-slate-500">{{ state.totalStrings }} strings detected</p>
            </div>
            <div class="rounded-xl border border-white/10 bg-white/5 p-4">
              <p class="text-xs uppercase tracking-wide text-slate-400 dark:text-slate-500">Source language</p>
              <div class="mt-1 flex items-center gap-2">
                <input
                  v-model="state.sourceLanguage"
                  class="w-full rounded-lg bg-midnight/50 px-3 py-2 text-sm text-white ring-1 ring-white/10 focus:ring-2 focus:ring-mint"
                  placeholder="e.g. en"
                />
                <button class="rounded-lg border border-white/20 px-3 py-2 text-xs" @click="autoDetectSource">Auto</button>
              </div>
              <p class="mt-2 text-xs text-slate-400 dark:text-slate-500">Hint: we use this to prompt the provider.</p>
            </div>
            <div class="rounded-xl border border-white/10 bg-white/5 p-4">
              <p class="text-xs uppercase tracking-wide text-slate-400 dark:text-slate-500">Available languages</p>
              <div class="mt-2 flex flex-wrap gap-2">
                <span
                  v-for="lang in state.availableLanguages"
                  :key="lang"
                  class="rounded-full border border-white/10 px-2 py-1 text-xs text-slate-200"
                >
                  {{ lang }}
                </span>
                <span v-if="!state.availableLanguages.length" class="text-xs text-slate-400 dark:text-slate-500">–</span>
              </div>
            </div>
          </div>

          <div class="mt-6 rounded-xl border border-white/10 bg-white/5 p-4">
            <div class="flex items-center justify-between">
              <div>
                <p class="text-xs uppercase tracking-wide text-slate-400 dark:text-slate-500">Target languages</p>
                <p class="text-sm text-slate-300 dark:text-slate-400">Add one per code, press enter to confirm.</p>
              </div>
              <button
                class="rounded-full border border-white/20 px-3 py-1 text-xs text-slate-200 hover:border-mint/60 hover:text-mint"
                @click="useAvailableTargets"
              >
                Use detected
              </button>
            </div>

            <div class="mt-3 flex flex-wrap gap-2">
              <span
                v-for="lang in state.targetLanguages"
                :key="lang"
                class="group inline-flex items-center gap-2 rounded-full bg-white/10 px-3 py-1 text-sm"
              >
                {{ lang }}
                <button class="text-slate-400 transition group-hover:text-white" @click="removeTarget(lang)">×</button>
              </span>
              <input
                v-model="targetInput"
                class="tag-input rounded-full border border-dashed border-white/30 bg-midnight/40 px-3 py-1 text-sm"
                placeholder="Add language"
                @keydown.enter.prevent="addTarget"
                @blur="addTarget"
              />
            </div>

            <div class="mt-3 flex flex-wrap gap-2 text-xs text-slate-400 dark:text-slate-500">
              <button
                v-for="preset in presets"
                :key="preset"
                class="rounded-full border border-white/20 px-2 py-1 hover:border-mint/60 hover:text-mint"
                @click="addPreset(preset)"
              >
                {{ preset }}
              </button>
            </div>

            <div class="mt-4 space-y-2">
              <div class="flex items-center justify-between text-xs text-slate-400 dark:text-slate-500">
                <span>Language library</span>
                <button class="text-mint underline" @click="showAllLanguages = !showAllLanguages">
                  {{ showAllLanguages ? 'Collapse' : 'Show all' }}
                </button>
              </div>
              <div class="flex flex-wrap gap-2">
                <button
                  v-for="lang in languageOptions"
                  :key="lang.code"
                  class="rounded-full border border-white/20 px-3 py-1 text-xs text-slate-100 transition hover:border-mint/60 hover:text-mint"
                  @click="addPreset(lang.code)"
                >
                  {{ lang.name }} ({{ lang.code }})
                </button>
              </div>
            </div>
          </div>
        </section>

        <section class="glass rounded-2xl p-6 space-y-4">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm uppercase tracking-[0.2em] text-slate-400 dark:text-slate-500">Provider</p>
              <h2 class="font-display text-xl font-semibold">Batch translate</h2>
            </div>
            <span class="text-xs text-slate-400 dark:text-slate-500">{{ providerLabel }}</span>
          </div>

          <div class="grid grid-cols-2 gap-3">
            <button
              v-for="p in providers"
              :key="p.id"
              class="rounded-xl border px-3 py-2 text-left transition"
              :class="p.id === state.provider ? 'border-mint/60 bg-mint/10 text-white' : 'border-white/10 bg-white/5 text-slate-300'"
              @click="state.provider = p.id"
            >
              <p class="font-semibold">{{ p.name }}</p>
              <p class="text-xs text-slate-400 dark:text-slate-500">{{ p.hint }}</p>
            </button>
          </div>

          <div class="space-y-3 text-sm text-slate-200">
            <template v-if="state.provider === 'openai'">
              <label class="block">
                <span class="text-xs text-slate-400 dark:text-slate-500">API key</span>
                <input v-model="state.openai.apiKey" class="mt-1 w-full rounded-lg bg-midnight/40 px-3 py-2 ring-1 ring-white/10 focus:ring-2 focus:ring-mint" type="password" placeholder="sk-…" />
              </label>
              <label class="block">
                <span class="text-xs text-slate-400 dark:text-slate-500">Model</span>
                <input v-model="state.openai.model" class="mt-1 w-full rounded-lg bg-midnight/40 px-3 py-2 ring-1 ring-white/10 focus:ring-2 focus:ring-mint" />
              </label>
              <label class="block">
                <span class="text-xs text-slate-400 dark:text-slate-500">API base (optional)</span>
                <input v-model="state.openai.apiBaseUrl" class="mt-1 w-full rounded-lg bg-midnight/40 px-3 py-2 ring-1 ring-white/10 focus:ring-2 focus:ring-mint" placeholder="https://api.openai.com" />
              </label>
            </template>
            <template v-else-if="state.provider === 'google'">
              <label class="block">
                <span class="text-xs text-slate-400 dark:text-slate-500">Google API key</span>
                <input v-model="state.google.apiKey" class="mt-1 w-full rounded-lg bg-midnight/40 px-3 py-2 ring-1 ring-white/10 focus:ring-2 focus:ring-mint" placeholder="AIza…" />
              </label>
            </template>
            <template v-else-if="state.provider === 'deepl'">
              <label class="block">
                <span class="text-xs text-slate-400 dark:text-slate-500">DeepL API key</span>
                <input v-model="state.deepl.apiKey" class="mt-1 w-full rounded-lg bg-midnight/40 px-3 py-2 ring-1 ring-white/10 focus:ring-2 focus:ring-mint" placeholder="auth-key" />
              </label>
              <label class="block">
                <span class="text-xs text-slate-400 dark:text-slate-500">Formality</span>
                <input v-model="state.deepl.formality" class="mt-1 w-full rounded-lg bg-midnight/40 px-3 py-2 ring-1 ring-white/10 focus:ring-2 focus:ring-mint" placeholder="default" />
              </label>
              <label class="inline-flex items-center gap-2 text-xs text-slate-300 dark:text-slate-400">
                <input v-model="state.deepl.isFree" type="checkbox" class="h-4 w-4 rounded border-white/30 bg-midnight/50" />
                Use free API
              </label>
            </template>
            <template v-else>
              <label class="block">
                <span class="text-xs text-slate-400 dark:text-slate-500">Baidu App ID</span>
                <input v-model="state.baidu.appId" class="mt-1 w-full rounded-lg bg-midnight/40 px-3 py-2 ring-1 ring-white/10 focus:ring-2 focus:ring-mint" />
              </label>
              <label class="block">
                <span class="text-xs text-slate-400 dark:text-slate-500">Baidu App Secret</span>
                <input v-model="state.baidu.appSecret" class="mt-1 w-full rounded-lg bg-midnight/40 px-3 py-2 ring-1 ring-white/10 focus:ring-2 focus:ring-mint" type="password" />
              </label>
            </template>

            <div class="grid grid-cols-2 gap-3 text-xs text-slate-300 dark:text-slate-400">
              <label class="block">
                <span>Concurrency</span>
                <input v-model.number="state.concurrency" type="number" min="1" class="mt-1 w-full rounded-lg bg-midnight/40 px-3 py-2 ring-1 ring-white/10 focus:ring-2 focus:ring-mint" />
              </label>
              <label class="block">
                <span>Timeout (sec)</span>
                <input v-model.number="state.timeoutSeconds" type="number" min="30" class="mt-1 w-full rounded-lg bg-midnight/40 px-3 py-2 ring-1 ring-white/10 focus:ring-2 focus:ring-mint" />
              </label>
            </div>

            <div class="grid grid-cols-2 gap-3">
              <button
                class="w-full rounded-xl bg-mint px-4 py-3 text-center text-sm font-semibold text-midnight shadow-lg shadow-mint/20 transition hover:shadow-mint/40 disabled:cursor-not-allowed disabled:opacity-60"
                :disabled="isTranslating || !project?.fileName || !state.targetLanguages.length"
                @click="batchTranslate"
              >
                <span v-if="isTranslating">Translating…</span>
                <span v-else>Run batch translate</span>
              </button>
              <button
                class="w-full rounded-xl border border-white/20 px-4 py-3 text-center text-sm font-semibold text-slate-100 transition hover:border-mint/60 hover:text-mint"
                @click="saveOptions"
              >
                Save options
              </button>
            </div>
            <p v-if="progress.id" class="text-xs text-slate-300 dark:text-slate-400">
              Progress: {{ progress.done }}/{{ progress.total || '…' }} • {{ progress.status }}
            </p>
            <p class="text-xs text-slate-400 dark:text-slate-500">We only translate missing entries for the selected targets.</p>
          </div>
        </section>
      </div>

      <section class="glass rounded-2xl p-6">
        <div class="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
          <div>
            <p class="text-sm uppercase tracking-[0.2em] text-slate-400 dark:text-slate-500">Strings</p>
            <h2 class="font-display text-xl font-semibold">Visual localisation</h2>
            <p class="text-slate-400 dark:text-slate-500">Browse source text next to target translations. Missing items are highlighted.</p>
          </div>
          <div class="flex items-center gap-2 text-xs text-slate-300 dark:text-slate-400">
            <label class="flex items-center gap-2">
              <span>Filter</span>
              <input v-model="filter" class="rounded-lg bg-midnight/50 px-3 py-2 ring-1 ring-white/10 focus:ring-2 focus:ring-mint" placeholder="Search key or text" />
            </label>
          </div>
        </div>

        <div class="mt-4 overflow-hidden rounded-xl border border-white/10 bg-white/5">
          <div class="max-h-[480px] overflow-auto scrollbar-thin">
            <table class="min-w-full divide-y divide-white/5 text-sm">
              <thead class="sticky top-0 bg-midnight">
                <tr>
                  <th class="px-4 py-3 text-left font-semibold text-slate-300 dark:text-slate-400">Key</th>
                  <th class="px-4 py-3 text-left font-semibold text-slate-300 dark:text-slate-400">{{ state.sourceLanguage || 'source' }}</th>
                  <th
                    v-for="lang in displayTargets"
                    :key="lang"
                    class="px-4 py-3 text-left font-semibold text-slate-300 dark:text-slate-400"
                  >
                    {{ lang }}
                  </th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="row in filteredEntries" :key="row.key" class="divide-y divide-white/5">
                  <td class="px-4 py-3 align-top font-mono text-xs text-slate-300 dark:text-slate-400">{{ row.key }}</td>
                  <td class="px-4 py-3 align-top text-slate-100 dark:text-slate-200">
                    <p class="whitespace-pre-line">{{ row.source || '—' }}</p>
                    <p class="mt-1 text-xs text-slate-400 dark:text-slate-500">{{ row.state }}</p>
                  </td>
                  <td
                    v-for="lang in displayTargets"
                    :key="lang"
                    class="px-4 py-3 align-top"
                    :class="row.missing.includes(lang) ? 'bg-orange-500/5 text-orange-200' : 'text-slate-100 dark:text-slate-200'"
                  >
                    <p class="whitespace-pre-line">{{ row.translations[lang] || '–' }}</p>
                  </td>
                </tr>
                <tr v-if="!filteredEntries.length">
                  <td class="px-4 py-6 text-center text-slate-400 dark:text-slate-500" :colspan="2 + displayTargets.length">No strings yet.</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </section>

      <!-- File upload section -->
      <section class="glass rounded-2xl p-6">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm uppercase tracking-[0.2em] text-slate-400 dark:text-slate-500">Upload</p>
            <h2 class="font-display text-xl font-semibold">Project File</h2>
          </div>
        </div>
        <div class="mt-4">
          <div class="flex items-center justify-center w-full">
            <label class="flex flex-col items-center justify-center w-full h-32 border-2 border-dashed rounded-2xl cursor-pointer border-white/20 hover:border-mint/60 bg-white/5 hover:bg-mint/5 transition">
              <div class="flex flex-col items-center justify-center pt-5 pb-6">
                <svg class="w-8 h-8 mb-4 text-slate-500 dark:text-slate-400" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 20 16">
                  <path stroke="currentColor" stroke-dasharray="3,3" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 13h3a3 3 0 0 0 0-6h-.025A5.56 5.56 0 0 0 16 6.5 5.5 5.5 0 0 0 5.207 5.021C5.137 5.017 5.071 5 5 5a4 4 0 0 0 0 8h2.167M10 15V6m0 0L8 8m2-2 2 2"/>
                </svg>
                <p class="mb-2 text-sm text-slate-500 dark:text-slate-400"><span class="font-semibold">Click to upload</span> or drag and drop</p>
                <p class="text-xs text-slate-500 dark:text-slate-500">.xcstrings files only</p>
              </div>
              <input 
                ref="fileInput" 
                type="file" 
                accept=".xcstrings,application/json" 
                class="hidden" 
                @change="onFileChange" 
              />
            </label>
          </div>
        </div>
      </section>

      <div v-if="statusMessage" class="rounded-xl border border-white/10 px-4 py-3 text-sm" :class="statusClass">
        {{ statusMessage }}
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const props = defineProps({
  id: {
    type: String,
    required: true
  }
})

// Project and state management
const project = ref<any>(null)
const fileInput = ref<HTMLInputElement | null>(null)

type LocalizationEntry = {
  key: string
  source: string
  state: string
  translations: Record<string, string>
  missing: string[]
}

type Payload = {
  fileName: string
  sourceLanguage: string
  availableLanguages: string[]
  totalStrings: number
  entries: LocalizationEntry[]
  warning?: string
}

type ProviderId = 'openai' | 'google' | 'deepl' | 'baidu'
type JobState = { id: string; status: string; done: number; total: number; message?: string }

type User = {
  id: number;
  username: string;
  email: string;
  isActive: boolean;
}

type ProjectForm = {
  name: string;
  description: string;
  sourceLanguage: string;
}

const presets = ['zh-Hans', 'ja', 'ko', 'de', 'fr', 'es', 'ar']
const languages = [
  { code: 'zh-Hans', name: 'Chinese (Simplified)' },
  { code: 'zh-Hant', name: 'Chinese (Traditional)' },
  { code: 'en', name: 'English' },
  { code: 'ja', name: 'Japanese' },
  { code: 'ko', name: 'Korean' },
  { code: 'de', name: 'German' },
  { code: 'fr', name: 'French' },
  { code: 'es', name: 'Spanish' },
  { code: 'ar', name: 'Arabic' },
  { code: 'pt', name: 'Portuguese' },
  { code: 'ru', name: 'Russian' },
  { code: 'it', name: 'Italian' },
  { code: 'nl', name: 'Dutch' },
  { code: 'pl', name: 'Polish' },
  { code: 'sv', name: 'Swedish' },
  { code: 'da', name: 'Danish' },
  { code: 'fi', name: 'Finnish' },
  { code: 'no', name: 'Norwegian' },
  { code: 'cs', name: 'Czech' },
  { code: 'tr', name: 'Turkish' },
  { code: 'he', name: 'Hebrew' },
  { code: 'hi', name: 'Hindi' },
  { code: 'id', name: 'Indonesian' },
  { code: 'th', name: 'Thai' },
  { code: 'vi', name: 'Vietnamese' }
]
const providers = [
  { id: 'openai' as ProviderId, name: 'OpenAI', hint: 'GPT style chat completion' },
  { id: 'google' as ProviderId, name: 'Google', hint: 'Google Cloud translation' },
  { id: 'deepl' as ProviderId, name: 'DeepL', hint: 'Great for EU languages' },
  { id: 'baidu' as ProviderId, name: 'Baidu', hint: 'China-friendly' }
]

const state = reactive({
  fileName: '',
  sourceLanguage: 'en',
  targetLanguages: [] as string[],
  availableLanguages: [] as string[],
  totalStrings: 0,
  entries: [] as LocalizationEntry[],
  provider: 'openai' as ProviderId,
  openai: { apiKey: '', apiBaseUrl: '', model: 'gpt-3.5-turbo', temperature: 0.3, maxTokens: 1024 },
  google: { apiKey: '' },
  deepl: { apiKey: '', isFree: true, formality: 'default' },
  baidu: { appId: '', appSecret: '' },
  concurrency: 6,
  timeoutSeconds: 300
})

const isTranslating = ref(false)
const progress = reactive<JobState>({ id: '', status: 'idle', done: 0, total: 0 })
let progressTimer: number | null = null
const statusMessage = ref('')
const statusTone = ref<'info' | 'error'>('info')
const filter = ref('')
const targetInput = ref('')
const showAllLanguages = ref(false)

const LOCAL_KEY = 'xcstrings-translator-ui'

const providerLabel = computed(() => providers.find((p) => p.id === state.provider)?.name ?? '')
const displayTargets = computed(() => state.targetLanguages)
const filteredEntries = computed(() => {
  const term = filter.value.trim().toLowerCase()
  if (!term) return state.entries
  return state.entries.filter((row) =>
    row.key.toLowerCase().includes(term) ||
    row.source.toLowerCase().includes(term) ||
    displayTargets.value.some((lang) => (row.translations[lang] || '').toLowerCase().includes(term))
  )
})

const statusClass = computed(() =>
  statusTone.value === 'error'
    ? 'text-red-200 border-red-400/30 bg-red-900/30'
    : 'text-mint border-mint/50 bg-mint/10'
)

const languageOptions = computed(() => (showAllLanguages.value ? languages : languages.slice(0, 20)))

// Helper functions
function showStatus(message: string, tone: 'info' | 'error' = 'info') {
  statusMessage.value = message
  statusTone.value = tone
}

function triggerUpload() {
  fileInput.value?.click()
}

function addTarget() {
  const trimmed = targetInput.value.trim()
  if (!trimmed) return
  if (!state.targetLanguages.includes(trimmed)) {
    state.targetLanguages.push(trimmed)
  }
  targetInput.value = ''
}

function addPreset(code: string) {
  if (!state.targetLanguages.includes(code)) {
    state.targetLanguages.push(code)
  }
}

function removeTarget(lang: string) {
  state.targetLanguages = state.targetLanguages.filter((l) => l !== lang)
}

function useAvailableTargets() {
  state.targetLanguages = state.availableLanguages.filter((lang) => lang !== state.sourceLanguage)
}

function autoDetectSource() {
  if (state.availableLanguages.includes(state.sourceLanguage)) return
  if (state.availableLanguages.length) {
    state.sourceLanguage = state.availableLanguages[0]
  }
}

function formatDate(dateString: string) {
  const date = new Date(dateString)
  return date.toLocaleDateString()
}

// API functions
async function loadProject() {
  const token = localStorage.getItem('token')
  if (!token) {
    router.push('/login')
    return
  }

  try {
    const res = await fetch(`/api/protected/projects/${props.id}`, {
      headers: { 'Authorization': `Bearer ${token}` }
    })

    if (!res.ok) {
      throw new Error(`Failed to load project: ${await res.text()}`)
    }

    const data = await res.json()
    project.value = data.project
  } catch (err) {
    showStatus('Failed to load project: ' + (err as Error).message, 'error')
  }
}

async function onFileChange(event: Event) {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return

  const token = localStorage.getItem('token')
  if (!token) {
    router.push('/login')
    return
  }

  const formData = new FormData()
  formData.append('file', file)
  if (state.sourceLanguage) {
    formData.append('sourceLanguage', state.sourceLanguage)
  }

  try {
    const res = await fetch(`/api/protected/projects/${props.id}/upload`, {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${token}` },
      body: formData
    })

    if (!res.ok) {
      throw new Error(`Upload failed: ${await res.text()}`)
    }

    const payload = (await res.json()) as Payload
    applyPayload(payload)
    if (!state.targetLanguages.length) {
      useAvailableTargets()
    }
    if (payload.warning) {
      showStatus(payload.warning, 'error')
    } else {
      showStatus('File uploaded to project.', 'info')
    }
  } catch (err) {
    showStatus('Upload failed: ' + (err as Error).message, 'error')
  } finally {
    input.value = ''
  }
}

function applyPayload(payload: Payload) {
  state.fileName = payload.fileName
  state.sourceLanguage = payload.sourceLanguage
  state.availableLanguages = payload.availableLanguages
  state.totalStrings = payload.totalStrings
  state.entries = payload.entries
}

async function batchTranslate() {
  if (!project.value?.fileName) {
    showStatus('Upload a file first.', 'error')
    return
  }
  if (!state.targetLanguages.length) {
    showStatus('Add at least one target language.', 'error')
    return
  }

  const token = localStorage.getItem('token')
  if (!token) {
    router.push('/login')
    return
  }

  const body = {
    provider: state.provider,
    targetLanguages: state.targetLanguages,
    sourceLanguage: state.sourceLanguage,
    concurrency: state.concurrency,
    timeoutSeconds: state.timeoutSeconds,
    config: {
      apiKey: state.provider === 'baidu' ? undefined : getApiKey(),
      apiBaseUrl: state.openai.apiBaseUrl,
      model: state.provider === 'openai' ? state.openai.model : undefined,
      temperature: state.provider === 'openai' ? state.openai.temperature : undefined,
      maxTokens: state.provider === 'openai' ? state.openai.maxTokens : undefined,
      appId: state.baidu.appId,
      appSecret: state.baidu.appSecret,
      formality: state.deepl.formality,
      isFree: state.deepl.isFree
    }
  }

  isTranslating.value = true
  showStatus('Running batch translation…')

  try {
    const res = await fetch(`/api/protected/projects/${props.id}/translate`, {
      method: 'POST',
      headers: { 
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json' 
      },
      body: JSON.stringify(body)
    })

    if (!res.ok) {
      throw new Error(`Translate failed: ${await res.text()}`)
    }

    const { jobId } = (await res.json()) as { jobId: string }
    if (!jobId) {
      throw new Error('Translate job did not start.')
    }
    startProgress(jobId)
  } catch (err) {
    isTranslating.value = false
    showStatus('Translate failed: ' + (err as Error).message, 'error')
  }
}

function getApiKey() {
  if (state.provider === 'openai') return state.openai.apiKey
  if (state.provider === 'google') return state.google.apiKey
  return state.deepl.apiKey
}

async function exportProject() {
  const token = localStorage.getItem('token')
  if (!token) {
    router.push('/login')
    return
  }

  try {
    const res = await fetch(`/api/protected/projects/${props.id}/export`, {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    if (!res.ok) {
      throw new Error('Export failed')
    }
    const blob = await res.blob()
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = project.value?.fileName || `project_${props.id}_translated.xcstrings`
    link.click()
    URL.revokeObjectURL(url)
    showStatus('Project exported.', 'info')
  } catch (err) {
    showStatus('Export failed: ' + (err as Error).message, 'error')
  }
}

// Progress tracking
function startProgress(id: string) {
  progress.id = id
  progress.status = 'running'
  progress.done = 0
  progress.total = 0
  isTranslating.value = true
  if (progressTimer) {
    clearInterval(progressTimer)
  }
  pollProgress()
  progressTimer = window.setInterval(pollProgress, 1200)
}

async function pollProgress() {
  const token = localStorage.getItem('token')
  if (!token) {
    router.push('/login')
    return
  }

  try {
    const res = await fetch(`/api/protected/projects/${props.id}/translate/progress`, {
      headers: { 
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    })
    
    if (!res.ok) return
    
    const data = (await res.json()) as { job?: JobState | null; payload?: Payload }

    if (data.payload) {
      applyPayload(data.payload)
    }

    if (data.job) {
      progress.id = data.job.id
      progress.status = data.job.status
      progress.done = data.job.done
      progress.total = data.job.total
      if (data.job.status !== 'running') {
        stopProgress()
        showStatus(data.job.status === 'done' ? 'Translations applied.' : data.job.message || 'Translation stopped.', data.job.status === 'done' ? 'info' : 'error')
      }
    } else if (progress.id) {
      stopProgress()
    }
  } catch (err) {
    console.error('Error polling progress:', err)
  }
}

function stopProgress() {
  if (progressTimer) {
    clearInterval(progressTimer)
    progressTimer = null
  }
  isTranslating.value = false
  progress.id = ''
}

function saveOptions() {
  saveLocalState(snapshotOptions())
  showStatus('Options saved locally.', 'info')
}

function snapshotOptions() {
  return {
    provider: state.provider,
    sourceLanguage: state.sourceLanguage,
    targetLanguages: [...state.targetLanguages],
    concurrency: state.concurrency,
    timeoutSeconds: state.timeoutSeconds,
    openai: { ...state.openai },
    google: { ...state.google },
    deepl: { ...state.deepl },
    baidu: { ...state.baidu }
  }
}

function loadLocalState() {
  try {
    const raw = localStorage.getItem(LOCAL_KEY)
    if (!raw) return null
    return JSON.parse(raw)
  } catch (err) {
    console.warn('Failed to load local state', err)
    return null
  }
}

function saveLocalState(val: unknown) {
  try {
    localStorage.setItem(LOCAL_KEY, JSON.stringify(val))
  } catch (err) {
    console.warn('Failed to persist local state', err)
  }
}

async function startTranslation() {
  await batchTranslate()
}

function logout() {
  localStorage.removeItem('token')
  router.push('/login')
}

onMounted(async () => {
  // Check authentication
  const token = localStorage.getItem('token')
  if (!token) {
    router.push('/login')
    return
  }

  // Load project
  await loadProject()
  
  // Load saved state
  const saved = loadLocalState()
  if (saved) {
    Object.assign(state, saved)
  }
})

watch(
  () => snapshotOptions(),
  (val) => saveLocalState(val),
  { deep: true }
)
</script>
