import API from '@/utils/request-api'
import QiniuOctokit from '@/http/qiniu'

const baseUrl = './twitter'

function getUsersData () {
  let jsonPath = '/json/userList.json'
  if (process.env.VUE_APP_PIC_BED === 'qiniu') {
    jsonPath += '?v=' + (new Date()).getTime()
    return new Promise((resolve, reject) => {
      QiniuOctokit.request('GET /' + process.env.VUE_APP_QINIU_RESOURCE_PREFIX + jsonPath).then(response => {
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
  let jsonPath = '/json/' + user + '/' + page + '.json'
  if (process.env.VUE_APP_PIC_BED === 'qiniu') {
    jsonPath += '?v=' + (new Date()).getTime()
    return new Promise((resolve, reject) => {
      QiniuOctokit.request('GET /' + process.env.VUE_APP_QINIU_RESOURCE_PREFIX + jsonPath).then(response => {
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
  let jsonPath = '/json/updateInfo.json'
  if (process.env.VUE_APP_PIC_BED === 'qiniu') {
    jsonPath += '?v=' + (new Date()).getTime()
    return new Promise((resolve, reject) => {
      QiniuOctokit.request('GET /' + process.env.VUE_APP_QINIU_RESOURCE_PREFIX + jsonPath).then(response => {
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
