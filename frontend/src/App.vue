<template>
  <div class="app">
    <b-navbar toggleable="lg" type="dark" variant="primary">
      <b-navbar-brand to="/">
        <img src="./assets/img/logo.svg" class="logo"> <span class="logo-text">Dapr eShop</span>
      </b-navbar-brand>

      <b-navbar-toggle target="nav-collapse" />

      <b-collapse id="nav-collapse" is-nav>
        <!-- <b-navbar-nav>
          <b-nav-item to="/catalog" active-class="active">
            <fa icon="shopping-basket" /> &nbsp; Catalog
          </b-nav-item>
          <b-nav-item to="/offers" active-class="active">
            <fa icon="trophy" /> &nbsp; Offers
          </b-nav-item>
        </b-navbar-nav> -->

        <b-navbar-nav class="ml-auto">
          <b-nav-item v-if="!user.userName" to="/login" variant="info">
            <fa icon="user" /> &nbsp; Login
          </b-nav-item>
          <b-nav-item v-if="user.userName" to="/cart" variant="info" active-class="active">
            <fa icon="shopping-cart" /> &nbsp; Cart
          </b-nav-item>
          <b-nav-item v-if="user.userName" to="/account" variant="info" active-class="active">
            <fa icon="id-card" /> &nbsp; Account
          </b-nav-item>
        </b-navbar-nav>
      </b-collapse>
    </b-navbar>


    <div class="container">
      <router-view />

      <footer>Dapr eShop v{{ version }} - (C) Ben Coleman, 2020</footer>
    </div>
  </div>
</template>

<script>
import { userProfile, msalApp, accessTokenRequest } from './main'
import User from './user'

export default {
  name: 'App',

  data() {
    return {
      user: userProfile,
      version: require('../package.json').version,
    }
  },

  async mounted() {
    // Try to refresh the token for the stored user
    // If it works great, if not we remove the stored local user
    // and the user will need to login again
    let storedUser = localStorage.getItem('user')
    if (storedUser) {
      try {
        let tokenResp = await msalApp.acquireTokenSilent(accessTokenRequest)

        if (tokenResp) {
          Object.assign(userProfile, new User(tokenResp.accessToken, msalApp.getAccount(), msalApp.getAccount().userName || msalApp.getAccount().preferred_username))
          console.log(`### App.vue: MSAL user ${userProfile.userName} is logged & has token`)
          localStorage.setItem('user', userProfile.userName)
          userProfile.cart = []

          try {
            if (localStorage.getItem('cart')) {
              userProfile.cart = JSON.parse(localStorage.getItem('cart'))
            }
          } catch (err) {
            userProfile.cart = []
          }
        } else {
          console.log('### acquireTokenSilent returned no token, removing stored user')
          Object.assign(userProfile, new User())
          localStorage.removeItem('user')
          localStorage.removeItem('cart')
        }
      } catch (err) {
        console.log(`### Error acquireTokenSilent ${err}, removing stored user`)
        Object.assign(userProfile, new User())
        localStorage.removeItem('user')
        localStorage.removeItem('cart')
      }
    }
  }
}
</script>

<style>
  * {
    font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
  }
  body, html, .app {
    margin: 0;
    padding: 0;
    height: 100%
  }
  .logo {
    height: 3rem;
    padding-right: 1rem;
  }
  .logo-text {
    font-family: 'Garamond', 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
    font-size: 2rem;
    padding-right: 3rem;
  }
  .navbar {
    padding: 0.4rem 1rem !important;
    margin-bottom: 2rem;
  }
  .nav-item {
    width: 8rem;
    text-align: center;
    font-size: 1.3rem;
    border-radius: 10px;
    margin-right: 2rem;
    /* padding: 1rem 2rem !important; */
  }
  .nav-item:hover {
    background-color: rgba(10, 10, 60, 0.1);
    border-radius: 10px;
  }
  .active {
    background-color:rgba(255, 255, 255, 0.1);
    border-radius: 10px;
  }
  .card-header {
    font-size: 150% !important;
  }
  footer {
    width: 100%;
    text-align: right;
    border-top: 2px solid lightgray;
    margin-top: 3rem;
  }
</style>
