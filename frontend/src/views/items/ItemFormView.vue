<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'

import ItemForm from '@/components/items/ItemForm.vue'
import { deleteItemImage, uploadItemImage } from '@/api/items'
import { useItemsStore } from '@/stores/items'
import type { ItemInput } from '@/types'

/**
 * Create/edit page for an Item. Edit mode is determined by the presence
 * of an `:inventoryNumber` route param. Image upload/removal happens
 * immediately against the dedicated image endpoint and only once the item
 * exists (i.e. in edit mode).
 */
const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const store = useItemsStore()

const inventoryNumber = computed(() =>
  typeof route.params.inventoryNumber === 'string' ? route.params.inventoryNumber : undefined,
)
const currentItem = computed(
  () => store.items.find((item) => item.inventoryNumber === inventoryNumber.value) ?? null,
)
const pageTitle = computed(() =>
  inventoryNumber.value ? t('items.editTitle') : t('items.createTitle'),
)

onMounted(async () => {
  if (inventoryNumber.value) {
    await store.fetchAll()
  }
})

async function handleSubmit(input: ItemInput): Promise<void> {
  if (inventoryNumber.value) {
    await store.update(inventoryNumber.value, input)
    router.push(`/items/${inventoryNumber.value}`)
  } else {
    const created = await store.create(input)
    router.push(`/items/${created.inventoryNumber}`)
  }
}

async function handleImageSelect(file: File): Promise<void> {
  if (!inventoryNumber.value) {
    return
  }
  await uploadItemImage(inventoryNumber.value, file)
  await store.fetchAll()
}

async function handleImageRemove(): Promise<void> {
  if (!inventoryNumber.value) {
    return
  }
  await deleteItemImage(inventoryNumber.value)
  await store.fetchAll()
}

function handleCancel(): void {
  router.push('/items')
}
</script>

<template>
  <div class="flex flex-col gap-4">
    <h1>{{ pageTitle }}</h1>
    <ItemForm
      :item="currentItem"
      @submit="handleSubmit"
      @image-select="handleImageSelect"
      @image-remove="handleImageRemove"
      @cancel="handleCancel"
    />
  </div>
</template>
