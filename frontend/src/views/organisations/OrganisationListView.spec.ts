import { describe, expect, it, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { axe } from 'vitest-axe'

import i18n from '@/locales'
import * as organisationsApi from '@/api/organisations'
import OrganisationListView from './OrganisationListView.vue'

vi.mock('@/api/organisations')

const organisation = { id: '1', name: 'Org 1', description: 'First org' }

describe('OrganisationListView', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
    vi.mocked(organisationsApi.listOrganisations).mockResolvedValue([organisation])
  })

  it('fetches and renders organisations on mount', async () => {
    const wrapper = mount(OrganisationListView, { global: { plugins: [i18n] } })
    await flushPromises()

    expect(wrapper.text()).toContain('Org 1')
  })

  it('opens the create modal and creates an organisation on submit', async () => {
    const created = { id: '2', name: 'New org', description: '' }
    vi.mocked(organisationsApi.createOrganisation).mockResolvedValue(created)
    const wrapper = mount(OrganisationListView, { global: { plugins: [i18n] } })
    await flushPromises()

    await wrapper.get('button').trigger('click')
    await flushPromises()

    const nameInput = document.body.querySelector('form input') as HTMLInputElement
    nameInput.value = 'New org'
    nameInput.dispatchEvent(new Event('input'))
    const form = document.body.querySelector('form') as HTMLFormElement
    form.dispatchEvent(new Event('submit', { cancelable: true }))
    await flushPromises()

    expect(organisationsApi.createOrganisation).toHaveBeenCalledWith({
      name: 'New org',
      description: '',
    })
    expect(wrapper.text()).toContain('New org')
  })

  it('deletes an organisation after confirmation', async () => {
    vi.mocked(organisationsApi.deleteOrganisation).mockResolvedValue(undefined)
    vi.spyOn(window, 'confirm').mockReturnValue(true)
    const wrapper = mount(OrganisationListView, { global: { plugins: [i18n] } })
    await flushPromises()

    const deleteButton = wrapper.findAll('button').find((b) => b.text() === 'Delete')
    await deleteButton?.trigger('click')
    await flushPromises()

    expect(organisationsApi.deleteOrganisation).toHaveBeenCalledWith('1')
  })

  it('does not delete when confirmation is declined', async () => {
    vi.spyOn(window, 'confirm').mockReturnValue(false)
    const wrapper = mount(OrganisationListView, { global: { plugins: [i18n] } })
    await flushPromises()

    const deleteButton = wrapper.findAll('button').find((b) => b.text() === 'Delete')
    await deleteButton?.trigger('click')
    await flushPromises()

    expect(organisationsApi.deleteOrganisation).not.toHaveBeenCalled()
  })

  it('has no detectable accessibility violations', async () => {
    const wrapper = mount(OrganisationListView, { global: { plugins: [i18n] } })
    await flushPromises()
    const results = await axe(wrapper.element, { rules: { region: { enabled: false } } })
    expect(results).toHaveNoViolations()
  })
})
