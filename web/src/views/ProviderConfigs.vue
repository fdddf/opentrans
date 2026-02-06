<template>
  <div class="space-y-6">
    <header class="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
      <div>
        <p class="text-xs uppercase tracking-[0.3em] text-slate-500">{{ t('nav.providerConfigs') }}</p>
        <h1 class="text-2xl font-semibold">{{ t('providerConfigs.title') }}</h1>
        <p class="text-sm text-slate-400">{{ t('providerConfigs.subtitle') }}</p>
      </div>
      <button class="rounded-lg bg-mint px-3 py-2 text-sm font-semibold text-midnight shadow" @click="showAddModal = true">
        {{ t('providerConfigs.add') }}
      </button>
    </header>

    <!-- Add/Edit Provider Config Modal -->
    <div v-if="showAddModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 px-4">
      <div class="w-full max-w-lg rounded-2xl border border-white/10 bg-midnight p-6 shadow-xl">
        <div class="flex items-center justify-between">
          <h2 class="text-lg font-semibold">{{ editingConfig ? t('providerConfigs.edit') : t('providerConfigs.add') }}</h2>
          <button class="text-slate-400 hover:text-white" @click="closeModal">×</button>
        </div>

        <form class="mt-4 space-y-4" @submit.prevent="saveConfig">
          <div>
            <label class="text-sm text-slate-400">{{ t('providerConfigs.providerType') }}</label>
            <select 
              v-model="form.providerType" 
              class="mt-1 w-full rounded-lg bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint" 
              :disabled="!!editingConfig"
              required
            >
              <option value="openai">OpenAI</option>
              <option value="google">Google</option>
              <option value="deepl">DeepL</option>
              <option value="baidu">Baidu</option>
              <option value="appleconnect">Apple Connect</option>
              <option value="llama">Hunyuan (Llama)</option>
            </select>
          </div>

          <!-- OpenAI Specific Fields -->
          <div v-if="form.providerType === 'openai'">
            <div>
              <label class="text-sm text-slate-400">API Key</label>
              <input 
                v-model="form.configData.apiKey" 
                type="password"
                class="mt-1 w-full rounded-lg bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint" 
                placeholder="sk-..."
                required 
              />
            </div>
            <div class="grid grid-cols-2 gap-4 mt-2">
              <div>
                <label class="text-sm text-slate-400">API Base URL (optional)</label>
                <input 
                  v-model="form.configData.apiBaseUrl" 
                  class="mt-1 w-full rounded-lg bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint" 
                  placeholder="https://api.openai.com" 
                />
              </div>
              <div>
                <label class="text-sm text-slate-400">Model</label>
                <input 
                  v-model="form.configData.model" 
                  class="mt-1 w-full rounded-lg bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint" 
                  placeholder="gpt-3.5-turbo" 
                />
              </div>
            </div>
          </div>

          <!-- Google Specific Fields -->
          <div v-if="form.providerType === 'google'">
            <div>
              <label class="text-sm text-slate-400">API Key</label>
              <input 
                v-model="form.configData.apiKey" 
                type="password"
                class="mt-1 w-full rounded-lg bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint" 
                placeholder="AIzaSy..."
                required 
              />
            </div>
          </div>

          <!-- DeepL Specific Fields -->
          <div v-if="form.providerType === 'deepl'">
            <div>
              <label class="text-sm text-slate-400">API Key</label>
              <input 
                v-model="form.configData.apiKey" 
                type="password"
                class="mt-1 w-full rounded-lg bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint" 
                placeholder="..."
                required 
              />
            </div>
            <div class="mt-2">
              <label class="flex items-center">
                <input 
                  v-model="form.configData.isFree" 
                  type="checkbox"
                  class="rounded bg-white/10 text-mint ring-1 ring-white/10 focus:ring-mint"
                />
                <span class="ml-2 text-sm text-slate-400">Use Free API</span>
              </label>
            </div>
          </div>

          <!-- Baidu Specific Fields -->
          <div v-if="form.providerType === 'baidu'">
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="text-sm text-slate-400">App ID</label>
                <input 
                  v-model="form.configData.appId" 
                  class="mt-1 w-full rounded-lg bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint" 
                  placeholder="xxx"
                  required 
                />
              </div>
              <div>
                <label class="text-sm text-slate-400">App Secret</label>
                <input 
                  v-model="form.configData.appSecret" 
                  type="password"
                  class="mt-1 w-full rounded-lg bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint" 
                  placeholder="xxx"
                  required 
                />
              </div>
            </div>
          </div>

          <!-- Apple Connect Specific Fields -->
          <div v-if="form.providerType === 'appleconnect'">
            <div>
              <label class="text-sm text-slate-400">Issuer ID</label>
              <input 
                v-model="form.configData.issuerID" 
                class="mt-1 w-full rounded-lg bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint" 
                placeholder="xxx"
                required 
              />
            </div>
            <div class="mt-2">
              <label class="text-sm text-slate-400">Key ID</label>
              <input 
                v-model="form.configData.keyID" 
                class="mt-1 w-full rounded-lg bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint" 
                placeholder="xxx"
                required 
              />
            </div>
            <div class="mt-2">
              <label class="text-sm text-slate-400">Private Key</label>
              <textarea 
                v-model="form.configData.privateKey" 
                rows="4"
                class="mt-1 w-full rounded-lg bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint" 
                placeholder="-----BEGIN PRIVATE KEY-----&#10;...&#10;-----END PRIVATE KEY-----"
                required 
              />
            </div>
          </div>

          <!-- Hunyuan/Llama Specific Fields -->
          <div v-if="form.providerType === 'llama'">
            <p class="text-xs text-slate-400">{{ t('providerConfigs.hunyuanHint') }}</p>
          </div>

          <div class="flex items-center">
            <label class="flex items-center">
              <input 
                v-model="form.isDefault" 
                type="checkbox"
                class="rounded bg-white/10 text-mint ring-1 ring-white/10 focus:ring-mint"
              />
              <span class="ml-2 text-sm text-slate-400">{{ t('providerConfigs.default') }}</span>
            </label>
          </div>

          <div class="flex justify-end gap-2 pt-4">
            <button type="button" class="rounded-lg border border-white/20 px-3 py-2 text-sm hover:border-mint/60 hover:text-mint" @click="closeModal">
              {{ t('common.cancel') || 'Cancel' }}
            </button>
            <button type="submit" class="rounded-lg bg-mint px-3 py-2 text-sm font-semibold text-midnight shadow">
              {{ editingConfig ? t('common.update') : t('common.save') }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- Provider Configs List -->
    <section class="rounded-2xl border border-white/10 bg-white/5 p-4">
      <div class="overflow-x-auto">
        <table class="min-w-full text-sm">
          <thead class="text-left text-slate-500">
            <tr>
              <th class="py-2 px-2">{{ t('providerConfigs.type') }}</th>
              <th class="py-2 px-2">{{ t('providerConfigs.name') }}</th>
              <th class="py-2 px-2">{{ t('providerConfigs.default') }}</th>
              <th class="py-2 px-2">{{ t('common.status') }}</th>
              <th class="py-2 px-2">{{ t('common.actions') }}</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-white/5">
            <tr v-for="config in configs" :key="config.id" class="hover:bg-white/5">
              <td class="py-2 px-2">
                <span class="rounded-full bg-white/10 px-2 py-1 text-xs uppercase">{{ config.providerType }}</span>
              </td>
              <td class="py-2 px-2">
                <span class="font-medium">{{ getProviderDisplayName(config.providerType) }}</span>
              </td>
              <td class="py-2 px-2">
                <span v-if="config.isDefault" class="rounded-full bg-mint/20 px-2 py-1 text-xs text-mint">
                  {{ t('providerConfigs.yes') }}
                </span>
                <span v-else class="rounded-full bg-slate-700/40 px-2 py-1 text-xs text-slate-400">
                  {{ t('providerConfigs.no') }}
                </span>
              </td>
              <td class="py-2 px-2">
                <span class="rounded-full bg-emerald-900/40 px-2 py-1 text-xs text-emerald-200">
                  {{ t('common.active') }}
                </span>
              </td>
              <td class="py-2 px-2 space-x-2 text-xs">
                <button 
                  class="rounded border border-white/20 px-2 py-1 hover:border-mint/60 hover:text-mint" 
                  @click="editConfig(config)"
                >
                  {{ t('common.edit') }}
                </button>
                <button 
                  class="rounded border border-white/20 px-2 py-1 hover:border-rose-600/60 hover:text-rose-500" 
                  @click="deleteConfig(config.id)"
                >
                  {{ t('common.delete') }}
                </button>
                <button 
                  v-if="!config.isDefault"
                  class="rounded border border-white/20 px-2 py-1 hover:border-amber-500/60 hover:text-amber-400" 
                  @click="setAsDefault(config.id)"
                >
                  {{ t('providerConfigs.makeDefault') }}
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useApi } from '../composables/useApi'
import type { ProviderConfig } from '../composables/useApi'

