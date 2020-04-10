<template>
  <div class="app">
    <b-navbar toggleable="lg" type="dark" variant="primary">
      <b-navbar-brand to="/">
        <img src="./assets/img/logo.svg" class="logo"> <span class="logo-text">Dapr eShop</span>
      </b-navbar-brand>

      <b-navbar-toggle target="nav-collapse" />

      <b-collapse id="nav-collapse" is-nav>
        <b-navbar-nav>
          <b-nav-item to="/catalog" active-class="active">
            <fa icon="shopping-basket" /> &nbsp; Catalog
          </b-nav-item>
          <b-nav-item to="/offers" active-class="active">
            <fa icon="trophy" /> &nbsp; Offers
          </b-nav-item>
        </b-navbar-nav>
      </b-collapse>
      <b-navbar-nav class="ml-auto">
        <b-nav-item to="/orders" active-class="active" variant="info">
          <fa icon="user-circle" /> &nbsp; My Orders
        </b-nav-item>
      </b-navbar-nav>
    </b-navbar>

    <div class="container">
      <router-view />
    </div>
  </div>
</template>

<script>
export default {
  name: 'App',

  data() {
    return {
      version: require('../package.json').version
    }
  },
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
    width: 10rem;
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
</style>
