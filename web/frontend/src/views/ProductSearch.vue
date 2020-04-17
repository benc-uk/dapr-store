<template>
  <div>
    <h1>Search results for: <i>{{ this.$route.params.query }}</i></h1>
    <product-list view-type="search" :query="this.$route.params.query" />
  </div>
</template>

<script>
import ProductList from '../components/ProductList'

export default {
  name: 'ProductSearch',

  components: {
    'product-list': ProductList
  }
}
</script>
