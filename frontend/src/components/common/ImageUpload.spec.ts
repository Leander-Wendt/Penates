import { describe, expect, it, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { axe } from 'vitest-axe'

import i18n from '@/locales'
import ImageUpload from './ImageUpload.vue'

beforeEach(() => {
  URL.createObjectURL = vi.fn(() => 'blob:preview-url')
})

describe('ImageUpload', () => {
  it('renders a file input associated with its label', () => {
    const wrapper = mount(ImageUpload, {
      props: { label: 'Item photo' },
      global: { plugins: [i18n] },
    })
    const label = wrapper.get('label')
    const input = wrapper.get('input[type="file"]')
    expect(label.attributes('for')).toBe(input.attributes('id'))
  })

  it('shows the existing image when imageUrl is provided', () => {
    const wrapper = mount(ImageUpload, {
      props: { label: 'Item photo', imageUrl: '/uploads/1.png' },
      global: { plugins: [i18n] },
    })
    expect(wrapper.get('img').attributes('src')).toBe('/uploads/1.png')
  })

  it('shows no image and no remove button when nothing is set', () => {
    const wrapper = mount(ImageUpload, {
      props: { label: 'Item photo' },
      global: { plugins: [i18n] },
    })
    expect(wrapper.find('img').exists()).toBe(false)
    expect(wrapper.find('button').exists()).toBe(false)
  })

  it('emits select with the chosen file and previews it', async () => {
    const wrapper = mount(ImageUpload, {
      props: { label: 'Item photo' },
      global: { plugins: [i18n] },
    })
    const file = new File(['x'], 'photo.png', { type: 'image/png' })
    const input = wrapper.get('input[type="file"]')
    Object.defineProperty(input.element, 'files', { value: [file] })
    await input.trigger('change')

    expect(wrapper.emitted('select')?.[0]).toEqual([file])
    expect(wrapper.get('img').attributes('src')).toBe('blob:preview-url')
  })

  it('emits remove when the remove button is clicked', async () => {
    const wrapper = mount(ImageUpload, {
      props: { label: 'Item photo', imageUrl: '/uploads/1.png' },
      global: { plugins: [i18n] },
    })
    await wrapper.get('button').trigger('click')
    expect(wrapper.emitted('remove')).toBeTruthy()
  })

  it('has no detectable accessibility violations', async () => {
    const wrapper = mount(ImageUpload, {
      props: { label: 'Item photo', imageUrl: '/uploads/1.png', imageAlt: 'A drill' },
      global: { plugins: [i18n] },
    })
    const results = await axe(wrapper.element, { rules: { region: { enabled: false } } })
    expect(results).toHaveNoViolations()
  })
})
