<template>
    <div class="flex flex-row bg-white/10 backdrop-blur-xl p-2 rounded-full">
        <div v-for="item in routes" :key="item.key" 
        :class="cx('flex flex-row text-sm px-4 py-2 font-bold rounded-full cursor-pointer', 
            item?.key === type && 'bg-black/10'
        )"
        @click="()=>{
            router.push({
                query: {
                    ...route.query,
                    type: item.key
                }
            })
        }"
        >
          {{ item.name }}
        </div>
    </div>
</template>

<script setup lang="ts">
import { cx } from 'class-variance-authority'
const routes = [
    { name: '文件', key: 'file' },
    { name: '文本', key: 'text' }
]
const route = useRoute()
const router = useRouter()
const type = computed(() => route?.query?.type)
</script>