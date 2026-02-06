<template>
  <div class="space-y-6">
    <header class="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
      <div>
        <p class="text-xs uppercase tracking-[0.3em] text-slate-500">{{ t('nav.projects') }}</p>
        <h1 class="text-2xl font-semibold">{{ t('projects.title') }}</h1>
        <p class="text-sm text-slate-400">{{ t('projects.subtitle') }}</p>
      </div>
      <div class="flex items-center gap-2">
        <button
          class="rounded-lg border border-white/20 px-3 py-2 text-sm hover:border-mint/60 hover:text-mint"
          @click="openCreateModal"
        >
          {{ t('projects.add') }}
        </button>
      </div>
    </header>

    <section class="grid grid-cols-1 lg:grid-cols-3 gap-4">
      <div
        v-for="project in projects"
        :key="project.id"
        class="rounded-2xl border border-white/10 bg-white/5 p-4"
      >
        <div class="flex items-center justify-between">
          <div>
            <h3 class="font-semibold">{{ project.name }}</h3>
            <p class="text-xs text-slate-500">{{ project.description || t('projects.noDescription') }}</p>
          </div>
          <span class="rounded-full bg-mint/20 px-2 py-1 text-xs text-mint">
            {{ t('projects.appGroup') }}
          </span>
        </div>
        <div class="mt-3 text-xs text-slate-400 space-y-1">
          <p>{{ t('projects.appsCount', { count: countApps(project.id) }) }}</p>
        </div>
        <div class="mt-4 flex gap-2 text-xs">
          <button
            class="rounded border border-white/20 px-2 py-1 hover:border-mint/60 hover:text-mint"
            @click="openEditModal(project)"
          >
            {{ t('common.edit') }}
          </button>
          <button
            class="rounded border border-white/20 px-2 py-1 hover:border-rose-600/60 hover:text-rose-500"
            @click="deleteProject(project.id)"
          >
            {{ t('common.delete') }}
          </button>
        </div>
      </div>
    </section>

    <div v-if="showModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 px-4">
      <div class="w-full max-w-lg rounded-2xl border border-white/10 bg-midnight p-6 shadow-xl">
        <div class="flex items-center justify-between">
          <h2 class="text-lg font-semibold">
            {{ editingProject ? t('projects.edit') : t('projects.add') }}
          </h2>
          <button class="text-slate-400 hover:text-white" @click="closeModal">×</button>
        </div>

        <form class="mt-4 space-y-4" @submit.prevent="saveProject">
          <div>
            <label class="text-sm text-slate-400">{{ t('projects.name') }}</label>
            <input
              v-model="form.name"
              class="mt-1 w-full rounded-lg bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint"
              required
            />
          </div>
          <div>
            <label class="text-sm text-slate-400">{{ t('projects.description') }}</label>
            <textarea
              v-model="form.description"
              rows="3"
              class="mt-1 w-full rounded-lg bg-white/5 px-3 py-2 text-sm ring-1 ring-white/10 focus:ring-2 focus:ring-mint"
            ></textarea>
          </div>

          <div class="flex justify-end gap-2 pt-2">
            <button
              type="button"
              class="rounded-lg border border-white/20 px-3 py-2 text-sm hover:border-mint/60 hover:text-mint"
              @click="closeModal"
            >
              {{ t('common.cancel') }}
            </button>
            <button
              type="submit"
              class="rounded-lg bg-mint px-3 py-2 text-sm font-semibold text-midnight shadow"
            >
              {{ editingProject ? t('common.update') : t('common.save') }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useApi } from '../composables/useApi'
import type { Project, App } from '../composables/useApi'

const { t } = useI18n()
const { api } = useApi()

const projects = ref<Project[]>([])
const apps = ref<App[]>([])
const showModal = ref(false)
const editingProject = ref<Project | null>(null)

const form = reactive({
  name: '',
  description: ''
})

function resetForm() {
  form.name = ''
  form.description = ''
}

function openCreateModal() {
  editingProject.value = null
  resetForm()
  showModal.value = true
}

function openEditModal(project: Project) {
  editingProject.value = project
  form.name = project.name
  form.description = project.description || ''
  showModal.value = true
}

function closeModal() {
  showModal.value = false
  editingProject.value = null
  resetForm()
}

function countApps(projectId: number) {
  return apps.value.filter((app) => app.projectId === projectId).length
}

async function fetchProjects() {
  const response = await api.getProjects('app_group')
  if (response.success) {
    projects.value = response.projects
  }
}

async function fetchApps() {
  const response = await api.getApps()
  if (response.success) {
    apps.value = response.apps
  }
}

async function saveProject() {
  if (editingProject.value) {
    await api.updateProject(editingProject.value.id, {
      name: form.name,
      description: form.description,
      projectType: 'app_group'
    })
  } else {
    await api.createProject({
      name: form.name,
      description: form.description,
      projectType: 'app_group'
    })
  }

  await fetchProjects()
  closeModal()
}

async function deleteProject(projectId: number) {
  if (!confirm(t('projects.confirmDelete'))) {
    return
  }
  await api.deleteProject(projectId)
  await fetchProjects()
}

onMounted(async () => {
  await Promise.all([fetchProjects(), fetchApps()])
})
</script>
