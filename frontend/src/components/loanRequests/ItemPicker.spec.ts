import { describe, expect, it } from 'vitest'
import { mount } from '@vue/test-utils'
import { axe } from 'vitest-axe'

import i18n from '@/locales'
import ItemPicker from './ItemPicker.vue'

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
  {
    inventoryNumber: '000002',
    name: 'Ladder',
    description: '',
    imagePath: null,
    category: 'Tools',
    location: 'A2',
    amount: 1,
    electricalAppliance: false,
    note: '',
  },
]

describe('ItemPicker', () => {
  it('renders a checkbox per item, each with a text label', () => {
    const wrapper = mount(ItemPicker, {
      props: { items, modelValue: [] },
      global: { plugins: [i18n] },
    })
    const checkboxes = wrapper.findAll('input[type="checkbox"]')
    expect(checkboxes).toHaveLength(2)
    expect(wrapper.text()).toContain('Drill')
    expect(wrapper.text()).toContain('Ladder')
  })

  it('reflects the selected items as checked', () => {
    const wrapper = mount(ItemPicker, {
      props: { items, modelValue: ['000001'] },
      global: { plugins: [i18n] },
    })
    const checkboxes = wrapper.findAll('input[type="checkbox"]')
    expect((checkboxes[0].element as HTMLInputElement).checked).toBe(true)
    expect((checkboxes[1].element as HTMLInputElement).checked).toBe(false)
  })

  it('emits update:modelValue adding an item when its checkbox is checked', async () => {
    const wrapper = mount(ItemPicker, {
      props: { items, modelValue: [] },
      global: { plugins: [i18n] },
    })
    await wrapper.findAll('input[type="checkbox"]')[0].setValue(true)
    expect(wrapper.emitted('update:modelValue')?.[0]).toEqual([['000001']])
  })

  it('emits update:modelValue removing an item when its checkbox is unchecked', async () => {
    const wrapper = mount(ItemPicker, {
      props: { items, modelValue: ['000001', '000002'] },
      global: { plugins: [i18n] },
    })
    await wrapper.findAll('input[type="checkbox"]')[0].setValue(false)
    expect(wrapper.emitted('update:modelValue')?.[0]).toEqual([['000002']])
  })

  it('is keyboard operable: each checkbox is a native, tabbable control with an associated label', () => {
    const wrapper = mount(ItemPicker, {
      props: { items, modelValue: [] },
      global: { plugins: [i18n] },
    })
    const checkboxes = wrapper.findAll('input[type="checkbox"]')
    const labels = wrapper.findAll('label')
    checkboxes.forEach((checkbox, index) => {
      expect(checkbox.attributes('tabindex')).not.toBe('-1')
      expect(labels[index].attributes('for')).toBe(checkbox.attributes('id'))
    })
  })

  it('has no detectable accessibility violations', async () => {
    const wrapper = mount(ItemPicker, {
      props: { items, modelValue: [] },
      global: { plugins: [i18n] },
    })
    const results = await axe(wrapper.element, { rules: { region: { enabled: false } } })
    expect(results).toHaveNoViolations()
  })
})
