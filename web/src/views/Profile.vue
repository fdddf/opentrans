<template>
  <div class="space-y-6">
    <header>
      <p class="text-xs uppercase tracking-[0.3em] text-slate-500">{{ t('profile.title') }}</p>
      <h1 class="text-2xl font-semibold">{{ t('profile.title') }}</h1>
      <p class="text-sm text-slate-400">{{ t('profile.subtitle') }}</p>
    </header>

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Profile Information -->
      <div class="lg:col-span-2 space-y-6">
        <section class="glass rounded-2xl p-6">
          <h2 class="text-lg font-semibold mb-4">{{ t('profile.personalInfo') }}</h2>
          <form @submit.prevent="updateProfile" class="space-y-4">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-slate-300 mb-1">{{ t('common.username') }}</label>
                <input
                  v-model="profileForm.username"
                  type="text"
                  class="w-full px-4 py-2 border border-white/10 rounded-xl bg-white/5 text-slate-200 focus:outline-none focus:ring-mint focus:border-mint"
                  disabled
                />
                <p class="text-xs text-slate-500 mt-1">{{ t('profile.usernameReadOnly') }}</p>
              </div>
              <div>
                <label class="block text-sm font-medium text-slate-300 mb-1">{{ t('common.email') }}</label>
                <input
                  v-model="profileForm.email"
                  type="email"
                  class="w-full px-4 py-2 border border-white/10 rounded-xl bg-white/5 text-slate-200 focus:outline-none focus:ring-mint focus:border-mint"
                />
              </div>
            </div>
            <div class="flex justify-end">
              <button
                type="submit"
                :disabled="updatingProfile"
                class="px-6 py-2 bg-mint text-midnight rounded-xl font-semibold hover:bg-mint/90 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {{ updatingProfile ? t('common.saving') + '...' : t('profile.updateProfile') }}
              </button>
            </div>
          </form>
        </section>

        <!-- Change Password -->
        <section class="glass rounded-2xl p-6">
          <h2 class="text-lg font-semibold mb-4">{{ t('profile.changePassword') }}</h2>
          <form @submit.prevent="changePassword" class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-slate-300 mb-1">{{ t('profile.currentPassword') }}</label>
              <input
                v-model="passwordForm.currentPassword"
                type="password"
                class="w-full px-4 py-2 border border-white/10 rounded-xl bg-white/5 text-slate-200 focus:outline-none focus:ring-mint focus:border-mint"
                required
              />
            </div>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-slate-300 mb-1">{{ t('profile.newPassword') }}</label>
                <input
                  v-model="passwordForm.newPassword"
                  type="password"
                  class="w-full px-4 py-2 border border-white/10 rounded-xl bg-white/5 text-slate-200 focus:outline-none focus:ring-mint focus:border-mint"
                  required
                />
              </div>
              <div>
                <label class="block text-sm font-medium text-slate-300 mb-1">{{ t('common.confirmPassword') }}</label>
                <input
                  v-model="passwordForm.confirmPassword"
                  type="password"
                  class="w-full px-4 py-2 border border-white/10 rounded-xl bg-white/5 text-slate-200 focus:outline-none focus:ring-mint focus:border-mint"
                  required
                />
              </div>
            </div>
            <div class="flex justify-end">
              <button
                type="submit"
                :disabled="changingPassword"
                class="px-6 py-2 bg-mint text-midnight rounded-xl font-semibold hover:bg-mint/90 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {{ changingPassword ? t('common.saving') + '...' : t('profile.updatePassword') }}
              </button>
            </div>
          </form>
        </section>
      </div>

      <!-- Account Info -->
      <div class="space-y-6">
        <section class="glass rounded-2xl p-6">
          <h2 class="text-lg font-semibold mb-4">{{ t('profile.accountInfo') }}</h2>
          <div class="space-y-4">
            <div>
              <p class="text-xs text-slate-500 mb-1">{{ t('profile.userId') }}</p>
              <p class="text-sm font-medium">{{ currentUser?.id }}</p>
            </div>
            <div>
              <p class="text-xs text-slate-500 mb-1">{{ t('profile.status') }}</p>
              <span class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium" :class="currentUser?.isActive ? 'bg-emerald-900/30 text-emerald-400' : 'bg-rose-900/30 text-rose-400'">
                {{ currentUser?.isActive ? t('common.active') : t('common.inactive') }}
              </span>
            </div>
            <div>
              <p class="text-xs text-slate-500 mb-1">{{ t('profile.subscriptionType') }}</p>
              <p class="text-sm font-medium">{{ currentUser?.subscriptionType || '-' }}</p>
            </div>
            <div>
              <p class="text-xs text-slate-500 mb-1">{{ t('profile.createdAt') }}</p>
              <p class="text-sm font-medium">{{ formatDate(currentUser?.createdAt) }}</p>
            </div>
          </div>
        </section>

        <!-- Usage Stats -->
        <section class="glass rounded-2xl p-6">
          <h2 class="text-lg font-semibold mb-4">{{ t('profile.usage') }}</h2>
          <div class="space-y-4">
            <div>
              <div class="flex justify-between text-sm mb-1">
                <span class="text-slate-400">{{ t('profile.appsUsed') }}</span>
                <span class="font-medium">{{ currentUser?.currentAppCount || 0 }} / {{ currentUser?.maxApps || 0 }}</span>
              </div>
              <div class="h-2 rounded-full bg-white/10">
                <div class="h-2 rounded-full bg-mint" :style="{ width: `${Math.min((currentUser?.currentAppCount || 0) / (currentUser?.maxApps || 1) * 100, 100)}%` }"></div>
              </div>
            </div>
            <div>
              <div class="flex justify-between text-sm mb-1">
                <span class="text-slate-400">{{ t('profile.translationsUsed') }}</span>
                <span class="font-medium">{{ currentUser?.currentUsage || 0 }} / {{ currentUser?.maxTranslations || 0 }}</span>
              </div>
              <div class="h-2 rounded-full bg-white/10">
                <div class="h-2 rounded-full bg-mint" :style="{ width: `${Math.min((currentUser?.currentUsage || 0) / (currentUser?.maxTranslations || 1) * 100, 100)}%` }"></div>
              </div>
            </div>
          </div>
        </section>
      </div>
    </div>

    <!-- Toast notification -->
    <div
      v-if="toast.show"
      class="fixed bottom-4 right-4 rounded-xl px-6 py-4 flex items-center gap-3 shadow-xl z-50 transition-all"
      :class="toast.type === 'success' ? 'bg-emerald-900/90 text-emerald-100' : 'bg-rose-900/90 text-rose-100'"
    >
      <span class="text-xl">{{ toast.type === 'success' ? '✓' : '✕' }}</span>
      <span class="text-sm font-medium">{{ toast.message }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useApi } from '@/composables/useApi'
