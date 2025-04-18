<template>  
    <Form as="" v-slot="slots">
        <slot v-bind="slots" />
    </Form>
</template>

<script setup lang="ts">
import { Form } from 'vee-validate'
</script>