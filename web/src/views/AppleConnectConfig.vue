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
          class="rounded-lg border border-white/20 px-3 py-2 text-sm hover:border-mint/60 hover:text-mint" 
          @click="testConnection"
          :disabled="testing"
        >
          <span v-if="!testing">{{ t('appleConnectConfig.testConnection') }}</span>
          <span v-else class="flex items-center gap-1">
            <span class="h-3 w-3 rounded-full bg-mint animate-pulse"></span>
            {{ t('appleConnectConfig.testing') }}
          </span>
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
    </div>

    <!-- Configuration Form -->
    <div class="rounded-2xl border border-white/10 bg-white/5 p-6">
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
              rows="6"
              class="mt-1 w-full rounded-lg bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint font-mono text-xs" 
              :placeholder="t('appleConnectConfig.privateKeyPlaceholder')"
              required
            ></textarea>
          </div>
        </div>
        
        <div class="flex justify-end">
          <button 
            type="submit" 
            class="rounded-lg bg-mint px-4 py-2 text-sm font-semibold text-midnight shadow"
            :disabled="saving"
          >
            <span v-if="!saving">{{ t('common.save') }}</span>
            <span v-else class="flex items-center gap-1">
              <span class="h-3 w-3 rounded-full bg-midnight animate-pulse"></span>
              {{ t('appleConnectConfig.saving') }}
            </span>
          </button>
        </div>
      </form>
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
  privateKey: ''
})

const saving = ref(false)
const testing = ref(false)
const connectionStatus = ref<{success: boolean; title: string; message: string} | null>(null)

// Load existing configuration if available
onMounted(async () => {
  try {
    const response = await api.getProviderConfigs()
    if (response.success) {
      const appleConnectConfig = response.configs.find((config: any) => config.providerType === 'appleconnect')
      if (appleConnectConfig) {
        form.issuerId = appleConnectConfig.configData.issuerID || ''
        form.keyId = appleConnectConfig.configData.keyID || ''
        // We don't load the actual private key for security reasons, only a placeholder
      }
    }
  } catch (error) {
    console.error('Failed to load existing configuration:', error)
  }
})

async function saveConfig() {
  saving.value = true
  connectionStatus.value = null
  
  try {
    // Prepare the configuration data in the expected format
    const configData = {
      issuerID: form.issuerId.trim(),
      keyID: form.keyId.trim(),
      privateKey: form.privateKey.trim()
    }
    
    // Check if we already have an apple connect config to update
    const response = await api.getProviderConfigs()
    if (response.success) {
      const existingConfig = response.configs.find((config: any) => config.providerType === 'appleconnect')
      
      if (existingConfig) {
        // Update existing configuration
        await api.updateProviderConfig(existingConfig.id, {
          providerType: 'appleconnect',
          configData: configData,
          isDefault: existingConfig.isDefault
        })
      } else {
        // Create new configuration
        await api.createProviderConfig({
          providerType: 'appleconnect',
          configData: configData,
          isDefault: true
        })
      }
      
      connectionStatus.value = {
        success: true,
        title: t('appleConnectConfig.saved.title'),
        message: t('appleConnectConfig.saved.message')
      }
    }
  } catch (error) {
    console.error('Failed to save configuration:', error)
    connectionStatus.value = {
      success: false,
      title: t('appleConnectConfig.saveError.title'),
      message: t('appleConnectConfig.saveError.message')
    }
  } finally {
    saving.value = false
  }
}

async function testConnection() {
  testing.value = true
  connectionStatus.value = null
  
  try {
    // Create a temporary config to test the connection
    const testCredentials = {
      issuerId: form.issuerId,
      keyId: form.keyId,
      privateKey: form.privateKey
    }
    
    // Test by attempting to sync apps (this will validate the credentials)
    const response = await api.syncAppleApps(testCredentials)
    
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
    console.error('Connection test failed:', error)
    connectionStatus.value = {
      success: false,
      title: t('appleConnectConfig.testError.title'),
      message: t('appleConnectConfig.testError.message')
    }
  } finally {
    testing.value = false
  }
}
</script>