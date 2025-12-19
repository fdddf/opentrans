<template>
  <div class="min-h-screen text-slate-50">
    <div class="mx-auto max-w-6xl px-6 py-10 space-y-8">
      <header class="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
        <div>
          <p class="text-sm uppercase tracking-[0.3em] text-slate-400">Visual Localisation</p>
          <h1 class="font-display text-3xl font-semibold tracking-tight text-white sm:text-4xl">
            XCStrings Translator Studio
          </h1>
          <p class="mt-2 max-w-2xl text-slate-400">
            Upload your <code class="font-mono text-mint">Localizable.xcstrings</code>, pick target languages and run batch
            translations with your favourite provider. Export the updated file in one click.
          </p>
        </div>
        <div class="flex flex-wrap items-center gap-3">
          <!-- Auth buttons -->
          <div v-if="!isAuthenticated" class="flex items-center gap-2">
            <button
              class="rounded-full border border-white/20 px-4 py-2 text-sm font-medium text-slate-100 transition hover:border-mint/60 hover:text-mint"
              @click="showLogin = true"
            >
              Login
            </button>
            <button
              class="rounded-full bg-mint px-4 py-2 text-sm font-semibold text-midnight shadow-glow transition hover:shadow-neon"
              @click="showRegister = true"
            >
              Register
            </button>
          </div>
          <div v-else class="flex items-center gap-2">
            <span class="text-sm text-slate-300">Hi, {{ user?.username }}!</span>
            <button
              class="rounded-full border border-white/20 px-4 py-2 text-sm font-medium text-slate-100 transition hover:border-mint/60 hover:text-mint"
              @click="logout"
            >
              Logout
            </button>
          </div>
          
          <button
            class="rounded-full border border-white/20 px-4 py-2 text-sm font-medium text-slate-100 transition hover:border-mint/60 hover:text-mint"
            @click="refreshState"
          >
            Refresh
          </button>
          
          <!-- Upload button with dropdown for new vs project upload -->
          <div class="relative" v-if="isAuthenticated">
            <button
              class="rounded-full bg-mint px-4 py-2 text-sm font-semibold text-midnight shadow-glow transition hover:shadow-neon"
              @click="showUploadMenu = !showUploadMenu"
            >
              Upload xcstrings
            </button>
            
            <div v-if="showUploadMenu" class="absolute right-0 mt-2 w-56 rounded-xl bg-midnight/90 border border-white/20 shadow-lg z-10">
              <button
                class="w-full text-left px-4 py-2 text-sm hover:bg-mint/10 hover:text-white rounded-t-xl"
                @click="triggerNewUpload"
              >
                New Project
              </button>
              <button
                class="w-full text-left px-4 py-2 text-sm hover:bg-mint/10 hover:text-white rounded-b-xl"
                @click="triggerProjectUpload"
              >
                To Project
              </button>
            </div>
          </div>
          <button
            v-else
            class="rounded-full bg-mint px-4 py-2 text-sm font-semibold text-midnight shadow-glow transition hover:shadow-neon"
            @click="triggerUpload"
          >
            Upload xcstrings
          </button>
          
          <input ref="fileInput" type="file" accept=".xcstrings,application/json" class="hidden" @change="onFileChange" />
          <input ref="projectFileInput" type="file" accept=".xcstrings,application/json" class="hidden" @change="onProjectFileChange" />
          
          <button
            class="rounded-full border border-white/20 px-4 py-2 text-sm font-semibold text-slate-900 transition hover:border-mint/60 hover:text-white disabled:cursor-not-allowed disabled:border-white/10 disabled:text-slate-500"
            :disabled="!hasFile"
            @click="exportFile"
          >
            Export
          </button>
        </div>
      </header>
      
      <!-- Auth Modals -->
      <div v-if="showLogin" class="fixed inset-0 bg-black/70 flex items-center justify-center z-50">
        <div class="bg-midnight rounded-2xl p-6 w-full max-w-md border border-white/20">
          <div class="flex justify-between items-center mb-4">
            <h2 class="text-xl font-semibold">Login</h2>
            <button @click="showLogin = false" class="text-slate-400 hover:text-white">×</button>
          </div>
          
          <form @submit.prevent="performLogin">
            <div class="mb-4">
              <label class="block text-sm text-slate-400 mb-2">Username</label>
              <input v-model="loginForm.username" class="w-full rounded-lg bg-white/5 px-3 py-2 ring-1 ring-white/10 focus:ring-2 focus:ring-mint" placeholder="Username" required />
            </div>
            <div class="mb-4">
              <label class="block text-sm text-slate-400 mb-2">Password</label>
              <input v-model="loginForm.password" type="password" class="w-full rounded-lg bg-white/5 px-3 py-2 ring-1 ring-white/10 focus:ring-2 focus:ring-mint" placeholder="Password" required />
            </div>
            <button type="submit" class="w-full rounded-xl bg-mint px-4 py-3 text-center text-sm font-semibold text-midnight shadow-lg shadow-mint/20 transition hover:shadow-mint/40">
              Login
            </button>
          </form>
        </div>
      </div>
      
      <div v-if="showRegister" class="fixed inset-0 bg-black/70 flex items-center justify-center z-50">
        <div class="bg-midnight rounded-2xl p-6 w-full max-w-md border border-white/20">
          <div class="flex justify-between items-center mb-4">
            <h2 class="text-xl font-semibold">Register</h2>
            <button @click="showRegister = false" class="text-slate-400 hover:text-white">×</button>
          </div>
          
          <form @submit.prevent="performRegister">
            <div class="mb-4">
              <label class="block text-sm text-slate-400 mb-2">Username</label>
              <input v-model="registerForm.username" class="w-full rounded-lg bg-white/5 px-3 py-2 ring-1 ring-white/10 focus:ring-2 focus:ring-mint" placeholder="Username" required />
            </div>
            <div class="mb-4">
              <label class="block text-sm text-slate-400 mb-2">Email</label>
              <input v-model="registerForm.email" type="email" class="w-full rounded-lg bg-white/5 px-3 py-2 ring-1 ring-white/10 focus:ring-2 focus:ring-mint" placeholder="Email" required />
            </div>
            <div class="mb-4">
              <label class="block text-sm text-slate-400 mb-2">Password</label>
              <input v-model="registerForm.password" type="password" class="w-full rounded-lg bg-white/5 px-3 py-2 ring-1 ring-white/10 focus:ring-2 focus:ring-mint" placeholder="Password" required />
            </div>
            <button type="submit" class="w-full rounded-xl bg-mint px-4 py-3 text-center text-sm font-semibold text-midnight shadow-lg shadow-mint/20 transition hover:shadow-mint/40">
              Register
            </button>
          </form>
        </div>
      </div>
      
      <!-- Project Selection Modal -->
      <div v-if="showProjectSelection" class="fixed inset-0 bg-black/70 flex items-center justify-center z-50">
        <div class="bg-midnight rounded-2xl p-6 w-full max-w-2xl border border-white/20 max-h-[80vh] overflow-auto">
          <div class="flex justify-between items-center mb-4">
            <h2 class="text-xl font-semibold">Select Project</h2>
            <button @click="closeProjectSelection" class="text-slate-400 hover:text-white">×</button>
          </div>
          
          <div class="mb-4">
            <button 
              @click="createNewProject" 
              class="w-full rounded-xl border border-white/20 px-4 py-3 text-left text-sm font-medium text-slate-100 transition hover:border-mint/60 hover:text-mint mb-2"
            >
              + Create New Project
            </button>
          </div>
          
          <div v-if="projects.length === 0" class="text-center py-8 text-slate-400">
            No projects found. Create one first.
          </div>
          <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div 
              v-for="project in projects" 
              :key="project.id"
              class="rounded-xl border border-white/20 bg-white/5 p-4 cursor-pointer hover:border-mint/60 transition"
              @click="selectProject(project.id)"
            >
              <h3 class="font-semibold">{{ project.name || 'Untitled Project' }}</h3>
              <p class="text-sm text-slate-400 mt-1">{{ project.description || 'No description' }}</p>
              <p class="text-xs text-slate-500 mt-2">{{ formatDate(project.createdAt) }}</p>
            </div>
          </div>
        </div>
      </div>
      
      <!-- Project Creation Modal -->
      <div v-if="showProjectCreation" class="fixed inset-0 bg-black/70 flex items-center justify-center z-50">
        <div class="bg-midnight rounded-2xl p-6 w-full max-w-md border border-white/20">
          <div class="flex justify-between items-center mb-4">
            <h2 class="text-xl font-semibold">Create New Project</h2>
            <button @click="showProjectCreation = false" class="text-slate-400 hover:text-white">×</button>
          </div>
          
          <form @submit.prevent="createProject">
            <div class="mb-4">
              <label class="block text-sm text-slate-400 mb-2">Project Name</label>
              <input v-model="newProject.name" class="w-full rounded-lg bg-white/5 px-3 py-2 ring-1 ring-white/10 focus:ring-2 focus:ring-mint" placeholder="My Localization Project" required />
            </div>
            <div class="mb-4">
              <label class="block text-sm text-slate-400 mb-2">Description</label>
              <textarea v-model="newProject.description" class="w-full rounded-lg bg-white/5 px-3 py-2 ring-1 ring-white/10 focus:ring-2 focus:ring-mint" placeholder="Project description"></textarea>
            </div>
            <div class="mb-4">
              <label class="block text-sm text-slate-400 mb-2">Source Language</label>
              <input v-model="newProject.sourceLanguage" class="w-full rounded-lg bg-white/5 px-3 py-2 ring-1 ring-white/10 focus:ring-2 focus:ring-mint" placeholder="e.g. en" />
            </div>
            <button type="submit" class="w-full rounded-xl bg-mint px-4 py-3 text-center text-sm font-semibold text-midnight shadow-lg shadow-mint/20 transition hover:shadow-mint/40">
              Create Project
            </button>
          </form>
        </div>
      </div>

      <div class="grid gap-4 lg:grid-cols-3">
        <section class="glass rounded-2xl p-6 lg:col-span-2">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm uppercase tracking-[0.2em] text-slate-400">Workspace</p>
              <h2 class="font-display text-xl font-semibold">File & Languages</h2>
            </div>
            <span v-if="hasFile" class="rounded-full bg-emerald-900/40 px-3 py-1 text-xs font-semibold text-mint">Loaded</span>
          </div>

          <div class="mt-5 grid gap-4 md:grid-cols-3">
            <div class="rounded-xl border border-white/10 bg-white/5 p-4">
              <p class="text-xs uppercase tracking-wide text-slate-400">File</p>
              <p class="mt-1 font-semibold" v-if="state.fileName">{{ state.fileName }}</p>
              <p class="mt-1 text-slate-400" v-else>Waiting for upload…</p>
              <p class="mt-2 text-xs text-slate-400">{{ state.totalStrings }} strings detected</p>
            </div>
            <div class="rounded-xl border border-white/10 bg-white/5 p-4">
              <p class="text-xs uppercase tracking-wide text-slate-400">Source language</p>
              <div class="mt-1 flex items-center gap-2">
                <input
                  v-model="state.sourceLanguage"
                  class="w-full rounded-lg bg-midnight/50 px-3 py-2 text-sm text-white ring-1 ring-white/10 focus:ring-2 focus:ring-mint"
                  placeholder="e.g. en"
                />
                <button class="rounded-lg border border-white/20 px-3 py-2 text-xs" @click="autoDetectSource">Auto</button>
              </div>
              <p class="mt-2 text-xs text-slate-400">Hint: we use this to prompt the provider.</p>
            </div>
            <div class="rounded-xl border border-white/10 bg-white/5 p-4">
              <p class="text-xs uppercase tracking-wide text-slate-400">Available languages</p>
              <div class="mt-2 flex flex-wrap gap-2">
                <span
                  v-for="lang in state.availableLanguages"
                  :key="lang"
                  class="rounded-full border border-white/10 px-2 py-1 text-xs text-slate-200"
                >
                  {{ lang }}
                </span>
                <span v-if="!state.availableLanguages.length" class="text-xs text-slate-400">–</span>
              </div>
            </div>
          </div>

          <div class="mt-6 rounded-xl border border-white/10 bg-white/5 p-4">
            <div class="flex items-center justify-between">
              <div>
                <p class="text-xs uppercase tracking-wide text-slate-400">Target languages</p>
                <p class="text-sm text-slate-300">Add one per code, press enter to confirm.</p>
              </div>
              <button
                class="rounded-full border border-white/20 px-3 py-1 text-xs text-slate-200 hover:border-mint/60 hover:text-mint"
                @click="useAvailableTargets"
              >
                Use detected
              </button>
            </div>

            <div class="mt-3 flex flex-wrap gap-2">
              <span
                v-for="lang in state.targetLanguages"
                :key="lang"
                class="group inline-flex items-center gap-2 rounded-full bg-white/10 px-3 py-1 text-sm"
              >
                {{ lang }}
                <button class="text-slate-400 transition group-hover:text-white" @click="removeTarget(lang)">×</button>
              </span>
              <input
                v-model="targetInput"
                class="tag-input rounded-full border border-dashed border-white/30 bg-midnight/40 px-3 py-1 text-sm"
                placeholder="Add language"
                @keydown.enter.prevent="addTarget"
                @blur="addTarget"
              />
            </div>

            <div class="mt-3 flex flex-wrap gap-2 text-xs text-slate-400">
              <button
                v-for="preset in presets"
                :key="preset"
                class="rounded-full border border-white/20 px-2 py-1 hover:border-mint/60 hover:text-mint"
                @click="addPreset(preset)"
              >
                {{ preset }}
              </button>
            </div>

            <div class="mt-4 space-y-2">
              <div class="flex items-center justify-between text-xs text-slate-400">
                <span>Language library</span>
                <button class="text-mint underline" @click="showAllLanguages = !showAllLanguages">
                  {{ showAllLanguages ? 'Collapse' : 'Show all' }}
                </button>
              </div>
              <div class="flex flex-wrap gap-2">
                <button
                  v-for="lang in languageOptions"
                  :key="lang.code"
                  class="rounded-full border border-white/20 px-3 py-1 text-xs text-slate-100 transition hover:border-mint/60 hover:text-mint"
                  @click="addPreset(lang.code)"
                >
                  {{ lang.name }} ({{ lang.code }})
                </button>
              </div>
            </div>
          </div>
        </section>

        <section class="glass rounded-2xl p-6 space-y-4">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm uppercase tracking-[0.2em] text-slate-400">Provider</p>
              <h2 class="font-display text-xl font-semibold">Batch translate</h2>
            </div>
            <span class="text-xs text-slate-400">{{ providerLabel }}</span>
          </div>

          <div class="grid grid-cols-2 gap-3">
            <button
              v-for="p in providers"
              :key="p.id"
              class="rounded-xl border px-3 py-2 text-left transition"
              :class="p.id === state.provider ? 'border-mint/60 bg-mint/10 text-white' : 'border-white/10 bg-white/5 text-slate-300'"
              @click="state.provider = p.id"
            >
              <p class="font-semibold">{{ p.name }}</p>
              <p class="text-xs text-slate-400">{{ p.hint }}</p>
            </button>
          </div>

          <div class="space-y-3 text-sm text-slate-200">
            <template v-if="state.provider === 'openai'">
              <label class="block">
                <span class="text-xs text-slate-400">API key</span>
                <input v-model="state.openai.apiKey" class="mt-1 w-full rounded-lg bg-midnight/40 px-3 py-2 ring-1 ring-white/10 focus:ring-2 focus:ring-mint" type="password" placeholder="sk-…" />
              </label>
              <label class="block">
                <span class="text-xs text-slate-400">Model</span>
                <input v-model="state.openai.model" class="mt-1 w-full rounded-lg bg-midnight/40 px-3 py-2 ring-1 ring-white/10 focus:ring-2 focus:ring-mint" />
              </label>
              <label class="block">
                <span class="text-xs text-slate-400">API base (optional)</span>
                <input v-model="state.openai.apiBaseUrl" class="mt-1 w-full rounded-lg bg-midnight/40 px-3 py-2 ring-1 ring-white/10 focus:ring-2 focus:ring-mint" placeholder="https://api.openai.com" />
              </label>
            </template>
            <template v-else-if="state.provider === 'google'">
              <label class="block">
                <span class="text-xs text-slate-400">Google API key</span>
                <input v-model="state.google.apiKey" class="mt-1 w-full rounded-lg bg-midnight/40 px-3 py-2 ring-1 ring-white/10 focus:ring-2 focus:ring-mint" placeholder="AIza…" />
              </label>
            </template>
            <template v-else-if="state.provider === 'deepl'">
              <label class="block">
                <span class="text-xs text-slate-400">DeepL API key</span>
                <input v-model="state.deepl.apiKey" class="mt-1 w-full rounded-lg bg-midnight/40 px-3 py-2 ring-1 ring-white/10 focus:ring-2 focus:ring-mint" placeholder="auth-key" />
              </label>
              <label class="block">
                <span class="text-xs text-slate-400">Formality</span>
                <input v-model="state.deepl.formality" class="mt-1 w-full rounded-lg bg-midnight/40 px-3 py-2 ring-1 ring-white/10 focus:ring-2 focus:ring-mint" placeholder="default" />
              </label>
              <label class="inline-flex items-center gap-2 text-xs text-slate-300">
                <input v-model="state.deepl.isFree" type="checkbox" class="h-4 w-4 rounded border-white/30 bg-midnight/50" />
                Use free API
              </label>
            </template>
            <template v-else>
              <label class="block">
                <span class="text-xs text-slate-400">Baidu App ID</span>
                <input v-model="state.baidu.appId" class="mt-1 w-full rounded-lg bg-midnight/40 px-3 py-2 ring-1 ring-white/10 focus:ring-2 focus:ring-mint" />
              </label>
              <label class="block">
                <span class="text-xs text-slate-400">Baidu App Secret</span>
                <input v-model="state.baidu.appSecret" class="mt-1 w-full rounded-lg bg-midnight/40 px-3 py-2 ring-1 ring-white/10 focus:ring-2 focus:ring-mint" type="password" />
              </label>
            </template>

            <div class="grid grid-cols-2 gap-3 text-xs text-slate-300">
              <label class="block">
                <span>Concurrency</span>
                <input v-model.number="state.concurrency" type="number" min="1" class="mt-1 w-full rounded-lg bg-midnight/40 px-3 py-2 ring-1 ring-white/10 focus:ring-2 focus:ring-mint" />
              </label>
              <label class="block">
                <span>Timeout (sec)</span>
                <input v-model.number="state.timeoutSeconds" type="number" min="30" class="mt-1 w-full rounded-lg bg-midnight/40 px-3 py-2 ring-1 ring-white/10 focus:ring-2 focus:ring-mint" />
              </label>
            </div>

            <div class="grid grid-cols-2 gap-3">
              <button
                class="w-full rounded-xl bg-mint px-4 py-3 text-center text-sm font-semibold text-midnight shadow-lg shadow-mint/20 transition hover:shadow-mint/40 disabled:cursor-not-allowed disabled:opacity-60"
                :disabled="isTranslating || !hasFile || !state.targetLanguages.length"
                @click="batchTranslate"
              >
                <span v-if="isTranslating">Translating…</span>
                <span v-else>Run batch translate</span>
              </button>
              <button
                class="w-full rounded-xl border border-white/20 px-4 py-3 text-center text-sm font-semibold text-slate-100 transition hover:border-mint/60 hover:text-mint"
                @click="saveOptions"
              >
                Save options
              </button>
            </div>
            <p v-if="progress.id" class="text-xs text-slate-300">
              Progress: {{ progress.done }}/{{ progress.total || '…' }} • {{ progress.status }}
            </p>
            <p class="text-xs text-slate-400">We only translate missing entries for the selected targets.</p>
          </div>
        </section>
      </div>

      <section class="glass rounded-2xl p-6">
        <div class="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
          <div>
            <p class="text-sm uppercase tracking-[0.2em] text-slate-400">Strings</p>
            <h2 class="font-display text-xl font-semibold">Visual localisation</h2>
            <p class="text-slate-400">Browse source text next to target translations. Missing items are highlighted.</p>
          </div>
          <div class="flex items-center gap-2 text-xs text-slate-300">
            <label class="flex items-center gap-2">
              <span>Filter</span>
              <input v-model="filter" class="rounded-lg bg-midnight/50 px-3 py-2 ring-1 ring-white/10 focus:ring-2 focus:ring-mint" placeholder="Search key or text" />
            </label>
          </div>
        </div>

        <div class="mt-4 overflow-hidden rounded-xl border border-white/10 bg-white/5">
          <div class="max-h-[480px] overflow-auto scrollbar-thin">
            <table class="min-w-full divide-y divide-white/5 text-sm">
              <thead class="sticky top-0 bg-midnight">
                <tr>
                  <th class="px-4 py-3 text-left font-semibold text-slate-300">Key</th>
                  <th class="px-4 py-3 text-left font-semibold text-slate-300">{{ state.sourceLanguage || 'source' }}</th>
                  <th
                    v-for="lang in displayTargets"
                    :key="lang"
                    class="px-4 py-3 text-left font-semibold text-slate-300"
                  >
                    {{ lang }}
                  </th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="row in filteredEntries" :key="row.key" class="divide-y divide-white/5">
                  <td class="px-4 py-3 align-top font-mono text-xs text-slate-300">{{ row.key }}</td>
                  <td class="px-4 py-3 align-top text-slate-100">
                    <p class="whitespace-pre-line">{{ row.source || '—' }}</p>
                    <p class="mt-1 text-xs text-slate-400">{{ row.state }}</p>
                  </td>
                  <td
                    v-for="lang in displayTargets"
                    :key="lang"
                    class="px-4 py-3 align-top"
                    :class="row.missing.includes(lang) ? 'bg-orange-500/5 text-orange-200' : 'text-slate-100'"
                  >
                    <p class="whitespace-pre-line">{{ row.translations[lang] || '–' }}</p>
                  </td>
                </tr>
                <tr v-if="!filteredEntries.length">
                  <td class="px-4 py-6 text-center text-slate-400" :colspan="2 + displayTargets.length">No strings yet.</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </section>

      <section v-if="statusMessage" class="rounded-xl border border-white/10 bg-white/5 px-4 py-3 text-sm" :class="statusClass">
        {{ statusMessage }}
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue'

