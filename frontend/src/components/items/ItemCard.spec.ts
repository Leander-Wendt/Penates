import { describe, expect, it } from 'vitest'
import { mount } from '@vue/test-utils'
import { createRouter, createWebHistory } from 'vue-router'
import { axe } from 'vitest-axe'

import i18n from '@/locales'
import ItemCard from './ItemCard.vue'

const item = {
  inventoryNumber: '000001',
  name: 'Drill',
  description: '',
  imagePath: '/uploads/000001.png',
  category: 'Tools',
  location: 'Shelf A1',
  amount: 1,
  electricalAppliance: true,
  note: '',
}

function makeRouter() {
  return createRouter({
    history: createWebHistory(),
    routes: [{ path: '/items/:inventoryNumber', component: { template: '<div />' } }],
  })
}

describe('ItemCard', () => {
  it('links to the item detail page', () => {
    const wrapper = mount(ItemCard, { props: { item }, global: { plugins: [i18n, makeRouter()] } })
    expect(wrapper.get('a').attributes('href')).toBe('/items/000001')
  })

  it('renders the image when imagePath is set', () => {
    const wrapper = mount(ItemCard, { props: { item }, global: { plugins: [i18n, makeRouter()] } })
    expect(wrapper.get('img').attributes('src')).toBe('/uploads/000001.png')
  })

  it('renders a no-image placeholder when imagePath is null', () => {
    const wrapper = mount(ItemCard, {
      props: { item: { ...item, imagePath: null } },
      global: { plugins: [i18n, makeRouter()] },
    })
    expect(wrapper.find('img').exists()).toBe(false)
    expect(wrapper.text()).toContain('No image')
  })

  it('shows the item name, category and location', () => {
    const wrapper = mount(ItemCard, { props: { item }, global: { plugins: [i18n, makeRouter()] } })
    expect(wrapper.text()).toContain('Drill')
    expect(wrapper.text()).toContain('Tools')
    expect(wrapper.text()).toContain('Shelf A1')
  })

  it('has no detectable accessibility violations', async () => {
    const wrapper = mount(ItemCard, { props: { item }, global: { plugins: [i18n, makeRouter()] } })
    const results = await axe(wrapper.element, { rules: { region: { enabled: false } } })
    expect(results).toHaveNoViolations()
  })
})
