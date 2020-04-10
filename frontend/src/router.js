import VueRouter from 'vue-router'

import ProductCatalog from './components/ProductCatalog'
import ProductOffers from './components/ProductOffers'
import Home from './components/Home'

const router = new VueRouter({
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
      component: Home
    }
  ]
})

export default router