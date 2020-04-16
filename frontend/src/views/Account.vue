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
      <img class="profile" :src="registeredUser.profileImage">
      Display Name: <b>{{ registeredUser.displayName }}</b>
      <br>
      Username: <b>{{ registeredUser.username }}</b>
      <br><br>
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
    <b-card v-for="order in orders" :key="order.id" class="order">
      <h2>{{ order.title }}</h2>
      <h2 class="float-right text-uppercase">
        {{ order.status }}
      </h2>
      <ul>
        <li>Amount: £{{ order.amount }}</li>
        <li>
          Order Id: <code>{{ order.id }}</code>
        </li>
        <li>Items:</li>
        <ul>
          <li v-for="(item, index) in order.itemsHydrated" :key="index">
            {{ item.name }} &mdash; £{{ item.cost }}
          </li>
        </ul>
      </ul>
    </b-card>
  </div>
</template>

<script>
import { userProfile, msalApp } from '../main'
import api from '../mixins/api'
import User from '../user'
import ErrorBox from '../components/ErrorBox'

export default {
  name: 'Account',

  components: {
    'error-box': ErrorBox
  },

  mixins: [ api ],

  data() {
    return {
      registeredUser: null,
      error: null,
      orders: null,
      ordersLoaded: false
    }
  },

  async created() {
    if (!userProfile.userName) {
      this.$router.replace({ path: '/' })
      return
    }

    try {
      let resp = await this.apiUserGet(userProfile.userName)
      if (resp.data) {
        this.registeredUser = resp.data
      }
    } catch (err) {
      this.error = this.apiDecodeError(err)
    }

    this.reloadOrders()

  },

  methods: {
    async logout() {
      Object.assign(userProfile, new User())
      localStorage.removeItem('user')
      localStorage.removeItem('cart')
      await msalApp.logout()

      this.$router.push({ name: 'home' })
    },

    async reloadOrders() {
      this.ordersLoaded = false
      this.orders = []
      let orderList = []

      try {
        let resp = await this.apiOrdersForUser(userProfile.userName)
        if (resp.data) {
          orderList = resp.data
        }
      } catch (err) {
        this.error = this.apiDecodeError(err)
      }

      // If you have no orders, skip it
      if (!orderList) {
        this.ordersLoaded = true
        return
      }

      // Load orders call the API to fetch details
      for (let orderId of orderList) {
        try {
          let resp = await this.apiOrderGet(orderId)

          if (resp.data) {
            let order = resp.data

            // Items on order are just the product ids, we can rehydrae with full product objects
            order.itemsHydrated = []
            for (let itemId of order.items) {
              let resp = await this.apiProductGet(itemId)
              if (resp.data) {
                order.itemsHydrated.push(resp.data)
              }
            }
            this.orders.push(order)
          }
        } catch (err) {
          this.error += JSON.stringify(this.apiDecodeError(err))+'\n\n'
          continue
        }
      }

      this.ordersLoaded = true
    }
  }
}
</script>

<style scoped>
.details {
  font-size: 140%;
}
.order {
  font-size: 140%;
  margin-bottom: 1rem;
}
.profile {
  float: right;
  width: 10rem;
  border-radius: 50%;
}
</style>