import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
import AppShell from '../components/AppShell.vue'
import Login from '../views/Login.vue'
import Register from '../views/Register.vue'
import Dashboard from '../views/Dashboard.vue'
import Apps from '../views/Apps.vue'
import AppWorkspace from '../views/AppWorkspace.vue'
import AppLocalizations from '../views/AppLocalizations.vue'
import Users from '../views/Users.vue'
import Subscriptions from '../views/Subscriptions.vue'
import Languages from '../views/Languages.vue'
import AppleConnectConfig from '../views/AppleConnectConfig.vue'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: Login
  },
  {
    path: '/register',
    name: 'Register',
    component: Register
  },
  {
    path: '/',
    component: AppShell,
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
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