type LocalizationEntry = {
  key: string
  source: string
  state: string
  translations: Record<string, string>
  missing: string[]
}

type Payload = {
  fileName: string
  sourceLanguage: string
  availableLanguages: string[]
  totalStrings: number
  entries: LocalizationEntry[]
  warning?: string
}

type ProviderId = 'openai' | 'google' | 'deepl' | 'baidu'
type JobState = { id: string; status: string; done: number; total: number; message?: string }

type User = {
  id: number;
  username: string;
  email: string;
  isActive: boolean;
}

type Project = {
  id: number;
  name: string;
  description: string;
  fileName: string;
  sourceLanguage: string;
  createdAt: string;
  updatedAt: string;
}

type ProjectForm = {
  name: string;
  description: string;
  sourceLanguage: string;
}

const presets = ['zh-Hans', 'ja', 'ko', 'de', 'fr', 'es', 'ar']
const languages = [
  { code: 'zh-Hans', name: 'Chinese (Simplified)' },
  { code: 'zh-Hant', name: 'Chinese (Traditional)' },
  { code: 'en', name: 'English' },
  { code: 'ja', name: 'Japanese' },
  { code: 'ko', name: 'Korean' },
  { code: 'de', name: 'German' },
  { code: 'fr', name: 'French' },
  { code: 'es', name: 'Spanish' },
  { code: 'ar', name: 'Arabic' },
  { code: 'pt', name: 'Portuguese' },
  { code: 'ru', name: 'Russian' },
  { code: 'it', name: 'Italian' },
  { code: 'nl', name: 'Dutch' },
  { code: 'pl', name: 'Polish' },
  { code: 'sv', name: 'Swedish' },
  { code: 'da', name: 'Danish' },
  { code: 'fi', name: 'Finnish' },
  { code: 'no', name: 'Norwegian' },
  { code: 'cs', name: 'Czech' },
  { code: 'tr', name: 'Turkish' },
  { code: 'he', name: 'Hebrew' },
  { code: 'hi', name: 'Hindi' },
  { code: 'id', name: 'Indonesian' },
  { code: 'th', name: 'Thai' },
  { code: 'vi', name: 'Vietnamese' }
]
const providers = [
  { id: 'openai' as ProviderId, name: 'OpenAI', hint: 'GPT style chat completion' },
  { id: 'google' as ProviderId, name: 'Google', hint: 'Google Cloud translation' },
  { id: 'deepl' as ProviderId, name: 'DeepL', hint: 'Great for EU languages' },
  { id: 'baidu' as ProviderId, name: 'Baidu', hint: 'China-friendly' }
]

