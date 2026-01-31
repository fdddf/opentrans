<template>
  <div class="space-y-6">
    <header class="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
      <div>
        <p class="text-xs uppercase tracking-[0.3em] text-slate-500">{{ t('nav.appleConnectConfig') }}</p>
        <h1 class="text-2xl font-semibold">{{ t('appleConnectConfig.title') }}</h1>
        <p class="text-sm text-slate-400">{{ t('appleConnectConfig.subtitle') }}</p>
      </div>
      <div class="flex items-center gap-2">
        <button 
          class="rounded-lg bg-mint px-4 py-2 text-sm font-semibold text-midnight shadow hover:bg-mint/90"
          @click="showForm = true; editingConfig = null; resetForm()"
        >
          {{ t('appleConnectConfig.addNew') }}
        </button>
      </div>
    </header>

    <!-- Connection Status Notification -->
    <div v-if="connectionStatus" 
         :class="[
           'p-4 rounded-2xl border border-white/10',
           connectionStatus.success ? 'bg-emerald-900/20 border-emerald-500/30' : 'bg-rose-900/20 border-rose-500/30'
         ]"
    >
      <div class="flex items-start justify-between">
        <div class="flex items-start">
          <span :class="connectionStatus.success ? 'text-emerald-400' : 'text-rose-400'">
            {{ connectionStatus.success ? '✓' : '✕' }}
          </span>
          <div class="ml-3">
            <p class="text-sm font-medium" 
               :class="connectionStatus.success ? 'text-emerald-100' : 'text-rose-100'">
              {{ connectionStatus.title }}
            </p>
            <p class="text-sm mt-1" 
               :class="connectionStatus.success ? 'text-emerald-200' : 'text-rose-200'">
              {{ connectionStatus.message }}
            </p>
          </div>
        </div>
        <button 
          @click="connectionStatus = null"
          class="text-slate-400 hover:text-white"
        >
          ×
        </button>
      </div>
    </div>

    <!-- Configurations List -->
    <div class="space-y-4">
      <div v-if="loading" class="flex items-center justify-center py-12">
        <div class="h-6 w-6 rounded-full border-2 border-mint border-t-transparent animate-spin"></div>
      </div>

      <div v-else-if="configs.length === 0" class="rounded-2xl border border-white/10 bg-white/5 p-12 text-center">
        <div class="text-4xl mb-4">🍎</div>
        <h3 class="text-lg font-semibold mb-2">{{ t('appleConnectConfig.noConfigs.title') }}</h3>
        <p class="text-sm text-slate-400 mb-4">{{ t('appleConnectConfig.noConfigs.message') }}</p>
        <button 
          class="rounded-lg bg-mint px-4 py-2 text-sm font-semibold text-midnight shadow hover:bg-mint/90"
          @click="showForm = true; editingConfig = null; resetForm()"
        >
          {{ t('appleConnectConfig.addFirst') }}
        </button>
      </div>

      <div v-else class="space-y-4">
        <div 
          v-for="config in configs" 
          :key="config.id"
          class="rounded-2xl border border-white/10 bg-white/5 p-6 hover:border-white/20 transition-colors"
        >
          <div class="flex flex-col gap-4 md:flex-row md:items-start md:justify-between">
            <div class="flex-1 space-y-3">
              <div class="flex items-center gap-2">
                <h3 class="font-semibold">{{ config.issuerId }}</h3>
                <span 
                  v-if="config.isDefault"
                  class="rounded-full bg-mint/20 px-2 py-0.5 text-xs text-mint"
                >
                  {{ t('common.default') }}
                </span>
              </div>
              <div class="space-y-1 text-sm">
                <div class="flex items-center gap-2 text-slate-400">
                  <span class="w-20">{{ t('appleConnectConfig.keyId') }}:</span>
                  <span class="font-mono">{{ config.keyId }}</span>
                </div>
                <div class="flex items-center gap-2 text-slate-400">
                  <span class="w-20">{{ t('appleConnectConfig.createdAt') }}:</span>
                  <span>{{ formatDate(config.createdAt) }}</span>
                </div>
              </div>
            </div>

            <div class="flex items-center gap-2">
              <button 
                class="rounded-lg border border-white/20 px-3 py-2 text-sm hover:border-mint/60 hover:text-mint"
                @click="testConnection(config.id)"
                :disabled="testingConfigId === config.id"
              >
                <span v-if="testingConfigId !== config.id">{{ t('appleConnectConfig.testConnection') }}</span>
                <span v-else class="flex items-center gap-1">
                  <span class="h-3 w-3 rounded-full bg-mint animate-pulse"></span>
                  {{ t('appleConnectConfig.testing') }}
                </span>
              </button>
              <button 
                class="rounded-lg border border-white/20 px-3 py-2 text-sm hover:border-white/40"
                @click="editConfig(config)"
              >
                {{ t('common.edit') }}
              </button>
              <button 
                class="rounded-lg border border-rose-500/30 px-3 py-2 text-sm text-rose-400 hover:border-rose-500/60 hover:bg-rose-500/10"
                @click="deleteConfig(config.id)"
                :disabled="deletingConfigId === config.id"
              >
                <span v-if="deletingConfigId !== config.id">{{ t('common.delete') }}</span>
                <span v-else class="flex items-center gap-1">
                  <span class="h-3 w-3 rounded-full bg-rose-400 animate-pulse"></span>
                  {{ t('common.deleting') }}
                </span>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Configuration Form Modal -->
    <div v-if="showForm" class="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div class="absolute inset-0 bg-black/50" @click="closeForm"></div>
      <div class="relative w-full max-w-2xl rounded-2xl border border-white/10 bg-slate-900 p-6 shadow-2xl">
        <div class="flex items-center justify-between mb-6">
          <h2 class="text-xl font-semibold">
            {{ editingConfig ? t('appleConnectConfig.editConfig') : t('appleConnectConfig.addNew') }}
          </h2>
          <button 
            @click="closeForm"
            class="text-slate-400 hover:text-white"
          >
            ×
          </button>
        </div>

        <form @submit.prevent="saveConfig" class="space-y-6">
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div>
              <label class="text-sm text-slate-400 flex items-center gap-1">
                <span>{{ t('appleConnectConfig.issuerId') }}</span>
                <span class="text-xs text-rose-400">*</span>
              </label>
              <input 
                v-model="form.issuerId" 
                type="text" 
                class="mt-1 w-full rounded-lg bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint" 
                :placeholder="t('appleConnectConfig.issuerIdPlaceholder')"
                required
              />
            </div>
            
            <div>
              <label class="text-sm text-slate-400 flex items-center gap-1">
                <span>{{ t('appleConnectConfig.keyId') }}</span>
                <span class="text-xs text-rose-400">*</span>
              </label>
              <input 
                v-model="form.keyId" 
                type="text" 
                class="mt-1 w-full rounded-lg bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint" 
                :placeholder="t('appleConnectConfig.keyIdPlaceholder')"
                required
              />
            </div>
            
            <div class="md:col-span-2">
              <label class="text-sm text-slate-400 flex items-center gap-1">
                <span>{{ t('appleConnectConfig.privateKey') }}</span>
                <span class="text-xs text-rose-400">*</span>
              </label>
              <textarea 
                v-model="form.privateKey" 
                rows="8"
                class="mt-1 w-full rounded-lg bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint font-mono text-xs" 
                :placeholder="t('appleConnectConfig.privateKeyPlaceholder')"
                required
              ></textarea>
            </div>

            <div class="md:col-span-2">
              <label class="flex items-center gap-2 cursor-pointer">
                <input 
                  v-model="form.isDefault" 
                  type="checkbox" 
                  class="rounded bg-white/5 ring-1 ring-white/10 text-mint focus:ring-2 focus:ring-mint"
                />
                <span class="text-sm text-slate-400">{{ t('appleConnectConfig.setAsDefault') }}</span>
              </label>
            </div>
          </div>
          
          <!-- Form-level test status -->
          <div v-if="formTestStatus" 
               :class="[
                 'p-3 rounded-lg border',
                 formTestStatus.success ? 'bg-emerald-900/20 border-emerald-500/30' : 'bg-rose-900/20 border-rose-500/30'
               ]"
          >
            <div class="flex items-center gap-2 text-sm">
              <span :class="formTestStatus.success ? 'text-emerald-400' : 'text-rose-400'">
                {{ formTestStatus.success ? '✓' : '✕' }}
              </span>
              <span :class="formTestStatus.success ? 'text-emerald-100' : 'text-rose-100'">
                {{ formTestStatus.message }}
              </span>
            </div>
          </div>

          <div class="flex justify-end gap-2">
            <button 
              type="button"
              class="rounded-lg border border-white/20 px-4 py-2 text-sm hover:border-white/40"
              @click="closeForm"
            >
              {{ t('common.cancel') }}
            </button>
            <button 
              type="button"
              class="rounded-lg border border-white/20 px-4 py-2 text-sm hover:border-mint/60 hover:text-mint"
              @click="testFormConnection"
              :disabled="testingForm"
            >
              <span v-if="!testingForm">{{ t('appleConnectConfig.testConnection') }}</span>
              <span v-else class="flex items-center gap-1">
                <span class="h-3 w-3 rounded-full bg-mint animate-pulse"></span>
                {{ t('appleConnectConfig.testing') }}
              </span>
            </button>
            <button 
              type="submit" 
              class="rounded-lg bg-mint px-4 py-2 text-sm font-semibold text-midnight shadow hover:bg-mint/90"
              :disabled="saving"
            >
              <span v-if="!saving">{{ editingConfig ? t('common.update') : t('common.save') }}</span>
              <span v-else class="flex items-center gap-1">
                <span class="h-3 w-3 rounded-full bg-midnight animate-pulse"></span>
                {{ t('appleConnectConfig.saving') }}
              </span>
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- Help Section -->
    <div class="rounded-2xl border border-white/10 bg-white/5 p-6">
      <h2 class="text-lg font-semibold mb-4">{{ t('appleConnectConfig.help.title') }}</h2>
      <div class="space-y-3 text-sm text-slate-300">
        <p>{{ t('appleConnectConfig.help.description') }}</p>
        <ol class="list-decimal list-inside space-y-2 mt-3">
          <li>{{ t('appleConnectConfig.help.step1') }}</li>
          <li>{{ t('appleConnectConfig.help.step2') }}</li>
          <li>{{ t('appleConnectConfig.help.step3') }}</li>
          <li>{{ t('appleConnectConfig.help.step4') }}</li>
        </ol>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useApi } from '../composables/useApi'

