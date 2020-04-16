<template>
  <div>
    <h1>On Sale</h1>
    <product-list view-type="offers" />
  </div>
</template>

<script>
import ProductList from '../components/ProductList'

export default {
  name: 'ProductOffers',

  components: {
    'product-list': ProductList
  }
}
</script>

<style scoped>

</style>
