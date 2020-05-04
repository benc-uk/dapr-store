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
      <b-button size="lg" variant="success" class="float-right" @click="loadOrder($route.params.id)">
        <fa icon="redo-alt" /> &nbsp; Refresh
      </b-button>
    </h1>

    <error-box :error="error" />

    <div v-if="!order && !error" class="text-center">
      <b-spinner variant="success" style="width: 5rem; height: 5rem;" />
    </div>

    <order v-if="order" :order="order" />
  </div>
</template>

<script>
import ErrorBox from '../components/ErrorBox'
import api from '../mixins/api'
import { userProfile } from '../main'
import Order from '../components/Order'

export default {
  name: 'ViewOrder',

  components: {
    'error-box': ErrorBox,
    'order': Order
  },

  mixins: [ api ],

  data() {
    return {
      order: null,
      error: null,
      user: userProfile
    }
  },

  mounted() {
    this.loadOrder(this.$route.params.id)
  },

  methods: {
    async loadOrder(id) {
      this.order = null
      try {
        let resp = await this.apiOrderGet(id)
        if (resp.data) {
          this.order = resp.data
        }
      } catch (err) {
        this.error = this.apiDecodeError(err)
      }
    }
  }
}
</script>

<style scoped>

</style>