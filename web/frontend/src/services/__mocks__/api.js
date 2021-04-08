// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Dapr Store frontend - Mock API helper
// ----------------------------------------------------------------------------

const fs = require('fs')

// Load mock data, which we put beside the tests
let mockDataDir = __dirname + '/../../../../../testing/mock-data'
let mockJson = fs.readFileSync(`${mockDataDir}/carts.json`)
let mockCarts = JSON.parse(mockJson)
mockJson = fs.readFileSync(`${mockDataDir}/orders.json`)
let mockOrders = JSON.parse(mockJson)
mockJson = fs.readFileSync(`${mockDataDir}/products.json`)
let mockProducts = JSON.parse(mockJson)
mockJson = fs.readFileSync(`${mockDataDir}/user-orders.json`)
let mockUserOrders = JSON.parse(mockJson)
mockJson = fs.readFileSync(`${mockDataDir}/users.json`)
let mockUsers = JSON.parse(mockJson)

export default {
  productOffers() {
    return mockProducts.filter((p) => p.onOffer == true)
  },

  productCatalog: function () {
    return mockProducts
  },

  productGet(productId) {
    return new Promise((resolve) => {
      resolve(mockProducts.find((p) => p.id == productId))
    })
  },

  productSearch(query) {
    return mockProducts.filter((p) => p.name.includes(query))
  },

  ordersForUser: function () {
    return mockUserOrders
  },

  userGet() {
    return mockUsers[0]
  },

  orderGet(orderId) {
    return mockOrders.find((o) => o.id == orderId)
  },

  cartGet() {
    return mockCarts[0]
  },

  cartAddAmount(_, productId, amount) {
    let count = mockCarts[0][productId]
    count += amount
    mockCarts[0][productId] = count
    return mockCarts[0]
  }
}
