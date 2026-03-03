<template>
  <div class="min-h-screen bg-midnight text-slate-100">
    <div class="flex min-h-screen">
      <!-- Mobile overlay -->
      <div
        v-if="isMobileMenuOpen"
        class="fixed inset-0 bg-black/50 z-40 lg:hidden"
        @click="isMobileMenuOpen = false"
      ></div>

      <!-- Sidebar -->
      <aside
        class="fixed lg:relative z-50 h-screen border-r border-white/10 bg-midnight/80 backdrop-blur-lg flex flex-col transition-all duration-300 ease-in-out overflow-hidden"
        :class="[
          isCollapsed ? 'lg:w-16' : 'lg:w-64',
          isMobileMenuOpen ? 'translate-x-0 w-64' : '-translate-x-full lg:translate-x-0'
        ]"
      >
        <!-- Header with toggle button -->
        <div class="p-4 border-b border-white/5 flex items-center justify-between min-h-[72px]">
          <div v-if="!isCollapsed" class="overflow-hidden">
            <div class="text-sm uppercase tracking-[0.25em] text-slate-500">XCStrings</div>
            <div class="mt-1 text-xl font-semibold whitespace-nowrap">{{ t('common.brand') }}</div>
          </div>
          <button
            @click="isCollapsed = !isCollapsed"
            class="hidden lg:flex h-8 w-8 items-center justify-center rounded-lg bg-white/5 hover:bg-white/10 transition-colors shrink-0"
            :class="isCollapsed ? 'mx-auto' : ''"
          >
            <svg
              class="w-4 h-4 text-slate-400 transition-transform duration-300"
              :class="isCollapsed ? 'rotate-180' : ''"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 19l-7-7 7-7m8 14l-7-7 7-7" />
            </svg>
          </button>
          <!-- Mobile close button -->
          <button
            @click="isMobileMenuOpen = false"
            class="lg:hidden h-8 w-8 items-center justify-center rounded-lg bg-white/5 hover:bg-white/10 transition-colors"
          >
            <svg class="w-4 h-4 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <!-- Navigation -->
        <nav class="flex-1 p-2 space-y-1 overflow-y-auto overflow-x-hidden">
          <RouterLink
            v-for="item in navItems"
            :key="item.to"
            :to="item.to"
            class="flex items-center gap-3 rounded-xl px-3 py-2 text-sm transition group relative"
            :class="[
              isActive(item.to) ? 'bg-mint/10 text-mint border border-mint/40' : 'text-slate-300 hover:text-white hover:bg-white/5',
              isCollapsed ? 'justify-center' : ''
            ]"
            @click="isMobileMenuOpen = false"
          >
            <span
              class="inline-flex shrink-0 items-center justify-center rounded-lg bg-white/5 group-hover:scale-110 transition-transform"
              :class="isCollapsed ? 'h-8 w-8 text-lg' : 'h-6 w-6 text-xs'"
            >{{ item.icon }}</span>
            <span v-if="!isCollapsed" class="whitespace-nowrap overflow-hidden">{{ typeof item.label === 'function' ? item.label() : item.label }}</span>
            <!-- Tooltip for collapsed state -->
            <div
              v-if="isCollapsed"
              class="absolute left-full ml-2 px-2 py-1 bg-midnight/95 border border-white/10 rounded-lg text-xs whitespace-nowrap opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all z-50"
            >
              {{ typeof item.label === 'function' ? item.label() : item.label }}
            </div>
          </RouterLink>
        </nav>

        <!-- User section -->
        <div class="p-2 border-t border-white/5">
          <div
            class="flex items-center gap-2 p-2 rounded-xl hover:bg-white/5 transition-colors"
            :class="isCollapsed ? 'justify-center' : ''"
          >
            <div
              class="shrink-0 rounded-full bg-mint/20 flex items-center justify-center transition-all"
              :class="isCollapsed ? 'h-10 w-10' : 'h-8 w-8'"
            >
              <span
                class="font-semibold text-mint transition-all"
                :class="isCollapsed ? 'text-base' : 'text-sm'"
              >{{ currentUserInitials }}</span>
            </div>
            <div v-if="!isCollapsed" class="flex-1 min-w-0 overflow-hidden">
              <p class="text-sm font-medium truncate">{{ currentUser?.username || 'User' }}</p>
              <p class="text-xs text-slate-500 truncate">{{ currentUser?.email || '' }}</p>
            </div>
          </div>
          <div v-if="!isCollapsed" class="text-xs text-slate-500 px-2 mt-1">
            {{ t('users.' + (currentUser?.role || 'userRegular')) }}
          </div>
        </div>
      </aside>

      <div class="flex-1 flex flex-col min-w-0">
        <header class="sticky top-0 z-10 bg-midnight/90 backdrop-blur-xl border-b border-white/10">
          <div class="mx-auto max-w-7xl px-4 py-4 flex items-center justify-between gap-4">
            <div class="flex items-center gap-3">
              <!-- Mobile menu toggle -->
              <button
                @click="isMobileMenuOpen = !isMobileMenuOpen"
                class="lg:hidden h-8 w-8 flex items-center justify-center rounded-lg bg-white/5 hover:bg-white/10 transition-colors"
              >
                <svg class="w-5 h-5 text-slate-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
                </svg>
              </button>
              <div>
                <div class="text-xs uppercase tracking-[0.2em] text-slate-500">Admin Console</div>
                <div class="text-lg font-semibold">Localization Ops</div>
              </div>
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
const isCollapsed = ref(false)
const isMobileMenuOpen = ref(false)

const navItems = [
  { to: '/dashboard', label: () => t('nav.dashboard'), icon: '🏠' },
  { to: '/apps', label: () => t('nav.apps'), icon: '📱' },
  { to: '/projects', label: () => t('nav.projects'), icon: '🗂️' },
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
