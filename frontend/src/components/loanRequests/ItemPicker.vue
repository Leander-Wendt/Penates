<script setup lang="ts">
import { useId } from 'vue'
import { useI18n } from 'vue-i18n'

import type { Item } from '@/types'

/**
 * Multi-select list of items for a loan request, built from native
 * checkboxes so it is fully operable by keyboard (Tab to move, Space to
 * toggle) without any custom widget behaviour.
 */
const props = defineProps<{
  items: Item[]
  modelValue: string[]
}>()

const emit = defineEmits<{ 'update:modelValue': [value: string[]] }>()

const { t } = useI18n()
const legendId = useId()

function isSelected(inventoryNumber: string): boolean {
  return props.modelValue.includes(inventoryNumber)
}

function toggle(inventoryNumber: string, checked: boolean): void {
  if (checked) {
    emit('update:modelValue', [...props.modelValue, inventoryNumber])
  } else {
    emit(
      'update:modelValue',
      props.modelValue.filter((value) => value !== inventoryNumber),
    )
  }
}
</script>

<template>
  <fieldset class="border-3 border-brut-black p-4">
    <legend :id="legendId" class="px-2 font-extrabold uppercase">
      {{ t('loanRequests.itemPickerLabel') }}
    </legend>
    <ul class="flex flex-col gap-2">
      <li v-for="item in items" :key="item.inventoryNumber" class="flex items-center gap-2">
        <input
          :id="`item-picker-${item.inventoryNumber}`"
          type="checkbox"
          :value="item.inventoryNumber"
          :checked="isSelected(item.inventoryNumber)"
          @change="toggle(item.inventoryNumber, ($event.target as HTMLInputElement).checked)"
        />
        <label :for="`item-picker-${item.inventoryNumber}`">
          {{ item.name }} ({{ item.inventoryNumber }})
        </label>
      </li>
    </ul>
  </fieldset>
</template>
