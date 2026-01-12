<template>
  <div class="min-h-screen text-slate-50 dark:text-slate-100">
    <div class="mx-auto max-w-6xl px-6 py-10 space-y-8">
      <header class="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
        <div>
          <p class="text-sm uppercase tracking-[0.3em] text-slate-400 dark:text-slate-500">Dashboard</p>
          <h1 class="font-display text-3xl font-semibold tracking-tight text-white dark:text-slate-200 sm:text-4xl">
            Welcome back, {{ user?.username }}!
          </h1>
          <p class="mt-2 max-w-2xl text-slate-400 dark:text-slate-500">
            Here's what's happening with your XCStrings translation projects.
          </p>
        </div>
        <div class="flex flex-wrap items-center gap-3">
          <router-link 
            to="/translate" 
            class="rounded-full bg-mint px-4 py-2 text-sm font-semibold text-midnight shadow-glow transition hover:shadow-neon"
          >
            Start New Translation
          </router-link>
          <router-link 
            to="/projects" 
            class="rounded-full border border-white/20 px-4 py-2 text-sm font-medium text-slate-100 transition hover:border-mint/60 hover:text-mint"
          >
            Manage Projects
          </router-link>
          <button
            class="rounded-full border border-white/20 px-4 py-2 text-sm font-medium text-slate-100 transition hover:border-mint/60 hover:text-mint"
            @click="logout"
          >
            Logout
          </button>
        </div>
      </header>

      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div class="glass rounded-2xl p-6">
          <div class="flex items-center">
            <div class="rounded-lg bg-mint/10 p-3">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-mint" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
              </svg>
            </div>
            <div class="ml-4">
              <h3 class="text-sm font-medium text-slate-400 dark:text-slate-500">Total Projects</h3>
              <p class="text-2xl font-semibold text-white dark:text-slate-200">{{ projects.length }}</p>
            </div>
          </div>
        </div>

        <div class="glass rounded-2xl p-6">
          <div class="flex items-center">
            <div class="rounded-lg bg-mint/10 p-3">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-mint" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
              </svg>
            </div>
            <div class="ml-4">
              <h3 class="text-sm font-medium text-slate-400 dark:text-slate-500">Active Translations</h3>
              <p class="text-2xl font-semibold text-white dark:text-slate-200">{{ activeTranslations }}</p>
            </div>
          </div>
        </div>

        <div class="glass rounded-2xl p-6">
          <div class="flex items-center">
            <div class="rounded-lg bg-mint/10 p-3">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-mint" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
            <div class="ml-4">
              <h3 class="text-sm font-medium text-slate-400 dark:text-slate-500">Completed Translations</h3>
              <p class="text-2xl font-semibold text-white dark:text-slate-200">{{ completedTranslations }}</p>
            </div>
          </div>
        </div>
      </div>

      <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
        <section class="glass rounded-2xl p-6">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm uppercase tracking-[0.2em] text-slate-400 dark:text-slate-500">Recent Projects</p>
              <h2 class="font-display text-xl font-semibold">Your Projects</h2>
            </div>
            <router-link to="/projects" class="text-sm text-mint hover:underline">View all</router-link>
          </div>

          <div class="mt-6 space-y-4">
            <div 
              v-for="project in recentProjects" 
              :key="project.id"
              class="rounded-xl border border-white/10 bg-white/5 p-4 cursor-pointer hover:border-mint/60 transition"
              @click="goToProject(project.id)"
            >
              <div class="flex justify-between">
                <h3 class="font-semibold">{{ project.name || 'Untitled Project' }}</h3>
                <span class="text-xs text-slate-400 dark:text-slate-500">{{ formatDate(project.createdAt) }}</span>
              </div>
              <p class="text-sm text-slate-400 dark:text-slate-500 mt-1">{{ project.description || 'No description' }}</p>
              <div class="mt-2 flex items-center text-xs text-slate-500 dark:text-slate-400">
                <span class="flex items-center">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-3 w-3 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                  </svg>
                  {{ project.fileName }}
                </span>
                <span class="mx-2">•</span>
                <span>{{ project.sourceLanguage }} source</span>
              </div>
            </div>
            <div v-if="projects.length === 0" class="text-center py-8 text-slate-400 dark:text-slate-500">
              No projects yet. 
              <router-link to="/projects" class="text-mint hover:underline">Create your first project</router-link>
            </div>
          </div>
        </section>

        <section class="glass rounded-2xl p-6">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm uppercase tracking-[0.2em] text-slate-400 dark:text-slate-500">Recent Activity</p>
              <h2 class="font-display text-xl font-semibold">Latest Actions</h2>
            </div>
          </div>

          <div class="mt-6 space-y-4">
            <div 
              v-for="activity in recentActivities" 
              :key="activity.id"
              class="flex items-start"
            >
              <div class="rounded-lg bg-mint/10 p-2">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-mint" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" :d="activity.icon" />
                </svg>
              </div>
              <div class="ml-3">
                <p class="text-sm font-medium">{{ activity.action }}</p>
                <p class="text-xs text-slate-400 dark:text-slate-500">{{ activity.timestamp }}</p>
              </div>
            </div>
            <div v-if="recentActivities.length === 0" class="text-center py-8 text-slate-400 dark:text-slate-500">
              No recent activity.
            </div>
          </div>
        </section>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

// Mock data - in a real app this would come from the API
const user = ref({
  id: 1,
  username: 'JohnDoe',
  email: 'john@example.com',
  isActive: true
})

const projects = ref([
  {
    id: 1,
    name: 'Mobile App Localization',
    description: 'Localization for our iOS mobile application',
    fileName: 'Localizable.xcstrings',
    sourceLanguage: 'en',
    createdAt: '2023-06-15T10:30:00Z',
    updatedAt: '2023-06-15T10:30:00Z'
  },
  {
    id: 2,
    name: 'Feature X Strings',
    description: 'Localization for new feature X',
    fileName: 'FeatureX.xcstrings',
    sourceLanguage: 'en',
    createdAt: '2023-06-10T14:22:00Z',
    updatedAt: '2023-06-10T14:22:00Z'
  }
])

const activeTranslations = ref(2)
const completedTranslations = ref(5)
const recentActivities = ref([
  {
    id: 1,
    action: 'Started translation for Mobile App Localization',
    timestamp: '2 hours ago',
    icon: 'M13 10V3L4 14h7v7l9-11h-7z'
  },
  {
    id: 2,
    action: 'Exported translated file for Feature X Strings',
    timestamp: 'Yesterday',
    icon: 'M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4'
  },
  {
    id: 3,
    action: 'Created new project: Feature X Strings',
    timestamp: '2 days ago',
    icon: 'M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2'
  }
])

const recentProjects = computed(() => {
  return [...projects.value]
    .sort((a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime())
    .slice(0, 3)
})

// Methods
function formatDate(dateString: string) {
  const date = new Date(dateString)
  return date.toLocaleDateString()
}

function goToProject(projectId: number) {
  router.push(`/projects/${projectId}`)
}

function logout() {
  localStorage.removeItem('token')
  router.push('/login')
}

onMounted(() => {
  // In a real app, we would fetch user data from the API
  // Check if user is authenticated
  const token = localStorage.getItem('token')
  if (!token) {
    router.push('/login')
  }
})
</script>
