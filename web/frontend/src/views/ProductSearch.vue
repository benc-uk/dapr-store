<!--
// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Dapr Store frontend - Show product results from a search
// ----------------------------------------------------------------------------
-->

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
