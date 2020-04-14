<template>
  <div>
    <error-box :error="error" />

    <div v-if="!products" class="text-center">
      <b-spinner variant="success" style="width: 5rem; height: 5rem;" />
    </div>

    <b-card v-for="product in products" :key="product.id" :img-src="product.image" img-height="200" img-width="200" :img-alt="product.name" img-right>
      <b-card-title>
        {{ product.name }}
      </b-card-title>
      <b-card-text>
        <h4>{{ product.description }}</h4>

        £{{ product.cost }}
      </b-card-text>

      <b-button :disabled="!user.token" href="#" variant="primary" @click="addToCart(product)">
        <fa icon="shopping-cart" />
        &nbsp; Add to Cart
      </b-button>
    </b-card>
  </div>
</template>

<script>
import api from '../mixins/api'
import ErrorBox from './ErrorBox'
import { userProfile } from '../main'

export default {
  name: 'ProductList',

  components: {
    'error-box': ErrorBox
  },

  mixins: [ api ],

  props: {
    viewType: {
      type: String,
      required: true
    }
  },

  data() {
    return {
      products: null,
      error: null,
      user: userProfile
    }
  },

  async mounted() {
    try {
      let resp
      if (this.viewType == 'all') {
        resp = await this.apiProductCatalog()
      } else {
        resp = await this.apiProductOffers()
      }
      this.products = resp.data
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
.card {
  margin: 1rem;
}
.card-title {
  font-size: 2rem;
}
</style>