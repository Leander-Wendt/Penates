import { describe, expect, it } from 'vitest'

import { generateStrongPassword } from './generatePassword'

describe('generateStrongPassword', () => {
  it('generates a password of the default length', () => {
    expect(generateStrongPassword()).toHaveLength(16)
  })

  it('generates a password of a custom length', () => {
    expect(generateStrongPassword(24)).toHaveLength(24)
  })

  it('only uses unambiguous letters, digits, and symbols', () => {
    const password = generateStrongPassword(64)
    expect(password).toMatch(/^[A-HJ-NP-Za-km-z2-9!@#$%^&*_=+?-]+$/)
  })

  it('includes at least one lowercase letter, uppercase letter, digit, and symbol', () => {
    const password = generateStrongPassword(64)
    expect(password).toMatch(/[a-z]/)
    expect(password).toMatch(/[A-Z]/)
    expect(password).toMatch(/[0-9]/)
    expect(password).toMatch(/[!@#$%^&*_=+?-]/)
  })

  it('generates different passwords across calls', () => {
    const a = generateStrongPassword()
    const b = generateStrongPassword()
    expect(a).not.toBe(b)
  })

  it('throws for a length too short to fit every character class', () => {
    expect(() => generateStrongPassword(2)).toThrow()
  })
})
