import { describe, expect, it, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { axe } from 'vitest-axe'

import BaseButton from './BaseButton.vue'

describe('BaseButton', () => {
  it('renders a native button with the default primary variant', () => {
    const wrapper = mount(BaseButton, { slots: { default: 'Save' } })
    const button = wrapper.get('button')
    expect(button.text()).toBe('Save')
    expect(button.classes()).toContain('bg-brut-yellow')
  })

  it('applies the requested variant classes', () => {
    const wrapper = mount(BaseButton, {
      props: { variant: 'danger' },
      slots: { default: 'Delete' },
    })
    expect(wrapper.get('button').classes()).toContain('bg-brut-red')
  })

  it('forwards click events to a handler passed via attrs', async () => {
    const onClick = vi.fn()
    const wrapper = mount(BaseButton, { attrs: { onClick }, slots: { default: 'Click me' } })
    await wrapper.get('button').trigger('click')
    expect(onClick).toHaveBeenCalled()
  })

  it('respects the disabled prop', () => {
    const wrapper = mount(BaseButton, { props: { disabled: true }, slots: { default: 'Save' } })
    expect(wrapper.get('button').attributes('disabled')).toBeDefined()
  })

  it('exposes an aria-label when provided, for icon-only usage', () => {
    const wrapper = mount(BaseButton, { props: { ariaLabel: 'Close' } })
    expect(wrapper.get('button').attributes('aria-label')).toBe('Close')
  })

  it('has no detectable accessibility violations', async () => {
    const wrapper = mount(BaseButton, { slots: { default: 'Save' } })
    const results = await axe(wrapper.element)
    expect(results).toHaveNoViolations()
  })
})
