<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'

import BaseButton from '@/components/common/BaseButton.vue'
import ItemCard from '@/components/items/ItemCard.vue'
import ItemSearchFilter from '@/components/items/ItemSearchFilter.vue'
import { useAuth } from '@/composables/useAuth'
import { useItemsStore } from '@/stores/items'

/**
 * Item list page: search/category filtering, a thumbnail grid linking to
 * item detail pages, and a create button for Admin/Logistics.
 */
const { t } = useI18n()
const { hasRole } = useAuth()
const store = useItemsStore()
const router = useRouter()

const categories = computed(() =>
  Array.from(new Set(store.items.map((item) => item.category).filter(Boolean))),
)

let debounceHandle: ReturnType<typeof setTimeout> | undefined

function updateSearchTerm(value: string): void {
  store.searchTerm = value
  if (debounceHandle) {
    clearTimeout(debounceHandle)
  }
  debounceHandle = setTimeout(() => store.fetchAll(), 300)
}

function updateCategoryFilter(value: string): void {
  store.categoryFilter = value
  store.fetchAll()
}

onMounted(() => {
  store.fetchAll()
})

onBeforeUnmount(() => {
  if (debounceHandle) {
    clearTimeout(debounceHandle)
  }
})

function goToCreate(): void {
  router.push('/items/new')
}
</script>

<template>
  <div class="flex flex-col gap-6">
    <div class="flex items-center justify-between">
      <h1>{{ t('items.title') }}</h1>
      <BaseButton v-if="hasRole(['admin', 'logistics'])" @click="goToCreate">
        {{ t('items.createButton') }}
      </BaseButton>
    </div>

    <ItemSearchFilter
      :search-term="store.searchTerm"
      :category-filter="store.categoryFilter"
      :categories="categories"
      @update:search-term="updateSearchTerm"
      @update:category-filter="updateCategoryFilter"
    />

    <p v-if="store.items.length === 0">{{ t('common.noResults') }}</p>
    <div v-else class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
      <ItemCard v-for="item in store.items" :key="item.inventoryNumber" :item="item" />
    </div>
  </div>
</template>
