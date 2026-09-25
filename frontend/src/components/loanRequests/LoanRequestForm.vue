<script setup lang="ts">
import { reactive } from 'vue'
import { useI18n } from 'vue-i18n'

import BaseButton from '@/components/common/BaseButton.vue'
import BaseInput from '@/components/common/BaseInput.vue'
import ItemPicker from './ItemPicker.vue'
import type { Item, LoanRequestCreateInput } from '@/types'

/**
 * Create form for a loan request: item multi-select, lending/return
 * dates, and a location. Available items to choose from are supplied by
 * the parent view.
 */
const props = defineProps<{
  items: Item[]
  preSelectedItems?: string[]
}>()
const emit = defineEmits<{ submit: [input: LoanRequestCreateInput]; cancel: [] }>()

const { t } = useI18n()

const form = reactive({
  itemInventoryNumbers: props.preSelectedItems || ([] as string[]),
  dateOfLending: '',
  dateOfReturn: '',
  locationOfItems: '',
})

const errors = reactive({ items: '', dateOfLending: '', dateOfReturn: '', locationOfItems: '' })

function handleSubmit(): void {
  errors.items = form.itemInventoryNumbers.length > 0 ? '' : t('common.required')
  errors.dateOfLending = form.dateOfLending ? '' : t('common.required')
  errors.dateOfReturn = form.dateOfReturn ? '' : t('common.required')
  errors.locationOfItems = form.locationOfItems.trim() ? '' : t('common.required')

  if (errors.items || errors.dateOfLending || errors.dateOfReturn || errors.locationOfItems) {
    return
  }

  emit('submit', { ...form })
}
</script>

<template>
  <form class="flex flex-col gap-4" @submit.prevent="handleSubmit">
    <ItemPicker v-model="form.itemInventoryNumbers" :items="props.items" />
    <p v-if="errors.items" role="alert" aria-live="assertive" class="font-bold text-brut-red">
      {{ errors.items }}
    </p>
    <BaseInput
      v-model="form.dateOfLending"
      type="date"
      :label="t('loanRequests.dateOfLending')"
      required
      :error="errors.dateOfLending"
    />
    <BaseInput
      v-model="form.dateOfReturn"
      type="date"
      :label="t('loanRequests.dateOfReturn')"
      required
      :error="errors.dateOfReturn"
    />
    <BaseInput
      v-model="form.locationOfItems"
      :label="t('loanRequests.locationOfItems')"
      required
      :error="errors.locationOfItems"
    />
    <div class="flex gap-2">
      <BaseButton type="submit">{{ t('common.save') }}</BaseButton>
      <BaseButton type="button" variant="secondary" @click="emit('cancel')">
        {{ t('common.cancel') }}
      </BaseButton>
    </div>
  </form>
</template>