const { t } = useI18n()
const { api } = useApi()

const form = reactive({
  issuerId: '',
  keyId: '',
  privateKey: '',
  isDefault: true
})

const configs = ref<any[]>([])
const editingConfig = ref<any | null>(null)
const showForm = ref(false)
const loading = ref(false)
const saving = ref(false)
const testingConfigId = ref<number | null>(null)
const deletingConfigId = ref<number | null>(null)
const connectionStatus = ref<{success: boolean; title: string; message: string} | null>(null)
const testingForm = ref(false)
const formTestStatus = ref<{success: boolean; message: string} | null>(null)

function resetForm() {
  form.issuerId = ''
  form.keyId = ''
  form.privateKey = ''
  form.isDefault = true
}

function closeForm() {
  showForm.value = false
  editingConfig.value = null
  resetForm()
}

function editConfig(config: any) {
  editingConfig.value = config
  form.issuerId = config.issuerId || ''
  form.keyId = config.keyId || ''
  form.privateKey = '' // Don't load private key for security
  form.isDefault = config.isDefault ?? false
  showForm.value = true
}

function formatDate(dateString: string) {
  return new Date(dateString).toLocaleDateString()
}

async function testFormConnection() {
  // Validate form data
  if (!form.issuerId.trim() || !form.keyId.trim() || !form.privateKey.trim()) {
    formTestStatus.value = {
      success: false,
      message: t('appleConnectConfig.fillAllFields')
    }
    return
  }

  testingForm.value = true
  formTestStatus.value = null

  try {
    // Test connection directly without creating a temporary config
    const response = await api.testAppleConnectCredentials({
      issuerId: form.issuerId.trim(),
      keyId: form.keyId.trim(),
      privateKey: form.privateKey.trim()
    })

    console.log('[testFormConnection] Test response:', response)

    if (response.success) {
      formTestStatus.value = {
        success: true,
        message: t('appleConnectConfig.testSuccess.message')
      }
    } else {
      formTestStatus.value = {
        success: false,
        message: response.message || t('appleConnectConfig.testError.message')
      }
    }
  } catch (error) {
    console.error('[testFormConnection] Test failed:', error)
    formTestStatus.value = {
      success: false,
      message: error instanceof Error ? error.message : t('appleConnectConfig.testError.message')
    }
  } finally {
    testingForm.value = false
  }
}

