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
    <h1><i class="fa-solid fa-cart-shopping"></i> &nbsp; Shopping Cart</h1>
    <br />
    <error-box :error="error" />

    <div v-if="!cart && !error" class="text-center">
      <div class="spinner-border text-success" role="status"><span class="visually-hidden">...</span></div>
    </div>

    <div v-if="newOrder" class="alert alert-success">
      <h4>Order accepted! 😄</h4>
      <router-link :to="`order/` + newOrder.id">
        <div>
          Your order for <b>£{{ newOrder.amount }}</b> has been accepted, the order ID is: <b>{{ newOrder.id }}</b
          ><br />
          Check your account for progress on your order(s)
        </div>
      </router-link>
    </div>

    <div v-if="cart && Object.keys(cart.products).length == 0" class="alert alert-dark">
      Your cart is empty! <br /><br /><router-link to="/catalog"> Go shopping! </router-link>
    </div>

    <div v-if="cart">
      <h2>Cart total: £{{ total.toFixed(2) }}</h2>
      <div v-for="product of cartProducts" :key="product.id" class="card m-3 p-1" header-bg-variant="primary" header-text-variant="white">
        <div class="card-header">{{ product.name }}</div>
        <div class="card-body">
          <h2>Count:</h2>
          <input :value="cart.products[product.id]" type="text" readonly class="m-3" />

          <button class="btn btn-sm btn-warning ml-5 mr-3" :disabled="!cart || cart.products.length == 0" @click="modifyProductAmmount(product.id, -1)">
            <i class="fa-solid fa-circle-minus"></i>
          </button>
          <button class="btn btn-sm btn-success ml-5 mr-3" :disabled="!cart || cart.products.length == 0" @click="modifyProductAmmount(product.id, 1)">
            <i class="fa-solid fa-circle-plus"></i>
          </button>
          <img :src="product.image" class="thumb" />
        </div>
      </div>
    </div>

    <button class="btn btn-lg btn-primary ml-5 mr-3" :disabled="!cart || cartProducts.length == 0" @click="submitOrder">
      <i class="fa-solid fa-basket-shopping"></i> &nbsp; CHECKOUT
    </button>
    &nbsp;
    <button class="btn btn-lg btn-warning ml-5 mr-3 float-end" :disabled="!cart || cartProducts.length == 0" @click="clearCart">
      <i class="fa-solid fa-trash-can"></i> &nbsp; EMPTY CART
    </button>
  </div>
</template>

<script>
import ErrorBox from '../components/ErrorBox'
import api from '../services/api'
import auth from '../services/auth'

export default {
  name: 'Cart',

  components: {
    'error-box': ErrorBox
  },

  data() {
    return {
      error: null,
      newOrder: null,
      cart: null,
      cartProducts: [],
      user: null
    }
  },

  computed: {
    total() {
      let tot = 0
      for (let product of this.cartProducts) {
        let count = this.cart.products[product.id]
        tot += count * product.cost
      }

      return tot
    }
  },

  async mounted() {
    try {
      this.user = auth.user()
      if (!this.user) {
        return
      }

      let resp = await api.cartGet(this.user.username)
      if (resp) {
        this.cart = resp
        this.cartProducts = []
        for (let productId in this.cart.products) {
          // Do this async helps speed it up when running locally, due to Dapr bug
          api.productGet(productId).then((resp) => {
            this.cartProducts.push(resp)
          })
        }
      }
    } catch (err) {
      this.error = err
    }
  },

  methods: {
    async submitOrder() {
      try {
        this.newOrder = await api.cartSubmit(this.user.username)
        this.cart = await api.cartClear(this.user.username)
        this.cartProducts = []
      } catch (err) {
        this.error = err
      }
    },

    async clearCart() {
      try {
        let resp = await api.cartClear(this.user.username)
        this.cart = resp
        this.cartProducts = []
      } catch (err) {
        this.error = err
      }
    },

    async modifyProductAmmount(productId, amount) {
      try {
        this.cart = await api.cartAddAmount(this.user.username, productId, amount)

        // Fiddly nonsense to remove from cartProducts if removed from products.cart
        // Check if productId is removed from cart object, then recreate cartProducts array
        if (!Object.prototype.hasOwnProperty.call(this.cart.products, productId)) {
          this.cartProducts = this.cartProducts.filter((p) => p.id != productId)
        }
      } catch (err) {
        this.error = err
      }
    }
  }
}
</script>

<style scoped>
input[type='text'] {
  border: none;
  width: 4rem;
  text-align: center;
}
.thumb {
  width: 80px;
  float: right;
}
</style>
