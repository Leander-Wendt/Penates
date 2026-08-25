<script setup lang="ts">
import { reactive, watch } from 'vue'
import { useI18n } from 'vue-i18n'

import BaseButton from '@/components/common/BaseButton.vue'
import BaseInput from '@/components/common/BaseInput.vue'
import type { Organisation, OrganisationInput } from '@/types'

/**
 * Create/edit form for an Organisation. When `organisation` is provided the
 * fields are pre-filled for editing; otherwise the form starts empty for
 * creation.
 */
const props = defineProps<{ organisation?: Organisation | null }>()
const emit = defineEmits<{ submit: [input: OrganisationInput]; cancel: [] }>()

const { t } = useI18n()

const form = reactive<OrganisationInput>({
  name: props.organisation?.name ?? '',
  description: props.organisation?.description ?? '',
})

const errors = reactive<{ name: string }>({ name: '' })

watch(
  () => props.organisation,
  (organisation) => {
    form.name = organisation?.name ?? ''
    form.description = organisation?.description ?? ''
  },
)

function handleSubmit(): void {
  errors.name = form.name.trim() ? '' : t('common.required')
  if (errors.name) {
    return
  }
  emit('submit', { name: form.name, description: form.description })
}
</script>

<template>
  <form class="flex flex-col gap-4" @submit.prevent="handleSubmit">
    <BaseInput v-model="form.name" :label="t('common.name')" required :error="errors.name" />
    <BaseInput v-model="form.description" :label="t('common.description')" />
    <div class="flex gap-2">
      <BaseButton type="submit">{{ t('common.save') }}</BaseButton>
      <BaseButton type="button" variant="secondary" @click="emit('cancel')">
        {{ t('common.cancel') }}
      </BaseButton>
    </div>
  </form>
</template>
