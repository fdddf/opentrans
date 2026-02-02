<template>
  <div class="space-y-6">
    <header class="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
      <div>
        <p class="text-xs uppercase tracking-[0.3em] text-slate-500">{{ t('nav.users') }}</p>
        <h1 class="text-2xl font-semibold">{{ t('users.title') }}</h1>
        <p class="text-sm text-slate-400">{{ t('users.subtitle') }}</p>
      </div>
      <div class="flex items-center gap-2">
        <button class="rounded-lg border border-white/20 px-3 py-2 text-sm hover:border-mint/60 hover:text-mint">{{ t('users.inviteAdmin') }}</button>
        <button class="rounded-lg bg-mint px-3 py-2 text-sm font-semibold text-midnight shadow">{{ t('users.createUser') }}</button>
      </div>
    </header>

    <section class="rounded-2xl border border-white/10 bg-white/5 p-4">
      <div class="flex items-center justify-between mb-3">
        <div class="text-sm text-slate-400">{{ t('users.userList') }}</div>
        <input class="rounded-lg bg-midnight/40 px-3 py-2 text-xs ring-1 ring-white/10" :placeholder="t('users.searchPlaceholder')" />
      </div>
      <div class="overflow-x-auto">
        <table class="min-w-full text-sm">
          <thead class="text-left text-slate-500">
            <tr>
              <th class="py-2">{{ t('users.user') }}</th>
              <th class="py-2">{{ t('users.role') }}</th>
              <th class="py-2">{{ t('users.subscription') }}</th>
              <th class="py-2">{{ t('users.status') }}</th>
              <th class="py-2">{{ t('common.actions') }}</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-white/5">
            <tr v-for="user in users" :key="user.id" class="hover:bg-white/5">
              <td class="py-2">
                <div class="font-semibold">{{ user.username }}</div>
                <div class="text-xs text-slate-500">{{ user.email }}</div>
              </td>
              <td class="py-2">
                <span class="rounded-full px-2 py-1 text-xs bg-white/10 text-slate-200">
                  {{ t('users.userRegular') }}
                </span>
              </td>
              <td class="py-2">{{ user.subscriptionType }}</td>
              <td class="py-2">
                <span class="rounded-full bg-emerald-900/40 px-2 py-1 text-xs text-emerald-200" v-if="user.isActive">{{ t('users.active') }}</span>
                <span class="rounded-full bg-rose-900/40 px-2 py-1 text-xs text-rose-200" v-else>{{ t('users.inactive') }}</span>
              </td>
              <td class="py-2 space-x-2 text-xs">
                <button class="rounded border border-white/20 px-2 py-1 hover:border-mint/60 hover:text-mint">{{ t('users.resetPassword') }}</button>
                <button class="rounded border border-white/20 px-2 py-1 hover:border-mint/60 hover:text-mint" @click="toggleUserStatus(user)">
                  {{ user.isActive ? t('users.deactivate') : t('users.activate') }}
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useApi } from '../composables/useApi'
import { useToast } from '../composables/useToast'
import type { User } from '../composables/useApi'

const { t } = useI18n()
const { api } = useApi()
const toast = useToast()

const users = ref<User[]>([])

async function fetchUsers() {
  try {
    const response = await api.getUsers()
    if (response.success) {
      users.value = response.users
    }
  } catch (error: any) {
    console.error('Failed to fetch users:', error)
    toast.error('Failed to fetch users: ' + (error.message || 'Unknown error'))
  }
}

async function toggleUserStatus(user: User) {
  try {
    if (user.isActive) {
      await api.deactivateUser(user.id)
      toast.success(`User ${user.username} deactivated`)
    } else {
      await api.activateUser(user.id)
      toast.success(`User ${user.username} activated`)
    }
    fetchUsers()
  } catch (error: any) {
    console.error('Failed to toggle user status:', error)
    toast.error('Failed to toggle user status: ' + (error.message || 'Unknown error'))
  }
}

onMounted(() => {
  fetchUsers()
})
</script>
