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
            <template v-slot:btn v-if="needUpdate">
              <a-button :disabled="triggerUpdate" @click="triggerUpdate = true">refresh</a-button>
            </template>
            <template v-slot:default v-if="usersList.length > 0">
              <div
                v-for="(user, i) in usersList"
                :key="i"
                @click="changeUser(i)"
              >
                <a :class="user.Username === currentUser ? 'current' : 'normal'">{{ user.Name }}</a>
              </div>
            </template>
          </aside-box>
        </div>
        <div v-if="usersList.length > 0 && Object.keys(usersData).length === usersList.length" class="tweets">
          <twitter
            v-if="currentUser === usersList[0].Username"
            class="detail"
            :isAll="true"
            :detail="usersData[currentUser]"
            :usersObj="usersListObj"
          />
          <twitter
            v-else
            class="detail"
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
import { ref, getCurrentInstance, onMounted, onUnmounted, watch } from 'vue'
import twitterApi from '@/api/twitter/index'
import AsideBox from '@/components/AsideBox/index.vue'
import Twitter from '@/components/Twitter/index.vue'
import Octokit from '@/http/qiniu'
import { arrToObj } from '@/utils/index'

export default {
  components: { AsideBox, Twitter },
  setup () {
    const { ctx } = getCurrentInstance()
    const updateUser = ref([])
    const usersList = ref([])
    const usersData = ref({})
    const usersListObj = ref({})
    const currentUser = ref('')
    const curPage = ref(1)
    const updateTime = ref(0)
    const needUpdate = ref(false)
    const triggerUpdate = ref(false)
    var timer

    const onExit = () => {
      ctx.$router.push({
        path: '/login'
      })
    }

    const webhook = () => {
      Octokit.request('POST /repos/LarchLiu/gh-twitter/dispatches', {
        event_type: 'scraper'
      })
      console.log(process.env.VUE_APP_GITHUB_TOKEN)
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
            updateUser.value = res.Users
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
    }

    const margeDetail = (tweet, profile) => {
      return { Avatar: profile.Avatar, Name: profile.Name, ...tweet }
    }

    onMounted(() => {
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
      webhook,
      margeDetail,
      usersListObj,
      updateUser,
      needUpdate,
      triggerUpdate,
      updateTime
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