import type { User } from '@/composables/useApi'

const { t } = useI18n()
const { api } = useApi()

const currentUser = ref<User | null>(null)

const profileForm = reactive({
  username: '',
  email: ''
})

const passwordForm = reactive({
  currentPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const updatingProfile = ref(false)
const changingPassword = ref(false)

const toast = reactive({
  show: false,
  message: '',
  type: 'success' as 'success' | 'error'
})

function showToast(message: string, type: 'success' | 'error' = 'success') {
  toast.message = message
  toast.type = type
  toast.show = true
  setTimeout(() => {
    toast.show = false
  }, 3000)
}

function formatDate(dateStr?: string): string {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleDateString()
}

async function loadCurrentUser() {
  try {
    const storedUser = localStorage.getItem('currentUser')
    if (storedUser) {
      currentUser.value = JSON.parse(storedUser)
      profileForm.username = currentUser.value.username
      profileForm.email = currentUser.value.email
    }
  } catch (error) {
    console.error('Failed to load current user:', error)
  }
}

async function updateProfile() {
  if (!currentUser.value) return

  updatingProfile.value = true
  try {
    const response = await api.request('/protected/user/profile', {
      method: 'PUT',
      body: JSON.stringify({
        email: profileForm.email
      })
    })

    if (response && response.success) {
      currentUser.value.email = profileForm.email
      localStorage.setItem('currentUser', JSON.stringify(currentUser.value))
      showToast(t('profile.updateSuccess'))
    } else {
      showToast(t('profile.updateFailed'), 'error')
    }
  } catch (error) {
    console.error('Failed to update profile:', error)
    showToast(t('profile.updateFailed'), 'error')
  } finally {
    updatingProfile.value = false
  }
}

async function changePassword() {
  if (passwordForm.newPassword !== passwordForm.confirmPassword) {
    showToast(t('profile.passwordMismatch'), 'error')
    return
  }

  if (passwordForm.newPassword.length < 6) {
    showToast(t('profile.passwordTooShort'), 'error')
    return
  }

  changingPassword.value = true
  try {
    const response = await fetch('/api/protected/change-password', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      },
      body: JSON.stringify({
        currentPassword: passwordForm.currentPassword,
        newPassword: passwordForm.newPassword
      })
    })

    const data = await response.json()

    if (data.success) {
      showToast(t('profile.passwordUpdateSuccess'))
      passwordForm.currentPassword = ''
      passwordForm.newPassword = ''
      passwordForm.confirmPassword = ''
    } else {
      showToast(data.message || t('profile.passwordUpdateFailed'), 'error')
    }
  } catch (error) {
    console.error('Failed to change password:', error)
    showToast(t('profile.passwordUpdateFailed'), 'error')
  } finally {
    changingPassword.value = false
  }
}

onMounted(() => {
  loadCurrentUser()
})
</script>