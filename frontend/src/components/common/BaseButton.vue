<script setup lang="ts">
/**
 * Neo-brutalist button primitive: thick black border, hard offset shadow
 * that "presses" flat on click, and a high-contrast color per variant.
 */
withDefaults(
  defineProps<{
    variant?: 'primary' | 'secondary' | 'danger'
    type?: 'button' | 'submit' | 'reset'
    disabled?: boolean
    ariaLabel?: string
  }>(),
  {
    variant: 'primary',
    type: 'button',
    disabled: false,
    ariaLabel: undefined,
  },
)

const variantClasses: Record<string, string> = {
  primary: 'bg-brut-yellow text-brut-black',
  secondary: 'bg-brut-white text-brut-black',
  danger: 'bg-brut-red text-brut-white',
}

defineOptions({ inheritAttrs: false })
</script>

<template>
  <button
    :type="type"
    :disabled="disabled"
    :aria-label="ariaLabel"
    class="inline-flex items-center justify-center gap-2 border-3 border-brut-black px-4 py-2 font-bold uppercase tracking-tight shadow-brut transition-transform duration-100 ease-out hover:-translate-x-0.5 hover:-translate-y-0.5 hover:shadow-brut-lg active:translate-x-0.5 active:translate-y-0.5 active:shadow-brut-pressed disabled:cursor-not-allowed disabled:opacity-50 disabled:hover:translate-x-0 disabled:hover:translate-y-0 disabled:hover:shadow-brut"
    :class="variantClasses[variant]"
    v-bind="$attrs"
  >
    <slot />
  </button>
</template>
