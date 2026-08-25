import { describe, expect, it } from 'vitest'
import { mount } from '@vue/test-utils'
import { axe } from 'vitest-axe'

import i18n from '@/locales'
import LoanRequestStatusBadge from './LoanRequestStatusBadge.vue'

describe('LoanRequestStatusBadge', () => {
  it('renders the status as visible text, not only as color', () => {
    const wrapper = mount(LoanRequestStatusBadge, {
      props: { status: 'approved' },
      global: { plugins: [i18n] },
    })
    expect(wrapper.text()).toContain('Approved')
  })

  it('renders an icon alongside the text for each status', () => {
    const statuses = ['pending', 'approved', 'rejected', 'returned', 'cancelled'] as const
    for (const status of statuses) {
      const wrapper = mount(LoanRequestStatusBadge, {
        props: { status },
        global: { plugins: [i18n] },
      })
      const icon = wrapper.find('span[aria-hidden="true"]')
      expect(icon.exists()).toBe(true)
      expect(icon.text().length).toBeGreaterThan(0)
    }
  })

  it('has no detectable accessibility violations', async () => {
    const wrapper = mount(LoanRequestStatusBadge, {
      props: { status: 'rejected' },
      global: { plugins: [i18n] },
    })
    const results = await axe(wrapper.element, { rules: { region: { enabled: false } } })
    expect(results).toHaveNoViolations()
  })
})
