import VueRouter from 'vue-router'

import ProductCatalog from './views/ProductCatalog'
import ProductOffers from './views/ProductOffers'
import Login from './views/Login'
import Account from './views/Account'
import Home from './views/Home'
//import Error from './views/Error'
import About from './views/About'
import Cart from './views/Cart'

const router = new VueRouter({
  mode: 'history',
  routes: [
    {
      path: '/catalog',
      component: ProductCatalog
    },
    {
      path: '/offers',
      component: ProductOffers
    },
    {
      path: '/',
      component: Home,
      name: 'home'
    },
    {
      path: '/login',
      component: Login
    },
    {
      path: '/account',
      component: Account
    },
    {
      path: '/cart',
      component: Cart
    },
    {
      path: '/about',
      component: About
    },
    // {
    //   path: '/error',
    //   name: 'error',
    //   component: Error,
    //   props: true
    // },
  ]
})

export default router