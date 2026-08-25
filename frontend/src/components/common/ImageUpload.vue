<script setup lang="ts">
import { computed, ref, useId } from 'vue'

/**
 * File picker with an image preview. Used for item photos: shows the
 * currently stored image (if any) or a live preview of a newly selected
 * file, and lets the user replace or remove it.
 */
const props = withDefaults(
  defineProps<{
    imageUrl?: string | null
    label: string
    imageAlt?: string
  }>(),
  {
    imageUrl: null,
    imageAlt: '',
  },
)

const emit = defineEmits<{ select: [file: File]; remove: [] }>()

const inputId = useId()
const localPreview = ref<string | null>(null)

const previewUrl = computed(() => localPreview.value ?? props.imageUrl)

function handleChange(event: Event): void {
  const file = (event.target as HTMLInputElement).files?.[0]
  if (!file) {
    return
  }
  localPreview.value = URL.createObjectURL(file)
  emit('select', file)
}

function handleRemove(): void {
  localPreview.value = null
  emit('remove')
}
</script>

<template>
  <div class="flex flex-col gap-2">
    <label :for="inputId" class="font-bold uppercase tracking-tight">{{ label }}</label>
    <img
      v-if="previewUrl"
      :src="previewUrl"
      :alt="imageAlt"
      class="h-40 w-40 border-3 border-brut-black object-cover"
    />
    <input :id="inputId" type="file" accept="image/*" @change="handleChange" />
    <button
      v-if="previewUrl"
      type="button"
      class="w-fit border-3 border-brut-black bg-brut-white px-3 py-1 font-bold"
      @click="handleRemove"
    >
      {{ $t('common.removeImage') }}
    </button>
  </div>
</template>
