import { useToast as useVueToastification } from 'vue-toastification'
import type { ToastOptions } from 'vue-toastification'

export function useToast() {
  const toast = useVueToastification()

  const defaultOptions: ToastOptions = {
    timeout: 3000,
    position: 'top-right',
    closeOnClick: true,
    pauseOnHover: true,
    draggable: true,
  }

  return {
    success: (message: string, options?: ToastOptions) => {
      toast.success(message, { ...defaultOptions, ...options })
    },
    error: (message: string, options?: ToastOptions) => {
      toast.error(message, { ...defaultOptions, timeout: 5000, ...options })
    },
    warning: (message: string, options?: ToastOptions) => {
      toast.warning(message, { ...defaultOptions, ...options })
    },
    info: (message: string, options?: ToastOptions) => {
      toast.info(message, { ...defaultOptions, ...options })
    },
    showToast: (message: string, options?: ToastOptions) => {
      toast(message, { ...defaultOptions, ...options })
    },
  }
}