async function loadConfigs() {
  loading.value = true
  try {
    const response = await api.getAppleConnectConfigs()
    if (response.success && Array.isArray(response.data)) {
      configs.value = response.data
    }
  } catch (error) {
    console.error('[loadConfigs] Failed to load configurations:', error)
  } finally {
    loading.value = false
  }
}

async function saveConfig() {
  saving.value = true
  connectionStatus.value = null

  try {
    const configData = {
      issuerId: form.issuerId.trim(),
      keyId: form.keyId.trim(),
      privateKey: form.privateKey.trim(),
      isDefault: form.isDefault
    }

    console.log('[saveConfig] Saving config:', { issuerId: configData.issuerId, keyId: configData.keyId, isDefault: configData.isDefault })

    if (editingConfig.value) {
      // Update existing configuration
      await api.updateAppleConnectConfig(editingConfig.value.id, configData)
    } else {
      // Create new configuration
      await api.createAppleConnectConfig(configData)
    }

    connectionStatus.value = {
      success: true,
      title: editingConfig.value ? t('appleConnectConfig.updated.title') : t('appleConnectConfig.saved.title'),
      message: editingConfig.value ? t('appleConnectConfig.updated.message') : t('appleConnectConfig.saved.message')
    }

    closeForm()
    await loadConfigs()
  } catch (error) {
    console.error('[saveConfig] Failed to save configuration:', error)
    const errorMessage = error instanceof Error ? error.message : 'Unknown error occurred'
    connectionStatus.value = {
      success: false,
      title: t('appleConnectConfig.saveError.title'),
      message: errorMessage
    }
  } finally {
    saving.value = false
  }
}

