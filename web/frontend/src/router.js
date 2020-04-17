import VueRouter from 'vue-router'

import { userProfile } from './main'

import ProductCatalog from './views/ProductCatalog'
import ProductOffers from './views/ProductOffers'
import ProductSearch from './views/ProductSearch'
import Login from './views/Login'
import Account from './views/Account'
import Home from './views/Home'
import About from './views/About'
import Cart from './views/Cart'

const router = new VueRouter({
  mode: 'history',
  routes: [
    {
      name: 'home',
      path: '/',
      component: Home
    },
    {
      name: 'login',
      path: '/login',
      component: Login
    },
    {
      name: 'account',
      path: '/account',
      component: Account,
      beforeEnter: signedInCheck
    },
    {
      name: 'cart',
      path: '/cart',
      component: Cart,
      beforeEnter: signedInCheck
    },
    {
      name: 'about',
      path: '/about',
      component: About
    },
    {
      name: 'search',
      path: '/search/:query',
      component: ProductSearch
    },
    {
      name: 'catalog',
      path: '/catalog',
      component: ProductCatalog
    },
    {
      name: 'offers',
      path: '/offers',
      component: ProductOffers
    },
  ]
})

function signedInCheck(to, from, next) {
  if (!userProfile.userName) {
    next('/login')
  } else {
    next()
  }
}

export default router