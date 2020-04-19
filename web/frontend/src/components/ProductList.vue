<template>
  <div>
    <error-box :error="error" />

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
    },
    query: {
      type: String,
      default: ''
    }
  },

  data() {
    return {
      products: null,
      error: null,
      user: userProfile
    }
  },

  watch: {
    query: async function (val) {
      try {
        if (this.query !== '') {
          this.products = null
          let resp = await this.apiProductSearch(val)
          this.products = resp.data
        }
      } catch (err) {
        this.error = this.apiDecodeError(err)
      }
    },
  },

  async mounted() {
    try {
      let resp
      if (this.viewType == 'offers') {
        resp = await this.apiProductOffers()
      } else if (this.query !== '') {
        resp = await this.apiProductSearch(this.query)
      } else {
        resp = await this.apiProductCatalog()
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
.product-img {
  float: right;
}
.product-img img {
  width: 12rem;
  border-radius: 0.5rem;
}
</style>