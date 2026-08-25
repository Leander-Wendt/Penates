<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'

import BaseButton from '@/components/common/BaseButton.vue'
import BaseTable from '@/components/common/BaseTable.vue'
import LoanRequestStatusBadge from '@/components/loanRequests/LoanRequestStatusBadge.vue'
import { useLoanRequestsStore } from '@/stores/loanRequests'
import type { LoanRequest } from '@/types'
import { formatDate } from '@/utils/formatDate'

/**
 * Loan request list. The backend already scopes results to the current
 * user's own requests for Students, and to all requests for
 * Admin/Logistics.
 */
const { t } = useI18n()
const store = useLoanRequestsStore()
const router = useRouter()

const columns = computed(() => [
  { id: 'items', label: t('loanRequests.itemsLabel') },
  { id: 'dateOfLending', label: t('loanRequests.dateOfLending') },
  { id: 'dateOfReturn', label: t('loanRequests.dateOfReturn') },
  { id: 'locationOfItems', label: t('loanRequests.locationOfItems') },
  { id: 'status', label: t('common.status') },
  { id: 'actions', label: t('common.actions') },
])

onMounted(() => {
  store.fetchAll()
})

function goToCreate(): void {
  router.push('/loan-requests/new')
}

function goToDetail(request: LoanRequest): void {
  router.push(`/loan-requests/${request.id}`)
}
</script>

<template>
  <div class="flex flex-col gap-4">
    <div class="flex items-center justify-between">
      <h1>{{ t('loanRequests.title') }}</h1>
      <BaseButton @click="goToCreate">{{ t('loanRequests.createButton') }}</BaseButton>
    </div>

    <BaseTable
      :columns="columns"
      :rows="store.loanRequests"
      row-key="id"
      :caption="t('loanRequests.title')"
      :empty-message="t('common.noResults')"
    >
      <template #items="{ row }">
        {{ (row as LoanRequest).items.map((item) => item.name).join(', ') }}
      </template>
      <template #dateOfLending="{ row }">
        {{ formatDate((row as LoanRequest).dateOfLending) }}
      </template>
      <template #dateOfReturn="{ row }">
        {{ formatDate((row as LoanRequest).dateOfReturn) }}
      </template>
      <template #status="{ row }">
        <LoanRequestStatusBadge :status="(row as LoanRequest).status" />
      </template>
      <template #actions="{ row }">
        <BaseButton variant="secondary" @click="goToDetail(row as LoanRequest)">
          {{ t('common.view') }}
        </BaseButton>
      </template>
    </BaseTable>
  </div>
</template>
