<!--
// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Dapr Store frontend - Show a single order
// ----------------------------------------------------------------------------
-->

<template>
  <div>
    <h1>
      View Order Status
      <button class="btn btn-lg btn-success float-end" @click="loadOrder($route.params.id)">
        <i class="fa-solid fa-rotate"></i> &nbsp; Refresh
      </button>
    </h1>

    <error-box :error="error" />

    <div v-if="!order && !error" class="text-center">
      <div class="spinner-border text-success" role="status"><span class="visually-hidden">Not Registered...</span></div>
    </div>

    <order v-if="order" :order="order" />
  </div>
</template>

<script>
import ErrorBox from '../components/ErrorBox'
import api from '../services/api'
import Order from '../components/Order'

export default {
  name: 'ViewOrder',

  components: {
    'error-box': ErrorBox,
    order: Order
  },

  data() {
    return {
      order: null,
      error: null
    }
  },

  mounted() {
    this.loadOrder(this.$route.params.id)
  },

  methods: {
    async loadOrder(id) {
      this.order = null
      try {
        this.order = await api.orderGet(id)
      } catch (err) {
        this.error = err
      }
    }
  }
}
</script>
