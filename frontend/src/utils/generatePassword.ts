const LOWERCASE = 'abcdefghijkmnopqrstuvwxyz'
const UPPERCASE = 'ABCDEFGHJKLMNPQRSTUVWXYZ'
const DIGITS = '23456789'
const SYMBOLS = '!@#$%^&*-_=+?'
const CHARACTER_CLASSES = [LOWERCASE, UPPERCASE, DIGITS, SYMBOLS]
const ALL_CHARACTERS = CHARACTER_CLASSES.join('')

/**
 * Generates a cryptographically random password (via `crypto.getRandomValues`)
 * of the given length (default 16), guaranteed to contain at least one
 * lowercase letter, one uppercase letter, one digit, and one symbol.
 * Visually ambiguous characters (`l`, `I`, `1`, `O`, `0`) are excluded so a
 * generated password can be read and typed back correctly. Throws if length
 * is too small to fit one character from each class.
 */
export function generateStrongPassword(length = 16): string {
  if (length < CHARACTER_CLASSES.length) {
    throw new Error(
      `generateStrongPassword requires a length of at least ${CHARACTER_CLASSES.length}`,
    )
  }

  let password: string
  do {
    password = Array.from({ length }, () => randomCharacter(ALL_CHARACTERS)).join('')
  } while (!CHARACTER_CLASSES.every((charClass) => containsAny(password, charClass)))

  return password
}

function randomCharacter(pool: string): string {
  const bytes = new Uint32Array(1)
  crypto.getRandomValues(bytes)
  return pool[bytes[0] % pool.length]
}

function containsAny(value: string, characters: string): boolean {
  return [...value].some((char) => characters.includes(char))
}
