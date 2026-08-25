<script setup lang="ts">
import { nextTick, ref, useId, watch } from 'vue'

/**
 * Neo-brutalist modal dialog. Traps keyboard focus while open, closes on
 * `Escape`, and restores focus to whichever element triggered it once
 * closed.
 */
const props = withDefaults(
  defineProps<{
    open: boolean
    title: string
  }>(),
  {},
)

const emit = defineEmits<{ close: [] }>()

const titleId = useId()
const dialogRef = ref<HTMLElement | null>(null)
let previouslyFocused: HTMLElement | null = null

function getFocusableElements(): HTMLElement[] {
  if (!dialogRef.value) {
    return []
  }
  const selector =
    'a[href], button:not([disabled]), textarea:not([disabled]), input:not([disabled]), select:not([disabled]), [tabindex]:not([tabindex="-1"])'
  return Array.from(dialogRef.value.querySelectorAll<HTMLElement>(selector))
}

function handleKeydown(event: KeyboardEvent): void {
  if (event.key === 'Escape') {
    emit('close')
    return
  }
  if (event.key !== 'Tab') {
    return
  }
  const focusable = getFocusableElements()
  if (focusable.length === 0) {
    return
  }
  const first = focusable[0]
  const last = focusable[focusable.length - 1]
  if (event.shiftKey && document.activeElement === first) {
    event.preventDefault()
    last.focus()
  } else if (!event.shiftKey && document.activeElement === last) {
    event.preventDefault()
    first.focus()
  }
}

watch(
  () => props.open,
  async (isOpen) => {
    if (isOpen) {
      previouslyFocused = document.activeElement as HTMLElement | null
      await nextTick()
      const focusable = getFocusableElements()
      ;(focusable[0] ?? dialogRef.value)?.focus()
    } else {
      previouslyFocused?.focus()
    }
  },
)
</script>

<template>
  <Teleport to="body">
    <div v-if="open" class="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div class="fixed inset-0 bg-brut-black/60" @click="emit('close')" />
      <div
        ref="dialogRef"
        role="dialog"
        aria-modal="true"
        :aria-labelledby="titleId"
        class="relative z-10 w-full max-w-lg border-3 border-brut-black bg-brut-white p-6 shadow-brut-lg"
        tabindex="-1"
        @keydown="handleKeydown"
      >
        <div class="mb-4 flex items-start justify-between gap-4">
          <h2 :id="titleId" class="text-xl">{{ title }}</h2>
          <button
            type="button"
            class="border-3 border-brut-black bg-brut-white px-2 py-1 font-bold"
            :aria-label="$t('common.close')"
            @click="emit('close')"
          >
            &times;
          </button>
        </div>
        <slot />
      </div>
    </div>
  </Teleport>
</template>
