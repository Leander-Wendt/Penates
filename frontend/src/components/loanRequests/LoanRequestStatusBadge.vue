<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

import type { LoanRequestStatus } from '@/types'

/**
 * Renders a loan request status as an icon plus text label, so the state
 * is never conveyed by color alone.
 */
const props = defineProps<{ status: LoanRequestStatus }>()

const { t } = useI18n()

const ICONS: Record<LoanRequestStatus, string> = {
  pending: '⏳',
  approved: '✔',
  rejected: '✕',
  returned: '↩',
  cancelled: '⊘',
}

const COLORS: Record<LoanRequestStatus, string> = {
  pending: 'bg-brut-yellow text-brut-black',
  approved: 'bg-brut-green text-brut-white',
  rejected: 'bg-brut-red text-brut-white',
  returned: 'bg-brut-blue text-brut-white',
  cancelled: 'bg-brut-white text-brut-black',
}

const icon = computed(() => ICONS[props.status])
const colorClass = computed(() => COLORS[props.status])
const label = computed(() => t(`loanRequests.status.${props.status}`))
</script>

<template>
  <span
    class="inline-flex items-center gap-1 border-3 border-brut-black px-2 py-1 text-sm font-bold uppercase"
    :class="colorClass"
  >
    <span aria-hidden="true">{{ icon }}</span>
    <span>{{ label }}</span>
  </span>
</template>
