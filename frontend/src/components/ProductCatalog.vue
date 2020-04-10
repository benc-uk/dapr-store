<template>
  <div>
    <h1>All Products</h1>
    <product-list view-type="main" />
  </div>
</template>

<script>
import ProductList from './ProductList'

export default {
  name: 'ProductCatalog',

  components: {
    'product-list': ProductList
  }
}
</script>

<style scoped>

</style>
