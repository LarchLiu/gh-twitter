import service from './request'

const base = ''

const postRequest = (url, payload = {}) => {
  return new Promise((resolve, reject) => {
    return service({
      method: 'post',
      url: base + url,
      data: payload
    }).then(response => {
      resolve(response)
    }).catch(error => {
      reject(error)
    })
  })
}
const uploadFileRequest = (url, payload = {}) => {
  return new Promise((resolve, reject) => {
    service({
      method: 'post',
      url: base + url,
      data: payload,
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    }).then(response => {
      resolve(response)
    }).catch(error => {
      reject(error)
    })
  })
}
const patchRequest = (url, payload = {}) => {
  return new Promise((resolve, reject) => {
    service({
      method: 'patch',
      url: base + url,
      data: payload
    }).then(response => {
      resolve(response)
    }).catch(error => {
      reject(error)
    })
  })
}
const deleteRequest = (url, payload = {}) => {
  return new Promise((resolve, reject) => {
    service({
      method: 'delete',
      url: base + url,
      data: payload
    }).then(response => {
      resolve(response)
    }).catch(error => {
      reject(error)
    })
  })
}
const getRequest = (url, payload = {}) => {
  return new Promise((resolve, reject) => {
    service({
      method: 'get',
      url: base + url,
      params: payload
    }).then(response => {
      resolve(response)
    }).catch(error => {
      reject(error)
    })
  })
}
export default {
  getRequest,
  postRequest,
  deleteRequest,
  patchRequest,
  uploadFileRequest
}
