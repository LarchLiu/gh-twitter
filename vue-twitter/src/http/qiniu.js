import { Octokit } from '@octokit/core'

const QiniuOctokit = Octokit.defaults({
  baseUrl: process.env.VUE_APP_QINIU_DOMAIN
})

const octokit = new QiniuOctokit()

export default octokit
