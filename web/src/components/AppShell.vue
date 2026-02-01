<template>
  <div class="min-h-screen bg-midnight text-slate-100">
    <div class="flex min-h-screen">
      <aside class="w-64 border-r border-white/10 bg-midnight/80 backdrop-blur-lg hidden lg:flex flex-col">
        <div class="p-6 border-b border-white/5">
          <div class="text-sm uppercase tracking-[0.25em] text-slate-500">XCStrings</div>
          <div class="mt-1 text-xl font-semibold">{{ t('common.brand') }}</div>
        </div>
        <nav class="flex-1 p-4 space-y-2">
          <RouterLink
            v-for="item in navItems"
            :key="item.to"
            :to="item.to"
            class="flex items-center gap-3 rounded-xl px-3 py-2 text-sm transition"
            :class="isActive(item.to) ? 'bg-mint/10 text-mint border border-mint/40' : 'text-slate-300 hover:text-white hover:bg-white/5'"
          >
            <span class="inline-flex h-6 w-6 items-center justify-center rounded-lg bg-white/5 text-xs">{{ item.icon }}</span>
            <span>{{ typeof item.label === 'function' ? item.label() : item.label }}</span>
          </RouterLink>
        </nav>
        <div class="p-4 border-t border-white/5">
          <div class="flex items-center gap-2 mb-2">
            <div class="h-8 w-8 rounded-full bg-mint/20 flex items-center justify-center">
              <span class="text-sm font-semibold text-mint">{{ currentUserInitials }}</span>
            </div>
            <div class="flex-1 min-w-0">
              <p class="text-sm font-medium truncate">{{ currentUser?.username || 'User' }}</p>
              <p class="text-xs text-slate-500 truncate">{{ currentUser?.email || '' }}</p>
            </div>
          </div>
          <div class="text-xs text-slate-500">
            {{ t('users.' + (currentUser?.role || 'userRegular')) }}
          </div>
        </div>
      </aside>

      <div class="flex-1 flex flex-col">
        <header class="sticky top-0 z-10 bg-midnight/90 backdrop-blur-xl border-b border-white/10">
          <div class="mx-auto max-w-7xl px-4 py-4 flex items-center justify-between">
            <div>
              <div class="text-xs uppercase tracking-[0.2em] text-slate-500">Admin Console</div>
              <div class="text-lg font-semibold">Localization Ops</div>
            </div>
            <div class="flex items-center gap-3">
              <div class="hidden sm:flex items-center gap-2 rounded-full bg-white/5 px-3 py-1 text-xs">
                <span class="h-2 w-2 rounded-full bg-emerald-400"></span>
                <span>Data isolated per tenant</span>
              </div>
              <select v-model="locale" class="rounded-lg bg-white/5 px-2 py-1 text-xs border border-white/10">
                <option value="en">English</option>
                <option value="zh">中文</option>
              </select>
              <router-link to="/profile" class="rounded-full border border-white/20 px-3 py-1 text-xs hover:border-mint/60 hover:text-mint">
                {{ t('nav.profile') }}
              </router-link>
              <button class="rounded-full bg-mint px-3 py-1 text-xs font-semibold text-midnight" @click="handleLogout">
                {{ t('common.logout') }}
              </button>
            </div>
          </div>
          <!-- Profile Dropdown -->
          <div v-if="showProfileMenu" class="absolute right-4 top-16 w-48 bg-midnight/95 backdrop-blur-lg rounded-lg border border-white/10 shadow-xl p-2 space-y-1 z-50">
            <div class="px-3 py-2 border-b border-white/5">
              <p class="text-sm font-medium">{{ currentUser?.username }}</p>
              <p class="text-xs text-slate-500">{{ currentUser?.email }}</p>
            </div>
            <router-link to="/profile" class="w-full text-left px-3 py-2 text-sm rounded hover:bg-white/5 text-slate-300 hover:text-white">
              {{ t('common.settings') }}
            </router-link>
            <button class="w-full text-left px-3 py-2 text-sm rounded hover:bg-white/5 text-rose-400 hover:text-rose-300" @click="handleLogout">
              {{ t('common.logout') }}
            </button>
          </div>
        </header>

        <main class="flex-1 mx-auto w-full max-w-7xl px-4 py-8">
          <router-view />
        </main>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { RouterLink, useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useApi } from '../composables/useApi'
import type { User } from '../composables/useApi'

const route = useRoute()
const router = useRouter()
const { t, locale } = useI18n()
const { api } = useApi()

const currentUser = ref<User | null>(null)
const showProfileMenu = ref(false)

const navItems = [
  { to: '/dashboard', label: () => t('nav.dashboard'), icon: '🏠' },
  { to: '/apps', label: () => t('nav.apps'), icon: '📱' },
  { to: '/languages', label: () => t('nav.languages'), icon: '🌐' },
  { to: '/users', label: () => t('nav.users'), icon: '👥' },
  { to: '/subscriptions', label: () => t('nav.subscriptions'), icon: '💳' },
  { to: '/apple-connect-config', label: () => t('nav.appleConnectConfig'), icon: '🍎' }
]

const isActive = (path: string) => route.path.startsWith(path)

const currentUserInitials = computed(() => {
  if (!currentUser.value?.username) return 'U'
  return currentUser.value.username.substring(0, 2).toUpperCase()
})

// Fetch current user from token
async function fetchCurrentUser() {
  try {
    // First try to get user from localStorage (stored during login)
    const storedUser = localStorage.getItem('currentUser')
    if (storedUser) {
      currentUser.value = JSON.parse(storedUser)
      return
    }

    // Fallback: Get user ID from token claims
    const token = localStorage.getItem('token')
    if (!token) {
      handleLogout()
      return
    }

    // Decode JWT to get user ID
    const payload = JSON.parse(atob(token.split('.')[1]))
    const userId = payload.user_id

    if (userId) {
      const response = await api.getUser(userId)
      if (response.success) {
        currentUser.value = response.user
        // Store for future use
        localStorage.setItem('currentUser', JSON.stringify(response.user))
      }
    }
  } catch (error) {
    console.error('Failed to fetch current user:', error)
  }
}

async function handleLogout() {
  try {
    await api.logout()
  } catch (error) {
    console.error('Logout error:', error)
  } finally {
    // Clear local storage
    localStorage.removeItem('token')
    localStorage.removeItem('currentUser')
    api.clearToken()
    currentUser.value = null

    // Redirect to login
    router.push('/login')
  }
}

// Close profile menu when clicking outside
function handleClickOutside(event: MouseEvent) {
  const target = event.target as HTMLElement
  if (!target.closest('.profile-menu-container')) {
    showProfileMenu.value = false
  }
}

onMounted(() => {
  fetchCurrentUser()
  document.addEventListener('click', handleClickOutside)

  // Listen for token expiration event
  window.addEventListener('token-expired', handleTokenExpired)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
  window.removeEventListener('token-expired', handleTokenExpired)
})

function handleTokenExpired() {
  // Clear current user
  currentUser.value = null
  // Redirect to login
  router.push('/login')
}
</script>
