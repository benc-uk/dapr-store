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
    <error-box :error="error" />
    <product-list v-if="!error" :products="products" />
  </div>
</template>

<script>
import api from '../services/api'
import ProductList from '../components/ProductList'
import ErrorBox from '../components/ErrorBox'

export default {
  name: 'ProductCatalog',

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
      let resp = await api.productCatalog()
      if (resp && typeof resp === 'object') {
        this.products = resp
      } else {
        throw new Error('Failed to fetch products')
      }
    } catch (err) {
      this.error = err
    }
  },
}
</script>
