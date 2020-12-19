import API from '@/utils/request-api'

const baseUrl = './twitter/'

function getUsersData () {
  return API.getRequest(baseUrl + 'userList.txt')
}

function getTweetsData (user) {
  return API.getRequest(baseUrl + user + '.json')
}

export default {
  getUsersData,
  getTweetsData
}
