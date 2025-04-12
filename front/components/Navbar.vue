<template>
    <div class="flex flex-row bg-white/10 backdrop-blur-xl p-2 rounded-full">
        <div v-for="item in routes" :key="item.key" 
        class="flex flex-row text-sm px-4 py-2 font-bold rounded-full cursor-pointer relative"
        @click="()=>{
            router.push({
                query: {
                    ...route.query,
                    type: item.key
                }
            })
        }"
        >
          <motion.div v-if="item?.key === type" layoutId="navbar-active" class="absolute inset-0 rounded-full w-full h-full bg-black/10"/>
          {{ item.name }}
        </div>
    </div>
</template>

<script setup lang="ts">
import { cx } from 'class-variance-authority'
import { motion } from "motion-v"
const routes = [
    { name: '文件', key: 'file' },
    { name: '文本', key: 'text' }
]
const route = useRoute()
const router = useRouter()
const type = computed(() => route?.query?.type)
</script>