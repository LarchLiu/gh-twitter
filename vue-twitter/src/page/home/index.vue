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
            <div v-if="usersListSort">
              <div
                v-for="(user, i) in usersListSort"
                :key="i"
                @click="changeUser(i)"
              >
                <a :class="i === currentUser ? 'current' : 'normal'">{{ user }}</a>
              </div>
            </div>
          </aside-box>
        </div>
        <twitter
          class="detail"
          :detail="usersData ? usersData[currentUser] : {}"
        />
      </div>
    </div>
  </div>
</template>
<script>
import { ref, getCurrentInstance, onMounted, watch } from 'vue'
import twitterApi from '@/api/twitter/index'
import AsideBox from '@/components/AsideBox/index.vue'
import Twitter from '@/components/Twitter/index.vue'
import Octokit from '@/http/github'

export default {
  components: { AsideBox, Twitter },
  setup () {
    const { ctx } = getCurrentInstance()
    const usersList = ref([])
    const usersData = ref([])
    const usersListSort = ref([])
    const currentUser = ref(0)

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

    const getUserList = () => {
      twitterApi.getUsersData().then(res => {
        usersList.value = res.replace(/\s*/g, '').split(',')
      }).catch(e => {
        console.log(e)
        usersList.value = []
      })
    }

    const getUserTweets = (user) => {
      twitterApi.getTweetsData(user).then(data => {
        usersData.value.push(data)
        usersListSort.value.push(data.Profile.Name)
        // console.log(data)
      }).catch(err => {
        console.log(err)
      })
    }

    const changeUser = (i) => {
      currentUser.value = i
    }

    onMounted(() => {
      getUserList()
    })

    watch(usersList, () => {
      usersData.value = []
      for (let i = 0; i < usersList.value.length; i++) {
        getUserTweets(usersList.value[i])
      }
    })

    return {
      currentUser,
      usersList,
      usersData,
      usersListSort,
      getUserTweets,
      getUserList,
      changeUser,
      onExit,
      webhook
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
