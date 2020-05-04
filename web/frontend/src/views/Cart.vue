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

    <div v-if="!cart && !error" class="text-center">
      <b-spinner variant="success" style="width: 5rem; height: 5rem;" />
    </div>

    <b-alert v-if="newOrder" show variant="success">
      <h4>Order accepted! 😄</h4>
      <b-link :to="`order/`+newOrder.id">
        <div>
          Your order for <b>£{{ newOrder.amount }}</b> has been accepted, the order ID is: <b>{{ newOrder.id }}</b><br>
          Check your account for progress on your order(s)
        </div>
      </b-link>
    </b-alert>

    <b-alert v-if="cart && Object.keys(cart.products).length==0" show variant="dark">
      Your cart is empty! <br><br><b-link href="/catalog">
        Go shopping!
      </b-link>
    </b-alert>

    <div v-if="cart">
      <h2>Cart total: £{{ total.toFixed(2) }} </h2>
      <b-card v-for="product of cartProducts" :key="product.id" class="m-3 p-1" header-bg-variant="primary" header-text-variant="white">
        <template v-slot:header>
          <span>{{ product.name }}</span>
        </template>
        <h2>Count:</h2> <input :value="cart.products[product.id]" type="text" readonly>

        <b-button class="ml-5 mr-3" :disabled="!cart || cart.products.length == 0" variant="warning" size="lg" @click="addSubProduct(product.id, -1)">
          <fa icon="minus-circle" />
        </b-button>
        <b-button :disabled="!cart || cart.products.length == 0" variant="success" size="lg" @click="addSubProduct(product.id, 1)">
          <fa icon="plus-circle" />
        </b-button>
        <img :src="product.image" class="thumb">
      </b-card>
    </div>

    <b-button :disabled="!cart || cartProducts.length == 0" variant="primary" size="lg" @click="submitOrder">
      <fa icon="shopping-basket" /> &nbsp; CHECKOUT
    </b-button>
    &nbsp;
    <b-button :disabled="!cart || cartProducts.length == 0" variant="warning" size="lg" class="float-right" @click="clearCart">
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
      user: userProfile,
      error: null,
      newOrder: null,
      cart: null,
      cartProducts: []
    }
  },

  computed: {
    total() {
      let tot = 0
      for (let product of this.cartProducts) {
        let count = this.cart.products[product.id]
        tot += (count * product.cost)
      }

      return tot
    }
  },

  async mounted() {
    try {
      let resp = await this.apiCartGet(userProfile.userName)
      if (resp.data) {
        this.cart = resp.data
        this.cartProducts = []
        for (let productId in this.cart.products) {
          // Do this async helps speed it up when running locally
          this.apiProductGet(productId)
            .then((resp) => {
              if (resp.data) {
                this.cartProducts.push(resp.data)
              }
            })
        }
      }
    } catch (err) {
      this.error = this.apiDecodeError(err)
    }
  },

  methods: {
    async submitOrder() {
      try {
        let resp = await this.apiCartSubmit(userProfile.userName)
        this.newOrder = resp.data
        resp = await this.apiCartClear(userProfile.userName)
        this.cart = resp.data
        this.cartProducts = []
      } catch (err) {
        this.error = this.apiDecodeError(err)
      }
    },

    async clearCart() {
      try {
        let resp = await this.apiCartClear(userProfile.userName)
        this.cart = resp.data
        this.cartProducts = []
      } catch (err) {
        this.error = this.apiDecodeError(err)
      }
    },

    async addSubProduct(productId, amount) {
      try {
        let resp = await this.apiCartAddAmount(userProfile.userName, productId, amount)

        this.cart = resp.data
        // Fiddly nonsense to remove from cartProducts if removed from products.cart
        if (!Object.prototype.hasOwnProperty.call(this.cart.products, productId)) {
          this.cartProducts = this.cartProducts.filter((p) => p.id != productId)
        }
      } catch (err) {
        this.error = this.apiDecodeError(err)
      }
    }
  }
}
</script>

<style scoped>
  input[type=text] {
    border: none;
    width: 4rem;
    text-align: center;
  }
  .thumb {
    width: 80px;
    float: right;
  }
</style>