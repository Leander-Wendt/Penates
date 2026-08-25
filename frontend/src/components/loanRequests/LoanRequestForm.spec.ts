import { describe, expect, it } from 'vitest'
import { mount } from '@vue/test-utils'
import { axe } from 'vitest-axe'

import i18n from '@/locales'
import LoanRequestForm from './LoanRequestForm.vue'

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

describe('LoanRequestForm', () => {
  it('blocks submission and shows errors when nothing is filled in', async () => {
    const wrapper = mount(LoanRequestForm, { props: { items }, global: { plugins: [i18n] } })
    await wrapper.get('form').trigger('submit')
    expect(wrapper.emitted('submit')).toBeUndefined()
    expect(wrapper.text()).toContain('Required')
  })

  it('emits submit with the selected items and field values once valid', async () => {
    const wrapper = mount(LoanRequestForm, { props: { items }, global: { plugins: [i18n] } })
    await wrapper.get('input[type="checkbox"]').setValue(true)
    const dateInputs = wrapper.findAll('input[type="date"]')
    await dateInputs[0].setValue('2026-09-01')
    await dateInputs[1].setValue('2026-09-10')
    await wrapper.get('input[type="text"]').setValue('Room 1')
    await wrapper.get('form').trigger('submit')

    expect(wrapper.emitted('submit')?.[0]).toEqual([
      {
        itemInventoryNumbers: ['000001'],
        dateOfLending: '2026-09-01',
        dateOfReturn: '2026-09-10',
        locationOfItems: 'Room 1',
      },
    ])
  })

  it('emits cancel when the cancel button is clicked', async () => {
    const wrapper = mount(LoanRequestForm, { props: { items }, global: { plugins: [i18n] } })
    const buttons = wrapper.findAll('button').filter((b) => b.text() === 'Cancel')
    await buttons[0].trigger('click')
    expect(wrapper.emitted('cancel')).toBeTruthy()
  })

  it('has no detectable accessibility violations', async () => {
    const wrapper = mount(LoanRequestForm, { props: { items }, global: { plugins: [i18n] } })
    const results = await axe(wrapper.element, { rules: { region: { enabled: false } } })
    expect(results).toHaveNoViolations()
  })
})
