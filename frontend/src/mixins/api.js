/* eslint-disable no-console */
import { userProfile } from '../main'
import axios from 'axios'

export default {
  methods: {
    //
    // ===== Users =====
    //
    apiUserRegister: function(user) {
      return this._apiRawCall('v1.0/invoke/users/method/register', 'POST', user)
    },

    apiUserGet: function(username) {
      return this._apiRawCall(`v1.0/invoke/users/method/get/${username}`)
    },

    //
    // ===== Products =====
    //
    apiProductCatalog: function() {
      return this._apiRawCall('v1.0/invoke/products/method/catalog')
    },

    apiProductOffers: function() {
      return this._apiRawCall('v1.0/invoke/products/method/offers')
    },




    _apiRawCall: function(apiPath, method = 'get', data = null) {
      let apiUrl = `/${apiPath}`
      console.log(`### API CALL ${method} ${apiUrl}`)

      let headers = {}

      // Send token as per the OAuth 2.0 bearer token scheme
      if (userProfile.token) {
        headers = {
          'Authorization': `Bearer ${userProfile.token}`
        }
      }

      return axios({
        method: method,
        url: apiUrl,
        data: data,
        headers: headers
      })
    },

    apiDecodeError(err) {
      if (err.response && err.response.data && err.response.headers['content-type'].includes('json')) {
        // err.response.data.httpStatus = err.response.status
        // err.response.data.httpStatusText = err.response.statusText
        return err.response.data
      }
      if (err.response && err.request) {
        return `HTTP ${err.response.status}: API call failed: ${err.request.responseURL}`
      }
      return err.toString()
    }
  }
}