const state = reactive({
  fileName: '',
  sourceLanguage: 'en',
  targetLanguages: [] as string[],
  availableLanguages: [] as string[],
  totalStrings: 0,
  entries: [] as LocalizationEntry[],
  provider: 'openai' as ProviderId,
  openai: { apiKey: '', apiBaseUrl: '', model: 'gpt-3.5-turbo', temperature: 0.3, maxTokens: 1024 },
  google: { apiKey: '' },
  deepl: { apiKey: '', isFree: true, formality: 'default' },
  baidu: { appId: '', appSecret: '' },
  concurrency: 6,
  timeoutSeconds: 300
})

// Authentication state
const isAuthenticated = ref(false)
const user = ref<User | null>(null)
const showLogin = ref(false)
const showRegister = ref(false)
const loginForm = reactive({
  username: '',
  password: ''
})
const registerForm = reactive({
  username: '',
  email: '',
  password: ''
})

// Project management state
const projects = ref<Project[]>([])
const showProjectSelection = ref(false)
const showProjectCreation = ref(false)
const showUploadMenu = ref(false)
const newProject = reactive<ProjectForm>({
  name: '',
  description: '',
  sourceLanguage: ''
})
const projectFileInput = ref<HTMLInputElement | null>(null)

const isTranslating = ref(false)
const progress = reactive<JobState>({ id: '', status: 'idle', done: 0, total: 0 })
let progressTimer: number | null = null
const statusMessage = ref('')
const statusTone = ref<'info' | 'error'>('info')
const filter = ref('')
const targetInput = ref('')
const fileInput = ref<HTMLInputElement | null>(null)
const showAllLanguages = ref(false)

