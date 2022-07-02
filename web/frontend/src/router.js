// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Dapr Store frontend - client side routes
// ----------------------------------------------------------------------------

import { createRouter, createWebHistory } from 'vue-router'

import auth from './services/auth'

import ProductCatalog from './views/ProductCatalog'
import ProductOffers from './views/ProductOffers'
import ProductSearch from './views/ProductSearch'
import ProductSingle from './views/ProductSingle'
import ViewOrder from './views/ViewOrder'
import Login from './views/Login'
import Account from './views/Account'
import Home from './views/Home'
import About from './views/About'
import Cart from './views/Cart'

const router = createRouter({
  history: createWebHistory(),
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
    {
      name: 'single-product',
      path: '/product/:id',
      component: ProductSingle
    },
    {
      name: 'view-order',
      path: '/order/:id',
      component: ViewOrder
      //beforeEnter: signedInCheck
    }
  ]
})

function signedInCheck(to, from, next) {
  const user = auth.user()
  if (!user || !user.username) {
    next('/login')
  } else {
    next()
  }
}

export default router
