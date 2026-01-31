<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-slate-900 via-midnight to-slate-900 py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full space-y-8">
      <div class="text-center">
        <h2 class="mt-6 text-center text-3xl font-extrabold text-slate-200">
          {{ t('auth.loginTitle') }}
        </h2>
        <p class="mt-2 text-center text-sm text-slate-400">
          {{ t('auth.noAccount') }}
          <router-link to="/register" class="font-medium text-mint hover:text-mint/80">
            {{ t('auth.registerPrompt') }}
          </router-link>
        </p>
      </div>
      <div class="mt-8 bg-midnight/50 backdrop-blur-lg rounded-2xl px-6 py-8 shadow-xl border border-white/10">
        <form class="space-y-6" @submit.prevent="performLogin">
          <div>
            <label for="username" class="block text-sm font-medium text-slate-300">{{ t('common.username') }}</label>
            <div class="mt-1">
              <input 
                id="username" 
                v-model="loginForm.username" 
                type="text" 
                class="appearance-none block w-full px-4 py-3 border border-white/10 rounded-xl bg-white/5 text-slate-200 placeholder-slate-500 focus:outline-none focus:ring-mint focus:border-mint focus:z-10 sm:text-sm"
                placeholder="Your username"
                required
              />
            </div>
          </div>

          <div>
            <label for="password" class="block text-sm font-medium text-slate-300">{{ t('common.password') }}</label>
            <div class="mt-1">
              <input 
                id="password" 
                v-model="loginForm.password" 
                type="password" 
                class="appearance-none block w-full px-4 py-3 border border-white/10 rounded-xl bg-white/5 text-slate-200 placeholder-slate-500 focus:outline-none focus:ring-mint focus:border-mint focus:z-10 sm:text-sm"
                placeholder="Your password"
                required
              />
            </div>
          </div>

          <div class="flex items-center justify-between">
            <div class="text-sm">
              <a href="#" class="font-medium text-mint hover:text-mint/80">
                {{ t('auth.forgot') }}
              </a>
            </div>
          </div>

          <div>
            <button
              type="submit"
              class="group relative w-full flex justify-center py-3 px-4 border border-transparent text-sm font-medium rounded-xl text-midnight bg-mint hover:bg-mint/90 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-mint shadow-lg shadow-mint/20 transition-all duration-200"
            >
              <span class="absolute left-0 inset-y-0 flex items-center pl-3">
                <svg class="h-5 w-5 text-mint/60 group-hover:text-mint/40" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                  <path fill-rule="evenodd" d="M5 9V7a5 5 0 0110 0v2a2 2 0 012 2v5a2 2 0 01-2 2H5a2 2 0 01-2-2v-5a2 2 0 012-2zm8-2v2H7V7a3 3 0 016 0z" clip-rule="evenodd" />
                </svg>
              </span>
              {{ t('auth.loginCta') }}
            </button>
          </div>
        </form>
      </div>
      <!-- Status Message with better styling -->
      <div v-if="statusMessage" class="rounded-2xl px-6 py-4 flex items-center justify-center gap-3" :class="statusTone === 'error' ? 'bg-rose-900/20 border border-rose-500/30' : 'bg-emerald-900/20 border border-emerald-500/30'">
        <span class="text-xl" :class="statusTone === 'error' ? 'text-rose-400' : 'text-emerald-400'">
          {{ statusTone === 'error' ? '✕' : '✓' }}
        </span>
        <span class="text-sm font-medium" :class="statusTone === 'error' ? 'text-rose-100' : 'text-emerald-100'">
          {{ statusMessage }}
        </span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useApi } from '../composables/useApi'

const router = useRouter()
const { t } = useI18n()
const { api } = useApi()

// Authentication state
const loginForm = reactive({
  username: '',
  password: ''
})

const statusMessage = ref('')
const statusTone = ref<'info' | 'error'>('info')

// Authentication functions
async function performLogin() {
  try {
    const response = await api.login(loginForm.username, loginForm.password)

    if (!response.success || !response.token) {
      // Extract message safely - handle case where message might be undefined
      const errorMsg = response.message || t('auth.loginFailed')
      showStatus(errorMsg, 'error')
      return
    }

    // ApiClient automatically sets the token via setToken
    // and stores it in localStorage
    // Also store the user object from the login response
    if (response.user) {
      localStorage.setItem('currentUser', JSON.stringify(response.user))
    }

    showStatus(t('auth.loginSuccess'), 'info')

    // Redirect to dashboard after a short delay
    setTimeout(() => {
      router.push('/dashboard')
    }, 1000)
  } catch (err) {
    // Handle network errors or API errors
    const error = err as Error
    let errorMsg = t('auth.loginFailed')

    // If error message is valid and doesn't look like JSON, use it
    if (error.message && !error.message.startsWith('{')) {
      errorMsg = error.message
    }

    showStatus(errorMsg, 'error')
  }
}

function showStatus(message: string, tone: 'info' | 'error' = 'info') {
  // Ensure message is a string and not an object
  statusMessage.value = String(message)
  statusTone.value = tone
}
</script>