const LOCAL_KEY = 'xcstrings-translator-ui'

const providerLabel = computed(() => providers.find((p) => p.id === state.provider)?.name ?? '')
const hasFile = computed(() => !!state.fileName)
const displayTargets = computed(() => state.targetLanguages)
const filteredEntries = computed(() => {
  const term = filter.value.trim().toLowerCase()
  if (!term) return state.entries
  return state.entries.filter((row) =>
    row.key.toLowerCase().includes(term) ||
    row.source.toLowerCase().includes(term) ||
    displayTargets.value.some((lang) => (row.translations[lang] || '').toLowerCase().includes(term))
  )
})

const statusClass = computed(() =>
  statusTone.value === 'error'
    ? 'text-red-200 border-red-400/30 bg-red-900/30'
    : 'text-mint border-mint/50 bg-mint/10'
)

const languageOptions = computed(() => (showAllLanguages.value ? languages : languages.slice(0, 20)))

// Helper functions
function showStatus(message: string, tone: 'info' | 'error' = 'info') {
  statusMessage.value = message
  statusTone.value = tone
}

function triggerUpload() {
  fileInput.value?.click()
}

function triggerNewUpload() {
  triggerUpload()
  showUploadMenu.value = false
}

function triggerProjectUpload() {
  loadProjects()
  showProjectSelection.value = true
  showUploadMenu.value = false
}

