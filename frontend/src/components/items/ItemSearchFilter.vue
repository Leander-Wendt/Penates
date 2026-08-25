<script setup lang="ts">
import { useId } from 'vue'
import { useI18n } from 'vue-i18n'

/**
 * Search box and category dropdown for filtering the item list. Both
 * fields are controlled via v-model-style props/events so the parent can
 * drive the actual API query.
 */
defineProps<{
  searchTerm: string
  categoryFilter: string
  categories: string[]
}>()

const emit = defineEmits<{
  'update:searchTerm': [value: string]
  'update:categoryFilter': [value: string]
}>()

const { t } = useI18n()

const searchId = useId()
const categoryId = useId()
</script>

<template>
  <div class="flex flex-wrap gap-4">
    <div class="flex flex-col gap-1">
      <label :for="searchId" class="font-bold uppercase tracking-tight">{{
        t('common.search')
      }}</label>
      <input
        :id="searchId"
        type="search"
        :value="searchTerm"
        :placeholder="t('items.searchPlaceholder')"
        class="border-3 border-brut-black bg-brut-white px-3 py-2"
        @input="emit('update:searchTerm', ($event.target as HTMLInputElement).value)"
      />
    </div>
    <div class="flex flex-col gap-1">
      <label :for="categoryId" class="font-bold uppercase tracking-tight">
        {{ t('items.categoryFilterLabel') }}
      </label>
      <select
        :id="categoryId"
        :value="categoryFilter"
        class="border-3 border-brut-black bg-brut-white px-3 py-2"
        @change="emit('update:categoryFilter', ($event.target as HTMLSelectElement).value)"
      >
        <option value="">{{ t('items.allCategories') }}</option>
        <option v-for="category in categories" :key="category" :value="category">
          {{ category }}
        </option>
      </select>
    </div>
  </div>
</template>
