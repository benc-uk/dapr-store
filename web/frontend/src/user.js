// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Dapr Store frontend - Simple data class to hold user profile data
// ----------------------------------------------------------------------------

export default class User {
  constructor(token, account, userName) {
    this.token = token
    this.account = account
    this.userName = userName
    this.cart = []
  }
}