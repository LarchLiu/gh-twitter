const GithubToken = 'gh-token'

export function getGHToken () {
  return window.localStorage.getItem(GithubToken)
}

export function setGHToken (token) {
  return window.localStorage.setItem(GithubToken, token)
}
