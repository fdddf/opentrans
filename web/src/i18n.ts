import { createI18n } from 'vue-i18n'

// Load locale messages synchronously
function loadLocaleMessages() {
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

// Load messages synchronously before creating i18n instance
const messages = loadLocaleMessages()

// Create i18n instance with pre-loaded messages
export const i18n = createI18n({
  legacy: false,
  locale: 'en',
  fallbackLocale: 'en',
  messages
})

export type MessageSchema = typeof import('./locales/en.json')
