import { describe, expect, it, vi, beforeEach, afterEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createRouter, createWebHistory } from 'vue-router'
import { createPinia, setActivePinia } from 'pinia'
import { axe } from 'vitest-axe'

import i18n from '@/locales'
import * as itemsApi from '@/api/items'
import { useAuthStore } from '@/stores/auth'
import ItemListView from './ItemListView.vue'

vi.mock('@/api/items')

const items = [
  {
    inventoryNumber: '000001',
    name: 'Drill',
    description: '',
    imagePath: null,
    category: 'Tools',
    location: 'A1',
    amount: 1,
    electricalAppliance: true,
    note: '',
  },
]

function makeRouter() {
  return createRouter({
    history: createWebHistory(),
    routes: [
      { path: '/items', component: ItemListView },
      { path: '/items/new', component: { template: '<div />' } },
      { path: '/items/:inventoryNumber', component: { template: '<div />' } },
    ],
  })
}

describe('ItemListView', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
    vi.mocked(itemsApi.listItems).mockResolvedValue(items)
  })

  afterEach(() => {
    vi.useRealTimers()
  })

  it('renders item cards for the fetched items', async () => {
    const authStore = useAuthStore()
    authStore.token = 'jwt'
    authStore.role = 'student'
    const router = makeRouter()
    await router.push('/items')
    await router.isReady()
    const wrapper = mount(ItemListView, { global: { plugins: [i18n, router] } })
    await flushPromises()

    expect(wrapper.text()).toContain('Drill')
  })

  it('shows the create button for admin/logistics but not for a student', async () => {
    const authStore = useAuthStore()
    authStore.token = 'jwt'
    authStore.role = 'student'
    const router = makeRouter()
    await router.push('/items')
    await router.isReady()
    const studentWrapper = mount(ItemListView, { global: { plugins: [i18n, router] } })
    await flushPromises()
    expect(studentWrapper.text()).not.toContain('New item')

    authStore.role = 'admin'
    const adminWrapper = mount(ItemListView, { global: { plugins: [i18n, router] } })
    await flushPromises()
    expect(adminWrapper.text()).toContain('New item')
  })

  it('debounces the search query and refetches with the new term', async () => {
    vi.useFakeTimers()
    const authStore = useAuthStore()
    authStore.token = 'jwt'
    authStore.role = 'student'
    const router = makeRouter()
    await router.push('/items')
    await router.isReady()
    const wrapper = mount(ItemListView, { global: { plugins: [i18n, router] } })
    await flushPromises()
    vi.mocked(itemsApi.listItems).mockClear()

    await wrapper.get('input[type="search"]').setValue('drill')
    vi.advanceTimersByTime(300)
    await flushPromises()

    expect(itemsApi.listItems).toHaveBeenCalledWith({ search: 'drill', category: '' })
  })

  it('refetches immediately when the category filter changes', async () => {
    const authStore = useAuthStore()
    authStore.token = 'jwt'
    authStore.role = 'student'
    const router = makeRouter()
    await router.push('/items')
    await router.isReady()
    const wrapper = mount(ItemListView, { global: { plugins: [i18n, router] } })
    await flushPromises()
    vi.mocked(itemsApi.listItems).mockClear()

    await wrapper.get('select').setValue('Tools')
    await flushPromises()

    expect(itemsApi.listItems).toHaveBeenCalledWith({ search: '', category: 'Tools' })
  })

  it('has no detectable accessibility violations', async () => {
    const authStore = useAuthStore()
    authStore.token = 'jwt'
    authStore.role = 'admin'
    const router = makeRouter()
    await router.push('/items')
    await router.isReady()
    const wrapper = mount(ItemListView, { global: { plugins: [i18n, router] } })
    await flushPromises()
    const results = await axe(wrapper.element, { rules: { region: { enabled: false } } })
    expect(results).toHaveNoViolations()
  })
})
