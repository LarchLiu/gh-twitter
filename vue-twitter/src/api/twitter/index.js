import API from '@/utils/request-api'
import QiniuOctokit from '@/http/qiniu'

const baseUrl = './twitter'

function getUsersData () {
  const jsonPath = '/json/userList.json'
  if (process.env.VUE_APP_PIC_BED === 'qiniu') {
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
  const jsonPath = '/json/' + user + '/' + page + '.json'
  if (process.env.VUE_APP_PIC_BED === 'qiniu') {
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
  getTweetsData
}
