<!--
// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Dapr Store frontend - Cart contents view
// ----------------------------------------------------------------------------
-->

<template>
  <div>
    <h1>
      <fa icon="shopping-cart" /> &nbsp; Shopping Cart
    </h1>
    <br>

    <error-box :error="error" />
    <b-alert v-if="newOrder" show variant="success">
      <h4>Order accepted! 😄</h4>
      <div>
        Your order for <b>£{{ newOrder.amount }}</b> has been accepted, the order ID is: <b>{{ newOrder.id }}</b><br>
        Check your account for progress on your order(s)
      </div>
    </b-alert>
    <b-table striped hover :items="user.cart" :fields="fields" />

    <b-button :disabled="user.cart.length == 0" variant="primary" size="lg" @click="submitOrder">
      <fa icon="shopping-basket" /> &nbsp; CHECKOUT
    </b-button>
    &nbsp;
    <b-button :disabled="user.cart.length == 0" variant="warning" size="lg" class="float-right" @click="clearCart">
      <fa icon="trash-alt" /> &nbsp; EMPTY CART
    </b-button>
  </div>
</template>

<script>
import { userProfile } from '../main'
import ErrorBox from '../components/ErrorBox'
import api from '../mixins/api'

export default {
  name: 'Cart',

  components: {
    'error-box': ErrorBox
  },

  mixins: [ api ],

  data() {
    return {
      fields: ['name', 'cost', 'description'],
      user: userProfile,
      error: null,
      newOrder: null
    }
  },

  methods: {
    async submitOrder() {
      try {
        let order = {
          forUser: userProfile.userName,
          amount: userProfile.cart.reduce((total, p) => { return total + parseFloat(p.cost) }, 0),
          items: userProfile.cart.map((p) => p.id),
          title: new Date().toLocaleString()
        }

        let resp = await this.apiOrderSubmit(order)
        if (resp.data) {
          this.newOrder = resp.data
          userProfile.cart = []
          localStorage.setItem('cart', JSON.stringify(userProfile.cart))
        } else {
          throw new Error('Some sort of problem!')
        }
      } catch (err) {
        this.error = this.apiDecodeError(err)
      }
    },

    clearCart() {
      userProfile.cart = []
      localStorage.setItem('cart', JSON.stringify(userProfile.cart))
    }
  }
}
</script>
