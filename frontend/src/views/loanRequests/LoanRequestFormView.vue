<script setup lang="ts">
import { onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'

import LoanRequestForm from '@/components/loanRequests/LoanRequestForm.vue'
import { useItemsStore } from '@/stores/items'
import { useLoanRequestsStore } from '@/stores/loanRequests'
import type { LoanRequestCreateInput } from '@/types'

/**
 * Create page for a new loan request: lets the current user pick from the
 * available items and submit a request for approval.
 */
const { t } = useI18n()
const router = useRouter()
const itemsStore = useItemsStore()
const loanRequestsStore = useLoanRequestsStore()

onMounted(() => {
  itemsStore.fetchAll()
})

async function handleSubmit(input: LoanRequestCreateInput): Promise<void> {
  const created = await loanRequestsStore.create(input)
  router.push(`/loan-requests/${created.id}`)
}

function handleCancel(): void {
  router.push('/loan-requests')
}
</script>

<template>
  <div class="flex flex-col gap-4">
    <h1>{{ t('loanRequests.createTitle') }}</h1>
    <LoanRequestForm :items="itemsStore.items" @submit="handleSubmit" @cancel="handleCancel" />
  </div>
</template>
