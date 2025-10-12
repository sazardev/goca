<template>
  <span :class="['custom-badge', `badge-${type}`]">
    <Icon v-if="icon" :icon="icon" class="badge-icon" />
    <slot />
  </span>
</template>

<script setup>
import { Icon } from '@iconify/vue'

defineProps({
  type: {
    type: String,
    default: 'info', // info, success, warning, error, tip
  },
  icon: {
    type: String,
    default: null,
  },
})
</script>

<style scoped>
.custom-badge {
  display: inline-flex;
  align-items: center;
  gap: 0.375rem;
  padding: 0.25rem 0.75rem;
  border-radius: 0.5rem;
  font-size: 0.875rem;
  font-weight: 500;
  line-height: 1.5;
  transition: all 0.2s ease;
}

.badge-icon {
  width: 1rem;
  height: 1rem;
  flex-shrink: 0;
}

.badge-info {
  background-color: rgba(56, 189, 248, 0.1);
  color: rgb(14, 165, 233);
  border: 1px solid rgba(56, 189, 248, 0.2);
}

.badge-success {
  background-color: rgba(34, 197, 94, 0.1);
  color: rgb(22, 163, 74);
  border: 1px solid rgba(34, 197, 94, 0.2);
}

.badge-warning {
  background-color: rgba(251, 191, 36, 0.1);
  color: rgb(245, 158, 11);
  border: 1px solid rgba(251, 191, 36, 0.2);
}

.badge-error {
  background-color: rgba(239, 68, 68, 0.1);
  color: rgb(220, 38, 38);
  border: 1px solid rgba(239, 68, 68, 0.2);
}

.badge-tip {
  background-color: rgba(168, 85, 247, 0.1);
  color: rgb(147, 51, 234);
  border: 1px solid rgba(168, 85, 247, 0.2);
}

.dark .badge-info {
  background-color: rgba(56, 189, 248, 0.15);
  color: rgb(125, 211, 252);
}

.dark .badge-success {
  background-color: rgba(34, 197, 94, 0.15);
  color: rgb(134, 239, 172);
}

.dark .badge-warning {
  background-color: rgba(251, 191, 36, 0.15);
  color: rgb(253, 224, 71);
}

.dark .badge-error {
  background-color: rgba(239, 68, 68, 0.15);
  color: rgb(252, 165, 165);
}

.dark .badge-tip {
  background-color: rgba(168, 85, 247, 0.15);
  color: rgb(196, 181, 253);
}
</style>
