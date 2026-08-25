<script setup lang="ts">
/**
 * Neo-brutalist data table primitive. Renders a semantic `<table>` with a
 * caption for assistive technology and thick-bordered rows; per-cell
 * content is supplied via scoped slots keyed by column id.
 */
export interface BaseTableColumn {
  id: string
  label: string
}

withDefaults(
  defineProps<{
    columns: BaseTableColumn[]
    rows: Record<string, unknown>[]
    rowKey: string
    caption?: string
    emptyMessage?: string
  }>(),
  {
    caption: '',
    emptyMessage: '',
  },
)
</script>

<template>
  <div class="overflow-x-auto border-3 border-brut-black">
    <table class="w-full border-collapse text-left">
      <caption v-if="caption" class="sr-only">
        {{
          caption
        }}
      </caption>
      <thead>
        <tr class="border-b-3 border-brut-black bg-brut-yellow">
          <th
            v-for="column in columns"
            :key="column.id"
            scope="col"
            class="border-r-3 border-brut-black px-3 py-2 font-extrabold uppercase tracking-tight last:border-r-0"
          >
            {{ column.label }}
          </th>
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="row in rows"
          :key="String(row[rowKey])"
          class="border-b-3 border-brut-black last:border-b-0 even:bg-brut-cream"
        >
          <td
            v-for="column in columns"
            :key="column.id"
            class="border-r-3 border-brut-black px-3 py-2 last:border-r-0"
          >
            <slot :name="column.id" :row="row">{{ row[column.id] }}</slot>
          </td>
        </tr>
        <tr v-if="rows.length === 0">
          <td :colspan="columns.length" class="px-3 py-4 text-center">
            {{ emptyMessage }}
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
