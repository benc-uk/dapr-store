<!--
// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Dapr Store frontend - Show a single product
// ----------------------------------------------------------------------------
-->

<template>
  <div>
    <h1>Product Details</h1>
    <error-box :error="error" />

    <div v-if="!product && !error" class="text-center">
      <div class="spinner-border text-success" role="status"><span class="visually-hidden">...</span></div>
    </div>

    <div v-if="product" class="card">
      <div class="card-body">
        <div class="row d-flex">
          <div class="col mb-3">
            <div class="card-title">
              {{ product.name }}
            </div>

            <div v-if="product.onOffer" class="onsale">On Sale</div>

            <div class="m-3">
              {{ product.description }}
            </div>

            <h3>
              £{{ product.cost }}
            </h3>

            <button id="addBut" class="btn btn-primary" :disabled="!isLoggedIn()" @click="addToCart">
              <i class="fa-solid fa-basket-shopping"></i>
              &nbsp; Add to Cart
            </button>
          </div>

          <div class="col">
            <div class="product-img">
              <img :src="product.image" />
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
import ErrorBox from '../components/ErrorBox'
import api from '../services/api'
import auth from '../services/auth'

import { Toast } from 'bootstrap/dist/js/bootstrap.min.js'
let toast

export default {
  name: 'ProductSingle',

  components: {
    'error-box': ErrorBox
  },

  data() {
    return {
      product: null,
      name: null,
      error: null
    }
  },

  watch: {
    product: function(prod) {
      if (prod) {
        this.name = prod.name
      }
    }
  },

  async mounted() {
    try {
      toast = new Toast(this.$refs.addedToast, {
        delay: 2000
      })
      this.product = await api.productGet(this.$route.params.id)
    } catch (err) {
      this.error = err
    }
  },

  methods: {
    async addToCart() {
      try {
        const user = auth.user()
        if (!user) {
          return
        }

        await api.cartAddAmount(user.username, this.product.id, +1)
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
.card-title {
  font-size: 2.5rem;
}
.card-body {
  font-size: 1.1rem;
}
.product-img img {
  width: 100%;
  max-width: 600px;
  border-radius: 1vw;
}
@media (max-width: 768px) {
  .row {
    flex-direction: column;
  }
}
.onsale {
  display: inline-block;
  width: 100%;
  height: 2rem;
  font-size: 1.2rem;
  line-height: 2rem;
  color: rgb(97, 9, 9);
  background-color: rgb(240, 216, 216);
  text-align: center;
  border-radius: 0.3rem;
}
</style>
