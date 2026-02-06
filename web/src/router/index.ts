import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
import AppShell from '../components/AppShell.vue'
import Login from '../views/Login.vue'
import Register from '../views/Register.vue'
import Dashboard from '../views/Dashboard.vue'
import Apps from '../views/Apps.vue'
import AppWorkspace from '../views/AppWorkspace.vue'
import AppLocalizations from '../views/AppLocalizations.vue'
import Projects from '../views/Projects.vue'
import Translator from '../views/Translator.vue'
import Users from '../views/Users.vue'
import Subscriptions from '../views/Subscriptions.vue'
import Languages from '../views/Languages.vue'
import AppleConnectConfig from '../views/AppleConnectConfig.vue'
import Profile from '../views/Profile.vue'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: Login,
    meta: { requiresAuth: false }
  },
  {
    path: '/register',
    name: 'Register',
    component: Register,
    meta: { requiresAuth: false }
  },
  {
    path: '/',
    component: AppShell,
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        redirect: '/dashboard'
      },
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: Dashboard
      },
      {
        path: 'apps',
        name: 'Apps',
        component: Apps
      },
      {
        path: 'projects',
        name: 'Projects',
        component: Projects
      },
      {
        path: 'translator',
        name: 'Translator',
        component: Translator
      },
      {
        path: 'apps/:id',
        name: 'AppWorkspace',
        component: AppWorkspace,
        props: true
      },
      {
        path: 'apps/:id/localizations',
        name: 'AppLocalizations',
        component: AppLocalizations,
        props: true
      },
      {
        path: 'users',
        name: 'Users',
        component: Users
      },
      {
        path: 'subscriptions',
        name: 'Subscriptions',
        component: Subscriptions
      },
      {
        path: 'languages',
        name: 'Languages',
        component: Languages
      },
      {
        path: 'apple-connect-config',
        name: 'AppleConnectConfig',
        component: AppleConnectConfig
      },
      {
        path: 'profile',
        name: 'Profile',
        component: Profile
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// Global navigation guard for authentication
router.beforeEach((to, _from, next) => {
  const token = localStorage.getItem('token')
  const isAuthenticated = !!token

  // Check if route requires authentication
  const requiresAuth = to.meta.requiresAuth !== false // Default to true

  if (requiresAuth && !isAuthenticated) {
    // Redirect to login page if not authenticated
    next('/login')
  } else if (!requiresAuth && isAuthenticated && (to.path === '/login' || to.path === '/register')) {
    // Redirect to dashboard if already authenticated and trying to access login/register
    next('/dashboard')
  } else {
    // Proceed to route
    next()
  }
})

export default router
