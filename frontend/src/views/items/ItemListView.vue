<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'

import BaseButton from '@/components/common/BaseButton.vue'
import ItemCard from '@/components/items/ItemCard.vue'
import ItemListRow from '@/components/items/ItemListRow.vue'
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

// View mode toggle - list view is now the default
const viewMode = ref<'card' | 'list'>('list')

// Selected items for loan request
const selectedItems = ref<string[]>([])

function toggleItemSelection(inventoryNumber: string, checked: boolean): void {
  if (checked) {
    selectedItems.value = [...selectedItems.value, inventoryNumber]
  } else {
    selectedItems.value = selectedItems.value.filter((item) => item !== inventoryNumber)
  }
}

function isItemSelected(inventoryNumber: string): boolean {
  return selectedItems.value.includes(inventoryNumber)
}

function createLoanRequest(): void {
  // Navigate to loan request creation with selected items pre-filled
  router.push({
    path: '/loan-requests/new',
    query: { items: selectedItems.value.join(',') },
  })
}

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
      <div class="flex gap-2">
        <BaseButton
          v-if="hasRole(['student']) && selectedItems.length > 0"
          @click="createLoanRequest"
        >
          {{ t('loanRequests.createButton') }}
        </BaseButton>
        <BaseButton v-if="hasRole(['admin', 'logistics'])" @click="goToCreate">
          {{ t('items.createButton') }}
        </BaseButton>
      </div>
    </div>

    <div class="flex flex-wrap items-center justify-between gap-4">
      <ItemSearchFilter
        :search-term="store.searchTerm"
        :category-filter="store.categoryFilter"
        :categories="categories"
        @update:search-term="updateSearchTerm"
        @update:category-filter="updateCategoryFilter"
      />
      <div class="flex items-center gap-2">
        <span class="font-bold uppercase tracking-tight">
          {{ t('items.viewToggleLabel') }}
        </span>
        <button
          type="button"
          class="flex items-center gap-2 border-3 border-brut-black px-4 py-2 font-bold uppercase tracking-tight shadow-brut transition-transform duration-100 ease-out hover:-translate-x-0.5 hover:-translate-y-0.5 hover:shadow-brut-lg active:translate-x-0.5 active:translate-y-0.5 active:shadow-brut-pressed disabled:cursor-not-allowed disabled:opacity-50 disabled:hover:translate-x-0 disabled:hover:translate-y-0 disabled:hover:shadow-brut"
          :class="viewMode === 'card' ? 'bg-brut-yellow' : 'bg-brut-white'"
          @click="() => (viewMode = 'card')"
        >
          <svg
            width="16"
            height="16"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="3"
            stroke-linecap="round"
            stroke-linejoin="round"
          >
            <rect x="3" y="3" width="18" height="18" rx="2" />
            <rect x="7" y="7" width="10" height="10" rx="1" />
          </svg>
          {{ t('items.cardView') }}
        </button>
        <button
          type="button"
          class="flex items-center gap-2 border-3 border-brut-black px-4 py-2 font-bold uppercase tracking-tight shadow-brut transition-transform duration-100 ease-out hover:-translate-x-0.5 hover:-translate-y-0.5 hover:shadow-brut-lg active:translate-x-0.5 active:translate-y-0.5 active:shadow-brut-pressed disabled:cursor-not-allowed disabled:opacity-50 disabled:hover:translate-x-0 disabled:hover:translate-y-0 disabled:hover:shadow-brut"
          :class="viewMode === 'list' ? 'bg-brut-yellow' : 'bg-brut-white'"
          @click="() => (viewMode = 'list')"
        >
          <svg
            width="16"
            height="16"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="3"
            stroke-linecap="round"
            stroke-linejoin="round"
          >
            <line x1="8" y1="6" x2="21" y2="6" />
            <line x1="8" y1="12" x2="21" y2="12" />
            <line x1="8" y1="18" x2="21" y2="18" />
            <line x1="3" y1="6" x2="3.01" y2="6" />
            <line x1="3" y1="12" x2="3.01" y2="12" />
            <line x1="3" y1="18" x2="3.01" y2="18" />
          </svg>
          {{ t('items.listView') }}
        </button>
      </div>
    </div>

    <p v-if="store.items.length === 0">{{ t('common.noResults') }}</p>
    <div
      v-else-if="viewMode === 'card'"
      class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3"
    >
      <ItemCard v-for="item in store.items" :key="item.inventoryNumber" :item="item" />
    </div>
    <div v-else class="border-3 border-brut-black">
      <div
        class="grid grid-cols-7 gap-3 border-b-3 border-brut-black bg-brut-yellow px-3 py-2 font-extrabold uppercase tracking-tight"
      >
        <p>{{ t('common.select') }}</p>
        <p>{{ t('items.inventoryNumber') }}</p>
        <p>{{ t('common.name') }}</p>
        <p>{{ t('common.category') }}</p>
        <p>{{ t('common.location') }}</p>
        <p>{{ t('common.amount') }}</p>
        <p>{{ t('common.description') }}</p>
      </div>
      <ItemListRow
        v-for="item in store.items"
        :key="item.inventoryNumber"
        :item="item"
        :is-selected="isItemSelected(item.inventoryNumber)"
        @toggle-selection="toggleItemSelection"
      />
    </div>
  </div>
</template>
