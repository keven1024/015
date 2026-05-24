<script setup lang="ts">
import { Button } from '@/components/ui/button'
const props = defineProps<{
    title?: string
    showBackButton?: boolean
}>()
const router = useRouter()
</script>

<template>
    <div class="rounded-xl p-5 bg-white/50 backdrop-blur-xl w-full lg:w-200">
        <div v-if="title" class="flex flex-row justify-between">
            <h1 class="text-xl font-normal">{{ title }}</h1>
            <Button
                v-if="!!showBackButton"
                variant="outline"
                class="bg-white/70"
                size="icon"
                @click="
                    () => {
                        router.push('/')
                    }
                "
            >
                <LucideHome class="size-4" />
            </Button>
        </div>
        <slot />
    </div>
</template>
