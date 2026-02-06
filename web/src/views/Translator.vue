<template>
  <div class="space-y-6">
    <header class="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
      <div>
        <p class="text-xs uppercase tracking-[0.3em] text-slate-500">{{ t('nav.translator') }}</p>
        <h1 class="text-2xl font-semibold">{{ t('translator.title') }}</h1>
        <p class="text-sm text-slate-400">{{ t('translator.subtitle') }}</p>
      </div>
    </header>

    <section class="rounded-2xl border border-white/10 bg-white/5 p-6 space-y-6">
      <div class="space-y-2">
        <label class="text-sm text-slate-400">{{ t('translator.file') }}</label>
        <input
          type="file"
          accept=".xcstrings,application/json"
          @change="handleFileChange"
          class="block w-full text-sm text-slate-200 file:mr-4 file:rounded-lg file:border-0 file:bg-mint file:px-4 file:py-2 file:text-sm file:font-semibold file:text-midnight hover:file:bg-mint/90"
        />
        <p class="text-xs text-slate-500">{{ selectedFileName }}</p>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
          <label class="text-sm text-slate-400">{{ t('translator.sourceLanguage') }}</label>
          <input
            v-model="form.sourceLanguage"
            class="mt-1 w-full rounded-lg bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint"
            placeholder="en"
          />
        </div>
        <div>
          <label class="text-sm text-slate-400">{{ t('translator.targetLanguages') }}</label>
          <input
            v-model="form.targetLanguages"
            class="mt-1 w-full rounded-lg bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint"
            :placeholder="t('translator.targetPlaceholder')"
          />
        </div>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div>
          <label class="text-sm text-slate-400">{{ t('translator.provider') }}</label>
          <select
            v-model="form.provider"
            class="mt-1 w-full rounded-lg bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint"
          >
            <option value="openai">OpenAI</option>
            <option value="llama">Hunyuan (Llama)</option>
          </select>
        </div>
        <div>
          <label class="text-sm text-slate-400">{{ t('translator.concurrency') }}</label>
          <input
            v-model.number="form.concurrency"
            type="number"
            min="1"
            class="mt-1 w-full rounded-lg bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint"
          />
        </div>
        <div>
          <label class="text-sm text-slate-400">{{ t('translator.timeout') }}</label>
          <input
            v-model.number="form.timeoutSeconds"
            type="number"
            min="30"
            class="mt-1 w-full rounded-lg bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint"
          />
        </div>
      </div>

      <div v-if="form.provider === 'openai'" class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
          <label class="text-sm text-slate-400">{{ t('translator.apiKey') }}</label>
          <input
            v-model="form.config.apiKey"
            type="password"
            class="mt-1 w-full rounded-lg bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint"
          />
        </div>
        <div>
          <label class="text-sm text-slate-400">{{ t('translator.model') }}</label>
          <input
            v-model="form.config.model"
            class="mt-1 w-full rounded-lg bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint"
            placeholder="gpt-4o-mini"
          />
        </div>
      </div>

      <div class="flex flex-wrap gap-2">
        <button
          class="rounded-lg bg-mint px-4 py-2 text-sm font-semibold text-midnight shadow hover:bg-mint/90"
          :disabled="!fileReady || translating"
          @click="startTranslation"
        >
          {{ translating ? t('translator.translating') : t('translator.translate') }}
        </button>
        <button
          class="rounded-lg border border-white/20 px-4 py-2 text-sm hover:border-mint/60 hover:text-mint"
          :disabled="!projectId"
          @click="exportResult"
        >
          {{ t('translator.export') }}
        </button>
      </div>

      <div v-if="statusMessage" class="text-sm text-slate-300">{{ statusMessage }}</div>
      <div v-if="jobId" class="text-xs text-slate-500">{{ t('translator.jobId') }}: {{ jobId }}</div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { computed, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useApi } from '../composables/useApi'

const { t } = useI18n()
const { api } = useApi()

const selectedFile = ref<File | null>(null)
const projectId = ref<number | null>(null)
const jobId = ref<string | null>(null)
const translating = ref(false)
const statusMessage = ref('')

const form = reactive({
  sourceLanguage: 'en',
  targetLanguages: 'zh-Hans,ja',
  provider: 'openai',
  concurrency: 4,
  timeoutSeconds: 300,
  config: {
    apiKey: '',
    apiBaseUrl: '',
    model: ''
  }
})

const selectedFileName = computed(() => selectedFile.value?.name || t('translator.noFile'))
const fileReady = computed(() => !!selectedFile.value)

function handleFileChange(event: Event) {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0]
  if (!file) {
    return
  }
  selectedFile.value = file
  statusMessage.value = ''
  jobId.value = null
  projectId.value = null
}

async function startTranslation() {
  if (!selectedFile.value) {
    return
  }

  translating.value = true
  statusMessage.value = t('translator.uploading')

  try {
    const uploadResponse = await api.createProjectFromFile(selectedFile.value, form.sourceLanguage)
    if (!uploadResponse.success || !uploadResponse.project) {
      throw new Error(uploadResponse.message || 'Upload failed')
    }
    projectId.value = uploadResponse.project.id
    statusMessage.value = t('translator.translating')

    const targets = form.targetLanguages
      .split(',')
      .map((lang) => lang.trim())
      .filter(Boolean)

    const translateResponse = await api.translateProject(projectId.value, {
      provider: form.provider,
      sourceLanguage: form.sourceLanguage,
      targetLanguages: targets,
      concurrency: form.concurrency,
      timeoutSeconds: form.timeoutSeconds,
      config: form.config
    })

    jobId.value = translateResponse.jobId
    statusMessage.value = t('translator.translateSubmitted')
  } catch (error) {
    statusMessage.value = error instanceof Error ? error.message : t('translator.translateError')
  } finally {
    translating.value = false
  }
}

async function exportResult() {
  if (!projectId.value) {
    return
  }

  const blob = await api.exportProject(projectId.value)
  const url = window.URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = 'Localizable_translated.xcstrings'
  link.click()
  window.URL.revokeObjectURL(url)
}
</script>
