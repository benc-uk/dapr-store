<!--
// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Dapr Store frontend - Show products that are on offer from the catalog
// ----------------------------------------------------------------------------
-->

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
