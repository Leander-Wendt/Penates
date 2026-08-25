/**
 * Formats an ISO date or datetime string (e.g. `2026-03-05` or
 * `2026-03-05T10:15:00Z`, as returned by the backend) as `dd/mm/yyyy` for
 * display. Returns an empty string for null/undefined/empty input, and the
 * original value unchanged if it doesn't match the expected ISO shape.
 */
export function formatDate(value: string | null | undefined): string {
  if (!value) {
    return ''
  }
  const match = /^(\d{4})-(\d{2})-(\d{2})/.exec(value)
  if (!match) {
    return value
  }
  const [, year, month, day] = match
  return `${day}/${month}/${year}`
}
