import { ref, watch } from 'vue'

const theme = ref<'light' | 'dark'>('dark')

// Initialize theme from localStorage or system preference
function initializeTheme() {
  const savedTheme = localStorage.getItem('theme')
  const systemPrefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches
  
  if (savedTheme) {
    theme.value = savedTheme as 'light' | 'dark'
  } else if (window.matchMedia('(prefers-color-scheme: light)').matches) {
    theme.value = 'light'
  } else {
    theme.value = 'dark'
  }
  
  updateDocumentTheme()
}

// Update the document's theme attribute
function updateDocumentTheme() {
  document.documentElement.setAttribute('data-theme', theme.value)
  localStorage.setItem('theme', theme.value)
}

// Toggle between light and dark themes
function toggleTheme() {
  theme.value = theme.value === 'light' ? 'dark' : 'light'
  updateDocumentTheme()
}

// Set a specific theme
function setTheme(newTheme: 'light' | 'dark') {
  theme.value = newTheme
  updateDocumentTheme()
}

// Watch for system theme changes
function watchSystemTheme() {
  const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
  mediaQuery.addEventListener('change', (e) => {
    if (!localStorage.getItem('theme')) {
      theme.value = e.matches ? 'dark' : 'light'
      updateDocumentTheme()
    }
  })
}

// Initialize on module load
initializeTheme()
watchSystemTheme()

// Watch for theme changes
watch(theme, updateDocumentTheme)

export { theme, toggleTheme, setTheme }
