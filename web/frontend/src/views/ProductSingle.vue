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
      <b-spinner variant="success" style="width: 5rem; height: 5rem;" />
    </div>

    <b-card v-if="product">
      <b-row class="d-flex">
        <b-col class="mb-3">
          <b-card-title>
            {{ product.name }}
          </b-card-title>
          <div v-if="product.onOffer" class="onsale">
            On Sale
          </div>
          <br><br><br>
          {{ product.description }}
          <br><br>
          £{{ product.cost }}
          <br><br>
          <b-button id="addBut" :disabled="!user().userName" variant="primary" @click="addToCart">
            <fa icon="shopping-cart" />
            &nbsp; Add to Cart
          </b-button>
        </b-col>

        <b-col>
          <div class="product-img">
            <img :src="product.image">
          </div>
        </b-col>
      </b-row>
    </b-card>
  </div>
</template>

<script>
import ErrorBox from '../components/ErrorBox'
import api from '../mixins/api'
import auth from '../mixins/auth'

export default {
  name: 'ProductSingle',

  components: {
    'error-box': ErrorBox
  },

  mixins: [ api, auth ],

  data() {
    return {
      product: null,
      error: null
    }
  },

  async mounted() {
    try {
      let resp = await this.apiProductGet(this.$route.params.id)
      if (resp.data) {
        this.product = resp.data
      }
    } catch (err) {
      this.error = this.apiDecodeError(err)
    }
  },

  methods: {
    async addToCart() {
      try {
        await this.apiCartAddAmount(this.user().userName, this.product.id, +1)
        this.showToast('Added to your cart!', 'success')
      } catch (err) {
        this.showToast('Error adding to cart 😫 '+err.toString(), 'danger')
      }
    },

    showToast(msg, variant) {
      console.log(msg, variant)

      this.$bvToast.toast(`${this.product.name}`, {
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
@media  (max-width: 768px) {
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