const { t } = useI18n()
const { api } = useApi()

const configs = ref<ProviderConfig[]>([])
const showAddModal = ref(false)
const editingConfig = ref<ProviderConfig | null>(null)

const form = reactive({
  providerType: 'openai',
  configData: {} as any,
  isDefault: false
})

function resetForm() {
  form.providerType = 'openai'
  form.configData = {}
  form.isDefault = false
}

function closeModal() {
  showAddModal.value = false
  editingConfig.value = null
  resetForm()
}

function getProviderDisplayName(providerType: string): string {
  const names: Record<string, string> = {
    'openai': 'OpenAI',
    'google': 'Google',
    'deepl': 'DeepL',
    'baidu': 'Baidu',
    'appleconnect': 'Apple Connect',
    'llama': 'Hunyuan (Llama)'
  }
  return names[providerType] || providerType
}

async function saveConfig() {
  try {
    let response
    if (editingConfig.value) {
      response = await api.updateProviderConfig(editingConfig.value.id, {
        providerType: form.providerType,
        configData: form.configData,
        isDefault: form.isDefault
      })
    } else {
      response = await api.createProviderConfig({
        providerType: form.providerType,
        configData: form.configData,
        isDefault: form.isDefault
      })
    }
    
    if (response.success) {
      fetchConfigs()
      closeModal()
    }
  } catch (error) {
    console.error('Failed to save provider config:', error)
  }
}

