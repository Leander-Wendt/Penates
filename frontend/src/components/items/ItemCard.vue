<script setup lang="ts">
import { useI18n } from 'vue-i18n'

import type { Item } from '@/types'

/**
 * Summary card for one item in the item list: thumbnail, name, category,
 * and a link through to its detail page.
 */
defineProps<{ item: Item }>()

const { t } = useI18n()
</script>

<template>
  <RouterLink
    :to="`/items/${item.inventoryNumber}`"
    class="flex flex-col gap-2 border-3 border-brut-black bg-brut-white p-3 no-underline shadow-brut-sm transition-transform hover:-translate-x-0.5 hover:-translate-y-0.5 hover:shadow-brut"
  >
    <img
      v-if="item.imagePath"
      :src="item.imagePath"
      :alt="t('items.thumbnailAlt', { name: item.name })"
      class="h-32 w-full border-3 border-brut-black object-cover"
    />
    <div
      v-else
      class="flex h-32 w-full items-center justify-center border-3 border-brut-black bg-brut-cream font-bold uppercase"
    >
      {{ t('items.noImage') }}
    </div>
    <p class="text-lg font-extrabold">{{ item.name }}</p>
    <p class="text-sm">{{ item.category }}</p>
    <p class="text-sm">{{ item.location }}</p>
  </RouterLink>
</template>
