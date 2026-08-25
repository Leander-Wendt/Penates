<script setup lang="ts">
import { computed, useId } from 'vue'

/**
 * Neo-brutalist labeled text input. Wires validation errors to
 * `aria-describedby` and announces them via an `aria-live` region.
 */
withDefaults(
  defineProps<{
    modelValue: string | number
    label: string
    type?: string
    error?: string
    required?: boolean
    placeholder?: string
    autocomplete?: string
  }>(),
  {
    type: 'text',
    error: '',
    required: false,
    placeholder: '',
    autocomplete: undefined,
  },
)

defineEmits<{ 'update:modelValue': [value: string] }>()

const inputId = useId()
const errorId = computed(() => `${inputId}-error`)
</script>

<template>
  <div class="flex flex-col gap-1">
    <label :for="inputId" class="font-bold uppercase tracking-tight">
      {{ label }}
      <span v-if="required" aria-hidden="true" class="text-brut-red">*</span>
    </label>
    <input
      :id="inputId"
      :type="type"
      :value="modelValue"
      :placeholder="placeholder"
      :required="required"
      :autocomplete="autocomplete"
      :aria-describedby="error ? errorId : undefined"
      :aria-invalid="error ? 'true' : undefined"
      class="border-3 border-brut-black bg-brut-white px-3 py-2 text-brut-black focus-visible:outline-none disabled:opacity-50"
      @input="$emit('update:modelValue', ($event.target as HTMLInputElement).value)"
    />
    <p
      v-if="error"
      :id="errorId"
      role="alert"
      aria-live="assertive"
      class="font-bold text-brut-red"
    >
      {{ error }}
    </p>
  </div>
</template>
