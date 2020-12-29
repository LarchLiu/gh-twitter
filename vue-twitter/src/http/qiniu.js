import { Octokit } from '@octokit/core'

const QiniuOctokit = Octokit.defaults({
  baseUrl: process.env.QINIU_DOMAIN
})

const octokit = new QiniuOctokit()

export default octokit
