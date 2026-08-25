import { describe, expect, it } from 'vitest'

import { formatDate } from './formatDate'

describe('formatDate', () => {
  it('formats an ISO date string as dd/mm/yyyy', () => {
    expect(formatDate('2026-03-05')).toBe('05/03/2026')
  })

  it('formats a full ISO datetime string using only its date portion', () => {
    expect(formatDate('2026-03-05T10:15:00Z')).toBe('05/03/2026')
  })

  it('returns an empty string for null, undefined, or empty input', () => {
    expect(formatDate(null)).toBe('')
    expect(formatDate(undefined)).toBe('')
    expect(formatDate('')).toBe('')
  })

  it('returns the original value unchanged when it does not match the expected format', () => {
    expect(formatDate('not-a-date')).toBe('not-a-date')
  })
})
