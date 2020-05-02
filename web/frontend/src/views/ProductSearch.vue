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
    <error-box :error="error" />
    <product-list v-if="!error" :products="products" />
  </div>
</template>

<script>
import api from '../mixins/api'
import ProductList from '../components/ProductList'
import ErrorBox from '../components/ErrorBox'

export default {
  name: 'ProductSearch',

  components: {
    'product-list': ProductList,
    'error-box': ErrorBox
  },

  mixins: [ api ],

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
        let resp = await this.apiProductSearch(this.$route.params.query)
        this.products = resp.data
      } catch (err) {
        this.error = this.apiDecodeError(err)
      }
    }
  }
}
</script>
