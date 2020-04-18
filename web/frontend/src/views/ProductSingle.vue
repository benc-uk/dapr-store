<template>
  <div>
    <h1>Product Details</h1>
    <error-box :error="error" />

    <div v-if="!product" class="text-center">
      <b-spinner variant="success" style="width: 5rem; height: 5rem;" />
    </div>

    <b-card v-if="product">
      <b-row class="d-flex">
        <b-col class="mb-3">
          <b-card-title>
            {{ product.name }}
          </b-card-title>
          <br>
          {{ product.description }}
          <br><br>
          £{{ product.cost }}
          <br><br>
          <b-button :disabled="!user.userName" href="#" variant="primary" @click="addToCart(product)">
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
import { userProfile } from '../main'

export default {
  name: 'ProductSingle',

  components: {
    'error-box': ErrorBox
  },

  mixins: [ api ],

  data() {
    return {
      product: null,
      error: null,
      user: userProfile
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
    addToCart(product) {
      userProfile.cart.push(product)
      localStorage.setItem('cart', JSON.stringify(userProfile.cart))

      this.$bvToast.toast(`${product.name}`, {
        title: 'Added to your cart!',
        variant: 'success',
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
</style>