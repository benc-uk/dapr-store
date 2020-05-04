// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Dapr Store frontend - Mock API helper
// ----------------------------------------------------------------------------


const fs = require('fs')

// Load mock data, which we put beside the tests
let mockJson = fs.readFileSync(__dirname+'/../../../tests/unit/mock-data.json')
let mockData = JSON.parse(mockJson)

export default {
  methods: {

    apiProductOffers: function() {
      return {
        data: mockData.products.filter((p) => p.onOffer == true)
      }
    },

    apiProductCatalog: function() {
      return {
        data: mockData.products
      }
    },

    apiProductGet: function(productId) {
      return new Promise((resolve) => {
        resolve({ data: mockData.products.find((p) => p.id == productId) } )
      })
    },

    apiProductSearch: function(query) {
      return {
        data: mockData.products.filter((p) => p.name.includes(query))
      }
    },

    apiOrdersForUser: function(username) {
      return {
        data: mockData.ordersForUser[username]
      }
    },

    apiUserGet: function() {
      return {
        data: mockData.users[0]
      }
    },

    apiOrderGet: function(orderId) {
      return {
        data: mockData.orders.find((o) => o.id == orderId)
      }
    },

    apiCartGet: function(username) {
      return {
        data: mockData.carts[username]
      }
    },

    apiCartAddAmount: function(username, productId, amount) {
      let count = mockData.carts[username][productId]
      count+=amount
      mockData.carts[username][productId] = count
      return {
        data: mockData.carts[username]
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