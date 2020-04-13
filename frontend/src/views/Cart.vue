<template>
  <div>
    <h1><fa icon="shopping-cart" /> &nbsp; Shopping Cart</h1><br>
    <b-table striped hover :items="user.cart" :fields="fields" />
    <b-button :disabled="user.cart.length == 0" variant="primary" size="lg">
      <fa icon="shopping-basket" /> &nbsp; CHECKOUT
    </b-button>
  </div>
</template>

<script>
import { userProfile } from '../main'

export default {
  name: 'Cart',

  data() {
    return {
      fields: ['name', 'cost', 'description'],
      user: userProfile
    }
  },
}
</script>

<style>

</style>