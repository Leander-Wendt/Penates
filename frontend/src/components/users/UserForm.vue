<script setup lang="ts">
import { reactive, ref, useId, watch } from 'vue'
import { useI18n } from 'vue-i18n'

import BaseButton from '@/components/common/BaseButton.vue'
import BaseInput from '@/components/common/BaseInput.vue'
import type { Organisation, Role, User, UserCreateInput, UserUpdateInput } from '@/types'
import { generateStrongPassword } from '@/utils/generatePassword'

/**
 * Create/edit form for a User: email, password (create only), role, and
 * organisation. Emits a create payload when `user` is absent, otherwise an
 * update payload without a password field.
 */
const props = defineProps<{
  user?: User | null
  organisations: Organisation[]
}>()
const emit = defineEmits<{
  submit: [input: UserCreateInput | UserUpdateInput]
  cancel: []
}>()

const { t } = useI18n()

const roles: Role[] = ['admin', 'logistics', 'student']

const roleSelectId = useId()
const organisationSelectId = useId()

const form = reactive({
  email: props.user?.email ?? '',
  password: '',
  role: props.user?.role ?? 'student',
  organisationId: props.user?.organisationId ?? props.organisations[0]?.id ?? '',
})

const errors = reactive({ email: '', password: '' })

const passwordVisible = ref(false)

/**
 * Fills the password field with a freshly generated strong password and
 * reveals it as plain text so it can be read and shared with the new user.
 */
function generatePassword(): void {
  form.password = generateStrongPassword()
  passwordVisible.value = true
}

watch(
  () => props.user,
  (user) => {
    form.email = user?.email ?? ''
    form.role = user?.role ?? 'student'
    form.organisationId = user?.organisationId ?? props.organisations[0]?.id ?? ''
  },
)

watch(
  () => props.organisations,
  (organisations) => {
    if (!props.user && !form.organisationId) {
      form.organisationId = organisations[0]?.id ?? ''
    }
  },
)

function handleSubmit(): void {
  errors.email = form.email.trim() ? '' : t('common.required')
  errors.password = !props.user && !form.password.trim() ? t('common.required') : ''
  if (errors.email || errors.password) {
    return
  }

  if (props.user) {
    emit('submit', { email: form.email, role: form.role, organisationId: form.organisationId })
  } else {
    emit('submit', {
      email: form.email,
      password: form.password,
      role: form.role,
      organisationId: form.organisationId,
    })
  }
}
</script>

<template>
  <form class="flex flex-col gap-4" @submit.prevent="handleSubmit">
    <BaseInput
      v-model="form.email"
      type="email"
      :label="t('common.email')"
      required
      :error="errors.email"
      autocomplete="email"
    />
    <div v-if="!user" class="flex flex-col gap-2">
      <BaseInput
        v-model="form.password"
        :type="passwordVisible ? 'text' : 'password'"
        :label="t('common.password')"
        required
        :error="errors.password"
        autocomplete="new-password"
      />
      <BaseButton
        type="button"
        variant="secondary"
        data-testid="generate-password"
        class="self-start"
        @click="generatePassword"
      >
        {{ t('users.generatePassword') }}
      </BaseButton>
      <p v-if="passwordVisible" class="text-sm">{{ t('users.generatedPasswordHint') }}</p>
    </div>
    <div class="flex flex-col gap-1">
      <label :for="roleSelectId" class="font-bold uppercase tracking-tight">{{
        t('common.role')
      }}</label>
      <select
        :id="roleSelectId"
        v-model="form.role"
        class="border-3 border-brut-black bg-brut-white px-3 py-2"
      >
        <option v-for="role in roles" :key="role" :value="role">{{ t(`roles.${role}`) }}</option>
      </select>
    </div>
    <div class="flex flex-col gap-1">
      <label :for="organisationSelectId" class="font-bold uppercase tracking-tight">
        {{ t('users.organisation') }}
      </label>
      <select
        :id="organisationSelectId"
        v-model="form.organisationId"
        class="border-3 border-brut-black bg-brut-white px-3 py-2"
      >
        <option v-for="org in organisations" :key="org.id" :value="org.id">{{ org.name }}</option>
      </select>
    </div>
    <div class="flex gap-2">
      <BaseButton type="submit">{{ t('common.save') }}</BaseButton>
      <BaseButton type="button" variant="secondary" @click="emit('cancel')">
        {{ t('common.cancel') }}
      </BaseButton>
    </div>
  </form>
</template>