async function testConnection(configId: number) {
  testingConfigId.value = configId
  connectionStatus.value = null

  try {
    console.log('[testConnection] Testing config:', configId)

    const response = await api.testAppleConnectConnection(configId)

    if (response.success) {
      connectionStatus.value = {
        success: true,
        title: t('appleConnectConfig.testSuccess.title'),
        message: response.message || t('appleConnectConfig.testSuccess.message')
      }
    } else {
      connectionStatus.value = {
        success: false,
        title: t('appleConnectConfig.testError.title'),
        message: response.message || t('appleConnectConfig.testError.message')
      }
    }
  } catch (error) {
    console.error('[testConnection] Connection test failed:', error)
    connectionStatus.value = {
      success: false,
      title: t('appleConnectConfig.testError.title'),
      message: error instanceof Error ? error.message : t('appleConnectConfig.testError.message')
    }
  } finally {
    testingConfigId.value = null
  }
}

async function deleteConfig(configId: number) {
  if (!confirm(t('appleConnectConfig.confirmDelete'))) {
    return
  }

  deletingConfigId.value = configId
  connectionStatus.value = null

  try {
    await api.deleteAppleConnectConfig(configId)

    connectionStatus.value = {
      success: true,
      title: t('appleConnectConfig.deleted.title'),
      message: t('appleConnectConfig.deleted.message')
    }

    await loadConfigs()
  } catch (error) {
    console.error('[deleteConfig] Failed to delete configuration:', error)
    connectionStatus.value = {
      success: false,
      title: t('appleConnectConfig.deleteError.title'),
      message: error instanceof Error ? error.message : t('appleConnectConfig.deleteError.message')
    }
  } finally {
    deletingConfigId.value = null
  }
}

onMounted(async () => {
  await loadConfigs()
})
</script>