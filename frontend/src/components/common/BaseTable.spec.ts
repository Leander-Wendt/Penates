import { describe, expect, it } from 'vitest'
import { mount } from '@vue/test-utils'
import { axe } from 'vitest-axe'

import BaseTable from './BaseTable.vue'

const columns = [
  { id: 'name', label: 'Name' },
  { id: 'category', label: 'Category' },
]
const rows = [
  { id: '1', name: 'Drill', category: 'Tools' },
  { id: '2', name: 'Ladder', category: 'Tools' },
]

describe('BaseTable', () => {
  it('renders a header cell per column', () => {
    const wrapper = mount(BaseTable, { props: { columns, rows, rowKey: 'id' } })
    const headers = wrapper.findAll('th')
    expect(headers.map((h) => h.text())).toEqual(['Name', 'Category'])
  })

  it('renders a row per data item with default cell content', () => {
    const wrapper = mount(BaseTable, { props: { columns, rows, rowKey: 'id' } })
    const dataRows = wrapper.findAll('tbody tr')
    expect(dataRows).toHaveLength(2)
    expect(dataRows[0].text()).toContain('Drill')
  })

  it('renders scoped slot content for a column when provided', () => {
    const wrapper = mount(BaseTable, {
      props: { columns, rows, rowKey: 'id' },
      slots: {
        name: `<template #name="{ row }"><strong>{{ row.name }}</strong></template>`,
      },
    })
    expect(wrapper.find('strong').exists()).toBe(true)
  })

  it('renders an empty message when there are no rows', () => {
    const wrapper = mount(BaseTable, {
      props: { columns, rows: [], rowKey: 'id', emptyMessage: 'Nothing here' },
    })
    expect(wrapper.text()).toContain('Nothing here')
  })

  it('has no detectable accessibility violations', async () => {
    const wrapper = mount(BaseTable, {
      props: { columns, rows, rowKey: 'id', caption: 'Items' },
    })
    const results = await axe(wrapper.element, { rules: { region: { enabled: false } } })
    expect(results).toHaveNoViolations()
  })
})
