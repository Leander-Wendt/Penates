import {
  createRouter,
  createWebHistory,
  type RouteLocationNormalized,
  type NavigationGuardNext,
} from 'vue-router'

import { useAuth } from '@/composables/useAuth'
import type { Role } from '@/types'

declare module 'vue-router' {
  interface RouteMeta {
    roles?: Role[]
    public?: boolean
  }
}

const routes = [
  {
    path: '/',
    redirect: '/items',
  },
  {
    path: '/login',
    name: 'login',
    component: () => import('@/views/LoginView.vue'),
    meta: { public: true },
  },
  {
    path: '/forbidden',
    name: 'forbidden',
    component: () => import('@/views/ForbiddenView.vue'),
    meta: { public: true },
  },
  {
    path: '/items',
    name: 'items-list',
    component: () => import('@/views/items/ItemListView.vue'),
  },
  {
    path: '/items/new',
    name: 'items-create',
    component: () => import('@/views/items/ItemFormView.vue'),
    meta: { roles: ['admin', 'logistics'] as Role[] },
  },
  {
    path: '/items/:inventoryNumber',
    name: 'items-detail',
    component: () => import('@/views/items/ItemDetailView.vue'),
    props: true,
  },
  {
    path: '/items/:inventoryNumber/edit',
    name: 'items-edit',
    component: () => import('@/views/items/ItemFormView.vue'),
    props: true,
    meta: { roles: ['admin', 'logistics'] as Role[] },
  },
  {
    path: '/organisations',
    name: 'organisations-list',
    component: () => import('@/views/organisations/OrganisationListView.vue'),
    meta: { roles: ['admin', 'logistics'] as Role[] },
  },
  {
    path: '/users',
    name: 'users-list',
    component: () => import('@/views/users/UserListView.vue'),
    meta: { roles: ['admin', 'logistics'] as Role[] },
  },
  {
    path: '/users/new',
    name: 'users-create',
    component: () => import('@/views/users/UserFormView.vue'),
    meta: { roles: ['admin', 'logistics'] as Role[] },
  },
  {
    path: '/users/:id/edit',
    name: 'users-edit',
    component: () => import('@/views/users/UserFormView.vue'),
    props: true,
    meta: { roles: ['admin', 'logistics'] as Role[] },
  },
  {
    path: '/loan-requests',
    name: 'loan-requests-list',
    component: () => import('@/views/loanRequests/LoanRequestListView.vue'),
  },
  {
    path: '/loan-requests/new',
    name: 'loan-requests-create',
    component: () => import('@/views/loanRequests/LoanRequestFormView.vue'),
  },
  {
    path: '/loan-requests/:id',
    name: 'loan-requests-detail',
    component: () => import('@/views/loanRequests/LoanRequestDetailView.vue'),
    props: true,
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'not-found',
    component: () => import('@/views/NotFoundView.vue'),
    meta: { public: true },
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

/**
 * Global navigation guard: redirects unauthenticated visitors of a
 * non-public route to `/login`, and redirects authenticated visitors who
 * lack a required role to `/forbidden`. A 403 from the API leaves the
 * session intact (handled separately in the API client); this guard only
 * decides route access before a request is even made.
 */
export function authGuard(
  to: RouteLocationNormalized,
  _from: RouteLocationNormalized,
  next: NavigationGuardNext,
): void {
  const { isAuthenticated, hasRole } = useAuth()

  if (!to.meta.public && !isAuthenticated.value) {
    next({ name: 'login', query: { redirect: to.fullPath } })
    return
  }

  if (to.meta.roles && !hasRole(to.meta.roles)) {
    next({ name: 'forbidden' })
    return
  }

  next()
}

router.beforeEach(authGuard)

export default router
