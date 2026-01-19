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
        <div class="text-sm text-slate-400">租户内用户列表</div>
        <input class="rounded-lg bg-midnight/40 px-3 py-2 text-xs ring-1 ring-white/10" placeholder="搜索用户名或邮箱" />
      </div>
      <div class="overflow-x-auto">
        <table class="min-w-full text-sm">
          <thead class="text-left text-slate-500">
            <tr>
              <th class="py-2">用户</th>
              <th class="py-2">角色</th>
              <th class="py-2">订阅</th>
              <th class="py-2">状态</th>
              <th class="py-2">操作</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-white/5">
            <tr v-for="user in users" :key="user.id" class="hover:bg-white/5">
              <td class="py-2">
                <div class="font-semibold">{{ user.name }}</div>
                <div class="text-xs text-slate-500">{{ user.email }}</div>
              </td>
              <td class="py-2">
                <span class="rounded-full px-2 py-1 text-xs" :class="user.role === 'admin' ? 'bg-indigo-900/40 text-indigo-200' : 'bg-white/10 text-slate-200'">
                  {{ user.role === 'admin' ? '管理员' : '普通用户' }}
                </span>
              </td>
              <td class="py-2">{{ user.subscription }}</td>
              <td class="py-2">
                <span class="rounded-full bg-emerald-900/40 px-2 py-1 text-xs text-emerald-200" v-if="user.active">启用</span>
                <span class="rounded-full bg-rose-900/40 px-2 py-1 text-xs text-rose-200" v-else>停用</span>
              </td>
              <td class="py-2 space-x-2 text-xs">
                <button class="rounded border border-white/20 px-2 py-1 hover:border-mint/60 hover:text-mint">重置密码</button>
                <button class="rounded border border-white/20 px-2 py-1 hover:border-mint/60 hover:text-mint">停用</button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const users = [
  { id: 1, name: 'Admin', email: 'admin@example.com', role: 'admin', subscription: 'Pro', active: true },
  { id: 2, name: 'User A', email: 'usera@example.com', role: 'user', subscription: 'Starter', active: true },
  { id: 3, name: 'User B', email: 'userb@example.com', role: 'user', subscription: 'None', active: false }
]
</script>
