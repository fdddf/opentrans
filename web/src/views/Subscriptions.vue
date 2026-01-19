<template>
  <div class="space-y-6">
    <header class="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
      <div>
        <p class="text-xs uppercase tracking-[0.3em] text-slate-500">{{ t('nav.subscriptions') }}</p>
        <h1 class="text-2xl font-semibold">{{ t('subscriptions.title') }}</h1>
        <p class="text-sm text-slate-400">{{ t('subscriptions.subtitle') }}</p>
      </div>
      <div class="flex items-center gap-2">
        <button class="rounded-lg border border-white/20 px-3 py-2 text-sm hover:border-mint/60 hover:text-mint">{{ t('subscriptions.new') }}</button>
        <button class="rounded-lg bg-mint px-3 py-2 text-sm font-semibold text-midnight shadow">{{ t('subscriptions.sync') }}</button>
      </div>
    </header>

    <section class="rounded-2xl border border-white/10 bg-white/5 p-4">
      <div class="flex items-center justify-between mb-3">
        <div class="text-sm text-slate-400">订阅列表</div>
        <input class="rounded-lg bg-midnight/40 px-3 py-2 text-xs ring-1 ring-white/10" placeholder="搜索用户或套餐" />
      </div>
      <div class="overflow-x-auto">
        <table class="min-w-full text-sm">
          <thead class="text-left text-slate-500">
            <tr>
              <th class="py-2">用户</th>
              <th class="py-2">套餐</th>
              <th class="py-2">状态</th>
              <th class="py-2">续订日期</th>
              <th class="py-2">操作</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-white/5">
            <tr v-for="item in subscriptions" :key="item.id" class="hover:bg-white/5">
              <td class="py-2">{{ item.user }}</td>
              <td class="py-2">{{ item.plan }}</td>
              <td class="py-2">
                <span class="rounded-full px-2 py-1 text-xs" :class="item.active ? 'bg-emerald-900/40 text-emerald-200' : 'bg-rose-900/40 text-rose-200'">
                  {{ item.active ? '生效中' : '已过期' }}
                </span>
              </td>
              <td class="py-2">{{ item.renewal }}</td>
              <td class="py-2 space-x-2 text-xs">
                <button class="rounded border border-white/20 px-2 py-1 hover:border-mint/60 hover:text-mint">续订</button>
                <button class="rounded border border-white/20 px-2 py-1 hover:border-mint/60 hover:text-mint">取消</button>
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

const subscriptions = [
  { id: 1, user: 'Admin', plan: 'Pro', active: true, renewal: '2024-12-01' },
  { id: 2, user: 'User A', plan: 'Starter', active: true, renewal: '2024-08-15' },
  { id: 3, user: 'User B', plan: 'Free', active: false, renewal: '-' }
]
</script>
