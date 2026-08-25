import { describe, expect, it } from 'vitest'
import { mount } from '@vue/test-utils'
import { axe } from 'vitest-axe'

import i18n from '@/locales'
import ItemSearchFilter from './ItemSearchFilter.vue'

describe('ItemSearchFilter', () => {
  it('emits update:searchTerm as the user types', async () => {
    const wrapper = mount(ItemSearchFilter, {
      props: { searchTerm: '', categoryFilter: '', categories: ['Tools'] },
      global: { plugins: [i18n] },
    })
    await wrapper.get('input[type="search"]').setValue('drill')
    expect(wrapper.emitted('update:searchTerm')?.[0]).toEqual(['drill'])
  })

  it('emits update:categoryFilter when a category is selected', async () => {
    const wrapper = mount(ItemSearchFilter, {
      props: { searchTerm: '', categoryFilter: '', categories: ['Tools', 'Electronics'] },
      global: { plugins: [i18n] },
    })
    await wrapper.get('select').setValue('Electronics')
    expect(wrapper.emitted('update:categoryFilter')?.[0]).toEqual(['Electronics'])
  })

  it('lists every provided category plus an "all" option', () => {
    const wrapper = mount(ItemSearchFilter, {
      props: { searchTerm: '', categoryFilter: '', categories: ['Tools', 'Electronics'] },
      global: { plugins: [i18n] },
    })
    const options = wrapper.findAll('option').map((o) => o.text())
    expect(options).toEqual(['All categories', 'Tools', 'Electronics'])
  })

  it('has no detectable accessibility violations', async () => {
    const wrapper = mount(ItemSearchFilter, {
      props: { searchTerm: '', categoryFilter: '', categories: ['Tools'] },
      global: { plugins: [i18n] },
    })
    const results = await axe(wrapper.element, { rules: { region: { enabled: false } } })
    expect(results).toHaveNoViolations()
  })
})
