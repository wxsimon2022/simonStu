import { createRouter, createWebHistory } from 'vue-router'
import Login from '../views/Login.vue'
import SystemAdmin from '../views/SystemAdmin.vue'
import SystemRole from '../views/SystemRole.vue'
import SystemPermission from '../views/SystemPermission.vue'

const routes = [
  { path: '/login', name: 'Login', component: Login },
  { path: '/', name: 'Dashboard', component: () => import('../views/Dashboard.vue'), meta: { requiresAuth: true } },
  { path: '/admin', name: 'SystemAdmin', component: SystemAdmin, meta: { requiresAuth: true } },
  { path: '/role', name: 'SystemRole', component: SystemRole, meta: { requiresAuth: true } },
  { path: '/permission', name: 'SystemPermission', component: SystemPermission, meta: { requiresAuth: true } },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('admin_token')
  if (to.meta.requiresAuth && !token) next('/login')
  else if (to.path === '/login' && token) next('/')
  else next()
})

export default router
