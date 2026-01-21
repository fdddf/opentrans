<template>
  <div class="space-y-6">
    <header class="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
      <div>
        <p class="text-xs uppercase tracking-[0.3em] text-slate-500">{{ t('nav.subscription') }}</p>
        <h1 class="text-2xl font-semibold">{{ t('subscription.title') }}</h1>
        <p class="text-sm text-slate-400">{{ t('subscription.subtitle') }}</p>
      </div>
    </header>

    <!-- User Subscription Info -->
    <section class="rounded-2xl border border-white/10 bg-white/5 p-6">
      <h2 class="text-xl font-semibold mb-4">{{ t('subscription.userInfo') }}</h2>
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <div class="rounded-xl border border-white/10 bg-midnight p-4">
          <p class="text-sm text-slate-400">{{ t('subscription.plan') }}</p>
          <p class="text-lg font-semibold">{{ user?.subscriptionType || 'N/A' }}</p>
        </div>
        <div class="rounded-xl border border-white/10 bg-midnight p-4">
          <p class="text-sm text-slate-400">{{ t('subscription.appsUsed') }}</p>
          <p class="text-lg font-semibold">{{ user?.currentAppCount || 0 }} / {{ user?.maxApps || 0 }}</p>
        </div>
        <div class="rounded-xl border border-white/10 bg-midnight p-4">
          <p class="text-sm text-slate-400">{{ t('subscription.translationsUsed') }}</p>
          <p class="text-lg font-semibold">{{ user?.currentUsage || 0 }} / {{ user?.maxTranslations || 0 }}</p>
        </div>
        <div class="rounded-xl border border-white/10 bg-midnight p-4">
          <p class="text-sm text-slate-400">{{ t('subscription.status') }}</p>
          <p class="text-lg font-semibold">{{ user?.isSubscribed ? t('subscription.active') : t('subscription.inactive') }}</p>
        </div>
      </div>
    </section>

    <!-- Subscription Plans -->
    <section class="rounded-2xl border border-white/10 bg-white/5 p-6">
      <h2 class="text-xl font-semibold mb-4">{{ t('subscription.plans') }}</h2>
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
        <!-- Free Plan -->
        <div class="rounded-2xl border border-white/10 bg-midnight p-6">
          <h3 class="text-lg font-semibold mb-2">{{ t('subscription.free') }}</h3>
          <p class="text-3xl font-bold mb-4">Free</p>
          
          <ul class="space-y-2 mb-6">
            <li class="flex items-center">
              <span class="text-emerald-400 mr-2">✓</span>
              <span>{{ t('subscription.features.upToApps', { count: 1 }) }}</span>
            </li>
            <li class="flex items-center">
              <span class="text-emerald-400 mr-2">✓</span>
              <span>{{ t('subscription.features.upToTranslations', { count: 1000 }) }}</span>
            </li>
            <li class="flex items-center">
              <span class="text-emerald-400 mr-2">✓</span>
              <span>{{ t('subscription.features.basicSupport') }}</span>
            </li>
          </ul>
          
          <button 
            v-if="user?.subscriptionType !== 'free'"
            class="w-full rounded-lg bg-slate-700 px-3 py-2 text-sm font-semibold text-slate-200"
            @click="switchToPlan('free')"
            :disabled="user?.subscriptionType === 'free'"
          >
            {{ t('subscription.current') }}
          </button>
          <div v-else class="w-full rounded-lg bg-mint/20 px-3 py-2 text-sm font-semibold text-mint text-center">
            {{ t('subscription.current') }}
          </div>
        </div>
        
        <!-- Basic Plan -->
        <div class="rounded-2xl border border-white/10 bg-midnight p-6 relative">
          <div v-if="user?.subscriptionType === 'basic'" class="absolute top-0 right-0 rounded-bl-2xl rounded-tr-2xl bg-mint px-3 py-1 text-xs font-semibold text-midnight">
            {{ t('subscription.current') }}
          </div>
          <h3 class="text-lg font-semibold mb-2">{{ t('subscription.basic') }}</h3>
          <p class="text-3xl font-bold mb-4">
            $9<span class="text-sm font-normal">/month</span>
          </p>
          
          <ul class="space-y-2 mb-6">
            <li class="flex items-center">
              <span class="text-emerald-400 mr-2">✓</span>
              <span>{{ t('subscription.features.upToApps', { count: 5 }) }}</span>
            </li>
            <li class="flex items-center">
              <span class="text-emerald-400 mr-2">✓</span>
              <span>{{ t('subscription.features.upToTranslations', { count: 10000 }) }}</span>
            </li>
            <li class="flex items-center">
              <span class="text-emerald-400 mr-2">✓</span>
              <span>{{ t('subscription.features.prioritySupport') }}</span>
            </li>
            <li class="flex items-center">
              <span class="text-emerald-400 mr-2">✓</span>
              <span>{{ t('subscription.features.teamAccess') }}</span>
            </li>
          </ul>
          
          <button 
            class="w-full rounded-lg bg-mint px-3 py-2 text-sm font-semibold text-midnight"
            @click="switchToPlan('basic')"
            :disabled="user?.subscriptionType === 'basic'"
          >
            <span v-if="user?.subscriptionType === 'basic'">{{ t('subscription.current') }}</span>
            <span v-else>{{ t('subscription.upgrade') }}</span>
          </button>
        </div>
        
        <!-- Premium Plan -->
        <div class="rounded-2xl border border-white/10 bg-midnight p-6 relative">
          <div v-if="user?.subscriptionType === 'premium'" class="absolute top-0 right-0 rounded-bl-2xl rounded-tr-2xl bg-mint px-3 py-1 text-xs font-semibold text-midnight">
            {{ t('subscription.current') }}
          </div>
          <h3 class="text-lg font-semibold mb-2">{{ t('subscription.premium') }}</h3>
          <p class="text-3xl font-bold mb-4">
            $29<span class="text-sm font-normal">/month</span>
          </p>
          
          <ul class="space-y-2 mb-6">
            <li class="flex items-center">
              <span class="text-emerald-400 mr-2">✓</span>
              <span>{{ t('subscription.features.upToApps', { count: 20 }) }}</span>
            </li>
            <li class="flex items-center">
              <span class="text-emerald-400 mr-2">✓</span>
              <span>{{ t('subscription.features.upToTranslations', { count: 50000 }) }}</span>
            </li>
            <li class="flex items-center">
              <span class="text-emerald-400 mr-2">✓</span>
              <span>{{ t('subscription.features.prioritySupport') }}</span>
            </li>
            <li class="flex items-center">
              <span class="text-emerald-400 mr-2">✓</span>
              <span>{{ t('subscription.features.teamAccess') }}</span>
            </li>
            <li class="flex items-center">
              <span class="text-emerald-400 mr-2">✓</span>
              <span>{{ t('subscription.features.advancedAnalytics') }}</span>
            </li>
          </ul>
          
          <button 
            class="w-full rounded-lg bg-mint px-3 py-2 text-sm font-semibold text-midnight"
            @click="switchToPlan('premium')"
            :disabled="user?.subscriptionType === 'premium'"
          >
            <span v-if="user?.subscriptionType === 'premium'">{{ t('subscription.current') }}</span>
            <span v-else>{{ t('subscription.upgrade') }}</span>
          </button>
        </div>
      </div>
    </section>

    <!-- Usage Information -->
    <section class="rounded-2xl border border-white/10 bg-white/5 p-6">
      <h2 class="text-xl font-semibold mb-4">{{ t('subscription.usage') }}</h2>
      <div class="space-y-4">
        <div>
          <div class="flex justify-between mb-1">
            <span class="text-sm">{{ t('subscription.appsUsed') }}</span>
            <span class="text-sm">{{ user?.currentAppCount || 0 }} / {{ user?.maxApps || 0 }}</span>
          </div>
          <div class="w-full bg-slate-700 rounded-full h-2.5">
            <div 
              class="bg-mint h-2.5 rounded-full" 
              :style="{ width: Math.min(100, ((user?.currentAppCount || 0) / (user?.maxApps || 1)) * 100) + '%' }"
            ></div>
          </div>
        </div>
        
        <div>
          <div class="flex justify-between mb-1">
            <span class="text-sm">{{ t('subscription.translationsUsed') }}</span>
            <span class="text-sm">{{ user?.currentUsage || 0 }} / {{ user?.maxTranslations || 0 }}</span>
          </div>
          <div class="w-full bg-slate-700 rounded-full h-2.5">
            <div 
              class="bg-mint h-2.5 rounded-full" 
              :style="{ width: Math.min(100, ((user?.currentUsage || 0) / (user?.maxTranslations || 1)) * 100) + '%' }"
            ></div>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useApi } from '../composables/useApi'
import type { User } from '../composables/useApi'

const { t } = useI18n()
const { api } = useApi()

const user = ref<User | null>(null)

async function fetchUserSubscription() {
  try {
    const response = await api.getUserSubscription()
    if (response.success) {
      user.value = response.user
    }
  } catch (error) {
    console.error('Failed to fetch user subscription:', error)
  }
}

async function switchToPlan(planType: string) {
  // In a real implementation, this would redirect to Stripe checkout
  // For now, we'll just log the action
  console.log(`Switching to plan: ${planType}`)
  
  // Show a message to the user
  alert(`In a real implementation, this would redirect you to the ${planType} plan checkout`)
}

onMounted(() => {
  fetchUserSubscription()
})
</script>
