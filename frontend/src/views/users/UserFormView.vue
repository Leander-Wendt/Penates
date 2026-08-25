<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'

import UserForm from '@/components/users/UserForm.vue'
import { useOrganisationsStore } from '@/stores/organisations'
import { useUsersStore } from '@/stores/users'
import type { UserCreateInput, UserUpdateInput } from '@/types'

/**
 * Create/edit page for a User. Edit mode is determined by the presence of
 * an `:id` route param.
 */
const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const usersStore = useUsersStore()
const organisationsStore = useOrganisationsStore()

const userId = computed(() => (typeof route.params.id === 'string' ? route.params.id : undefined))
const currentUser = computed(
  () => usersStore.users.find((user) => user.id === userId.value) ?? null,
)
const pageTitle = computed(() => (userId.value ? t('users.editTitle') : t('users.createTitle')))

onMounted(async () => {
  await organisationsStore.fetchAll()
  if (userId.value) {
    await usersStore.fetchAll()
  }
})

async function handleSubmit(input: UserCreateInput | UserUpdateInput): Promise<void> {
  if (userId.value) {
    await usersStore.update(userId.value, input as UserUpdateInput)
  } else {
    await usersStore.create(input as UserCreateInput)
  }
  router.push('/users')
}

function handleCancel(): void {
  router.push('/users')
}
</script>

<template>
  <div class="flex flex-col gap-4">
    <h1>{{ pageTitle }}</h1>
    <UserForm
      :user="currentUser"
      :organisations="organisationsStore.organisations"
      @submit="handleSubmit"
      @cancel="handleCancel"
    />
  </div>
</template>