function addTarget() {
  const trimmed = targetInput.value.trim()
  if (!trimmed) return
  if (!state.targetLanguages.includes(trimmed)) {
    state.targetLanguages.push(trimmed)
  }
  targetInput.value = ''
}

function addPreset(code: string) {
  if (!state.targetLanguages.includes(code)) {
    state.targetLanguages.push(code)
  }
}

function removeTarget(lang: string) {
  state.targetLanguages = state.targetLanguages.filter((l) => l !== lang)
}

function useAvailableTargets() {
  state.targetLanguages = state.availableLanguages.filter((lang) => lang !== state.sourceLanguage)
}

function autoDetectSource() {
  if (state.availableLanguages.includes(state.sourceLanguage)) return
  if (state.availableLanguages.length) {
    state.sourceLanguage = state.availableLanguages[0]
  }
}

// Authentication functions
async function performLogin() {
  try {
    const res = await fetch('/api/auth/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(loginForm)
    })

    if (!res.ok) {
      const error = await res.json()
      showStatus(`Login failed: ${error.error || 'Unknown error'}`, 'error')
      return
    }

    const data = await res.json()
    localStorage.setItem('token', data.token)
    user.value = data.user
    isAuthenticated.value = true
    showLogin.value = false
    
    // Reset form
    loginForm.username = ''
    loginForm.password = ''
    
    showStatus('Login successful.', 'info')
  } catch (err) {
    showStatus('Login failed: ' + (err as Error).message, 'error')
  }
}

