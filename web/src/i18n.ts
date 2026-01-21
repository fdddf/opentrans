import { createI18n } from 'vue-i18n'

const messages = {
  en: {
    nav: {
      dashboard: 'Dashboard',
      apps: 'Apps',
      applocalizations: 'App Localizations',
      languages: 'Languages',
      users: 'Users',
      subscriptions: 'Subscriptions'
    },
    common: {
      brand: 'XCStrings Translator',
      login: 'Sign in',
      register: 'Create account',
      username: 'Username',
      email: 'Email',
      password: 'Password',
      confirmPassword: 'Confirm Password',
      logout: 'Logout',
      dashboard: 'Dashboard'
    },
    dashboard: {
      title: 'Multi-tenant localization overview',
      subtitle: 'Track users, subscriptions, apps and translation progress.',
      activeUsers: 'Active users',
      activeSubscriptions: 'Active subscriptions',
      apps: 'Managed apps',
      pendingStrings: 'Pending strings',
      recentActivity: 'Recent activity',
      opsLog: 'Ops log',
      translationProgress: 'Translation progress',
      viewDetails: 'View details'
    },
    auth: {
      loginTitle: 'Sign in to your account',
      loginCta: 'Sign in',
      registerPrompt: 'Create a new account',
      registerTitle: 'Create your account',
      haveAccount: 'Already have an account?',
      noAccount: "Don't have an account?",
      forgot: 'Forgot your password?'
    },
    workspace: {
      title: 'Localization Workspace',
      upload: 'Upload xcstrings',
      addLang: 'Add target language',
      export: 'Export',
      translateAll: 'Translate all',
      translateLang: 'Translate language',
      editEntry: 'Edit entry',
      search: 'Search key or comment'
    },
    subscriptions: {
      title: 'Subscription management',
      subtitle: 'Users can access app and translation services after subscribing.',
      new: 'New subscription',
      sync: 'Sync billing'
    },
    languages: {
      title: 'Language library',
      subtitle: 'Maintain global languages, aliases, and direction (LTR/RTL).',
      readonly: 'Fixed dataset; editing disabled.'
    },
    users: {
      title: 'User management',
      subtitle: 'Admins and regular users with isolated data.',
      inviteAdmin: 'Invite admin',
      createUser: 'Create user'
    },
    apps: {
      title: 'App management (sync/manual)',
      subtitle: 'Sync existing apps or add manually. Only manual apps can be deleted.',
      sync: 'Sync apps',
      manual: 'Add manually'
    }
  },
  zh: {
    nav: {
      dashboard: '总览',
      apps: '应用管理',
      languages: '多语言库',
      users: '用户管理',
      subscriptions: '订阅管理'
    },
    common: {
      brand: 'XCStrings 翻译',
      login: '登录',
      register: '注册账号',
      username: '用户名',
      email: '邮箱',
      password: '密码',
      confirmPassword: '确认密码',
      logout: '退出登录',
      dashboard: '总览'
    },
    dashboard: {
      title: '多租户多语言运营总览',
      subtitle: '查看用户、订阅、应用与翻译进度汇总。',
      activeUsers: '活跃用户',
      activeSubscriptions: '订阅中',
      apps: '管理的应用',
      pendingStrings: '待翻译字符串',
      recentActivity: '最近活动',
      opsLog: '运维日志',
      translationProgress: '翻译进度',
      viewDetails: '查看详情'
    },
    auth: {
      loginTitle: '登录您的账号',
      loginCta: '登录',
      registerPrompt: '创建新账号',
      registerTitle: '创建您的账号',
      haveAccount: '已有账号？',
      noAccount: '还没有账号？',
      forgot: '忘记密码？'
    },
    workspace: {
      title: '本地化工作区',
      upload: '上传 xcstrings',
      addLang: '添加目标语言',
      export: '导出',
      translateAll: '全部翻译',
      translateLang: '逐语言翻译',
      editEntry: '单条编辑',
      search: '搜索 key 或备注'
    },
    subscriptions: {
      title: '订阅管理',
      subtitle: '用户购买/续订后可使用 App 管理与翻译服务。',
      new: '新建订阅',
      sync: '同步账单'
    },
    languages: {
      title: '多语言库',
      subtitle: '维护全局语言列表、别名、方向（LTR/RTL）和可用性。',
      readonly: '固定数据，禁止编辑。'
    },
    users: {
      title: '用户管理',
      subtitle: '管理员、普通用户，数据隔离，权限控制。',
      inviteAdmin: '邀请管理员',
      createUser: '创建用户'
    },
    apps: {
      title: '应用管理（同步/手动添加）',
      subtitle: '同步现有应用或手动添加，删除仅限手动添加的应用。',
      sync: '同步应用',
      manual: '手动添加'
    }
  }
}

export const i18n = createI18n({
  legacy: false,
  locale: 'en',
  fallbackLocale: 'en',
  messages
})

export type MessageSchema = typeof messages.en
