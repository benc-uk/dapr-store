<!--
// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Dapr Store frontend - Show all products in the catalog
// ----------------------------------------------------------------------------
-->

<template>
  <div>
    <h1>All Products</h1>
    <product-list view-type="all" />
  </div>
</template>

<script>
import ProductList from '../components/ProductList'

export default {
  name: 'ProductCatalog',

  components: {
    'product-list': ProductList
  }
}
</script>
