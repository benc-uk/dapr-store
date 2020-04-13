<template>
  <div>
    <h1>
      User Account
      <b-button size="lg" variant="danger" class="float-right" @click="logout">
        <fa icon="sign-out-alt" /> &nbsp; LOGOUT
      </b-button>
    </h1>

    <table v-if="user.account && user.userName" class="table">
      <tr><td>Display Name</td><td>{{ user.account.name }}</td></tr>
      <tr><td>Username</td><td>{{ user.userName }}</td></tr>
      <tr><td>Unique ID</td><td>{{ user.account.accountIdentifier }}</td></tr>
    </table>

    <div />
  </div>
</template>

<script>
import { userProfile } from '../main'
import User from '../user'

export default {
  name: 'Account',

  data() {
    return {
      user: userProfile
    }
  },

  created() {
    if (!userProfile.token) {
      this.$router.replace({ path: '/' })
    }
  },

  methods: {
    async logout() {
      Object.assign(userProfile, new User())
      localStorage.removeItem('user')
      localStorage.removeItem('cart')
      //await msalApp.logout()

      this.$router.push({ name: 'home' })
    }
  }
}
</script>

<style scoped>
.table {
  font-size: 140%;
}
</style>