<template>
  <div>
    <h1>Special Offers!</h1>
    <product-list view-type="offers" />
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
