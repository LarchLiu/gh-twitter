<template>
  <div class="container">
    <div class="w1210">
      <div class="bd">
        <div class="aside-wrap">
          <aside-box
            title="GH Twitter"
            :need-fixed="true"
            id-name="header"
          >
            <template v-slot:btn>
              <span style="float: right">
                <PlusOutlined v-if="ghToken && dbInit" style="margin-right: 4px" @click="addUserVisible=true" />
                <SyncOutlined v-if="ghToken" style="margin-right: 4px" @click="actionScraper" />
                <SettingFilled @click="visible = true"/>
                <a-modal
                  :visible="visible"
                  title="Add Token"
                  ok-text="确认"
                  cancel-text="取消"
                  @cancel="cancelInput"
                  @ok="setGHToken"
                >
                  <p>请输入具有 repo 权限的 github personal access token, 切勿泄露 token</p>
                  <!-- eslint-disable-next-line -->
                  <a-input v-model:value="inputToken" placeholder="github personal access token" />
                </a-modal>
                <a-modal
                  :visible="addUserVisible"
                  title="Add Users"
                  ok-text="确认"
                  cancel-text="取消"
                  @cancel="cancelUserInput"
                  @ok="addUsersAction"
                >
                  <p>请输入用户名，添加多个用户以空格分割</p>
                  <!-- eslint-disable-next-line -->
                  <a-input v-model:value="inputUsers" placeholder="@username" />
                </a-modal>
              </span>
            </template>
            <template v-slot:default v-if="usersList.length > 0">
              <div
                v-for="(user, i) in usersList"
                :key="i"
                @click="changeUser(i)"
              >
                <a-badge :dot="needUpdate && (updateUser.findIndex(e => e === user.Username) >= 0)">
                  <a :class="user.Username === currentUser ? 'current' : 'normal'">{{ user.Name }}</a>
                </a-badge>
              </div>
            </template>
          </aside-box>
        </div>
        <div v-if="usersList.length > 0 && Object.keys(usersData).length === usersList.length" class="tweets">
          <twitter
            v-if="currentUser === usersList[0].Username"
            class="detail"
            :updateTime="updateTime"
            :isAll="true"
            :detail="usersData[currentUser]"
            :usersObj="usersListObj"
          />
          <twitter
            v-else
            class="detail"
            :updateTime="updateTime"
            :detail="usersData[currentUser]"
          />
        </div>
        <twitter
            v-else
            class="detail"
          />
      </div>
    </div>
  </div>
</template>
<script>
import { ref, getCurrentInstance, computed, onMounted, onUnmounted, watch } from 'vue'
import { useStore } from 'vuex'
import twitterApi from '@/api/twitter/index'
import AsideBox from '@/components/AsideBox/index.vue'
import Twitter from '@/components/Twitter/index.vue'
import { Octokit } from '@octokit/core'
import { arrToObj, uniqueArr } from '@/utils/index'
import { message } from 'ant-design-vue'
import {
  PlusOutlined,
  SettingFilled,
  SyncOutlined
} from '@ant-design/icons-vue'

