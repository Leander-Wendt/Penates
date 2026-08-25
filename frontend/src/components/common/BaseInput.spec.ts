import { describe, expect, it } from 'vitest'
import { mount } from '@vue/test-utils'
import { axe } from 'vitest-axe'

import BaseInput from './BaseInput.vue'

describe('BaseInput', () => {
  it('renders a label associated with the input via a matching id/for pair', () => {
    const wrapper = mount(BaseInput, { props: { modelValue: '', label: 'Email' } })
    const label = wrapper.get('label')
    const input = wrapper.get('input')
    expect(label.attributes('for')).toBe(input.attributes('id'))
    expect(label.text()).toContain('Email')
  })

  it('emits update:modelValue on input', async () => {
    const wrapper = mount(BaseInput, { props: { modelValue: '', label: 'Email' } })
    await wrapper.get('input').setValue('a@b.com')
    expect(wrapper.emitted('update:modelValue')?.[0]).toEqual(['a@b.com'])
  })

  it('links an error message to the input via aria-describedby and renders it in a live region', () => {
    const wrapper = mount(BaseInput, {
      props: { modelValue: '', label: 'Email', error: 'Email is required' },
    })
    const input = wrapper.get('input')
    const describedBy = input.attributes('aria-describedby')
    expect(describedBy).toBeDefined()
    const errorEl = wrapper.get(`#${describedBy}`)
    expect(errorEl.text()).toBe('Email is required')
    expect(errorEl.attributes('aria-live')).toBe('assertive')
    expect(input.attributes('aria-invalid')).toBe('true')
  })

  it('has no aria-describedby when there is no error', () => {
    const wrapper = mount(BaseInput, { props: { modelValue: '', label: 'Email' } })
    expect(wrapper.get('input').attributes('aria-describedby')).toBeUndefined()
  })

  it('has no detectable accessibility violations', async () => {
    const wrapper = mount(BaseInput, {
      props: { modelValue: '', label: 'Email', error: 'Email is required' },
    })
    const results = await axe(wrapper.element, { rules: { region: { enabled: false } } })
    expect(results).toHaveNoViolations()
  })
})
