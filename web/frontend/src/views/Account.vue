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

    <div class="d-flex justify-content-between mb-2">
      <h1>User Account</h1>
      <button class="btn btn-danger" style="height: 3rem" @click="logout"><i class="fa-solid fa-right-from-bracket"></i> &nbsp; LOGOUT</button>
    </div>

    <div v-if="!registeredUser" class="text-center">
      <div class="spinner-border text-success" role="status"><span class="visually-hidden">Not Registered...</span></div>
    </div>

    <div v-if="registeredUser" class="card mb-5">
      <div class="card-body">
        <div class="d-flex justify-content-between mb-2">
          <table class="details">
            <tr>
              <td>Display Name</td>
              <td>{{ registeredUser.displayName }}</td>
            </tr>
            <tr>
              <td>Username</td>
              <td>{{ registeredUser.username }}</td>
            </tr>
          </table>
          <img class="profile d-md-block" :src="photo" />
        </div>
      </div>
    </div>

    <div class="d-flex justify-content-between">
      <h1>Orders</h1>
      <div class="btn btn-success" style="height: 3rem" @click="reloadOrders"><i class="fa-solid fa-rotate"></i> &nbsp; Refresh</div>
    </div>

    <div v-if="!ordersLoaded" class="text-center">
      <div class="spinner-border text-success" role="status">
        <span class="visually-hidden">Not Registered...</span>
      </div>
    </div>

    <order v-for="order in orders" :key="order.id" :order="order" :hideDetailsButton="false" />
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
    order: Order
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
        let resp = await api.userGet(auth.user().username)
        if (resp) {
          this.registeredUser = resp
        }
        graph.getPhoto().then((photo) => {
          if (photo) {
            this.photo = photo
          }
        })
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
        orderList = await api.ordersForUser(auth.user().username)
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
td {
  width: 200px;
}
td + td {
  font-weight: 700;
}
code {
  color: rgb(23, 38, 173);
  font-size: 1.2rem;
}

.details {
  font-size: 140%;
}

.profile {
  width: 10rem;
  border-radius: 50%;
}
</style>
