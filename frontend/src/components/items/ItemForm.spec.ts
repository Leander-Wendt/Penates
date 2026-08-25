import { describe, expect, it, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { axe } from 'vitest-axe'

import i18n from '@/locales'
import { useAuthStore } from '@/stores/auth'
import ItemForm from './ItemForm.vue'

const item = {
  inventoryNumber: '000001',
  name: 'Drill',
  description: 'A drill',
  imagePath: '/uploads/000001.png',
  category: 'Tools',
  location: 'Shelf A1',
  amount: 2,
  electricalAppliance: true,
  note: 'note',
  manufacturer: 'Acme',
  serialNumber: 'SN1',
  lastTechnicalInspectionDate: '2026-01-01',
  lastElectricalInspectionDate: '2026-01-02',
  dateOfPurchase: '2025-01-01',
  price: 99,
  resolutionNumber: 'R1',
}

function mountForm(role: 'admin' | 'logistics' | 'student', props: Record<string, unknown> = {}) {
  const authStore = useAuthStore()
  authStore.token = 'jwt'
  authStore.role = role
  return mount(ItemForm, { props, global: { plugins: [i18n] } })
}

beforeEach(() => {
  URL.createObjectURL = vi.fn(() => 'blob:preview')
  setActivePinia(createPinia())
})

describe('ItemForm', () => {
  it('renders the public fields', () => {
    const wrapper = mountForm('admin')
    expect(wrapper.text()).toContain('Name')
    expect(wrapper.text()).toContain('Category')
    expect(wrapper.text()).toContain('Location')
  })

  it('pre-fills fields when editing an existing item', () => {
    const wrapper = mountForm('admin', { item })
    const nameInput = wrapper.get('input').element as HTMLInputElement
    expect(nameInput.value).toBe('Drill')
  })

  it('hides the restricted fields section for a student', () => {
    const wrapper = mountForm('student', { item })
    expect(wrapper.text()).not.toContain('Manufacturer')
    expect(wrapper.text()).not.toContain('Restricted information')
  })

  it('shows the restricted fields section for admin and logistics', () => {
    const wrapperAdmin = mountForm('admin', { item })
    expect(wrapperAdmin.text()).toContain('Manufacturer')

    setActivePinia(createPinia())
    const wrapperLogistics = mountForm('logistics', { item })
    expect(wrapperLogistics.text()).toContain('Manufacturer')
  })

  it('validates required fields before emitting submit', async () => {
    const wrapper = mountForm('admin')
    await wrapper.get('form').trigger('submit')
    expect(wrapper.emitted('submit')).toBeUndefined()
    expect(wrapper.text()).toContain('Required')
  })

  it('emits submit with the full form payload once required fields are valid', async () => {
    const wrapper = mountForm('student')
    const inputs = wrapper.findAll('input')
    await inputs[0].setValue('Drill')
    await inputs[2].setValue('Tools')
    await inputs[3].setValue('Shelf A1')
    await wrapper.get('form').trigger('submit')

    const emitted = wrapper.emitted('submit')?.[0]?.[0] as {
      name: string
      category: string
      location: string
    }
    expect(emitted.name).toBe('Drill')
    expect(emitted.category).toBe('Tools')
    expect(emitted.location).toBe('Shelf A1')
  })

  it('emits imageSelect and imageRemove from the embedded ImageUpload', async () => {
    const wrapper = mountForm('admin', { item })
    const file = new File(['x'], 'photo.png', { type: 'image/png' })
    const fileInput = wrapper.get('input[type="file"]')
    Object.defineProperty(fileInput.element, 'files', { value: [file] })
    await fileInput.trigger('change')
    expect(wrapper.emitted('imageSelect')?.[0]).toEqual([file])

    const removeButton = wrapper.findAll('button').find((b) => b.text() === 'Remove image')
    await removeButton?.trigger('click')
    expect(wrapper.emitted('imageRemove')).toBeTruthy()
  })

  it('emits cancel when the cancel button is clicked', async () => {
    const wrapper = mountForm('admin')
    const buttons = wrapper.findAll('button').filter((b) => b.text() === 'Cancel')
    await buttons[0].trigger('click')
    expect(wrapper.emitted('cancel')).toBeTruthy()
  })

  it('has no detectable accessibility violations', async () => {
    const wrapper = mountForm('admin', { item })
    const results = await axe(wrapper.element, { rules: { region: { enabled: false } } })
    expect(results).toHaveNoViolations()
  })
})
