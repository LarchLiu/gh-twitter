import API from '@/utils/request-api'
import QiniuOctokit from '@/http/qiniu'

const baseUrl = './twitter'
const qiniuPrefix = process.env.VUE_APP_QINIU_RESOURCE_PREFIX ? process.env.VUE_APP_QINIU_RESOURCE_PREFIX.replace(/^\/*|\/*$/g, '') : ''

function getUsersData () {
  const jsonPath = '/json/userList.json?v=' + (new Date()).getTime()
  if (process.env.VUE_APP_PIC_BED === 'qiniu') {
    // jsonPath += '?v=' + (new Date()).getTime()
    return new Promise((resolve, reject) => {
      const path = (qiniuPrefix + jsonPath).replace(/^\/*|\/*$/g, '')
      QiniuOctokit.request('GET /' + path).then(response => {
        const { data } = response
        resolve(data)
      }).catch(_ => {
        // reject(error)
        return API.getRequest(baseUrl + jsonPath)
      })
    })
  }
  return API.getRequest(baseUrl + jsonPath)
}

function getTweetsData (user, page) {
  const jsonPath = '/json/' + user + '/' + page + '.json?v=' + (new Date()).getTime()
  if (process.env.VUE_APP_PIC_BED === 'qiniu') {
    // jsonPath += '?v=' + (new Date()).getTime()
    return new Promise((resolve, reject) => {
      const path = (qiniuPrefix + jsonPath).replace(/^\/*|\/*$/g, '')
      QiniuOctokit.request('GET /' + path).then(response => {
        const { data } = response
        resolve(data)
      }).catch(_ => {
        // reject(error)
        return API.getRequest(baseUrl + jsonPath)
      })
    })
  }
  return API.getRequest(baseUrl + jsonPath)
}

function getUpdateInfo () {
  const jsonPath = '/json/updateInfo.json?v=' + (new Date()).getTime()
  if (process.env.VUE_APP_PIC_BED === 'qiniu') {
    // jsonPath += '?v=' + (new Date()).getTime()
    return new Promise((resolve, reject) => {
      const path = (qiniuPrefix + jsonPath).replace(/^\/*|\/*$/g, '')
      QiniuOctokit.request('GET /' + path).then(response => {
        const { data } = response
        resolve(data)
      }).catch(_ => {
        // reject(error)
        return API.getRequest(baseUrl + jsonPath)
      })
    })
  }
  return API.getRequest(baseUrl + jsonPath)
}

export default {
  getUsersData,
  getTweetsData,
  getUpdateInfo
}
