export default class User {
  constructor(token, account, userName) {
    this.token = token
    this.account = account
    this.userName = userName
    this.cart = []
  }
}