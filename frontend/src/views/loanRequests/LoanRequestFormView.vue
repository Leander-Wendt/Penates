<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'

import LoanRequestForm from '@/components/loanRequests/LoanRequestForm.vue'
import { useItemsStore } from '@/stores/items'
import { useLoanRequestsStore } from '@/stores/loanRequests'
import type { LoanRequestCreateInput } from '@/types'

/**
 * Create page for a new loan request: lets the current user pick from the
 * available items and submit a request for approval.
 */
const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const itemsStore = useItemsStore()
const loanRequestsStore = useLoanRequestsStore()

// Get pre-selected items from query parameters
const preSelectedItems = ref<string[]>([])

onMounted(() => {
  itemsStore.fetchAll()

  // Check for pre-selected items in query params
  const itemsQuery = route.query.items
  if (itemsQuery) {
    if (Array.isArray(itemsQuery)) {
      preSelectedItems.value = itemsQuery
    } else {
      preSelectedItems.value = itemsQuery.split(',').filter(Boolean)
    }
  }
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
    <LoanRequestForm
      :items="itemsStore.items"
      :pre-selected-items="preSelectedItems"
      @submit="handleSubmit"
      @cancel="handleCancel"
    />
  </div>
</template>
