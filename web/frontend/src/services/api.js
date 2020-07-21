// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Dapr Store frontend - API helper to call the various backend microservices
// ----------------------------------------------------------------------------

import auth from './auth'

const API_SCOPE = 'store-api'
let accessToken
let clientId
let apiEndpoint

export default {
  //
  // Must be called at startup
  //
  configure(appApiEndpoint, appClientId) {
    clientId = appClientId
    apiEndpoint = appApiEndpoint
  },

  //
  // ===== Users =====
  //
  userRegister: function(user) {
    return apiCall('v1.0/invoke/users/method/register', 'POST', user)
  },

  userGet: function(username) {
    return apiCall(`v1.0/invoke/users/method/get/${username}`)
  },

  userCheckReg: function(username) {
    return apiCall(`v1.0/invoke/users/method/isregistered/${username}`)
  },

  //
  // ===== Products =====
  //
  productCatalog: function() {
    return apiCall('v1.0/invoke/products/method/catalog')
  },

  productOffers: function() {
    return apiCall('v1.0/invoke/products/method/offers')
  },

  productGet: function(productId) {
    return apiCall(`v1.0/invoke/products/method/get/${productId}`)
  },

  productSearch: function(query) {
    return apiCall(`v1.0/invoke/products/method/search/${query}`)
  },

  //
  // ===== Cart =====
  //
  cartProductSet: function(username, productId, count) {
    return apiCall(`v1.0/invoke/cart/method/setProduct/${username}/${productId}/${count}`, 'PUT')
  },

  cartGet: function(username) {
    return apiCall(`v1.0/invoke/cart/method/get/${username}`, 'GET')
  },

  cartSubmit: function(username) {
    return apiCall('v1.0/invoke/cart/method/submit', 'POST', `${username}`)
  },

  cartClear: function(username) {
    return apiCall(`v1.0/invoke/cart/method/clear/${username}`, 'PUT', `${username}`)
  },

  cartAddAmount: async function(username, productId, amount) {
    let count = 0

    let cartResp = await this.cartGet(username)
    if (cartResp) {
      let cart = cartResp
      let productCount = cart.products[productId]
      if (productCount) {
        count = productCount + amount
      } else {
        count = amount
      }
      if (count < 0) { count = 0 }
    }

    return this.cartProductSet(username, productId, count)
  },

  //
  // ===== Orders =====
  //
  orderGet: function(orderId) {
    return apiCall(`v1.0/invoke/orders/method/get/${orderId}`)
  },

  ordersForUser: function(username) {
    return apiCall(`v1.0/invoke/orders/method/getForUser/${username}`)
  }
}

//
// ===== Base fetch wrapper =====
//
async function apiCall(apiPath, method = 'get', data = null) {
  let headers = {}
  let url = `${(apiEndpoint || '/')}${apiPath}`
  console.log(`### API CALL ${method} ${url}`)

  // Only get a token if logged in & using real auth (i.e AUTH_CLIENT_ID set)
  if (auth.user() && clientId) {
    const scopes = [ `api://${clientId}/${API_SCOPE}` ]

    // Try to get an access token with our API scope
    if (!accessToken) {
      accessToken = await auth.acquireToken(scopes)
    }

    // Send token as per the OAuth 2.0 bearer token scheme
    if (accessToken) {
      headers = {
        'Authorization': `Bearer ${accessToken}`
      }
    }
  }

  // Build request
  const request = {
    method,
    headers,
  }

  // Add payload if required
  if (data) {
    request.body = JSON.stringify(data)
  }

  // Make the HTTP request
  let resp = await fetch(url, request)

  // Decode error message when non-HTTP OK (200~299) is received
  if (!resp.ok) {
    let error = `API call to ${url} failed with ${resp.status} ${resp.statusText}`
    if (resp.headers.get('Content-Type') === 'application/json') {
      error = ''
      let errorObj = await resp.json()
      for (const [key, value] of Object.entries(errorObj)) {
        error += `${key}: '${value}', `
      }
    }
    throw new Error(error)
  }

  // Attempt to return response body as data object if JSON
  if (resp.headers.get('Content-Type') === 'application/json') {
    return resp.json()
  } else {
    return resp.text()
  }
}