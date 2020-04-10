<template>
  <h1>PRODUCT LIST {{ viewType }}</h1>
</template>

<script>
export default {
  name: 'ProductList',

  props: {
    viewType: {
      type: String,
      required: true
    }
  }
}
</script>

<style scoped>

</style>