import { describe, expect, it } from 'vitest'
import { mount } from '@vue/test-utils'

import i18n from '@/locales'
import OrganisationForm from './OrganisationForm.vue'

describe('OrganisationForm', () => {
  it('starts empty when no organisation is given', () => {
    const wrapper = mount(OrganisationForm, { global: { plugins: [i18n] } })
    const nameInput = wrapper.get('input')
    expect((nameInput.element as HTMLInputElement).value).toBe('')
  })

  it('pre-fills fields when editing an existing organisation', () => {
    const wrapper = mount(OrganisationForm, {
      props: { organisation: { id: '1', name: 'Org', description: 'desc' } },
      global: { plugins: [i18n] },
    })
    const inputs = wrapper.findAll('input')
    expect((inputs[0].element as HTMLInputElement).value).toBe('Org')
    expect((inputs[1].element as HTMLInputElement).value).toBe('desc')
  })

  it('emits submit with the form values on valid submission', async () => {
    const wrapper = mount(OrganisationForm, { global: { plugins: [i18n] } })
    const inputs = wrapper.findAll('input')
    await inputs[0].setValue('New org')
    await inputs[1].setValue('New description')
    await wrapper.get('form').trigger('submit')

    expect(wrapper.emitted('submit')?.[0]).toEqual([
      { name: 'New org', description: 'New description' },
    ])
  })

  it('shows a validation error and does not submit when the name is blank', async () => {
    const wrapper = mount(OrganisationForm, { global: { plugins: [i18n] } })
    await wrapper.get('form').trigger('submit')

    expect(wrapper.emitted('submit')).toBeUndefined()
    expect(wrapper.text()).toContain('Required')
  })

  it('emits cancel when the cancel button is clicked', async () => {
    const wrapper = mount(OrganisationForm, { global: { plugins: [i18n] } })
    const buttons = wrapper.findAll('button')
    await buttons[buttons.length - 1].trigger('click')

    expect(wrapper.emitted('cancel')).toBeTruthy()
  })
})
