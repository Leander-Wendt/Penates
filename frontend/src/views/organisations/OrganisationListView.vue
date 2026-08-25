<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'

import BaseButton from '@/components/common/BaseButton.vue'
import BaseModal from '@/components/common/BaseModal.vue'
import BaseTable from '@/components/common/BaseTable.vue'
import OrganisationForm from '@/components/organisations/OrganisationForm.vue'
import { useOrganisationsStore } from '@/stores/organisations'
import type { Organisation, OrganisationInput } from '@/types'

/**
 * Admin/Logistics-only organisation management: list, create, edit and
 * delete, all driven through a shared modal form.
 */
const { t } = useI18n()
const store = useOrganisationsStore()

const isModalOpen = ref(false)
const editingOrganisation = ref<Organisation | null>(null)

const columns = computed(() => [
  { id: 'name', label: t('common.name') },
  { id: 'description', label: t('common.description') },
  { id: 'actions', label: t('common.actions') },
])

const modalTitle = computed(() =>
  editingOrganisation.value ? t('organisations.editTitle') : t('organisations.createTitle'),
)

onMounted(() => {
  store.fetchAll()
})

function openCreateModal(): void {
  editingOrganisation.value = null
  isModalOpen.value = true
}

function openEditModal(organisation: Organisation): void {
  editingOrganisation.value = organisation
  isModalOpen.value = true
}

function closeModal(): void {
  isModalOpen.value = false
}

async function handleSubmit(input: OrganisationInput): Promise<void> {
  if (editingOrganisation.value) {
    await store.update(editingOrganisation.value.id, input)
  } else {
    await store.create(input)
  }
  closeModal()
}

async function handleDelete(organisation: Organisation): Promise<void> {
  if (window.confirm(t('organisations.deleteConfirm'))) {
    await store.remove(organisation.id)
  }
}
</script>

<template>
  <div class="flex flex-col gap-4">
    <div class="flex items-center justify-between">
      <h1>{{ t('organisations.title') }}</h1>
      <BaseButton @click="openCreateModal">{{ t('organisations.createButton') }}</BaseButton>
    </div>

    <BaseTable
      :columns="columns"
      :rows="store.organisations"
      row-key="id"
      :caption="t('organisations.title')"
      :empty-message="t('common.noResults')"
    >
      <template #actions="{ row }">
        <div class="flex gap-2">
          <BaseButton variant="secondary" @click="openEditModal(row as Organisation)">
            {{ t('common.edit') }}
          </BaseButton>
          <BaseButton variant="danger" @click="handleDelete(row as Organisation)">
            {{ t('common.delete') }}
          </BaseButton>
        </div>
      </template>
    </BaseTable>

    <BaseModal :open="isModalOpen" :title="modalTitle" @close="closeModal">
      <OrganisationForm
        :organisation="editingOrganisation"
        @submit="handleSubmit"
        @cancel="closeModal"
      />
    </BaseModal>
  </div>
</template>
