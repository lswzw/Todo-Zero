<template>
  <div class="locale-float">
    <button class="locale-float-btn" :title="t('common.switchLanguage')" @click="handleToggle">
      {{ currentLocaleAlias }}
    </button>
  </div>
</template>

<script setup lang="ts">
import { useLocale } from '@/composables/useLocale'

const { currentLocaleAlias, toggleLocale, t } = useLocale()

function handleToggle() {
  toggleLocale()
}
</script>

<style scoped>
.locale-float {
  position: fixed;
  right: 24px;
  bottom: 24px;
  z-index: var(--z-fixed);
}

.locale-float-btn {
  width: 52px;
  height: 52px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea, #764ba2);
  border: none;
  color: #fff;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  box-shadow: 0 4px 20px rgba(102, 126, 234, 0.4);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  display: flex;
  align-items: center;
  justify-content: center;
}

.locale-float-btn:hover {
  transform: scale(1.1);
  box-shadow: 0 8px 30px rgba(102, 126, 234, 0.5);
}

.locale-float-btn:active {
  transform: scale(0.95);
}

@media (max-width: 768px) {
  .locale-float {
    right: 16px;
    bottom: 16px;
  }

  .locale-float-btn {
    width: 48px;
    height: 48px;
    font-size: 13px;
  }
}
</style>
