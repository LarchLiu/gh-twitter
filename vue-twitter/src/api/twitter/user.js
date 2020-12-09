import API from '/@/utils/request-api'

const base_url = '/twitter/'

function getUsersData() {
  return API.getRequest(base_url + 'userList.txt')
}

function getTweetsData(user) {
  return API.getRequest(base_url + user + '.json')
}

export default {
  getUsersData,
  getTweetsData
}
