<!--
// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Dapr Store frontend - Reusable component that lists products
// ----------------------------------------------------------------------------
-->

<template>
  <div>
    <div v-if="!products" class="text-center">
      <b-spinner variant="success" style="width: 5rem; height: 5rem" />
    </div>

    <b-card v-for="product in products" :key="product.id">
      <b-row>
        <b-col>
          <b-link :to="`/product/` + product.id">
            <b-card-title>
              {{ product.name }}
            </b-card-title>
          </b-link>
          <b-card-text>
            {{ product.description }}
            <br /><br />
            <h4>£{{ product.cost }}</h4>
          </b-card-text>

          <b-button :disabled="!isLoggedIn()" href="#" variant="primary" class="d-none d-md-inline" @click="addToCart(product)">
            <fa icon="shopping-cart" />
            &nbsp; Add to Cart
          </b-button>
        </b-col>

        <b-col class="flex-grow-0 d-none d-md-block">
          <div class="product-img">
            <span v-if="product.onOffer" class="onsale">On Sale</span>
            <b-link :to="`/product/` + product.id">
              <img :src="product.image" />
            </b-link>
          </div>
        </b-col>
      </b-row>
    </b-card>
  </div>
</template>

<script>
import api from '../services/api'
import auth from '../services/auth'

export default {
  name: 'ProductList',

  props: {
    products: {
      type: Array,
      default: () => []
    }
  },

  methods: {
    async addToCart(product) {
      try {
        if (!auth.user()) {
          return
        }

        await api.cartAddAmount(auth.user().username, product.id, +1)
        this.showToast('Added to your cart!', 'success', product)
      } catch (err) {
        this.showToast('Error adding to cart 😫 ' + err.toString(), 'danger', product)
      }
    },

    isLoggedIn() {
      if (auth.user()) {
        return true
      }
      return false
    },

    showToast(msg, variant, product) {
      this.$bvToast.toast(`${product.name}`, {
        title: msg,
        variant: variant,
        autoHideDelay: 3000,
        appendToast: true,
        toaster: 'b-toaster-top-center',
        solid: true
      })
    }
  }
}
</script>

<style scoped>
.card {
  margin: 1rem;
}
.card-title {
  font-size: 2rem;
}
.product-img {
  float: right;
}
.product-img img {
  width: 12rem;
  border-radius: 0.5rem;
}
.card-text {
  font-size: 130%;
}
.onsale {
  display: inline-block;
  position: absolute;
  bottom: 0;
  width: 12rem;
  height: 1.5rem;
  font-size: 1rem;
  color: white;
  background-color: rgb(185, 17, 17);
  text-align: center;
  border-bottom-left-radius: 0.5rem;
  border-bottom-right-radius: 0.5rem;
}
</style>
