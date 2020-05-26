// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Dapr Store frontend - Mock API helper
// ----------------------------------------------------------------------------


const fs = require('fs')

// Load mock data, which we put beside the tests
let mockDataDir = __dirname+'/../../../../../testing/mock-data'
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
  methods: {

    apiProductOffers: function() {
      return {
        data: mockProducts.filter((p) => p.onOffer == true)
      }
    },

    apiProductCatalog: function() {
      return {
        data: mockProducts
      }
    },

    apiProductGet: function(productId) {
      return new Promise((resolve) => {
        resolve({ data: mockProducts.find((p) => p.id == productId) } )
      })
    },

    apiProductSearch: function(query) {
      return {
        data: mockProducts.filter((p) => p.name.includes(query))
      }
    },

    apiOrdersForUser: function() {
      return {
        data: mockUserOrders
      }
    },

    apiUserGet: function() {
      return {
        data: mockUsers[0]
      }
    },

    apiOrderGet: function(orderId) {
      return {
        data: mockOrders.find((o) => o.id == orderId)
      }
    },

    apiCartGet: function() {
      return {
        data: mockCarts[0]
      }
    },

    apiCartAddAmount: function(_, productId, amount) {
      let count = mockCarts[0][productId]
      count+=amount
      mockCarts[0][productId] = count
      return {
        data: mockCarts[0]
      }
    },

    //
    // Helper to decode error messages
    //
    apiDecodeError(err) {
      if (err.response && err.response.data && err.response.headers['content-type'].includes('json')) {
        return err.response.data
      }
      if (err.response && err.request) {
        return `HTTP ${err.response.status}: API call failed: ${err.request.responseURL}`
      }
      return err.toString()
    }
  }
}