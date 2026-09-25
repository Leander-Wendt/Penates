<script setup lang="ts">
import { useId } from 'vue'

import type { Item } from '@/types'

/**
 * Table row for one item in the item list view: displays item details
 * in a tabular format for the list view with a checkbox for selection.
 */
const props = defineProps<{
  item: Item
  isSelected: boolean
}>()

const emit = defineEmits<{
  'toggle-selection': [inventoryNumber: string, checked: boolean]
}>()

const checkboxId = useId()

function handleCheckboxChange(checked: boolean): void {
  emit('toggle-selection', props.item.inventoryNumber, checked)
}
</script>

<template>
  <div
    class="grid grid-cols-7 gap-3 border-b-3 border-brut-black px-3 py-2 last:border-b-0 hover:bg-brut-cream transition-colors"
  >
    <div class="flex h-full items-center justify-center">
      <input
        :id="checkboxId"
        type="checkbox"
        :checked="isSelected"
        :aria-label="`Select ${item.name}`"
        class="h-5 w-5 border-3 border-brut-black bg-brut-white appearance-none checked:bg-brut-yellow checked:border-brut-black focus:outline-none focus:ring-0 transition-colors cursor-pointer"
        @change="handleCheckboxChange(($event.target as HTMLInputElement).checked)"
      />
    </div>
    <RouterLink :to="`/items/${item.inventoryNumber}`" class="font-extrabold hover:underline">
      {{ item.inventoryNumber }}
    </RouterLink>
    <RouterLink :to="`/items/${item.inventoryNumber}`" class="font-extrabold hover:underline">
      {{ item.name }}
    </RouterLink>
    <p>{{ item.category }}</p>
    <p>{{ item.location }}</p>
    <p>{{ item.amount }}</p>
    <p class="truncate" :title="item.description">{{ item.description || '-' }}</p>
  </div>
</template>
