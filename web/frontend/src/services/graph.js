// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Set of methods to call the Graph API, using REST and fetch
// ----------------------------------------------------------------------------

import auth from './auth'

const GRAPH_BASE = 'https://graph.microsoft.com/beta'
const GRAPH_SCOPES = [ 'user.read', 'user.readbasic.all' ]

let accessToken

export default {
  //
  // Get details of user, and return as JSON
  // https://docs.microsoft.com/en-us/graph/api/user-get?view=graph-rest-1.0&tabs=http#response-1
  //
  async getSelf() {
    let resp = await callGraph('/me')
    if (resp) {
      let data = await resp.json()
      return data
    }
  },

  //
  // Get user's photo and return as a blob object URL
  // https://developer.mozilla.org/en-US/docs/Web/API/URL/createObjectURL
  //
  async getPhoto() {
    let resp = await callGraph('/me/photos/240x240/$value')
    if (resp) {
      let blob = await resp.blob()
      return URL.createObjectURL(blob)
    }
  },

  //
  // Search for users
  // https://developer.mozilla.org/en-US/docs/Web/API/URL/createObjectURL
  //
  async searchUsers(searchString, max = 50) {
    let resp = await callGraph(`/users?$filter=startswith(displayName, '${searchString}') or startswith(userPrincipalName, '${searchString}')&$top=${max}`)
    if (resp) {
      let data = await resp.json()
      return data
    }
  },

  //
  // Accessor for access token, only included for demo purposes
  //
  getAccessToken() {
    return accessToken
  }
}

//
// Common fetch wrapper (private)
//
async function callGraph(apiPath) {
  if (!auth.clientId() || !auth.user()) { return }

  // Acquire an access token to call APIs (like Graph)
  // Safe to call repeatedly as MSAL caches tokens locally
  accessToken = await auth.acquireToken(GRAPH_SCOPES)

  let resp = await fetch(
    `${GRAPH_BASE}${apiPath}`,
    {
      headers: { authorization: `bearer ${accessToken}` }
    }
  )

  if (!resp.ok) { throw new Error(`Call to ${GRAPH_BASE}${apiPath} failed: ${resp.statusText}`) }

  return resp
}