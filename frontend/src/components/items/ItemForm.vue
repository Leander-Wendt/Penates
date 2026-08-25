<script setup lang="ts">
import { reactive, useId, watch } from 'vue'
import { useI18n } from 'vue-i18n'

import BaseButton from '@/components/common/BaseButton.vue'
import BaseInput from '@/components/common/BaseInput.vue'
import ImageUpload from '@/components/common/ImageUpload.vue'
import { useAuth } from '@/composables/useAuth'
import type { Item, ItemInput } from '@/types'

/**
 * Create/edit form for an Item. Public fields are always shown; the
 * restricted fields section (manufacturer, serial number, inspection
 * dates, purchase date, price, resolution number) only renders for
 * Admin/Logistics, mirroring the backend's field-level restriction for
 * Students.
 */
const props = defineProps<{ item?: Item | null }>()
const emit = defineEmits<{
  submit: [input: ItemInput]
  imageSelect: [file: File]
  imageRemove: []
  cancel: []
}>()

const { t } = useI18n()
const { hasRole } = useAuth()

const electricalId = useId()

const form = reactive<ItemInput>({
  name: props.item?.name ?? '',
  description: props.item?.description ?? '',
  category: props.item?.category ?? '',
  location: props.item?.location ?? '',
  amount: props.item?.amount ?? 1,
  electricalAppliance: props.item?.electricalAppliance ?? false,
  note: props.item?.note ?? '',
  manufacturer: props.item?.manufacturer ?? '',
  serialNumber: props.item?.serialNumber ?? '',
  lastTechnicalInspectionDate: props.item?.lastTechnicalInspectionDate ?? '',
  lastElectricalInspectionDate: props.item?.lastElectricalInspectionDate ?? '',
  dateOfPurchase: props.item?.dateOfPurchase ?? '',
  price: props.item?.price ?? 0,
  resolutionNumber: props.item?.resolutionNumber ?? '',
})

const errors = reactive({ name: '', category: '', location: '' })

watch(
  () => props.item,
  (item) => {
    form.name = item?.name ?? ''
    form.description = item?.description ?? ''
    form.category = item?.category ?? ''
    form.location = item?.location ?? ''
    form.amount = item?.amount ?? 1
    form.electricalAppliance = item?.electricalAppliance ?? false
    form.note = item?.note ?? ''
    form.manufacturer = item?.manufacturer ?? ''
    form.serialNumber = item?.serialNumber ?? ''
    form.lastTechnicalInspectionDate = item?.lastTechnicalInspectionDate ?? ''
    form.lastElectricalInspectionDate = item?.lastElectricalInspectionDate ?? ''
    form.dateOfPurchase = item?.dateOfPurchase ?? ''
    form.price = item?.price ?? 0
    form.resolutionNumber = item?.resolutionNumber ?? ''
  },
)

function handleSubmit(): void {
  errors.name = form.name.trim() ? '' : t('common.required')
  errors.category = form.category.trim() ? '' : t('common.required')
  errors.location = form.location.trim() ? '' : t('common.required')
  if (errors.name || errors.category || errors.location) {
    return
  }
  emit('submit', { ...form })
}
</script>

<template>
  <form class="flex flex-col gap-4" @submit.prevent="handleSubmit">
    <BaseInput v-model="form.name" :label="t('common.name')" required :error="errors.name" />
    <BaseInput v-model="form.description" :label="t('common.description')" />
    <BaseInput
      v-model="form.category"
      :label="t('common.category')"
      required
      :error="errors.category"
    />
    <BaseInput
      v-model="form.location"
      :label="t('common.location')"
      required
      :error="errors.location"
    />
    <BaseInput v-model.number="form.amount" type="number" :label="t('common.amount')" />
    <div class="flex items-center gap-2">
      <input :id="electricalId" v-model="form.electricalAppliance" type="checkbox" />
      <label :for="electricalId" class="font-bold uppercase tracking-tight">
        {{ t('items.electricalAppliance') }}
      </label>
    </div>
    <BaseInput v-model="form.note" :label="t('common.note')" />

    <ImageUpload
      :label="t('common.uploadImage')"
      :image-url="item?.imagePath ?? null"
      :image-alt="item?.name ?? ''"
      @select="emit('imageSelect', $event)"
      @remove="emit('imageRemove')"
    />

    <fieldset v-if="hasRole(['admin', 'logistics'])" class="border-3 border-brut-black p-4">
      <legend class="px-2 font-extrabold uppercase">{{ t('items.restrictedSection') }}</legend>
      <div class="flex flex-col gap-4">
        <BaseInput v-model="form.manufacturer" :label="t('items.manufacturer')" />
        <BaseInput v-model="form.serialNumber" :label="t('items.serialNumber')" />
        <BaseInput
          v-model="form.lastTechnicalInspectionDate"
          type="date"
          :label="t('items.lastTechnicalInspectionDate')"
        />
        <BaseInput
          v-model="form.lastElectricalInspectionDate"
          type="date"
          :label="t('items.lastElectricalInspectionDate')"
        />
        <BaseInput v-model="form.dateOfPurchase" type="date" :label="t('items.dateOfPurchase')" />
        <BaseInput v-model.number="form.price" type="number" :label="t('items.price')" />
        <BaseInput v-model="form.resolutionNumber" :label="t('items.resolutionNumber')" />
      </div>
    </fieldset>

    <div class="flex gap-2">
      <BaseButton type="submit">{{ t('common.save') }}</BaseButton>
      <BaseButton type="button" variant="secondary" @click="emit('cancel')">
        {{ t('common.cancel') }}
      </BaseButton>
    </div>
  </form>
</template>
