<template>
  <div class="space-y-6">
    <header class="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
      <div>
        <p class="text-xs uppercase tracking-[0.3em] text-slate-500">{{ t('nav.languages') }}</p>
        <h1 class="text-2xl font-semibold">{{ t('languages.title') }}</h1>
        <p class="text-sm text-slate-400">{{ t('languages.subtitle') }}</p>
      </div>
      <div class="text-xs text-slate-500">{{ t('languages.readonly') }}</div>
    </header>

    <section class="rounded-2xl border border-white/10 bg-white/5 p-4">
      <div class="flex items-center justify-between mb-3">
        <div class="text-sm text-slate-400">语言列表</div>
        <input class="rounded-lg bg-midnight/40 px-3 py-2 text-xs ring-1 ring-white/10" placeholder="搜索语言或代码" />
      </div>
      <div class="overflow-x-auto">
        <table class="min-w-full text-sm">
          <thead class="text-left text-slate-500">
            <tr>
              <th class="py-2">名称</th>
              <th class="py-2">代码</th>
              <th class="py-2">别名</th>
              <th class="py-2">方向</th>
              <th class="py-2">可用</th>
              <th class="py-2">操作</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-white/5">
            <tr v-for="lang in languages" :key="lang.code" class="hover:bg-white/5">
              <td class="py-2">{{ lang.name }}</td>
              <td class="py-2">{{ lang.code }}</td>
              <td class="py-2">{{ lang.alias || '-' }}</td>
              <td class="py-2">{{ lang.rtl ? 'RTL' : 'LTR' }}</td>
              <td class="py-2">
                <span class="rounded-full px-2 py-1 text-xs" :class="lang.enabled ? 'bg-emerald-900/40 text-emerald-200' : 'bg-rose-900/40 text-rose-200'">
                  {{ lang.enabled ? '可用' : '禁用' }}
                </span>
              </td>
              <td class="py-2 space-x-2 text-xs">
                <button class="rounded border border-white/20 px-2 py-1 hover:border-mint/60 hover:text-mint">编辑</button>
                <button class="rounded border border-white/20 px-2 py-1 hover:border-mint/60 hover:text-mint">禁用</button>
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

const languages = [
  { code: 'en', name: 'English', alias: 'EN', rtl: false, enabled: true },
  { code: 'zh-CN', name: '简体中文', alias: 'Chinese Simplified', rtl: false, enabled: true },
  { code: 'ar', name: 'Arabic', alias: 'العربية', rtl: true, enabled: false }
]
</script>
