import { describe, expect, it } from 'vitest'
import { mount } from '@vue/test-utils'
import { axe } from 'vitest-axe'

import i18n from '@/locales'
import NotFoundView from './NotFoundView.vue'

describe('NotFoundView', () => {
  it('renders the not-found message', () => {
    const wrapper = mount(NotFoundView, { global: { plugins: [i18n] } })
    expect(wrapper.text()).toContain('Page not found')
  })

  it('has no detectable accessibility violations', async () => {
    const wrapper = mount(NotFoundView, { global: { plugins: [i18n] } })
    const results = await axe(wrapper.element, { rules: { region: { enabled: false } } })
    expect(results).toHaveNoViolations()
  })
})
