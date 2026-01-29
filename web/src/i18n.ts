import { createI18n } from 'vue-i18n'

// Load locale messages
async function loadLocaleMessages() {
  const locales = import.meta.glob('/src/locales/*.json', { eager: true })
  const messages: Record<string, any> = {}

  for (const [path, module] of Object.entries(locales)) {
    const locale = path.match(/([\w-]+)\.json$/)?.[1]
    if (locale) {
      messages[locale] = (module as { default: any }).default
    }
  }

  return messages
}

// Create i18n instance with dynamic loading
export const i18n = createI18n({
  legacy: false,
  locale: 'en',
  fallbackLocale: 'en',
  messages: {}
})

// Load and set locale messages
loadLocaleMessages().then(messages => {
  // Set the loaded messages
  Object.entries(messages).forEach(([locale, messagesForLocale]) => {
    i18n.global.setLocaleMessage(locale, messagesForLocale)
  })
})

export type MessageSchema = typeof import('./locales/en.json')
