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
    <error-box :error="error" />
    <product-list v-if="!error" :products="products" />
  </div>
</template>

<script>
import api from '../services/api'
import ProductList from '../components/ProductList'
import ErrorBox from '../components/ErrorBox'

export default {
  name: 'ProductOffers',

  components: {
    'product-list': ProductList,
    'error-box': ErrorBox
  },

  data() {
    return {
      products: null,
      error: null,
    }
  },

  async mounted() {
    try {
      this.products = await api.productOffers()
    } catch (err) {
      this.error = err
    }
  },
}
</script>
