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
      <div class="spinner-border text-success" role="status"><span class="visually-hidden">...</span></div>
    </div>

    <error-box-list :error="error" />

    <div v-for="product in products" :key="product.id" class="card">
      <div class="card-body">
        <div class="row">
          <div class="col">
            <router-link :to="`/product/` + product.id">
              <div class="card-title">
                {{ product.name }}
              </div>
            </router-link>
            <div class="card-text">
              {{ product.description }}
              <h4 class="mt-4">£{{ product.cost }}</h4>
            </div>

            <button :disabled="!isLoggedIn()" href="#" class="btn btn-primary d-none d-md-inline" @click="addToCart(product)">
              <i class="fa-solid fa-basket-shopping"></i>
              &nbsp; Add to Cart
            </button>
          </div>

          <div class="col flex-grow-0 d-none d-md-block">
            <div class="product-img">
              <span v-if="product.onOffer" class="onsale">On Sale</span>
              <router-link :to="`/product/` + product.id">
                <img :src="product.image" />
              </router-link>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="position-fixed top-0 start-50 translate-middle-x p-3" style="z-index: 11">
      <div ref="addedToast" class="toast hide" role="alert">
        <div class="d-flex">
          <div class="toast-body fs-4">
            {{ name }} was added to your cart!
          </div>
          <button type="button" class="btn-close me-2 m-auto" data-bs-dismiss="toast" aria-label="Close"></button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import api from '../services/api'
import auth from '../services/auth'
import ErrorBox from '../components/ErrorBox'

import { Toast } from 'bootstrap/dist/js/bootstrap.min.js'
let toast

export default {
  name: 'ProductList',

  components: {
    'error-box-list': ErrorBox
  },

  props: {
    products: {
      type: Array,
      default: () => []
    }
  },

  data() {
    return {
      error: null,
      name: ''
    }
  },
  
  async mounted() {
    toast = new Toast(this.$refs.addedToast, {
      delay: 2000
    })
  },

  methods: {
    async addToCart(product) {
      try {
        if (!auth.user()) {
          return
        }

        await api.cartAddAmount(auth.user().username, product.id, +1)
        this.name = product.name
        toast.show()
      } catch (err) {
        this.error = err
      }
    },

    isLoggedIn() {
      if (auth.user()) {
        return true
      }
      return false
    },
  }
}
</script>

<style scoped>
a {
  text-decoration: none;
}
.card {
  margin: 1rem;
}
.card-title {
  font-size: 2rem;
}
.product-img {
  float: right;
  position: relative;
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
  bottom: 0px;
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
