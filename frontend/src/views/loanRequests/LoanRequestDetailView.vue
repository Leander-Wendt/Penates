<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'

import BaseButton from '@/components/common/BaseButton.vue'
import LoanRequestStatusBadge from '@/components/loanRequests/LoanRequestStatusBadge.vue'
import { useAuth } from '@/composables/useAuth'
import { useLoanRequestsStore } from '@/stores/loanRequests'
import type { LoanRequestStatus } from '@/types'
import { formatDate } from '@/utils/formatDate'

/**
 * Detail page for a single loan request, with status-change actions
 * (approve/reject/mark returned/cancel) available to Admin/Logistics
 * depending on the request's current status.
 */
const props = defineProps<{ id: string }>()

const { t } = useI18n()
const { hasRole } = useAuth()
const store = useLoanRequestsStore()

const request = computed(() => store.loanRequests.find((candidate) => candidate.id === props.id))

const canManage = computed(() => hasRole(['admin', 'logistics']))

const availableTransitions = computed<LoanRequestStatus[]>(() => {
  if (!request.value) {
    return []
  }
  if (request.value.status === 'pending') {
    return ['approved', 'rejected']
  }
  if (request.value.status === 'approved') {
    return ['returned', 'cancelled']
  }
  return []
})

onMounted(async () => {
  if (!request.value) {
    await store.fetchAll()
  }
})

async function setStatus(status: LoanRequestStatus): Promise<void> {
  await store.updateStatus(props.id, status)
}

const transitionLabels: Record<LoanRequestStatus, string> = {
  pending: '',
  approved: 'loanRequests.approve',
  rejected: 'loanRequests.reject',
  returned: 'loanRequests.markReturned',
  cancelled: 'loanRequests.cancel',
}

const transitionVariants: Record<LoanRequestStatus, 'primary' | 'secondary' | 'danger'> = {
  pending: 'secondary',
  approved: 'primary',
  rejected: 'danger',
  returned: 'primary',
  cancelled: 'danger',
}
</script>

<template>
  <div v-if="request" class="flex flex-col gap-6">
    <div class="flex items-center justify-between">
      <h1>{{ t('loanRequests.detailTitle') }}</h1>
      <LoanRequestStatusBadge :status="request.status" />
    </div>

    <dl class="grid grid-cols-1 gap-x-6 gap-y-2 sm:grid-cols-2">
      <div class="sm:col-span-2">
        <dt class="font-bold uppercase">{{ t('loanRequests.itemsLabel') }}</dt>
        <dd>{{ request.items.map((item) => item.name).join(', ') }}</dd>
      </div>
      <div>
        <dt class="font-bold uppercase">{{ t('loanRequests.dateOfLending') }}</dt>
        <dd>{{ formatDate(request.dateOfLending) }}</dd>
      </div>
      <div>
        <dt class="font-bold uppercase">{{ t('loanRequests.dateOfReturn') }}</dt>
        <dd>{{ formatDate(request.dateOfReturn) }}</dd>
      </div>
      <div>
        <dt class="font-bold uppercase">{{ t('loanRequests.locationOfItems') }}</dt>
        <dd>{{ request.locationOfItems }}</dd>
      </div>
      <div>
        <dt class="font-bold uppercase">{{ t('loanRequests.requestingUser') }}</dt>
        <dd>{{ request.requestingUserId }}</dd>
      </div>
    </dl>

    <div v-if="canManage && availableTransitions.length > 0" class="flex flex-wrap gap-2">
      <BaseButton
        v-for="status in availableTransitions"
        :key="status"
        :variant="transitionVariants[status]"
        @click="setStatus(status)"
      >
        {{ t(transitionLabels[status]) }}
      </BaseButton>
    </div>
  </div>
</template>
