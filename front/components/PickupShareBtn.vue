<script setup lang="ts">
import AsyncButton from '@/components/ui/button/AsyncButton.vue';
import showDrawer from '@/lib/showDrawer';
import PickupShareDrawer from './Drawer/PickupShareDrawer.vue';
</script>

<template>
    <AsyncButton variant="outline" class="gap-2" @click="async()=>{
        showDrawer({ render: ({ hide }) => h(PickupShareDrawer, { hide }) })
    }">
        <LucideArchive class="size-4" />
        取件
    </AsyncButton>
</template>