export default {
  components: { AsideBox, PlusOutlined, SettingFilled, SyncOutlined, Twitter },
  setup () {
    const { ctx } = getCurrentInstance()
    const updateUser = ref([])
    const usersList = ref([])
    const usersData = ref({})
    const usersListObj = ref({})
    const currentUser = ref('')
    const inputToken = ref('')
    const inputUsers = ref('')
    const curPage = ref(1)
    const updateTime = ref(0)
    const needUpdate = ref(false)
    const triggerUpdate = ref(false)
    const visible = ref(false)
    const addUserVisible = ref(false)
    const store = useStore()
    const ghToken = computed(() => store.getters.gh_token)
    const repoUrl = process.env.VUE_APP_REPO_URL
    const dbInit = process.env.VUE_APP_DB_INIT
    let timer
    let ghApi

    const onExit = () => {
      ctx.$router.push({
        path: '/login'
      })
    }

    const actionScraper = () => {
      ghApi.request(`POST ${repoUrl}`, {
        event_type: 'scraper'
      }).then(res => {
        message.success({
          content: '更新请求已发出，请等待3-4分钟',
          duration: 3
        })
      }).catch(err => {
        console.log(err)
      })
    }

    const actionChangeUser = (type, users) => {
      ghApi.request(`POST ${repoUrl}`, {
        event_type: type,
        client_payload: {
          users: users
        }
      }).then(res => {
        message.success({
          content: '添加用户请求已发出，请等待3-4分钟',
          duration: 3
        })
      }).catch(err => {
        console.log(err)
      })
    }

    const setGHToken = () => {
      visible.value = false
      store.dispatch('setGHToken', inputToken.value)
      inputToken.value = ''
    }

    const cancelInput = () => {
      visible.value = false
      inputToken.value = ''
    }

    const addUsersAction = () => {
      addUserVisible.value = false
      const users = inputUsers.value.replace(/@/g, '').replace(/\s+/g, ',')
      console.log(users)
      actionChangeUser('addusers', users)
      inputUsers.value = ''
    }

    const cancelUserInput = () => {
      addUserVisible.value = false
      inputUsers.value = ''
    }

    const dataInit = () => {
      twitterApi.getUsersData().then(res => {
        updateUser.value = res.map(item => {
          return item.Username
        })
        usersList.value = res
        usersListObj.value = arrToObj(res, 'Username')
        currentUser.value = res[0].Username
        twitterApi.getUpdateInfo().then(info => {
          updateTime.value = info.UpdateTime
        }).catch(e => {
          console.log(e)
        })
        for (let i = 0; i < updateUser.value.length; i++) {
          getUserTweets(updateUser.value[i], 1)
        }
      }).catch(e => {
        console.log(e)
        usersList.value = []
      })
    }

    const getUpdateInfo = () => {
      twitterApi.getUpdateInfo().then(res => {
        if (updateTime.value < res.UpdateTime) {
          updateTime.value = res.UpdateTime
          if (res.IsUpdate) {
            needUpdate.value = res.IsUpdate
            updateUser.value = updateUser.value.concat(res.Users)
            updateUser.value = uniqueArr(updateUser.value)
            // updateUser.value = res.Users
          }
        }
      }).catch(e => {
        console.log(e)
        usersList.value = []
      })
    }

    const getUserList = () => {
      twitterApi.getUsersData().then(res => {
        usersList.value = res
        usersListObj.value = arrToObj(res, 'Username')
        // currentUser.value = res[0].Username
      }).catch(e => {
        console.log(e)
        usersList.value = []
      })
    }

    const getUserTweets = (user, page) => {
      twitterApi.getTweetsData(user, page).then(data => {
        usersData.value[user] = data
        updateUser.value.splice(updateUser.value.findIndex(e => e === user), 1)
        if (updateUser.value.length === 0) {
          needUpdate.value = false
          triggerUpdate.value = false
        }
        // usersListSort.value.push(data.Profile.Name)
        // console.log(data)
      }).catch(err => {
        console.log(err)
      })
    }

    const changeUser = (i) => {
      currentUser.value = usersList.value[i].Username
      curPage.value = 1
      if (updateUser.value.findIndex(e => e === currentUser.value) >= 0) {
        getUserTweets(currentUser.value, 1)
      }
    }

    const margeDetail = (tweet, profile) => {
      return { Avatar: profile.Avatar, Name: profile.Name, ...tweet }
    }

    onMounted(() => {
      ghApi = new Octokit({
        auth: ghToken.value
      })
      dataInit()
      timer = window.setInterval(() => {
        setTimeout(() => {
          getUpdateInfo()
        }, 0)
      }, 60000)
    })

    onUnmounted(() => {
      window.clearInterval(timer)
    })

    // watch(usersList, () => {
    //   usersData.value = {}
    //   for (let i = 0; i < usersList.value.length; i++) {
    //     getUserTweets(usersList.value[i].Username, curPage.value)
    //   }
    // })

    watch([needUpdate, triggerUpdate], () => {
      if (needUpdate.value && triggerUpdate.value) {
        getUserList()
        for (let i = 0; i < updateUser.value.length; i++) {
          getUserTweets(updateUser.value[i], 1)
        }
      }
    })

    watch([needUpdate, updateUser], () => {
      if (needUpdate.value) {
        let hasUser = true
        for (let i = 0; i < updateUser.value.length; i++) {
          let flag = false
          for (let j = 0; j < usersList.value.length; j++) {
            if (updateUser.value[i] === usersList.value[j].Username) {
              flag = true
              break
            }
          }
          if (!flag) {
            hasUser = false
            break
          }
        }
        if (!hasUser) {
          getUserList()
        }
      }
    })

    watch(ghToken, () => {
      ghApi = new Octokit({
        auth: ghToken.value
      })
    })

    return {
      currentUser,
      curPage,
      usersList,
      usersData,
      getUserTweets,
      getUserList,
      getUpdateInfo,
      changeUser,
      onExit,
      actionScraper,
      actionChangeUser,
      margeDetail,
      usersListObj,
      updateUser,
      needUpdate,
      triggerUpdate,
      updateTime,
      ghToken,
      setGHToken,
      visible,
      inputToken,
      cancelInput,
      addUserVisible,
      inputUsers,
      addUsersAction,
      cancelUserInput,
      dbInit
    }
  }
}
</script>
<style lang="less" scoped>
  .container {
    font-size: 12px;
    color: #666;

    .w1210 {
      width: 820px;
      margin: 0 auto;
    }

    .bd {
      margin-top: 2px;
      overflow: hidden;
      .aside-wrap {
        float: left;
        width: 210px;
        font-size: 15px;
        font-weight: 500;
        margin-bottom: 20px;
        :deep(.current) {
          color: rgb(29, 161, 242);
        }
        :deep(.normal) {
          color: black;
        }
      }
      .detail {
        float: right;
        width: 600px;
        min-height: 220px;
      }
    }
  }
</style>
