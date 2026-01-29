<template>
  <div class="p-6">
    <h2 class="text-xl font-bold mb-4">Translation Test Component</h2>
    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div class="border border-white/20 p-4 rounded-lg">
        <h3 class="text-lg font-semibold mb-2">Common Translations Used in the App</h3>
        <ul class="space-y-2">
          <li><strong>common.save:</strong> {{ t('common.save') }}</li>
          <li><strong>common.edit:</strong> {{ t('common.edit') }}</li>
          <li><strong>common.update:</strong> {{ t('common.update') }}</li>
          <li><strong>common.delete:</strong> {{ t('common.delete') }}</li>
          <li><strong>common.cancel:</strong> {{ t('common.cancel') }}</li>
          <li><strong>common.status:</strong> {{ t('common.status') }}</li>
          <li><strong>common.actions:</strong> {{ t('common.actions') }}</li>
          <li><strong>common.active:</strong> {{ t('common.active') }}</li>
          <li><strong>common.inactive:</strong> {{ t('common.inactive') }}</li>
        </ul>
      </div>
      <div class="border border-white/20 p-4 rounded-lg">
        <h3 class="text-lg font-semibold mb-2">Navigation Translations</h3>
        <ul class="space-y-2">
          <li><strong>nav.appleConnectConfig:</strong> {{ t('nav.appleConnectConfig') }}</li>
          <li><strong>nav.subscriptions:</strong> {{ t('nav.subscriptions') }}</li>
          <li><strong>nav.apps:</strong> {{ t('nav.apps') }}</li>
          <li><strong>nav.applocalizations:</strong> {{ t('nav.applocalizations') }}</li>
        </ul>
      </div>
    </div>
    <div class="mt-4 flex gap-2">
      <button @click="changeLocale('en')" class="px-3 py-1 border border-white/20 rounded hover:bg-mint/10">English</button>
      <button @click="changeLocale('zh')" class="px-3 py-1 border border-white/20 rounded hover:bg-mint/10">中文</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'

const { t, locale } = useI18n()

const changeLocale = (lang: string) => {
  locale.value = lang
}
</script>