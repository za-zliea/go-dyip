import { createRouter, createWebHashHistory, type RouteRecordRaw } from 'vue-router'
import { useUserStore } from '@/store/user'

// Hash mode: URLs look like `/#/ddns/update`. This requires zero server-side
// cooperation (the Go server only ever serves index.html at `/`), so it works
// flawlessly with the existing atreugo.StaticCustom setup.
const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    {
      path: '/login',
      name: 'Login',
      component: () => import('@/views/Login/index.vue'),
      meta: { titleKey: 'route.login', requiresAuth: false }
    },
    {
      path: '/',
      component: () => import('@/layout/index.vue'),
      redirect: '/ddns/view',
      children: [
        {
          path: '/ddns/view',
          name: 'DdnsView',
          component: () => import('@/views/Ddns/View/index.vue'),
          meta: { titleKey: 'route.ddnsView', icon: 'View', requiresAuth: true }
        },
        {
          path: '/ddns/update',
          name: 'DdnsUpdate',
          component: () => import('@/views/Ddns/Update/index.vue'),
          meta: { titleKey: 'route.ddnsUpdate', icon: 'EditPen', requiresAuth: true }
        }
      ]
    },
    {
      path: '/:pathMatch(.*)*',
      redirect: '/ddns/view'
    }
  ]
})

// Global guard: every auth-required route needs a token. Unauthed users are
// sent to /login with the original target as a `redirect` query.
router.beforeEach((to) => {
  const userStore = useUserStore()
  const requiresAuth = to.meta.requiresAuth !== false

  if (requiresAuth && !userStore.isLogin) {
    return { name: 'Login', query: { redirect: to.fullPath } }
  }
  if (to.name === 'Login' && userStore.isLogin) {
    return { path: '/ddns/view' }
  }
  return true
})

// document.title is set in App.vue, which watches the active locale — the
// title is intentionally fixed to the console title and does not change with
// the active route.

export default router
