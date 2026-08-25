<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'

import BaseButton from '@/components/common/BaseButton.vue'
import BaseTable from '@/components/common/BaseTable.vue'
import { useAuth } from '@/composables/useAuth'
import { useOrganisationsStore } from '@/stores/organisations'
import { useUsersStore } from '@/stores/users'
import type { User } from '@/types'

/**
 * Admin/Logistics-only user management list. The delete action is only
 * shown to Admins; Logistics can create and edit but not delete.
 */
const { t } = useI18n()
const { hasRole } = useAuth()
const usersStore = useUsersStore()
const organisationsStore = useOrganisationsStore()
const router = useRouter()

const columns = computed(() => [
  { id: 'email', label: t('common.email') },
  { id: 'role', label: t('common.role') },
  { id: 'organisation', label: t('users.organisation') },
  { id: 'actions', label: t('common.actions') },
])

const rows = computed(() =>
  usersStore.users.map((user) => ({
    ...user,
    organisation:
      organisationsStore.organisations.find((org) => org.id === user.organisationId)?.name ?? '',
  })),
)

onMounted(async () => {
  await Promise.all([usersStore.fetchAll(), organisationsStore.fetchAll()])
})

function goToCreate(): void {
  router.push('/users/new')
}

function goToEdit(user: User): void {
  router.push(`/users/${user.id}/edit`)
}

async function handleDelete(user: User): Promise<void> {
  if (window.confirm(t('users.deleteConfirm'))) {
    await usersStore.remove(user.id)
  }
}
</script>

<template>
  <div class="flex flex-col gap-4">
    <div class="flex items-center justify-between">
      <h1>{{ t('users.title') }}</h1>
      <BaseButton @click="goToCreate">{{ t('users.createButton') }}</BaseButton>
    </div>

    <BaseTable
      :columns="columns"
      :rows="rows"
      row-key="id"
      :caption="t('users.title')"
      :empty-message="t('common.noResults')"
    >
      <template #role="{ row }">
        {{ t(`roles.${(row as User).role}`) }}
      </template>
      <template #actions="{ row }">
        <div class="flex gap-2">
          <BaseButton variant="secondary" @click="goToEdit(row as User)">
            {{ t('common.edit') }}
          </BaseButton>
          <BaseButton v-if="hasRole(['admin'])" variant="danger" @click="handleDelete(row as User)">
            {{ t('common.delete') }}
          </BaseButton>
        </div>
      </template>
    </BaseTable>
  </div>
</template>
