<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-slate-900 via-midnight to-slate-900 py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full space-y-8">
      <div class="text-center">
        <h2 class="mt-6 text-center text-3xl font-extrabold text-slate-200">
          {{ t('auth.registerTitle') }}
        </h2>
        <p class="mt-2 text-center text-sm text-slate-400">
          {{ t('auth.haveAccount') }}
          <router-link to="/login" class="font-medium text-mint hover:text-mint/80">
            {{ t('auth.loginTitle') }}
          </router-link>
        </p>
      </div>
      <div class="mt-8 bg-midnight/50 backdrop-blur-lg rounded-2xl px-6 py-8 shadow-xl border border-white/10">
        <form class="space-y-6" @submit.prevent="performRegister">
          <div>
            <label for="username" class="block text-sm font-medium text-slate-300">{{ t('common.username') }}</label>
            <div class="mt-1">
              <input 
                id="username" 
                v-model="registerForm.username" 
                type="text" 
                class="appearance-none block w-full px-4 py-3 border border-white/10 rounded-xl bg-white/5 text-slate-200 placeholder-slate-500 focus:outline-none focus:ring-mint focus:border-mint focus:z-10 sm:text-sm"
                placeholder="Choose a username"
                required
              />
            </div>
          </div>

          <div>
            <label for="email" class="block text-sm font-medium text-slate-300">{{ t('common.email') }}</label>
            <div class="mt-1">
              <input 
                id="email" 
                v-model="registerForm.email" 
                type="email" 
                class="appearance-none block w-full px-4 py-3 border border-white/10 rounded-xl bg-white/5 text-slate-200 placeholder-slate-500 focus:outline-none focus:ring-mint focus:border-mint focus:z-10 sm:text-sm"
                placeholder="your@email.com"
                required
              />
            </div>
          </div>

          <div>
            <label for="password" class="block text-sm font-medium text-slate-300">{{ t('common.password') }}</label>
            <div class="mt-1">
              <input 
                id="password" 
                v-model="registerForm.password" 
                type="password" 
                class="appearance-none block w-full px-4 py-3 border border-white/10 rounded-xl bg-white/5 text-slate-200 placeholder-slate-500 focus:outline-none focus:ring-mint focus:border-mint focus:z-10 sm:text-sm"
                placeholder="Create a password"
                required
              />
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
              {{ t('common.register') }}
            </button>
          </div>
        </form>
      </div>
      <div v-if="statusMessage" class="text-center text-sm" :class="statusClass">
        {{ statusMessage }}
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'

const router = useRouter()
const { t } = useI18n()

// Authentication state
const registerForm = reactive({
  username: '',
  email: '',
  password: ''
})

const statusMessage = ref('')
const statusTone = ref<'info' | 'error'>('info')

const statusClass = computed(() =>
  statusTone.value === 'error'
    ? 'text-red-400'
    : 'text-mint'
)

// Authentication functions
async function performRegister() {
  try {
    // Basic email validation
    if (!validateEmail(registerForm.email)) {
      showStatus('Please enter a valid email address', 'error')
      return
    }

    const res = await fetch('/api/auth/register', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(registerForm)
    })

    let data;
    // For consistency, always parse the response body in case of error or success
    if (!res.ok) {
      data = await res.json()
      showStatus(data.message || 'Registration failed', 'error')
      return
    }

    data = await res.json()
    
    // Check if the response has the expected structure for success
    if (!data.success) {
      showStatus('Registration failed: Invalid response format from server', 'error')
      return
    }
    
    showStatus(data.message || 'Registration successful. Please check your email for activation.', 'info')
    
    // Redirect to login after a short delay
    setTimeout(() => {
      router.push('/login')
    }, 2000)
  } catch (err) {
    showStatus('Registration failed: ' + (err as Error).message, 'error')
  }
}

function validateEmail(email: string) {
  const re = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  return re.test(email)
}

function showStatus(message: string, tone: 'info' | 'error' = 'info') {
  statusMessage.value = message
  statusTone.value = tone
}
</script>