async function performRegister() {
  try {
    const res = await fetch('/api/auth/register', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(registerForm)
    })

    if (!res.ok) {
      const error = await res.json()
      showStatus(`Registration failed: ${error.error || 'Unknown error'}`, 'error')
      return
    }

    const data = await res.json()
    showStatus('Registration successful. Please login.', 'info')
    showRegister.value = false
    
    // Reset form
    registerForm.username = ''
    registerForm.email = ''
    registerForm.password = ''
  } catch (err) {
    showStatus('Registration failed: ' + (err as Error).message, 'error')
  }
}

async function logout() {
  localStorage.removeItem('token')
  isAuthenticated.value = false
  user.value = null
  showStatus('Logged out successfully.', 'info')
}

// Project-related functions
async function loadProjects() {
  const token = localStorage.getItem('token')
  if (!token) {
    showStatus('Please log in to access projects.', 'error')
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

function createNewProject() {
  showProjectCreation.value = true
  showProjectSelection.value = false
}

function closeProjectSelection() {
  showProjectSelection.value = false
}

async function createProject() {
  const token = localStorage.getItem('token')
  if (!token) {
    showStatus('Please log in to create projects.', 'error')
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
        name: newProject.name,
        description: newProject.description,
        fileName: '',
        fileContent: '',
        sourceLanguage: newProject.sourceLanguage
      })
    })

    if (!res.ok) {
      throw new Error(`Failed to create project: ${await res.text()}`)
    }

    const data = await res.json()
    projects.value.unshift(data.project)
    showProjectCreation.value = false
    
    // Reset form
    newProject.name = ''
    newProject.description = ''
    newProject.sourceLanguage = ''
    
    showStatus('Project created successfully.', 'info')
  } catch (err) {
    showStatus('Failed to create project: ' + (err as Error).message, 'error')
  }
}

