<!--
// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Dapr Store frontend - User / account details view
// ----------------------------------------------------------------------------
-->

<template>
  <div>
    <error-box :error="error" />
    <h1>
      User Account
      <b-button size="lg" variant="danger" class="float-right" @click="logout">
        <fa icon="sign-out-alt" /> &nbsp; LOGOUT
      </b-button>
    </h1>
    <br>

    <div v-if="!registeredUser" class="text-center">
      <b-spinner variant="success" style="width: 5rem; height: 5rem;" />
    </div>
    <b-card v-if="registeredUser" class="details">
      <img class="profile d-none d-md-block" :src="photo">

      Display Name: <b>{{ registeredUser.displayName }}</b>
      <br>
      Username: <b>{{ registeredUser.username }}</b>
    </b-card>

    <br>
    <h1>
      Orders
      <b-button size="lg" variant="success" class="float-right" @click="reloadOrders">
        <fa icon="redo-alt" /> &nbsp; Refresh
      </b-button>
    </h1>

    <div v-if="!ordersLoaded" class="text-center">
      <b-spinner variant="success" style="width: 5rem; height: 5rem;" />
    </div>

    <order v-for="order in orders" :key="order.id" :order="order" />
  </div>
</template>

<script>
import api from '../services/api'
import auth from '../services/auth'
import graph from '../services/graph'
import ErrorBox from '../components/ErrorBox'
import Order from '../components/Order'

export default {
  name: 'Account',

  components: {
    'error-box': ErrorBox,
    'order': Order
  },

  data() {
    return {
      registeredUser: auth.user(),
      photo: 'img/placeholder-profile.jpg',
      error: null,
      orders: null,
      ordersLoaded: false
    }
  },

  async created() {
    try {
      if (auth.user()) {
        let resp = await api.userGet(auth.user().userName)
        if (resp) {
          this.registeredUser = resp
        }
        graph.getPhoto()
          .then((photo) => { if (photo){ this.photo = photo } })
      }
    } catch (err) {
      this.error = err
    }

    this.reloadOrders()
  },

  methods: {
    async logout() {
      await auth.logout()

      this.$router.push({ name: 'home' })
    },

    async reloadOrders() {
      this.ordersLoaded = false
      this.orders = []
      let orderList = []

      try {
        orderList = await api.ordersForUser(auth.user().userName)
      } catch (err) {
        this.error = err
      }

      // If you have no orders, skip it
      if (!orderList) {
        this.ordersLoaded = true
        return
      }

      // Load orders call the API to fetch details
      for (let orderId of orderList.reverse()) {
        try {
          let order = await api.orderGet(orderId)
          this.orders.push(order)
        } catch (err) {
          this.error += err
          continue
        }
      }

      this.ordersLoaded = true
    }
  }
}
</script>

<style scoped>
  code {
    color:rgb(23, 38, 173);
    font-size: 1.2rem;
  }

  .details {
    font-size: 140%;
  }

  .profile {
    float: right;
    width: 10rem;
    border-radius: 50%;
  }
</style>