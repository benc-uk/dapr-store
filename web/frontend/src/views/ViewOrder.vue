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
    <div class="d-flex justify-content-between mb-2">
      <h1>View Order Status</h1>
      <!--button class="btn btn-lg btn-success" @click="loadOrder()"><i class="fa-solid fa-rotate"></i> &nbsp; Refresh</button-->
    </div>

    <error-box :error="error" />

    <div v-if="!order && !error" class="text-center">
      <div class="spinner-border text-success" role="status"><span class="visually-hidden">...</span></div>
    </div>

    <order v-if="order" :order="order" />
  </div>
</template>

<script>
import ErrorBox from '../components/ErrorBox'
import api from '../services/api'
import Order from '../components/Order'
var timerId = null

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

  async mounted() {
    this.loadOrder()

    // Refresh in the background every 5 seconds
    timerId = setInterval(this.loadOrder, 5000) 
  },

  unmounted() {
    clearInterval(timerId)
  },

  methods: {
    async loadOrder() {
      try {
        this.order = await api.orderGet(this.$route.params.id)
      } catch (err) {
        this.error = err
      }
    }
  }
}
</script>
