import { describe, expect, it, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createRouter, createWebHistory } from 'vue-router'
import { createPinia, setActivePinia } from 'pinia'
import { axe } from 'vitest-axe'

import i18n from '@/locales'
import * as itemsApi from '@/api/items'
import { useAuthStore } from '@/stores/auth'
import ItemFormView from './ItemFormView.vue'

vi.mock('@/api/items')

const existingItem = {
  inventoryNumber: '000001',
  name: 'Drill',
  description: '',
  imagePath: null,
  category: 'Tools',
  location: 'A1',
  amount: 1,
  electricalAppliance: true,
  note: '',
}

function makeRouter() {
  return createRouter({
    history: createWebHistory(),
    routes: [
      { path: '/items', component: { template: '<div />' } },
      { path: '/items/new', component: ItemFormView },
      { path: '/items/:inventoryNumber', component: { template: '<div />' } },
      { path: '/items/:inventoryNumber/edit', component: ItemFormView, props: true },
    ],
  })
}

describe('ItemFormView', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
    const authStore = useAuthStore()
    authStore.token = 'jwt'
    authStore.role = 'admin'
  })

  it('creates an item and redirects to its detail page', async () => {
    vi.mocked(itemsApi.createItem).mockResolvedValue(existingItem)
    const router = makeRouter()
    await router.push('/items/new')
    await router.isReady()
    const wrapper = mount(ItemFormView, { global: { plugins: [i18n, router] } })
    await flushPromises()

    const inputs = wrapper.findAll('input')
    await inputs[0].setValue('Drill')
    await inputs[2].setValue('Tools')
    await inputs[3].setValue('A1')
    await wrapper.get('form').trigger('submit')
    await flushPromises()

    expect(itemsApi.createItem).toHaveBeenCalled()
    expect(router.currentRoute.value.path).toBe('/items/000001')
  })

  it('loads and pre-fills an existing item in edit mode', async () => {
    vi.mocked(itemsApi.listItems).mockResolvedValue([existingItem])
    const router = makeRouter()
    await router.push('/items/000001/edit')
    await router.isReady()
    const wrapper = mount(ItemFormView, { global: { plugins: [i18n, router] } })
    await flushPromises()

    expect(wrapper.text()).toContain('Edit item')
    const nameInput = wrapper.get('input').element as HTMLInputElement
    expect(nameInput.value).toBe('Drill')
  })

  it('uploads a new image immediately when editing an existing item', async () => {
    vi.mocked(itemsApi.listItems).mockResolvedValue([existingItem])
    vi.mocked(itemsApi.uploadItemImage).mockResolvedValue({
      ...existingItem,
      imagePath: '/uploads/x.png',
    })
    URL.createObjectURL = vi.fn(() => 'blob:preview')
    const router = makeRouter()
    await router.push('/items/000001/edit')
    await router.isReady()
    const wrapper = mount(ItemFormView, { global: { plugins: [i18n, router] } })
    await flushPromises()

    const file = new File(['x'], 'photo.png', { type: 'image/png' })
    const fileInput = wrapper.get('input[type="file"]')
    Object.defineProperty(fileInput.element, 'files', { value: [file] })
    await fileInput.trigger('change')
    await flushPromises()

    expect(itemsApi.uploadItemImage).toHaveBeenCalledWith('000001', file)
  })

  it('has no detectable accessibility violations', async () => {
    const router = makeRouter()
    await router.push('/items/new')
    await router.isReady()
    const wrapper = mount(ItemFormView, { global: { plugins: [i18n, router] } })
    await flushPromises()
    const results = await axe(wrapper.element, { rules: { region: { enabled: false } } })
    expect(results).toHaveNoViolations()
  })
})
