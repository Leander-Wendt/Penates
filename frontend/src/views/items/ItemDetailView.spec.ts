import { describe, expect, it, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createRouter, createWebHistory } from 'vue-router'
import { createPinia, setActivePinia } from 'pinia'
import { axe } from 'vitest-axe'

import i18n from '@/locales'
import * as itemsApi from '@/api/items'
import { useAuthStore } from '@/stores/auth'
import ItemDetailView from './ItemDetailView.vue'

vi.mock('@/api/items')

const fullItem = {
  inventoryNumber: '000001',
  name: 'Drill',
  description: 'A drill',
  imagePath: '/uploads/000001.png',
  category: 'Tools',
  location: 'Shelf A1',
  amount: 1,
  electricalAppliance: true,
  note: '',
  manufacturer: 'Acme',
  serialNumber: 'SN1',
  lastTechnicalInspectionDate: '2026-01-01',
  lastElectricalInspectionDate: '2026-01-02',
  dateOfPurchase: '2025-01-01',
  price: 99,
  resolutionNumber: 'R1',
}

const studentItem = {
  inventoryNumber: '000001',
  name: 'Drill',
  description: 'A drill',
  imagePath: '/uploads/000001.png',
  category: 'Tools',
  location: 'Shelf A1',
  amount: 1,
  electricalAppliance: true,
  note: '',
}

function makeRouter() {
  return createRouter({
    history: createWebHistory(),
    routes: [
      { path: '/items/:inventoryNumber', component: ItemDetailView, props: true },
      { path: '/items/:inventoryNumber/edit', component: { template: '<div />' } },
    ],
  })
}

describe('ItemDetailView', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  it('renders the item once fetched', async () => {
    vi.mocked(itemsApi.listItems).mockResolvedValue([fullItem])
    const authStore = useAuthStore()
    authStore.token = 'jwt'
    authStore.role = 'admin'
    const router = makeRouter()
    await router.push('/items/000001')
    await router.isReady()
    const wrapper = mount(ItemDetailView, {
      props: { inventoryNumber: '000001' },
      global: { plugins: [i18n, router] },
    })
    await flushPromises()

    expect(wrapper.text()).toContain('Drill')
    expect(wrapper.text()).toContain('Shelf A1')
  })

  it('shows the restricted section for admin/logistics', async () => {
    vi.mocked(itemsApi.listItems).mockResolvedValue([fullItem])
    const authStore = useAuthStore()
    authStore.token = 'jwt'
    authStore.role = 'admin'
    const router = makeRouter()
    await router.push('/items/000001')
    await router.isReady()
    const wrapper = mount(ItemDetailView, {
      props: { inventoryNumber: '000001' },
      global: { plugins: [i18n, router] },
    })
    await flushPromises()

    expect(wrapper.text()).toContain('Restricted information')
    expect(wrapper.text()).toContain('Acme')
  })

  it('hides the restricted section for a student', async () => {
    vi.mocked(itemsApi.listItems).mockResolvedValue([studentItem])
    const authStore = useAuthStore()
    authStore.token = 'jwt'
    authStore.role = 'student'
    const router = makeRouter()
    await router.push('/items/000001')
    await router.isReady()
    const wrapper = mount(ItemDetailView, {
      props: { inventoryNumber: '000001' },
      global: { plugins: [i18n, router] },
    })
    await flushPromises()

    expect(wrapper.text()).not.toContain('Restricted information')
  })

  it('shows an edit button only for admin/logistics', async () => {
    vi.mocked(itemsApi.listItems).mockResolvedValue([studentItem])
    const authStore = useAuthStore()
    authStore.token = 'jwt'
    authStore.role = 'student'
    const router = makeRouter()
    await router.push('/items/000001')
    await router.isReady()
    const wrapper = mount(ItemDetailView, {
      props: { inventoryNumber: '000001' },
      global: { plugins: [i18n, router] },
    })
    await flushPromises()

    expect(wrapper.findAll('button').some((b) => b.text() === 'Edit')).toBe(false)
  })

  it('has no detectable accessibility violations', async () => {
    vi.mocked(itemsApi.listItems).mockResolvedValue([fullItem])
    const authStore = useAuthStore()
    authStore.token = 'jwt'
    authStore.role = 'admin'
    const router = makeRouter()
    await router.push('/items/000001')
    await router.isReady()
    const wrapper = mount(ItemDetailView, {
      props: { inventoryNumber: '000001' },
      global: { plugins: [i18n, router] },
    })
    await flushPromises()
    const results = await axe(wrapper.element, { rules: { region: { enabled: false } } })
    expect(results).toHaveNoViolations()
  })
})
