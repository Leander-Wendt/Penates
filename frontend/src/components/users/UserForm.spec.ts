import { describe, expect, it } from 'vitest'
import { mount } from '@vue/test-utils'
import { axe } from 'vitest-axe'

import i18n from '@/locales'
import type { UserCreateInput } from '@/types'
import UserForm from './UserForm.vue'

const organisations = [
  { id: 'o1', name: 'Org 1', description: '' },
  { id: 'o2', name: 'Org 2', description: '' },
]

describe('UserForm', () => {
  it('shows a password field when creating a new user', () => {
    const wrapper = mount(UserForm, { props: { organisations }, global: { plugins: [i18n] } })
    expect(wrapper.find('input[type="password"]').exists()).toBe(true)
  })

  it('hides the password field when editing an existing user', () => {
    const wrapper = mount(UserForm, {
      props: {
        organisations,
        user: { id: '1', email: 'a@b.com', role: 'student', organisationId: 'o1' },
      },
      global: { plugins: [i18n] },
    })
    expect(wrapper.find('input[type="password"]').exists()).toBe(false)
  })

  it('emits a create payload including the password', async () => {
    const wrapper = mount(UserForm, { props: { organisations }, global: { plugins: [i18n] } })
    await wrapper.get('input[type="email"]').setValue('new@b.com')
    await wrapper.get('input[type="password"]').setValue('secret123')
    await wrapper.get('form').trigger('submit')

    expect(wrapper.emitted('submit')?.[0]).toEqual([
      { email: 'new@b.com', password: 'secret123', role: 'student', organisationId: 'o1' },
    ])
  })

  it('emits an update payload without a password when editing', async () => {
    const wrapper = mount(UserForm, {
      props: {
        organisations,
        user: { id: '1', email: 'a@b.com', role: 'student', organisationId: 'o1' },
      },
      global: { plugins: [i18n] },
    })
    await wrapper.get('form').trigger('submit')

    expect(wrapper.emitted('submit')?.[0]).toEqual([
      { email: 'a@b.com', role: 'student', organisationId: 'o1' },
    ])
  })

  it('requires an email and blocks submission when blank', async () => {
    const wrapper = mount(UserForm, { props: { organisations }, global: { plugins: [i18n] } })
    await wrapper.get('input[type="password"]').setValue('secret123')
    await wrapper.get('form').trigger('submit')

    expect(wrapper.emitted('submit')).toBeUndefined()
  })

  it('has no detectable accessibility violations', async () => {
    const wrapper = mount(UserForm, { props: { organisations }, global: { plugins: [i18n] } })
    const results = await axe(wrapper.element, { rules: { region: { enabled: false } } })
    expect(results).toHaveNoViolations()
  })

  it('shows a button to generate a strong password when creating a new user', () => {
    const wrapper = mount(UserForm, { props: { organisations }, global: { plugins: [i18n] } })
    expect(wrapper.find('[data-testid="generate-password"]').exists()).toBe(true)
  })

  it('hides the generate-password button when editing an existing user', () => {
    const wrapper = mount(UserForm, {
      props: {
        organisations,
        user: { id: '1', email: 'a@b.com', role: 'student', organisationId: 'o1' },
      },
      global: { plugins: [i18n] },
    })
    expect(wrapper.find('[data-testid="generate-password"]').exists()).toBe(false)
  })

  it('fills and reveals the password field with a strong password when generated', async () => {
    const wrapper = mount(UserForm, { props: { organisations }, global: { plugins: [i18n] } })
    await wrapper.get('[data-testid="generate-password"]').trigger('click')

    const input = wrapper.get('input[type="text"][autocomplete="new-password"]')
    const value = (input.element as HTMLInputElement).value
    expect(value.length).toBeGreaterThanOrEqual(16)
    expect(value).toMatch(/[a-z]/)
    expect(value).toMatch(/[A-Z]/)
    expect(value).toMatch(/[0-9]/)
  })

  it('submits the generated password as part of the create payload', async () => {
    const wrapper = mount(UserForm, { props: { organisations }, global: { plugins: [i18n] } })
    await wrapper.get('input[type="email"]').setValue('new@b.com')
    await wrapper.get('[data-testid="generate-password"]').trigger('click')
    await wrapper.get('form').trigger('submit')

    const payload = wrapper.emitted('submit')?.[0]?.[0] as UserCreateInput
    expect(payload.password.length).toBeGreaterThanOrEqual(16)
  })
})
