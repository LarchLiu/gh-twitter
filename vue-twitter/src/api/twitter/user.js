import API from '/@/utils/request-api'

const base_url = '/twitter/userList.txt'

function getUsersData() {
  return API.getRequest(base_url)
}

export default {
  getUsersData
}
