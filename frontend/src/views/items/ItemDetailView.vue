<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'

import BaseButton from '@/components/common/BaseButton.vue'
import { useAuth } from '@/composables/useAuth'
import { useItemsStore } from '@/stores/items'
import { formatDate } from '@/utils/formatDate'

/**
 * Read-only detail page for a single item. Restricted fields are only
 * rendered for Admin/Logistics; for Students the backend never sends them
 * in the first place, so they are simply absent from `item`.
 */
const props = defineProps<{ inventoryNumber: string }>()

const { t } = useI18n()
const { hasRole } = useAuth()
const store = useItemsStore()
const route = useRoute()
const router = useRouter()

const item = computed(() =>
  store.items.find((candidate) => candidate.inventoryNumber === props.inventoryNumber),
)

onMounted(async () => {
  if (!item.value) {
    await store.fetchAll()
  }
})

function goToEdit(): void {
  router.push(`/items/${route.params.inventoryNumber}/edit`)
}
</script>

<template>
  <div v-if="item" class="flex flex-col gap-6">
    <div class="flex items-center justify-between">
      <h1>{{ item.name }}</h1>
      <BaseButton v-if="hasRole(['admin', 'logistics'])" @click="goToEdit">
        {{ t('common.edit') }}
      </BaseButton>
    </div>

    <img
      v-if="item.imagePath"
      :src="item.imagePath"
      :alt="t('items.thumbnailAlt', { name: item.name })"
      class="h-48 w-48 border-3 border-brut-black object-cover"
    />

    <dl class="grid grid-cols-1 gap-x-6 gap-y-2 sm:grid-cols-2">
      <div>
        <dt class="font-bold uppercase">{{ t('items.inventoryNumber') }}</dt>
        <dd>{{ item.inventoryNumber }}</dd>
      </div>
      <div>
        <dt class="font-bold uppercase">{{ t('common.category') }}</dt>
        <dd>{{ item.category }}</dd>
      </div>
      <div>
        <dt class="font-bold uppercase">{{ t('common.location') }}</dt>
        <dd>{{ item.location }}</dd>
      </div>
      <div>
        <dt class="font-bold uppercase">{{ t('common.amount') }}</dt>
        <dd>{{ item.amount }}</dd>
      </div>
      <div>
        <dt class="font-bold uppercase">{{ t('items.electricalAppliance') }}</dt>
        <dd>{{ item.electricalAppliance ? t('common.yes') : t('common.no') }}</dd>
      </div>
      <div class="sm:col-span-2">
        <dt class="font-bold uppercase">{{ t('common.description') }}</dt>
        <dd>{{ item.description }}</dd>
      </div>
      <div class="sm:col-span-2">
        <dt class="font-bold uppercase">{{ t('common.note') }}</dt>
        <dd>{{ item.note }}</dd>
      </div>
    </dl>

    <section v-if="hasRole(['admin', 'logistics'])" class="border-3 border-brut-black p-4">
      <h2>{{ t('items.restrictedSection') }}</h2>
      <dl class="grid grid-cols-1 gap-x-6 gap-y-2 sm:grid-cols-2">
        <div>
          <dt class="font-bold uppercase">{{ t('items.manufacturer') }}</dt>
          <dd>{{ item.manufacturer }}</dd>
        </div>
        <div>
          <dt class="font-bold uppercase">{{ t('items.serialNumber') }}</dt>
          <dd>{{ item.serialNumber }}</dd>
        </div>
        <div>
          <dt class="font-bold uppercase">{{ t('items.lastTechnicalInspectionDate') }}</dt>
          <dd>{{ formatDate(item.lastTechnicalInspectionDate) }}</dd>
        </div>
        <div>
          <dt class="font-bold uppercase">{{ t('items.lastElectricalInspectionDate') }}</dt>
          <dd>{{ formatDate(item.lastElectricalInspectionDate) }}</dd>
        </div>
        <div>
          <dt class="font-bold uppercase">{{ t('items.dateOfPurchase') }}</dt>
          <dd>{{ formatDate(item.dateOfPurchase) }}</dd>
        </div>
        <div>
          <dt class="font-bold uppercase">{{ t('items.price') }}</dt>
          <dd>{{ item.price }}</dd>
        </div>
        <div>
          <dt class="font-bold uppercase">{{ t('items.resolutionNumber') }}</dt>
          <dd>{{ item.resolutionNumber }}</dd>
        </div>
      </dl>
    </section>
  </div>
</template>
