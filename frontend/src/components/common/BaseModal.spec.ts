import { describe, expect, it } from 'vitest'
import { DOMWrapper, mount } from '@vue/test-utils'
import { axe } from 'vitest-axe'

import i18n from '@/locales'
import BaseModal from './BaseModal.vue'

function mountModal(open: boolean) {
  return mount(BaseModal, {
    props: { open, title: 'Edit item' },
    slots: {
      default: '<button id="first">First</button><button id="second">Second</button>',
    },
    attachTo: document.body,
    global: { plugins: [i18n] },
  })
}

function getDialog(): DOMWrapper<HTMLElement> {
  const el = document.body.querySelector('[role="dialog"]')
  if (!el) {
    throw new Error('dialog not found')
  }
  return new DOMWrapper(el as HTMLElement)
}

describe('BaseModal', () => {
  it('renders nothing when closed', () => {
    const wrapper = mountModal(false)
    expect(document.body.querySelector('[role="dialog"]')).toBeNull()
    wrapper.unmount()
  })

  it('renders the dialog with an accessible name when open', () => {
    const wrapper = mountModal(true)
    const dialog = document.body.querySelector('[role="dialog"]') as HTMLElement
    expect(dialog).not.toBeNull()
    expect(dialog.getAttribute('aria-modal')).toBe('true')
    const labelledBy = dialog.getAttribute('aria-labelledby')
    expect(document.getElementById(labelledBy as string)?.textContent).toBe('Edit item')
    wrapper.unmount()
  })

  it('emits close on Escape', async () => {
    const wrapper = mountModal(true)
    await getDialog().trigger('keydown', { key: 'Escape' })
    expect(wrapper.emitted('close')).toBeTruthy()
    wrapper.unmount()
  })

  it('emits close when the backdrop is clicked', async () => {
    const wrapper = mountModal(true)
    const backdrop = new DOMWrapper(
      document.body.querySelector('.fixed.inset-0.bg-brut-black\\/60') as HTMLElement,
    )
    await backdrop.trigger('click')
    expect(wrapper.emitted('close')).toBeTruthy()
    wrapper.unmount()
  })

  it('moves focus into the dialog when opened, and restores it on close', async () => {
    const trigger = document.createElement('button')
    trigger.textContent = 'Open'
    document.body.appendChild(trigger)
    trigger.focus()

    const wrapper = mount(BaseModal, {
      props: { open: false, title: 'Edit item' },
      slots: { default: '<button id="first">First</button>' },
      attachTo: document.body,
      global: { plugins: [i18n] },
    })
    await wrapper.setProps({ open: true })
    await new Promise((resolve) => setTimeout(resolve, 0))

    const dialog = document.body.querySelector('[role="dialog"]') as HTMLElement
    const firstFocusable = dialog.querySelector('button, [href], input, select, textarea')
    expect(document.activeElement).toBe(firstFocusable)

    await wrapper.setProps({ open: false })
    await new Promise((resolve) => setTimeout(resolve, 0))
    expect(document.activeElement).toBe(trigger)

    wrapper.unmount()
    trigger.remove()
  })

  it('wraps Tab focus from the last to the first focusable element', async () => {
    mountModal(true)
    const dialog = getDialog()
    const buttons = dialog.findAll('button')
    const first = buttons[0].element as HTMLElement
    const last = buttons[buttons.length - 1].element as HTMLElement
    last.focus()
    await dialog.trigger('keydown', { key: 'Tab' })
    expect(document.activeElement).toBe(first)
  })

  it('has no detectable accessibility violations when open', async () => {
    mountModal(true)
    const results = await axe(getDialog().element)
    expect(results).toHaveNoViolations()
  })
})
