<template>
  <div class="app">
    <nav class="navbar navbar-expand-lg navbar-dark bg-primary">
      <div class="container-fluid">
        <router-link class="navbar-brand" to="/">
          <img src="./assets/img/logo.svg" class="logo" />
          <span class="logo-text">Dapr eShop</span>
        </router-link>
        <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navContent">
          <span class="navbar-toggler-icon"></span>
        </button>
        <div id="navContent" class="collapse navbar-collapse">
          <div class="d-flex">
            <input v-model="query" class="form-control me-2" type="search" placeholder="Search products" @keyup.enter="search" />
            <button class="btn btn-success" @click="search"><i class="fa-solid fa-magnifying-glass"></i></button>
          </div>
          <div class="filler" style="flex-grow: 1"></div>
          <ul class="navbar-nav ml-auto">
            <li v-if="!user" class="nav-item btn-success">
              <router-link class="nav-link" to="/login"><i class="fa-solid fa-right-to-bracket"></i> Login</router-link>
            </li>
            <template v-else>
              <li class="nav-item btn-success">
                <router-link class="nav-link" to="/cart"><i class="fa-solid fa-cart-shopping"></i> Cart</router-link>
              </li>
              <li class="nav-item btn-success">
                <router-link class="nav-link" to="/account"><i class="fa-solid fa-user"></i> Account</router-link>
              </li>
            </template>
          </ul>
        </div>
      </div>
    </nav>

    <div class="container">
      <!-- Views are injected here -->
      <router-view @loginComplete="refreshUser" />

      <footer>Dapr eShop v{{ version }} - (C) Ben Coleman, 2020</footer>
    </div>
  </div>
</template>

<script>
import auth from './services/auth'

export default {
  name: 'App',

  data() {
    return {
      version: require('../package.json').version,
      query: '',
      user: {}
    }
  },

  async created() {
    // Restore any cached or saved local user
    this.refreshUser()
  },

  methods: {
    search() {
      if (this.query) {
        this.$router.push({ name: 'search', params: { query: this.query } }).catch(() => {})
      }
    },

    refreshUser() {
      this.user = auth.user()
    }
  }
}
</script>

<style>
* {
  font-family: 'Josefin Sans', 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
}

body,
html,
.app {
  margin: 0;
  padding: 0;
  height: 100%;
}

.logo {
  height: 3rem;
  padding-right: 1rem;
  vertical-align: middle;
}

.logo-text {
  font-family: 'Josefin Sans', 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
  font-size: 2.3rem;
  padding-right: 1rem;
  line-height: 1.7rem;
  vertical-align: middle;
}

@media (max-width: 500px) {
  .logo-text {
    font-size: 1rem;
  }
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
}

.nav-item:hover {
  background-color: rgba(10, 10, 60, 0.1);
  border-radius: 10px;
}

.active {
  background-color: rgba(255, 255, 255, 0.1);
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

.alert {
  font-size: 140% !important;
}

.alert h4 {
  color: #222;
  padding-bottom: 0.3rem;
  border-bottom: 2px solid rgba(0, 0, 0, 0.2);
}

input[type='text'] {
  background-color: rgba(10, 10, 60, 0.1);
  padding: 0px 0px 0px 10px !important;
  border-radius: 5px !important;
  font-size: 1.2rem !important;
  height: 45px;
}

input[type='text']:focus {
  background-color: rgba(200, 200, 255, 0.2);
}
</style>
