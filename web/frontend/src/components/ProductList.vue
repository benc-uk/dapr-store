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
      <b-spinner variant="success" style="width: 5rem; height: 5rem;" />
    </div>

    <b-card v-for="product in products" :key="product.id">
      <b-row>
        <b-col>
          <b-link :to="`/product/`+product.id">
            <b-card-title>
              {{ product.name }}
            </b-card-title>
          </b-link>
          <b-card-text>
            <h4>{{ product.description }}</h4>

            £{{ product.cost }}
          </b-card-text>

          <b-button :disabled="!user.userName" href="#" variant="primary" class="d-none d-md-inline" @click="addToCart(product)">
            <fa icon="shopping-cart" />
            &nbsp; Add to Cart
          </b-button>
        </b-col>

        <b-col class="flex-grow-0 d-none d-md-block">
          <div class="product-img">
            <b-link :to="`/product/`+product.id">
              <img :src="product.image">
            </b-link>
          </div>
        </b-col>
      </b-row>
    </b-card>
  </div>
</template>

<script>
// import api from '../mixins/api'
import { userProfile } from '../main'
import api from '../mixins/api'

export default {
  name: 'ProductList',

  mixins: [ api ],

  props: {
    products: {
      type: Array,
      default: () => []
    }
  },

  data() {
    return {
      user: userProfile
    }
  },

  methods: {
    async addToCart(product) {
      try {
        await this.apiCartAddAmount(userProfile.userName, product.id, +1)
        this.showToast('Added to your cart!', 'success', product)
      } catch (err) {
        this.showToast('Error adding to cart 😫 '+err.toString(), 'danger', product)
      }
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
</style>