async function selectProject(projectId: number) {
  // Load project data
  const token = localStorage.getItem('token')
  if (!token) {
    showStatus('Please log in to access projects.', 'error')
    return
  }

  try {
    const res = await fetch(`/api/protected/projects/${projectId}`, {
      headers: { 'Authorization': `Bearer ${token}` }
    })

    if (!res.ok) {
      throw new Error(`Failed to load project: ${await res.text()}`)
    }

    const data = await res.json()
    showProjectSelection.value = false
    
    // Since we can't directly load the project into local state,
    // we will just close the modal and show a message
    showStatus(`Project selected. File: ${data.project.fileName}`, 'info')
  } catch (err) {
    showStatus('Failed to load project: ' + (err as Error).message, 'error')
  }
}

async function onProjectFileChange(event: Event) {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return

  // Show project selection to upload to
  loadProjects()
  showProjectSelection.value = true
  input.value = ''
}

function formatDate(dateString: string) {
  const date = new Date(dateString)
  return date.toLocaleDateString()
}

// Existing functions remain the same
async function onFileChange(event: Event) {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return
  await uploadFile(file)
  input.value = ''
}

async function uploadFile(file: File) {
  const form = new FormData()
  form.append('file', file)
  if (state.sourceLanguage) {
    form.append('sourceLanguage', state.sourceLanguage)
  }

  showStatus('Uploading…')
  const res = await fetch('/api/upload', {
    method: 'POST',
    body: form
  })

  if (!res.ok) {
    showStatus(`Upload failed: ${await res.text()}`, 'error')
    return
  }

  const payload = (await res.json()) as Payload
  applyPayload(payload)
  if (!state.targetLanguages.length) {
    useAvailableTargets()
  }
  if (payload.warning) {
    showStatus(payload.warning, 'error')
  } else {
    showStatus('File loaded. Ready to translate.', 'info')
  }
}

function applyPayload(payload: Payload) {
  state.fileName = payload.fileName
  state.sourceLanguage = payload.sourceLanguage
  state.availableLanguages = payload.availableLanguages
  state.totalStrings = payload.totalStrings
  state.entries = payload.entries
}

function saveOptions() {
  saveLocalState(snapshotOptions())
  showStatus('Options saved locally.', 'info')
}