function editConfig(config: ProviderConfig) {
  editingConfig.value = config
  form.providerType = config.providerType
  form.configData = { ...config.configData }
  // Remove redacted values and replace with empty strings
  for (const key in form.configData) {
    if (form.configData[key] === '***REDACTED***') {
      form.configData[key] = ''
    }
  }
  form.isDefault = config.isDefault
  showAddModal.value = true
}

async function deleteConfig(configId: number) {
  if (!confirm(t('providerConfigs.confirmDelete') || 'Are you sure you want to delete this configuration?')) {
    return
  }
  
  try {
    await api.deleteProviderConfig(configId)
    fetchConfigs()
  } catch (error) {
    console.error('Failed to delete provider config:', error)
  }
}

async function setAsDefault(configId: number) {
  try {
    const config = configs.value.find(c => c.id === configId)
    if (!config) return
    
    await api.updateProviderConfig(configId, {
      isDefault: true
    })
    
    fetchConfigs()
  } catch (error) {
    console.error('Failed to set as default:', error)
  }
}

async function fetchConfigs() {
  try {
    const response = await api.getProviderConfigs()
    if (response.success) {
      configs.value = response.configs
    }
  } catch (error) {
    console.error('Failed to fetch provider configs:', error)
  }
}

onMounted(() => {
  fetchConfigs()
})
</script>
