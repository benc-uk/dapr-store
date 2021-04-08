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
    <h1>
      Search results for: <i>{{ $route.params.query }}</i>
    </h1>
    <error-box :error="error" />
    <product-list v-if="!error" :products="products" />
  </div>
</template>

<script>
import api from '../services/api'
import ProductList from '../components/ProductList'
import ErrorBox from '../components/ErrorBox'

export default {
  name: 'ProductSearch',

  components: {
    'product-list': ProductList,
    'error-box': ErrorBox
  },

  data() {
    return {
      products: null,
      error: null
    }
  },

  watch: {
    async $route() {
      this.doSearch()
    }
  },

  async mounted() {
    this.doSearch()
  },

  methods: {
    async doSearch() {
      try {
        this.products = null
        this.products = await api.productSearch(this.$route.params.query)
      } catch (err) {
        this.error = err
      }
    }
  }
}
</script>
