import API from '@/utils/request-api'

const baseUrl = './twitter/'

function getUsersData () {
  return API.getRequest(baseUrl + 'userList.txt')
}

function getTweetsData (user, page) {
  return API.getRequest(baseUrl + 'json/' + user + '/' + page + '.json')
}

export default {
  getUsersData,
  getTweetsData
}
