<template>
  <div class="min-h-screen text-slate-50 dark:text-slate-100">
    <div class="mx-auto max-w-6xl px-6 py-10 space-y-8">
      <header class="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
        <div>
          <p class="text-sm uppercase tracking-[0.3em] text-slate-400 dark:text-slate-500">PROJECTS</p>
          <h1 class="font-display text-3xl font-semibold tracking-tight text-white dark:text-slate-200 sm:text-4xl">
            Your XCStrings Projects
          </h1>
          <p class="mt-2 max-w-2xl text-slate-400 dark:text-slate-500">
            Manage all your localization projects in one place.
          </p>
        </div>
        <div class="flex flex-wrap items-center gap-3">
          <button
            class="rounded-full bg-mint px-4 py-2 text-sm font-semibold text-midnight shadow-glow transition hover:shadow-neon"
            @click="showProjectCreation = true"
          >
            New Project
          </button>
          <button
            class="rounded-full border border-white/20 px-4 py-2 text-sm font-medium text-slate-100 transition hover:border-mint/60 hover:text-mint"
            @click="loadProjects"
          >
            Refresh
          </button>
          <button
            class="rounded-full border border-white/20 px-4 py-2 text-sm font-medium text-slate-100 transition hover:border-mint/60 hover:text-mint"
            @click="logout"
          >
            Logout
          </button>
        </div>
      </header>

      <!-- Project Creation Modal -->
      <div v-if="showProjectCreation" class="fixed inset-0 bg-black/70 flex items-center justify-center z-50">
        <div class="bg-midnight rounded-2xl p-6 w-full max-w-md border border-white/20">
          <div class="flex justify-between items-center mb-4">
            <h2 class="text-xl font-semibold">Create New Project</h2>
            <button @click="showProjectCreation = false" class="text-slate-400 hover:text-white">×</button>
          </div>
          
          <form @submit.prevent="createProject">
            <div class="mb-4">
              <label class="block text-sm text-slate-400 dark:text-slate-500 mb-2">Project Name</label>
              <input v-model="newProject.name" class="w-full rounded-lg bg-white/5 px-3 py-2 ring-1 ring-white/10 focus:ring-2 focus:ring-mint text-slate-200" placeholder="My Localization Project" required />
            </div>
            <div class="mb-4">
              <label class="block text-sm text-slate-400 dark:text-slate-500 mb-2">Description</label>
              <textarea v-model="newProject.description" class="w-full rounded-lg bg-white/5 px-3 py-2 ring-1 ring-white/10 focus:ring-2 focus:ring-mint text-slate-200" placeholder="Project description"></textarea>
            </div>
            <div class="mb-4">
              <label class="block text-sm text-slate-400 dark:text-slate-500 mb-2">Source Language</label>
              <input v-model="newProject.sourceLanguage" class="w-full rounded-lg bg-white/5 px-3 py-2 ring-1 ring-white/10 focus:ring-2 focus:ring-mint text-slate-200" placeholder="e.g. en" />
            </div>
            <button type="submit" class="w-full rounded-xl bg-mint px-4 py-3 text-center text-sm font-semibold text-midnight shadow-lg shadow-mint/20 transition hover:shadow-mint/40">
              Create Project
            </button>
          </form>
        </div>
      </div>

      <section class="glass rounded-2xl p-6">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm uppercase tracking-[0.2em] text-slate-400 dark:text-slate-500">Projects</p>
            <h2 class="font-display text-xl font-semibold">Manage Projects</h2>
          </div>
          <div class="flex items-center gap-2 text-xs text-slate-300">
            <label class="flex items-center gap-2">
              <span>Filter</span>
              <input v-model="filter" class="rounded-lg bg-midnight/50 px-3 py-2 ring-1 ring-white/10 focus:ring-2 focus:ring-mint text-slate-200" placeholder="Search projects" />
            </label>
          </div>
        </div>

        <div class="mt-6">
          <div v-if="projects.length === 0" class="text-center py-12">
            <svg class="mx-auto h-12 w-12 text-slate-400 dark:text-slate-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
            <h3 class="mt-2 text-sm font-medium text-slate-200 dark:text-slate-300">No projects</h3>
            <p class="mt-1 text-sm text-slate-400 dark:text-slate-500">Get started by creating a new project.</p>
            <div class="mt-6">
              <button
                type="button"
                class="inline-flex items-center rounded-xl border border-transparent bg-mint px-4 py-2 text-sm font-medium text-midnight shadow-sm hover:bg-mint/90 focus:outline-none"
                @click="showProjectCreation = true"
              >
                <svg class="-ml-1 mr-2 h-5 w-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                  <path fill-rule="evenodd" d="M10 3a1 1 0 011 1v5h5a1 1 0 110 2h-5v5a1 1 0 11-2 0v-5H4a1 1 0 110-2h5V4a1 1 0 011-1z" clip-rule="evenodd" />
                </svg>
                New Project
              </button>
            </div>
          </div>

          <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            <div 
              v-for="project in filteredProjects" 
              :key="project.id"
              class="rounded-xl border border-white/10 bg-white/5 p-5 cursor-pointer hover:border-mint/60 hover:bg-mint/5 transition"
              @click="goToProject(project.id)"
            >
              <div class="flex justify-between">
                <h3 class="font-semibold text-slate-100 dark:text-slate-200">{{ project.name || 'Untitled Project' }}</h3>
                <span class="text-xs text-slate-400 dark:text-slate-500">{{ formatDate(project.createdAt) }}</span>
              </div>
              <p class="text-sm text-slate-400 dark:text-slate-500 mt-2">{{ project.description || 'No description' }}</p>
              <div class="mt-4 flex items-center text-xs text-slate-500 dark:text-slate-400">
                <span class="flex items-center">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-3 w-3 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                  </svg>
                  {{ project.fileName }}
                </span>
                <span class="mx-2">•</span>
                <span>{{ project.sourceLanguage }} source</span>
              </div>
              <div class="mt-3 flex space-x-2">
                <button 
                  class="text-xs rounded border border-white/20 px-2 py-1 hover:border-mint/60 hover:text-mint"
                  @click.stop="goToProject(project.id)"
                >
                  Open
                </button>
                <button 
                  class="text-xs rounded border border-white/20 px-2 py-1 hover:border-mint/60 hover:text-mint"
                  @click.stop="exportProject(project.id)"
                >
                  Export
                </button>
              </div>
            </div>
          </div>
        </div>
      </section>

      <div v-if="statusMessage" class="rounded-xl border px-4 py-3 text-sm" :class="statusClass">
        {{ statusMessage }}
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