async function batchTranslate() {
  if (!state.fileName) {
    showStatus('Upload a file first.', 'error')
    return
  }
  if (!state.targetLanguages.length) {
    showStatus('Add at least one target language.', 'error')
    return
  }

  const body = {
    provider: state.provider,
    targetLanguages: state.targetLanguages,
    sourceLanguage: state.sourceLanguage,
    concurrency: state.concurrency,
    timeoutSeconds: state.timeoutSeconds,
    config: {
      apiKey: state.provider === 'baidu' ? undefined : getApiKey(),
      apiBaseUrl: state.openai.apiBaseUrl,
      model: state.provider === 'openai' ? state.openai.model : undefined,
      temperature: state.provider === 'openai' ? state.openai.temperature : undefined,
      maxTokens: state.provider === 'openai' ? state.openai.maxTokens : undefined,
      appId: state.baidu.appId,
      appSecret: state.baidu.appSecret,
      formality: state.deepl.formality,
      isFree: state.deepl.isFree
    }
  }

  isTranslating.value = true
  showStatus('Running batch translation…')

  const res = await fetch('/api/translate', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body)
  })

  if (!res.ok) {
    isTranslating.value = false
    showStatus(`Translate failed: ${await res.text()}`, 'error')
    return
  }

  const { jobId } = (await res.json()) as { jobId: string }
  if (!jobId) {
    isTranslating.value = false
    showStatus('Translate job did not start.', 'error')
    return
  }
  startProgress(jobId)
}

function getApiKey() {
  if (state.provider === 'openai') return state.openai.apiKey
  if (state.provider === 'google') return state.google.apiKey
  return state.deepl.apiKey
}

async function exportFile() {
  const res = await fetch('/api/export')
  if (!res.ok) {
    showStatus('No file loaded to export.', 'error')
    return
  }
  const blob = await res.blob()
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = state.fileName || 'Localizable_translated.xcstrings'
  link.click()
  URL.revokeObjectURL(url)
  showStatus('Exported xcstrings file.', 'info')
}

async function refreshState() {
  const res = await fetch('/api/strings')
  if (!res.ok) return
  const payload = (await res.json()) as Payload
  applyPayload(payload)
  if (!state.targetLanguages.length) {
    useAvailableTargets()
  }
}

// Check authentication on mount
onMounted(async () => {
  const token = localStorage.getItem('token')
  if (token) {
    try {
      // Verify token by making a request to protected endpoint
      const res = await fetch('/api/protected/projects', {
        headers: { 'Authorization': `Bearer ${token}` }
      })
      
      if (res.ok) {
        isAuthenticated.value = true
        // Get user info - in a real implementation we'd have a /me endpoint
        // For now, we'll just assume the user is valid
      } else {
        localStorage.removeItem('token')
      }
    } catch (err) {
      console.error('Token verification failed:', err)
      localStorage.removeItem('token')
    }
  }

  refreshState().catch(() => null)

  const saved = loadLocalState()
  if (saved) {
    Object.assign(state, saved)
  }
})

watch(
  () => snapshotOptions(),
  (val) => saveLocalState(val),
  { deep: true }
)

function startProgress(id: string) {
  progress.id = id
  progress.status = 'running'
  progress.done = 0
  progress.total = 0
  isTranslating.value = true
  if (progressTimer) {
    clearInterval(progressTimer)
  }
  pollProgress()
  progressTimer = window.setInterval(pollProgress, 1200)
}

async function pollProgress() {
  const res = await fetch('/api/progress')
  if (!res.ok) return
  const data = (await res.json()) as { job?: JobState | null; payload?: Payload }

  if (data.payload) {
    applyPayload(data.payload)
  }

  if (data.job) {
    progress.id = data.job.id
    progress.status = data.job.status
    progress.done = data.job.done
    progress.total = data.job.total
    if (data.job.status !== 'running') {
      stopProgress()
      showStatus(data.job.status === 'done' ? 'Translations applied.' : data.job.message || 'Translation stopped.', data.job.status === 'done' ? 'info' : 'error')
    }
  } else if (progress.id) {
    stopProgress()
  }
}

function stopProgress() {
  if (progressTimer) {
    clearInterval(progressTimer)
    progressTimer = null
  }
  isTranslating.value = false
  progress.id = ''
}

function snapshotOptions() {
  return {
    provider: state.provider,
    sourceLanguage: state.sourceLanguage,
    targetLanguages: [...state.targetLanguages],
    concurrency: state.concurrency,
    timeoutSeconds: state.timeoutSeconds,
    openai: { ...state.openai },
    google: { ...state.google },
    deepl: { ...state.deepl },
    baidu: { ...state.baidu }
  }
}

function loadLocalState() {
  try {
    const raw = localStorage.getItem(LOCAL_KEY)
    if (!raw) return null
    return JSON.parse(raw)
  } catch (err) {
    console.warn('Failed to load local state', err)
    return null
  }
}

function saveLocalState(val: unknown) {
  try {
    localStorage.setItem(LOCAL_KEY, JSON.stringify(val))
  } catch (err) {
    console.warn('Failed to persist local state', err)
  }
}
</script>