// Project management state
const projects = ref<any[]>([])
const showProjectCreation = ref(false)
const newProject = ref({
  name: '',
  description: '',
  sourceLanguage: ''
})
const filter = ref('')
const statusMessage = ref('')
const statusTone = ref<'info' | 'error'>('info')

const statusClass = computed(() =>
  statusTone.value === 'error'
    ? 'text-red-200 border-red-400/30 bg-red-900/30'
    : 'text-mint border-mint/50 bg-mint/10'
)

const filteredProjects = computed(() => {
  const term = filter.value.trim().toLowerCase()
  if (!term) return projects.value
  return projects.value.filter((project) =>
    project.name.toLowerCase().includes(term) ||
    project.description.toLowerCase().includes(term) ||
    project.fileName.toLowerCase().includes(term)
  )
})

// Methods
async function loadProjects() {
  const token = localStorage.getItem('token')
  if (!token) {
    router.push('/login')
    return
  }

  try {
    const res = await fetch('/api/protected/projects', {
      headers: { 
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    })

    if (!res.ok) {
      throw new Error(`Failed to load projects: ${await res.text()}`)
    }

    const data = await res.json()
    projects.value = data.projects
  } catch (err) {
    showStatus('Failed to load projects: ' + (err as Error).message, 'error')
  }
}

async function createProject() {
  const token = localStorage.getItem('token')
  if (!token) {
    router.push('/login')
    return
  }

  try {
    const res = await fetch('/api/protected/projects', {
      method: 'POST',
      headers: { 
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        name: newProject.value.name,
        description: newProject.value.description,
        fileName: '',
        fileContent: '',
        sourceLanguage: newProject.value.sourceLanguage
      })
    })

    if (!res.ok) {
      throw new Error(`Failed to create project: ${await res.text()}`)
    }

    const data = await res.json()
    projects.value.unshift(data.project)
    showProjectCreation.value = false
    
    // Reset form
    newProject.value.name = ''
    newProject.value.description = ''
    newProject.value.sourceLanguage = ''
    
    showStatus('Project created successfully.', 'info')
  } catch (err) {
    showStatus('Failed to create project: ' + (err as Error).message, 'error')
  }
}

function goToProject(projectId: number) {
  router.push(`/projects/${projectId}`)
}

function formatDate(dateString: string) {
  const date = new Date(dateString)
  return date.toLocaleDateString()
}

function showStatus(message: string, tone: 'info' | 'error' = 'info') {
  statusMessage.value = message
  statusTone.value = tone
}

function logout() {
  localStorage.removeItem('token')
  router.push('/login')
}

async function exportProject(projectId: number) {
  const token = localStorage.getItem('token')
  if (!token) {
    router.push('/login')
    return
  }

  try {
    const res = await fetch(`/api/protected/projects/${projectId}/export`, {
      headers: { 
        'Authorization': `Bearer ${token}`,
      }
    })

    if (!res.ok) {
      throw new Error(`Export failed: ${await res.text()}`)
    }

    const blob = await res.blob()
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `project_${projectId}_translated.xcstrings`
    link.click()
    URL.revokeObjectURL(url)
    showStatus('Project exported successfully.', 'info')
  } catch (err) {
    showStatus('Export failed: ' + (err as Error).message, 'error')
  }
}

onMounted(() => {
  // Check if user is authenticated
  const token = localStorage.getItem('token')
  if (!token) {
    router.push('/login')
    return
  }
  
  loadProjects()
})
</script>
